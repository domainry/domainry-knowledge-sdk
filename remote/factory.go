// Package remote implements the official Knowledge SaaS client. It owns HTTP
// transport only; Knowledge business behavior and persistence remain in the
// independently deployed Knowledge service.
package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	agentsdk "github.com/domainry/domainry-agent-sdk"
	"github.com/domainry/domainry-agent-sdk/persistence"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	"github.com/domainry/domainry-knowledge-sdk/saashost"
	lifecyclecontract "github.com/domainry/domainry-lifecycle-sdk/contract"
	lifecyclemodel "github.com/domainry/domainry-lifecycle-sdk/model"
)

type Factory struct{ config Config }

func NewFactory(config Config) *Factory { return &Factory{config: config} }

func (factory *Factory) OpenSaaS(ctx context.Context, ref knowledgecontract.ApplicationRef) (saashost.Binding, error) {
	if err := ref.Validate(); err != nil {
		return nil, err
	}
	client, err := newClient(factory.config, ref.RuntimeID)
	if err != nil {
		return nil, err
	}
	var descriptor knowledgecontract.Descriptor
	if err := client.request(ctx, http.MethodGet, knowledgecontract.SaaSDiscoveryPath, nil, &descriptor); err != nil {
		return nil, err
	}
	if err := descriptor.Validate(); err != nil {
		return nil, remoteError("unavailable", "knowledge.protocol_incompatible", false, err)
	}
	if descriptor.Audience != ref.RuntimeID {
		return nil, remoteError("forbidden", "knowledge.runtime_mismatch", false, nil)
	}
	runtime := &runtimeClient{client: client}
	return &binding{descriptor: descriptor, client: client, runtime: runtime}, nil
}

type binding struct {
	descriptor knowledgecontract.Descriptor
	client     *client
	runtime    *runtimeClient
}

func (value *binding) Descriptor() knowledgecontract.Descriptor { return value.descriptor }
func (value *binding) Runtime() knowledgecontract.Runtime       { return value.runtime }
func (value *binding) ArtifactMutations() knowledgecontract.ArtifactMutationService {
	return artifactMutations{client: value.client}
}
func (value *binding) SubjectLifecycle() lifecyclecontract.SubjectExecutionHandler {
	return subjectLifecycle{client: value.client}
}
func (*binding) Close(context.Context) error { return nil }

type runtimeClient struct {
	client  *client
	options knowledgecontract.Options
}

func (value *runtimeClient) Validate(options *knowledgecontract.Options) error {
	if options == nil {
		return fmt.Errorf("Knowledge options are required")
	}
	return nil
}
func (value *runtimeClient) Prepare(_ string, options knowledgecontract.Options) (knowledgecontract.ConversationKnowledge, error) {
	value.options = options
	return conversationKnowledge{client: value.client, options: options}, nil
}
func (value *runtimeClient) Activate(runtimeID string, options knowledgecontract.Options) error {
	if runtimeID != value.client.runtimeID {
		return remoteError("forbidden", "knowledge.runtime_mismatch", false, nil)
	}
	value.options = options
	return value.client.invoke(context.Background(), "runtime.activate", struct{}{}, nil)
}
func (value *runtimeClient) NewService(runtimeID string, options knowledgecontract.Options) knowledgecontract.Service {
	value.options = options
	return &serviceClient{client: value.client, options: options, runtimeID: runtimeID}
}

type conversationKnowledge struct {
	client  *client
	options knowledgecontract.Options
}

func (value conversationKnowledge) Search(ctx context.Context, query string, authority agentsdk.ConversationAuthority) (json.RawMessage, error) {
	var output json.RawMessage
	err := value.client.invokeWithOptions(ctx, "runtime.knowledge.search", struct {
		Query     string                         `json:"query"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{query, authority}, &output, value.options)
	return output, err
}

type artifactMutations struct{ client *client }

func (value artifactMutations) ArtifactRecord(ctx context.Context, id string, version int64, authority agentsdk.ConversationAuthority) (persistence.ConversationArtifactRecord, error) {
	var output persistence.ConversationArtifactRecord
	err := value.client.invoke(ctx, "artifacts.record", struct {
		ID        string                         `json:"id"`
		Version   int64                          `json:"version"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, version, authority}, &output)
	return output, err
}
func (value artifactMutations) SaveArtifact(ctx context.Context, input persistence.ConversationArtifactWrite, authority agentsdk.ConversationAuthority) (persistence.ConversationArtifactRecord, error) {
	var output persistence.ConversationArtifactRecord
	err := value.client.invoke(ctx, "artifacts.save", struct {
		Input     persistence.ConversationArtifactWrite `json:"input"`
		Authority agentsdk.ConversationAuthority        `json:"authority"`
	}{input, authority}, &output)
	return output, err
}
func (value artifactMutations) SaveArtifactExport(ctx context.Context, input persistence.ConversationArtifactExportWrite, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactExport, error) {
	var output agentsdk.ConversationArtifactExport
	type exportWrite struct {
		RequestSHA256 string                              `json:"request_sha256,omitempty"`
		TTLSeconds    int64                               `json:"ttl_seconds"`
		ClientID      string                              `json:"client_id"`
		Export        agentsdk.ConversationArtifactExport `json:"export"`
		Content       []byte                              `json:"content"`
	}
	err := value.client.invoke(ctx, "artifacts.save_export", struct {
		Input     exportWrite                    `json:"input"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{exportWrite{input.RequestSHA256, input.TTLSeconds, input.ClientID, input.Export, input.Content}, authority}, &output)
	return output, err
}

type subjectLifecycle struct{ client *client }

func (value subjectLifecycle) Owner(context.Context) string { return "knowledge" }
func (value subjectLifecycle) PreviewSubject(ctx context.Context, workspaceID, subjectID string) (json.RawMessage, error) {
	return value.invoke(ctx, "subjects.preview", "", workspaceID, subjectID, nil)
}
func (value subjectLifecycle) ExportSubjectForRequest(ctx context.Context, requestID, workspaceID, subjectID string) (json.RawMessage, error) {
	return value.invoke(ctx, "subjects.export", requestID, workspaceID, subjectID, nil)
}
func (value subjectLifecycle) EraseSubjectForRequest(ctx context.Context, requestID, workspaceID, subjectID string, holds []lifecyclemodel.LegalHold) (json.RawMessage, error) {
	return value.invoke(ctx, "subjects.erase", requestID, workspaceID, subjectID, holds)
}
func (value subjectLifecycle) invoke(ctx context.Context, operation, requestID, workspaceID, subjectID string, holds []lifecyclemodel.LegalHold) (json.RawMessage, error) {
	var output json.RawMessage
	err := value.client.invoke(ctx, operation, struct {
		RequestID   string                     `json:"request_id,omitempty"`
		WorkspaceID string                     `json:"workspace_id"`
		SubjectID   string                     `json:"subject_id"`
		LegalHolds  []lifecyclemodel.LegalHold `json:"legal_holds,omitempty"`
	}{requestID, workspaceID, subjectID, holds}, &output)
	return output, err
}

var _ saashost.Factory = (*Factory)(nil)
var _ saashost.Binding = (*binding)(nil)
var _ knowledgecontract.Runtime = (*runtimeClient)(nil)
