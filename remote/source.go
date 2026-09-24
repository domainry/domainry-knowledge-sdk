package remote

import (
	"context"
	"fmt"

	agentsdk "github.com/domainry/domainry-agent-sdk"
)

type remoteSourceDescriptor struct {
	Handle             string `json:"handle"`
	Identity           string `json:"identity"`
	PermissionID       string `json:"permission_id,omitempty"`
	AccessPolicySHA256 string `json:"access_policy_sha256,omitempty"`
	MaxBytes           int64  `json:"max_bytes,omitempty"`
	Ready              bool   `json:"ready"`
}

type managedSource struct {
	client     *client
	descriptor remoteSourceDescriptor
}

func (value managedSource) KnowledgeDocumentSourceIdentity() string { return value.descriptor.Identity }
func (value managedSource) AttachmentKnowledgeSourceIdentity() string {
	return value.descriptor.Identity
}
func (value managedSource) KnowledgeDocumentAccessPolicySHA256() string {
	return value.descriptor.AccessPolicySHA256
}
func (value managedSource) KnowledgeDocumentMaxBytes() int64 { return value.descriptor.MaxBytes }
func (value managedSource) KnowledgeDocumentManagementReady() error {
	if !value.descriptor.Ready {
		return fmt.Errorf("Knowledge document source is unavailable")
	}
	return nil
}
func (value managedSource) PutKnowledgeDocument(ctx context.Context, input agentsdk.KnowledgeDocumentContent, authority agentsdk.ConversationAuthority) error {
	return value.client.invoke(ctx, "source.document.put", struct {
		Handle    string                            `json:"handle"`
		Input     agentsdk.KnowledgeDocumentContent `json:"input"`
		Authority agentsdk.ConversationAuthority    `json:"authority"`
	}{value.descriptor.Handle, input, authority}, nil)
}
func (value managedSource) InspectKnowledgeDocument(ctx context.Context, id string, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocumentState, error) {
	var output agentsdk.KnowledgeDocumentState
	err := value.client.invoke(ctx, "source.document.inspect", struct {
		Handle    string                         `json:"handle"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{value.descriptor.Handle, id, authority}, &output)
	return output, err
}
func (value managedSource) DeleteKnowledgeDocument(ctx context.Context, id string, authority agentsdk.ConversationAuthority) error {
	return value.client.invoke(ctx, "source.document.delete", struct {
		Handle    string                         `json:"handle"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{value.descriptor.Handle, id, authority}, nil)
}
func (value managedSource) RecoverKnowledgeDocumentDelete(ctx context.Context, id string, authority agentsdk.ConversationAuthority) error {
	return value.client.invoke(ctx, "source.document.recover_delete", struct {
		Handle    string                         `json:"handle"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{value.descriptor.Handle, id, authority}, nil)
}
func (value managedSource) SearchKnowledge(ctx context.Context, query string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationKnowledgeResult, error) {
	var output agentsdk.ConversationKnowledgeResult
	err := value.client.invoke(ctx, "source.knowledge.search", struct {
		Handle    string                         `json:"handle"`
		Query     string                         `json:"query"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{value.descriptor.Handle, query, authority}, &output)
	return output, err
}
func (value managedSource) ReadKnowledge(ctx context.Context, id string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationKnowledgeResult, error) {
	var output agentsdk.ConversationKnowledgeResult
	err := value.client.invoke(ctx, "source.knowledge.read", struct {
		Handle    string                         `json:"handle"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{value.descriptor.Handle, id, authority}, &output)
	return output, err
}
func (value managedSource) RevalidateKnowledge(ctx context.Context, input agentsdk.ConversationKnowledgeResult, authority agentsdk.ConversationAuthority) error {
	return value.client.invoke(ctx, "source.knowledge.revalidate", struct {
		Handle    string                               `json:"handle"`
		Input     agentsdk.ConversationKnowledgeResult `json:"input"`
		Authority agentsdk.ConversationAuthority       `json:"authority"`
	}{value.descriptor.Handle, input, authority}, nil)
}
func (value managedSource) SearchKnowledgeDocumentPassages(ctx context.Context, query string, authority agentsdk.ConversationAuthority) ([]agentsdk.KnowledgeDocumentPassage, error) {
	var output []agentsdk.KnowledgeDocumentPassage
	err := value.client.invoke(ctx, "source.passages.search", struct {
		Handle    string                         `json:"handle"`
		Query     string                         `json:"query"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{value.descriptor.Handle, query, authority}, &output)
	return output, err
}
func (value managedSource) ReadKnowledgeDocumentPassages(ctx context.Context, id string, authority agentsdk.ConversationAuthority) ([]agentsdk.KnowledgeDocumentPassage, error) {
	var output []agentsdk.KnowledgeDocumentPassage
	err := value.client.invoke(ctx, "source.passages.read", struct {
		Handle    string                         `json:"handle"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{value.descriptor.Handle, id, authority}, &output)
	return output, err
}

var _ agentsdk.ConversationAttachmentKnowledgeSource = managedSource{}
var _ agentsdk.KnowledgeDocumentDeleteRecoverySource = managedSource{}
