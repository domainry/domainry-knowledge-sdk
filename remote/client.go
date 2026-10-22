package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/domainry/domainry-foundation/telemetry"
	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	"github.com/domainry/domainry-knowledge-sdk/saashost"
)

type client struct {
	baseURL       *url.URL
	runtimeID     string
	token         string
	httpClient    *http.Client
	requestLimit  int64
	responseLimit int64
}

func newClient(config Config, runtimeID string) (*client, error) {
	config = normalizeConfig(config)
	base, err := url.Parse(strings.TrimSpace(config.Endpoint))
	if err != nil || base == nil || (base.Scheme != "http" && base.Scheme != "https") || strings.TrimSpace(base.Host) == "" || base.RawQuery != "" || base.Fragment != "" {
		return nil, fmt.Errorf("Knowledge SaaS endpoint is invalid")
	}
	base.Path = strings.TrimRight(base.Path, "/")
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: config.RequestTimeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	}
	return &client{baseURL: base, runtimeID: runtimeID, token: strings.TrimSpace(config.ServiceAccessToken), httpClient: telemetry.InstrumentHTTPClient(httpClient, "github.com/domainry/domainry-knowledge-sdk/remote"), requestLimit: config.MaxRequestBytes, responseLimit: config.MaxResponseBytes}, nil
}

func (c *client) request(ctx context.Context, method, path string, input any, output any) error {
	if ctx == nil {
		return remoteError("unavailable", "knowledge.context_required", false, nil)
	}
	var body io.Reader
	if input != nil {
		raw, err := json.Marshal(input)
		if err != nil {
			return remoteError("bad_request", "knowledge.request_invalid", false, err)
		}
		if int64(len(raw)) > c.requestLimit {
			return remoteError("bad_request", "knowledge.request_too_large", false, nil)
		}
		body = bytes.NewReader(raw)
	}
	endpoint := *c.baseURL
	endpoint.Path = strings.TrimRight(c.baseURL.Path, "/") + "/" + strings.TrimLeft(path, "/")
	ctx = telemetry.WithHTTPClientOperation(ctx, "knowledge.remote")
	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return remoteError("unavailable", "knowledge.request_invalid", false, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "domainry-knowledge-sdk-go")
	req.Header.Set(saashost.RuntimeIDHeader, c.runtimeID)
	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return remoteError("unavailable", "knowledge.remote_unavailable", true, err)
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, c.responseLimit+1))
	if err != nil {
		return remoteError("unavailable", "knowledge.response_read_failed", true, err)
	}
	if int64(len(raw)) > c.responseLimit {
		return remoteError("unavailable", "knowledge.response_too_large", false, nil)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var payload struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if json.Unmarshal(raw, &payload) == nil && strings.TrimSpace(payload.Code) != "" {
			return &knowledgesdk.Error{Class: statusClass(response.StatusCode), Code: safeCode(payload.Code), Message: safeMessage(payload.Message), Retryable: response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= 500}
		}
		return remoteError(statusClass(response.StatusCode), "knowledge.remote_request_failed", response.StatusCode >= 500, nil)
	}
	if output == nil {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(output); err != nil {
		return remoteError("unavailable", "knowledge.response_invalid", false, err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return remoteError("unavailable", "knowledge.response_invalid", false, nil)
	}
	return nil
}

func (c *client) invoke(ctx context.Context, operation string, input any, output any) error {
	return c.invokeWithOptions(ctx, operation, input, output, knowledgecontract.Options{})
}

func (c *client) invokeWithOptions(ctx context.Context, operation string, input any, output any, options knowledgecontract.Options) error {
	_, err := c.invokeWithGrants(ctx, operation, input, output, options)
	return err
}

func (c *client) invokeWithGrants(ctx context.Context, operation string, input any, output any, options knowledgecontract.Options) ([]knowledgecontract.SaaSGrant, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, remoteError("bad_request", "knowledge.request_invalid", false, err)
	}
	request := knowledgecontract.SaaSRequest{Operation: operation, Input: raw}
	for attempts := 0; attempts < 12; attempts++ {
		var response knowledgecontract.SaaSResponse
		if err := c.request(ctx, http.MethodPost, knowledgecontract.SaaSInvokePath, request, &response); err != nil {
			return nil, err
		}
		if response.Challenge != nil {
			grant, err := answerChallenge(ctx, *response.Challenge, options)
			if err != nil {
				return nil, err
			}
			request.Grants = append(request.Grants, grant)
			continue
		}
		if response.Error != nil {
			return nil, decodeRemoteError(response.Error)
		}
		if output == nil || len(response.Result) == 0 {
			return append([]knowledgecontract.SaaSGrant(nil), request.Grants...), nil
		}
		decoder := json.NewDecoder(bytes.NewReader(response.Result))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(output); err != nil {
			return nil, remoteError("unavailable", "knowledge.response_invalid", false, err)
		}
		if decoder.Decode(&struct{}{}) != io.EOF {
			return nil, remoteError("unavailable", "knowledge.response_invalid", false, nil)
		}
		return append([]knowledgecontract.SaaSGrant(nil), request.Grants...), nil
	}
	return nil, remoteError("unavailable", "knowledge.challenge_limit_exceeded", false, nil)
}

func answerChallenge(ctx context.Context, challenge knowledgecontract.SaaSChallenge, options knowledgecontract.Options) (knowledgecontract.SaaSGrant, error) {
	result := any(struct{}{})
	switch challenge.Kind {
	case "authorize.library":
		var input struct {
			Operation string                        `json:"operation"`
			Item      knowledgesdk.KnowledgeLibrary `json:"item"`
			Authority knowledgesdk.Authority        `json:"authority"`
		}
		if options.LibraryAuthorizer == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.library_authorizer_required", false, nil)
		}
		if err := options.LibraryAuthorizer.AuthorizeKnowledgeLibrary(ctx, input.Operation, input.Item, input.Authority); err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
	case "authorize.library_member":
		var input struct {
			User      string                 `json:"user"`
			Authority knowledgesdk.Authority `json:"authority"`
		}
		if options.LibraryAuthorizer == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.library_authorizer_required", false, nil)
		}
		if err := options.LibraryAuthorizer.ValidateKnowledgeLibraryMember(ctx, input.User, input.Authority); err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
	case "authorize.library_organization_member":
		var input struct {
			OrganizationID string                 `json:"organization_id"`
			Authority      knowledgesdk.Authority `json:"authority"`
		}
		if options.LibraryAuthorizer == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.library_authorizer_required", false, nil)
		}
		if err := options.LibraryAuthorizer.ValidateKnowledgeLibraryOrganizationMember(ctx, input.OrganizationID, input.Authority); err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
	case "authorize.document_source":
		var input struct {
			Source    knowledgesdk.KnowledgeDocumentSourceAccess `json:"source"`
			Authority knowledgesdk.Authority                     `json:"authority"`
		}
		if options.SourceAuthorizer == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.document_source_authorizer_required", false, nil)
		}
		if err := options.SourceAuthorizer.AuthorizeKnowledgeDocumentSource(ctx, input.Source, input.Authority); err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
	default:
		return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.challenge_unsupported", false, nil)
	}
	raw, err := json.Marshal(result)
	if err != nil {
		return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.challenge_result_invalid", false, err)
	}
	return knowledgecontract.SaaSGrant{Token: challenge.Token, Result: raw}, nil
}

func remoteError(class knowledgesdk.ErrorClass, code string, retryable bool, cause error) error {
	return &knowledgesdk.Error{Class: class, Code: code, Retryable: retryable, Cause: cause}
}

func decodeRemoteError(input *knowledgecontract.SaaSError) error {
	if input == nil {
		return remoteError("unavailable", "knowledge.remote_request_failed", false, nil)
	}
	classified := &knowledgesdk.Error{Class: safeClass(input.Class), Code: safeCode(input.Code), Message: safeMessage(input.Message), Retryable: input.Retryable}
	if classified.Class == knowledgesdk.ErrorClassConflict && len(input.Details) != 0 {
		var conflict knowledgesdk.KnowledgePageConflict
		if json.Unmarshal(input.Details, &conflict) == nil && conflict.CurrentDraftRevision > 0 && conflict.BasePublishedVersion >= 0 && len(conflict.CurrentContentSHA256) == 64 && !conflict.CurrentUpdatedAt.IsZero() {
			return &knowledgesdk.KnowledgePageConflictError{Err: classified, Conflict: conflict}
		}
	}
	return classified
}

func safeClass(value string) knowledgesdk.ErrorClass {
	switch strings.TrimSpace(value) {
	case string(knowledgesdk.ErrorClassBadRequest):
		return knowledgesdk.ErrorClassBadRequest
	case string(knowledgesdk.ErrorClassForbidden):
		return knowledgesdk.ErrorClassForbidden
	case string(knowledgesdk.ErrorClassNotFound):
		return knowledgesdk.ErrorClassNotFound
	case string(knowledgesdk.ErrorClassConflict):
		return knowledgesdk.ErrorClassConflict
	case string(knowledgesdk.ErrorClassRateLimited):
		return knowledgesdk.ErrorClassRateLimited
	case string(knowledgesdk.ErrorClassUnavailable):
		return knowledgesdk.ErrorClassUnavailable
	case string(knowledgesdk.ErrorClassInternal):
		return knowledgesdk.ErrorClassInternal
	default:
		return knowledgesdk.ErrorClassInternal
	}
}

func safeCode(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 128 {
		return "knowledge.remote_request_failed"
	}
	for _, r := range value {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '.' || r == '_') {
			return "knowledge.remote_request_failed"
		}
	}
	return value
}

func safeMessage(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 256 {
		return ""
	}
	return value
}

func statusClass(status int) knowledgesdk.ErrorClass {
	switch status {
	case http.StatusBadRequest, http.StatusRequestEntityTooLarge:
		return knowledgesdk.ErrorClassBadRequest
	case http.StatusUnauthorized, http.StatusForbidden:
		return knowledgesdk.ErrorClassForbidden
	case http.StatusNotFound:
		return knowledgesdk.ErrorClassNotFound
	case http.StatusConflict:
		return knowledgesdk.ErrorClassConflict
	case http.StatusTooManyRequests:
		return knowledgesdk.ErrorClassRateLimited
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return knowledgesdk.ErrorClassUnavailable
	default:
		return knowledgesdk.ErrorClassInternal
	}
}
