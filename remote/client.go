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

	agentsdk "github.com/domainry/domainry-agent-sdk"
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
	return &client{baseURL: base, runtimeID: runtimeID, token: strings.TrimSpace(config.ServiceAccessToken), httpClient: httpClient, requestLimit: config.MaxRequestBytes, responseLimit: config.MaxResponseBytes}, nil
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
		var payload knowledgecontract.SaaSResponse
		if json.Unmarshal(raw, &payload) == nil && payload.Error != nil {
			return &agentsdk.Error{Class: safeClass(payload.Error.Class), Code: safeCode(payload.Error.Code), Message: safeMessage(payload.Error.Message), Retryable: payload.Error.Retryable}
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
	raw, err := json.Marshal(input)
	if err != nil {
		return remoteError("bad_request", "knowledge.request_invalid", false, err)
	}
	request := knowledgecontract.SaaSRequest{Operation: operation, Input: raw}
	for attempts := 0; attempts < 12; attempts++ {
		var response knowledgecontract.SaaSResponse
		if err := c.request(ctx, http.MethodPost, knowledgecontract.SaaSInvokePath, request, &response); err != nil {
			return err
		}
		if response.Challenge != nil {
			grant, err := answerChallenge(ctx, *response.Challenge, options)
			if err != nil {
				return err
			}
			request.Grants = append(request.Grants, grant)
			continue
		}
		if response.Error != nil {
			return &agentsdk.Error{Class: safeClass(response.Error.Class), Code: safeCode(response.Error.Code), Message: safeMessage(response.Error.Message), Retryable: response.Error.Retryable}
		}
		if output == nil || len(response.Result) == 0 {
			return nil
		}
		decoder := json.NewDecoder(bytes.NewReader(response.Result))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(output); err != nil {
			return remoteError("unavailable", "knowledge.response_invalid", false, err)
		}
		return nil
	}
	return remoteError("unavailable", "knowledge.challenge_limit_exceeded", false, nil)
}

func answerChallenge(ctx context.Context, challenge knowledgecontract.SaaSChallenge, options knowledgecontract.Options) (knowledgecontract.SaaSGrant, error) {
	result := any(struct{}{})
	switch challenge.Kind {
	case "context.conversation":
		var input struct {
			ConversationID string                         `json:"conversation_id"`
			Authority      agentsdk.ConversationAuthority `json:"authority"`
		}
		if options.Context == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.context_reader_required", false, nil)
		}
		value, err := options.Context.Conversation(ctx, input.ConversationID, input.Authority)
		if err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
		result = value
	case "context.run":
		var input struct {
			ConversationID string                         `json:"conversation_id"`
			RunID          string                         `json:"run_id"`
			Authority      agentsdk.ConversationAuthority `json:"authority"`
		}
		if options.Context == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.context_reader_required", false, nil)
		}
		value, err := options.Context.Run(ctx, input.ConversationID, input.RunID, input.Authority)
		if err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
		result = value
	case "authorize.attachment":
		var input struct {
			Action    string                         `json:"action"`
			Authority agentsdk.ConversationAuthority `json:"authority"`
		}
		if options.AttachmentAuthorizer == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.attachment_authorizer_required", false, nil)
		}
		if err := options.AttachmentAuthorizer.AuthorizeConversationAttachment(ctx, input.Action, input.Authority); err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
	case "authorize.library":
		var input struct {
			Operation string                         `json:"operation"`
			Item      agentsdk.KnowledgeLibrary      `json:"item"`
			Authority agentsdk.ConversationAuthority `json:"authority"`
		}
		if options.LibraryAuthorizer == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.library_authorizer_required", false, nil)
		}
		if err := options.LibraryAuthorizer.AuthorizeKnowledgeLibrary(ctx, input.Operation, input.Item, input.Authority); err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
	case "authorize.library_member":
		var input struct {
			User      string                         `json:"user"`
			Authority agentsdk.ConversationAuthority `json:"authority"`
		}
		if options.LibraryAuthorizer == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.library_authorizer_required", false, nil)
		}
		if err := options.LibraryAuthorizer.ValidateKnowledgeLibraryMember(ctx, input.User, input.Authority); err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
	case "authorize.personal":
		var input agentsdk.ConversationToolRequest
		if options.PersonalAuthorizer == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.personal_authorizer_required", false, nil)
		}
		value, err := options.PersonalAuthorizer.AuthorizeConversationTool(ctx, input)
		if err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
		result = value
	case "sources.check":
		var input struct {
			Consumer  string                         `json:"consumer"`
			Sources   *agentsdk.ConversationSources  `json:"sources"`
			Authority agentsdk.ConversationAuthority `json:"authority"`
		}
		if options.Sources == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.source_policy_required", false, nil)
		}
		value, err := options.Sources.CheckSources(ctx, input.Authority, input.Consumer, input.Sources)
		if err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
		result = value
	case "sources.run":
		var input struct {
			Consumer  string                            `json:"consumer"`
			Ref       agentsdk.ConversationRunReference `json:"ref"`
			Authority agentsdk.ConversationAuthority    `json:"authority"`
		}
		if options.Sources == nil || json.Unmarshal(challenge.Input, &input) != nil {
			return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.source_policy_required", false, nil)
		}
		value, err := options.Sources.CheckRun(ctx, input.Authority, input.Consumer, input.Ref)
		if err != nil {
			return knowledgecontract.SaaSGrant{}, err
		}
		result = value
	default:
		return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.challenge_unsupported", false, nil)
	}
	raw, err := json.Marshal(result)
	if err != nil {
		return knowledgecontract.SaaSGrant{}, remoteError("unavailable", "knowledge.challenge_result_invalid", false, err)
	}
	return knowledgecontract.SaaSGrant{Token: challenge.Token, Result: raw}, nil
}

func remoteError(class, code string, retryable bool, cause error) error {
	return &agentsdk.Error{Class: class, Code: code, Retryable: retryable, Cause: cause}
}

func safeClass(value string) string {
	switch strings.TrimSpace(value) {
	case "bad_request", "forbidden", "not_found", "conflict", "unavailable":
		return strings.TrimSpace(value)
	default:
		return "unavailable"
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

func statusClass(status int) string {
	switch status {
	case http.StatusBadRequest, http.StatusRequestEntityTooLarge:
		return "bad_request"
	case http.StatusUnauthorized, http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusConflict:
		return "conflict"
	default:
		return "unavailable"
	}
}
