package remote

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/domainry/domainry-foundation/telemetry"
	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	"github.com/domainry/domainry-knowledge-sdk/saashost"
)

const (
	streamMetadataLimit = 256 << 10
	remoteRangeWindow   = 8 << 20
)

func (c *client) endpoint(path string) string {
	endpoint := *c.baseURL
	endpoint.Path = strings.TrimRight(c.baseURL.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return endpoint.String()
}

func (c *client) authorizeRequest(request *http.Request) {
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "domainry-knowledge-sdk-go")
	request.Header.Set(saashost.RuntimeIDHeader, c.runtimeID)
	if c.token != "" {
		request.Header.Set("Authorization", "Bearer "+c.token)
	}
}

func (c *client) streamHTTPClient() *http.Client {
	copy := *c.httpClient
	copy.Timeout = 0
	return &copy
}

func (c *client) uploadStream(ctx context.Context, path string, metadata any, content io.Reader, expectedBytes int64, output any) error {
	if ctx == nil || content == nil || expectedBytes < 1 || expectedBytes > knowledgesdk.KnowledgeDocumentMaximumBytes {
		return remoteError("bad_request", "knowledge.document_content_invalid", false, nil)
	}
	raw, err := json.Marshal(metadata)
	if err != nil || len(raw) == 0 || len(raw) > streamMetadataLimit {
		return remoteError("bad_request", "knowledge.request_invalid", false, err)
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
		if _, partErr = metadataPart.Write(raw); partErr != nil {
			writeErr = partErr
			return
		}
		contentPart, partErr := multipartWriter.CreateFormFile("content", "document")
		if partErr != nil {
			writeErr = partErr
			return
		}
		limited := &io.LimitedReader{R: content, N: expectedBytes + 1}
		written, partErr := io.CopyBuffer(contentPart, limited, make([]byte, 128*1024))
		if partErr != nil {
			writeErr = partErr
			return
		}
		if written != expectedBytes || limited.N != 1 {
			writeErr = remoteError("bad_request", "knowledge.document_size_mismatch", false, nil)
		}
	}()

	ctx = telemetry.WithHTTPClientOperation(ctx, "knowledge.document_upload")
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(path), reader)
	if err != nil {
		_ = reader.CloseWithError(err)
		return remoteError("unavailable", "knowledge.request_invalid", false, err)
	}
	c.authorizeRequest(request)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	response, err := c.streamHTTPClient().Do(request)
	if err != nil {
		_ = reader.CloseWithError(err)
		select {
		case <-writeDone:
		default:
		}
		return remoteError("unavailable", "knowledge.remote_unavailable", true, err)
	}
	defer response.Body.Close()
	if err = decodeStreamResponse(response, c.responseLimit, output); err != nil {
		return err
	}
	if writeErr := <-writeDone; writeErr != nil {
		return writeErr
	}
	return nil
}

func decodeStreamResponse(response *http.Response, limit int64, output any) error {
	raw, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return remoteError("unavailable", "knowledge.response_read_failed", true, err)
	}
	if int64(len(raw)) > limit {
		return remoteError("unavailable", "knowledge.response_too_large", false, nil)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var failure struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if json.Unmarshal(raw, &failure) == nil && strings.TrimSpace(failure.Code) != "" {
			return &knowledgesdk.Error{Class: statusClass(response.StatusCode), Code: safeCode(failure.Code), Message: safeMessage(failure.Message), Retryable: response.StatusCode >= 500}
		}
		return remoteError(statusClass(response.StatusCode), "knowledge.remote_request_failed", response.StatusCode >= 500, nil)
	}
	var envelope knowledgecontract.SaaSResponse
	if json.Unmarshal(raw, &envelope) != nil {
		return remoteError("unavailable", "knowledge.response_invalid", false, nil)
	}
	if envelope.Error != nil {
		return &knowledgesdk.Error{Class: safeClass(envelope.Error.Class), Code: safeCode(envelope.Error.Code), Message: safeMessage(envelope.Error.Message), Retryable: envelope.Error.Retryable}
	}
	if output == nil {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(envelope.Result))
	decoder.DisallowUnknownFields()
	if decoder.Decode(output) != nil || decoder.Decode(&struct{}{}) != io.EOF {
		return remoteError("unavailable", "knowledge.response_invalid", false, nil)
	}
	return nil
}

func encodeTransferHeader(value knowledgecontract.SaaSDocumentDownloadRequest) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil || len(raw) == 0 || len(raw) > streamMetadataLimit {
		return "", remoteError("bad_request", "knowledge.request_invalid", false, err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

type rangeDocumentReader struct {
	client         *client
	ctx            context.Context
	document       knowledgesdk.KnowledgeDocument
	transferHeader string
	offset         int64
	responseEnd    int64
	body           io.ReadCloser
	closed         bool
}

func (reader *rangeDocumentReader) Read(buffer []byte) (int, error) {
	if reader.closed {
		return 0, remoteError("unavailable", "knowledge.document_reader_closed", false, nil)
	}
	if len(buffer) == 0 {
		return 0, nil
	}
	if reader.offset >= reader.document.Bytes {
		return 0, io.EOF
	}
	if reader.body == nil {
		if err := reader.openRange(); err != nil {
			return 0, err
		}
	}
	n, err := reader.body.Read(buffer)
	reader.offset += int64(n)
	if err == io.EOF || reader.offset >= reader.responseEnd {
		_ = reader.body.Close()
		reader.body = nil
		if n > 0 {
			err = nil
		} else if reader.offset < reader.document.Bytes {
			err = nil
		}
	}
	return n, err
}

func (reader *rangeDocumentReader) openRange() error {
	end := reader.offset + remoteRangeWindow - 1
	if end >= reader.document.Bytes {
		end = reader.document.Bytes - 1
	}
	requestContext := telemetry.WithHTTPClientOperation(reader.ctx, "knowledge.document_download")
	request, err := http.NewRequestWithContext(requestContext, http.MethodGet, reader.client.endpoint(knowledgecontract.SaaSDocumentDownloadPath), nil)
	if err != nil {
		return remoteError("unavailable", "knowledge.request_invalid", false, err)
	}
	reader.client.authorizeRequest(request)
	request.Header.Set(knowledgecontract.SaaSTransferHeader, reader.transferHeader)
	request.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", reader.offset, end))
	response, err := reader.client.streamHTTPClient().Do(request)
	if err != nil {
		return remoteError("unavailable", "knowledge.remote_unavailable", true, err)
	}
	if response.StatusCode != http.StatusPartialContent {
		defer response.Body.Close()
		return decodeStreamResponse(response, reader.client.responseLimit, nil)
	}
	start, actualEnd, total, valid := parseContentRange(response.Header.Get("Content-Range"))
	if !valid || start != reader.offset || actualEnd != end || total != reader.document.Bytes || response.ContentLength != end-start+1 {
		_ = response.Body.Close()
		return remoteError("unavailable", "knowledge.document_range_invalid", false, nil)
	}
	reader.body = response.Body
	reader.responseEnd = actualEnd + 1
	return nil
}

func parseContentRange(value string) (int64, int64, int64, bool) {
	if !strings.HasPrefix(value, "bytes ") {
		return 0, 0, 0, false
	}
	parts := strings.Split(strings.TrimPrefix(value, "bytes "), "/")
	if len(parts) != 2 {
		return 0, 0, 0, false
	}
	bounds := strings.Split(parts[0], "-")
	if len(bounds) != 2 {
		return 0, 0, 0, false
	}
	start, errStart := strconv.ParseInt(bounds[0], 10, 64)
	end, errEnd := strconv.ParseInt(bounds[1], 10, 64)
	total, errTotal := strconv.ParseInt(parts[1], 10, 64)
	return start, end, total, errStart == nil && errEnd == nil && errTotal == nil && start >= 0 && end >= start && total > end
}

func (reader *rangeDocumentReader) Seek(offset int64, whence int) (int64, error) {
	if reader.closed {
		return 0, remoteError("unavailable", "knowledge.document_reader_closed", false, nil)
	}
	next := offset
	switch whence {
	case io.SeekStart:
	case io.SeekCurrent:
		next = reader.offset + offset
	case io.SeekEnd:
		next = reader.document.Bytes + offset
	default:
		return 0, fmt.Errorf("invalid seek origin")
	}
	if next < 0 {
		return 0, fmt.Errorf("negative seek offset")
	}
	if reader.body != nil {
		_ = reader.body.Close()
		reader.body = nil
	}
	reader.offset = next
	reader.responseEnd = next
	return next, nil
}

func (reader *rangeDocumentReader) Close() error {
	if reader.closed {
		return nil
	}
	reader.closed = true
	if reader.body != nil {
		return reader.body.Close()
	}
	return nil
}

var _ io.ReadSeekCloser = (*rangeDocumentReader)(nil)
