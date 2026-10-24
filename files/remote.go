package files

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/domainry/domainry-foundation/telemetry"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
)

const (
	discoveryPath  = "/knowledge/v1/discovery"
	invokePath     = "/knowledge/v1/invoke"
	runtimeHeader  = "X-Domainry-Runtime-ID"
	protocolV1     = "knowledge.saas.v1"
	deploymentSaaS = "saas"
)

type RemoteConfig struct {
	Endpoint           string
	ServiceAccessToken string
	HTTPClient         *http.Client
	RequestTimeout     time.Duration
	MaxRequestBytes    int64
	MaxResponseBytes   int64
}

func OpenRemote(ctx context.Context, config RemoteConfig, runtimeID string) (Service, error) {
	client, err := newRemoteClient(config, runtimeID)
	if err != nil {
		return nil, err
	}
	var descriptor remoteDescriptor
	if err = client.request(ctx, http.MethodGet, discoveryPath, nil, &descriptor); err != nil {
		return nil, err
	}
	if descriptor.ProtocolVersion != protocolV1 || descriptor.Mode != deploymentSaaS || descriptor.Audience != runtimeID || !slices.Contains(descriptor.Capabilities, "files.read") || !slices.Contains(descriptor.Capabilities, "files.write") {
		return nil, &Error{Class: "unavailable", Code: "knowledge.protocol_incompatible"}
	}
	return client, nil
}

type remoteDescriptor struct {
	ProtocolVersion string   `json:"protocol_version"`
	Mode            string   `json:"mode"`
	Audience        string   `json:"audience"`
	Capabilities    []string `json:"capabilities"`
}

type remoteRequest struct {
	Operation string          `json:"operation"`
	Input     json.RawMessage `json:"input,omitempty"`
}

type remoteResponse struct {
	Result json.RawMessage `json:"result,omitempty"`
	Error  *remoteError    `json:"error,omitempty"`
}

type remoteError struct {
	Class     string `json:"class"`
	Code      string `json:"code"`
	Message   string `json:"message,omitempty"`
	Retryable bool   `json:"retryable,omitempty"`
}

type remoteClient struct {
	baseURL       *url.URL
	runtimeID     string
	token         string
	httpClient    *http.Client
	requestLimit  int64
	responseLimit int64
}

func newRemoteClient(config RemoteConfig, runtimeID string) (*remoteClient, error) {
	runtimeID = strings.TrimSpace(runtimeID)
	base, err := url.Parse(strings.TrimSpace(config.Endpoint))
	if err != nil || base == nil || base.Scheme != "http" && base.Scheme != "https" || strings.TrimSpace(base.Host) == "" || base.RawQuery != "" || base.Fragment != "" || runtimeID == "" {
		return nil, fmt.Errorf("Knowledge file endpoint or Runtime ID is invalid")
	}
	base.Path = strings.TrimRight(base.Path, "/")
	if config.RequestTimeout <= 0 {
		config.RequestTimeout = 20 * time.Second
	}
	if config.MaxRequestBytes <= 0 {
		config.MaxRequestBytes = 40 << 20
	}
	if config.MaxResponseBytes <= 0 {
		config.MaxResponseBytes = 40 << 20
	}
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: config.RequestTimeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	}
	return &remoteClient{baseURL: base, runtimeID: runtimeID, token: strings.TrimSpace(config.ServiceAccessToken), httpClient: telemetry.InstrumentHTTPClient(httpClient, "github.com/domainry/domainry-knowledge-sdk/files"), requestLimit: config.MaxRequestBytes, responseLimit: config.MaxResponseBytes}, nil
}

func (client *remoteClient) Upload(ctx context.Context, authority Authority, input Upload) (File, error) {
	digest := sha256.Sum256(input.Data)
	return client.UploadStream(ctx, authority, StreamUpload{
		ClientID: input.ClientID, Filename: input.Filename, ContentType: input.ContentType,
		SizeBytes: int64(len(input.Data)), SHA256: hex.EncodeToString(digest[:]), Content: bytes.NewReader(input.Data),
	})
}

func (client *remoteClient) UploadStream(ctx context.Context, authority Authority, input StreamUpload) (File, error) {
	if ctx == nil || input.Content == nil || input.SizeBytes < 1 || input.SizeBytes > MaxContentBytes || len(input.SHA256) != 64 {
		return File{}, &Error{Class: ErrorClassBadRequest, Code: "knowledge.file_invalid"}
	}
	if rawDigest, err := hex.DecodeString(input.SHA256); err != nil || len(rawDigest) != sha256.Size || input.SHA256 != strings.ToLower(input.SHA256) {
		return File{}, &Error{Class: ErrorClassBadRequest, Code: "knowledge.file_invalid"}
	}
	metadata, err := json.Marshal(struct {
		Authority Authority    `json:"authority"`
		Input     StreamUpload `json:"input"`
	}{authority, input})
	if err != nil || len(metadata) == 0 || len(metadata) > 256<<10 {
		return File{}, &Error{Class: ErrorClassBadRequest, Code: "knowledge.request_invalid", Cause: err}
	}

	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)
	writeDone := make(chan error, 1)
	go func() {
		var writeErr error
		defer func() {
			if writeErr == nil {
				writeErr = multipartWriter.Close()
			}
			_ = writer.CloseWithError(writeErr)
			writeDone <- writeErr
		}()
		metadataPart, partErr := multipartWriter.CreateFormField("metadata")
		if partErr != nil {
			writeErr = partErr
			return
		}
		if _, partErr = metadataPart.Write(metadata); partErr != nil {
			writeErr = partErr
			return
		}
		contentPart, partErr := multipartWriter.CreateFormFile("content", input.Filename)
		if partErr != nil {
			writeErr = partErr
			return
		}
		limited := &io.LimitedReader{R: input.Content, N: input.SizeBytes + 1}
		written, partErr := io.CopyBuffer(contentPart, limited, make([]byte, 128*1024))
		if partErr != nil {
			writeErr = partErr
			return
		}
		if written != input.SizeBytes || limited.N != 1 {
			writeErr = &Error{Class: ErrorClassBadRequest, Code: "knowledge.file_size_mismatch"}
		}
	}()

	endpoint := *client.baseURL
	endpoint.Path = strings.TrimRight(client.baseURL.Path, "/") + knowledgecontract.SaaSFileUploadPath
	requestContext := telemetry.WithHTTPClientOperation(ctx, "knowledge.files.upload_stream")
	request, err := http.NewRequestWithContext(requestContext, http.MethodPost, endpoint.String(), reader)
	if err != nil {
		_ = reader.CloseWithError(err)
		return File{}, &Error{Class: ErrorClassUnavailable, Code: "knowledge.request_invalid", Cause: err}
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "domainry-knowledge-files-sdk-go")
	request.Header.Set(runtimeHeader, client.runtimeID)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	if client.token != "" {
		request.Header.Set("Authorization", "Bearer "+client.token)
	}
	httpClient := *client.httpClient
	httpClient.Timeout = 0
	response, err := httpClient.Do(request)
	if err != nil {
		_ = reader.CloseWithError(err)
		select {
		case <-writeDone:
		default:
		}
		return File{}, &Error{Class: ErrorClassUnavailable, Code: "knowledge.remote_unavailable", Retryable: true, Cause: err}
	}
	defer response.Body.Close()
	var output File
	responseErr := client.decodeStreamResponse(response, &output)
	writeErr := <-writeDone
	if responseErr != nil {
		return File{}, responseErr
	}
	if writeErr != nil {
		return File{}, writeErr
	}
	return output, nil
}

func (client *remoteClient) decodeStreamResponse(response *http.Response, output any) error {
	raw, err := io.ReadAll(io.LimitReader(response.Body, client.responseLimit+1))
	if err != nil {
		return &Error{Class: ErrorClassUnavailable, Code: "knowledge.response_read_failed", Retryable: true, Cause: err}
	}
	if int64(len(raw)) > client.responseLimit {
		return &Error{Class: ErrorClassUnavailable, Code: "knowledge.response_too_large"}
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var failure struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if json.Unmarshal(raw, &failure) == nil && strings.TrimSpace(failure.Code) != "" {
			return &Error{Class: statusClass(response.StatusCode), Code: failure.Code, Message: failure.Message, Retryable: response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= 500}
		}
		return &Error{Class: statusClass(response.StatusCode), Code: "knowledge.remote_request_failed", Retryable: response.StatusCode >= 500}
	}
	var envelope remoteResponse
	if json.Unmarshal(raw, &envelope) != nil || envelope.Error != nil {
		if envelope.Error != nil {
			return &Error{Class: safeErrorClass(envelope.Error.Class), Code: envelope.Error.Code, Message: envelope.Error.Message, Retryable: envelope.Error.Retryable}
		}
		return &Error{Class: ErrorClassUnavailable, Code: "knowledge.response_invalid"}
	}
	decoder := json.NewDecoder(bytes.NewReader(envelope.Result))
	decoder.DisallowUnknownFields()
	if decoder.Decode(output) != nil || decoder.Decode(&struct{}{}) != io.EOF {
		return &Error{Class: ErrorClassUnavailable, Code: "knowledge.response_invalid"}
	}
	return nil
}

func (client *remoteClient) Get(ctx context.Context, authority Authority, id string) (File, error) {
	var output File
	err := client.invoke(ctx, "files.get", struct {
		Authority Authority `json:"authority"`
		ID        string    `json:"id"`
	}{authority, id}, &output)
	return output, err
}

func (client *remoteClient) Download(ctx context.Context, authority Authority, id string) (Download, error) {
	var output Download
	err := client.invoke(ctx, "files.download", struct {
		Authority Authority `json:"authority"`
		ID        string    `json:"id"`
	}{authority, id}, &output)
	return output, err
}

func (client *remoteClient) Bind(ctx context.Context, authority Authority, id string, binding Binding) error {
	return client.invoke(ctx, "files.bind", struct {
		Authority Authority `json:"authority"`
		ID        string    `json:"id"`
		Binding   Binding   `json:"binding"`
	}{authority, id, binding}, nil)
}

func (client *remoteClient) Delete(ctx context.Context, authority Authority, id string) error {
	return client.invoke(ctx, "files.delete", struct {
		Authority Authority `json:"authority"`
		ID        string    `json:"id"`
	}{authority, id}, nil)
}

func (client *remoteClient) invoke(ctx context.Context, operation string, input, output any) error {
	raw, err := json.Marshal(input)
	if err != nil {
		return &Error{Class: "bad_request", Code: "knowledge.request_invalid", Cause: err}
	}
	var response remoteResponse
	if err = client.request(ctx, http.MethodPost, invokePath, remoteRequest{Operation: operation, Input: raw}, &response); err != nil {
		return err
	}
	if response.Error != nil {
		return &Error{Class: safeErrorClass(response.Error.Class), Code: response.Error.Code, Message: response.Error.Message, Retryable: response.Error.Retryable}
	}
	if output == nil || len(response.Result) == 0 {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(response.Result))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(output); err != nil || decoder.Decode(&struct{}{}) != io.EOF {
		return &Error{Class: "unavailable", Code: "knowledge.response_invalid", Cause: err}
	}
	return nil
}

func (client *remoteClient) request(ctx context.Context, method, requestPath string, input, output any) error {
	if ctx == nil {
		return &Error{Class: "unavailable", Code: "knowledge.context_required"}
	}
	var body io.Reader
	if input != nil {
		raw, err := json.Marshal(input)
		if err != nil || int64(len(raw)) > client.requestLimit {
			return &Error{Class: "bad_request", Code: "knowledge.request_too_large", Cause: err}
		}
		body = bytes.NewReader(raw)
	}
	endpoint := *client.baseURL
	endpoint.Path = strings.TrimRight(client.baseURL.Path, "/") + "/" + strings.TrimLeft(requestPath, "/")
	ctx = telemetry.WithHTTPClientOperation(ctx, "knowledge.files.remote")
	request, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return &Error{Class: "unavailable", Code: "knowledge.request_invalid", Cause: err}
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "domainry-knowledge-files-sdk-go")
	request.Header.Set(runtimeHeader, client.runtimeID)
	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if client.token != "" {
		request.Header.Set("Authorization", "Bearer "+client.token)
	}
	response, err := client.httpClient.Do(request)
	if err != nil {
		return &Error{Class: "unavailable", Code: "knowledge.remote_unavailable", Retryable: true, Cause: err}
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, client.responseLimit+1))
	if err != nil {
		return &Error{Class: "unavailable", Code: "knowledge.response_read_failed", Retryable: true, Cause: err}
	}
	if int64(len(raw)) > client.responseLimit {
		return &Error{Class: "unavailable", Code: "knowledge.response_too_large"}
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var payload struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if json.Unmarshal(raw, &payload) == nil && strings.TrimSpace(payload.Code) != "" {
			return &Error{Class: statusClass(response.StatusCode), Code: payload.Code, Message: payload.Message, Retryable: response.StatusCode >= 500}
		}
		return &Error{Class: "unavailable", Code: "knowledge.remote_request_failed", Retryable: response.StatusCode >= 500}
	}
	if output == nil {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(output); err != nil || decoder.Decode(&struct{}{}) != io.EOF {
		return &Error{Class: "unavailable", Code: "knowledge.response_invalid", Cause: err}
	}
	return nil
}

func statusClass(status int) ErrorClass {
	switch status {
	case http.StatusBadRequest, http.StatusRequestEntityTooLarge:
		return ErrorClassBadRequest
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrorClassForbidden
	case http.StatusNotFound:
		return ErrorClassNotFound
	case http.StatusConflict:
		return ErrorClassConflict
	case http.StatusTooManyRequests:
		return ErrorClassRateLimited
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return ErrorClassUnavailable
	default:
		return ErrorClassInternal
	}
}

func safeErrorClass(class string) ErrorClass {
	switch strings.TrimSpace(class) {
	case string(ErrorClassBadRequest):
		return ErrorClassBadRequest
	case string(ErrorClassForbidden):
		return ErrorClassForbidden
	case string(ErrorClassNotFound):
		return ErrorClassNotFound
	case string(ErrorClassConflict):
		return ErrorClassConflict
	case string(ErrorClassRateLimited):
		return ErrorClassRateLimited
	case string(ErrorClassUnavailable):
		return ErrorClassUnavailable
	case string(ErrorClassInternal):
		return ErrorClassInternal
	default:
		return ErrorClassInternal
	}
}

var _ Service = (*remoteClient)(nil)
