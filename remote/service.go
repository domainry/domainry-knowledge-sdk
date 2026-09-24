package remote

import (
	"context"
	"encoding/json"
	"time"

	agentsdk "github.com/domainry/domainry-agent-sdk"
	"github.com/domainry/domainry-agent-sdk/persistence"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
)

type serviceClient struct {
	client    *client
	options   knowledgecontract.Options
	runtimeID string
}

func (value *serviceClient) Close()                  {}
func (value *serviceClient) Start(context.Context)   {}
func (value *serviceClient) WakeAttachmentIndex()    {}
func (value *serviceClient) WakeKnowledgeDocuments() {}
func (value *serviceClient) AttachmentIndexWorker(context.Context, persistence.ConversationAttachmentIndexRepository) {
}
func (value *serviceClient) KnowledgeDocumentWorker(context.Context, persistence.KnowledgeDocumentRepository) {
}
func (value *serviceClient) ProcessAttachmentIndex(context.Context, persistence.ConversationAttachmentIndexRepository, persistence.ConversationAttachmentIndexLease) {
}
func (value *serviceClient) ProcessKnowledgeDocument(context.Context, persistence.KnowledgeDocumentRepository, persistence.KnowledgeDocumentLease) {
}

func (value *serviceClient) invoke(ctx context.Context, operation string, input any, output any) error {
	return value.client.invokeWithOptions(ctx, operation, input, output, value.options)
}

func (value *serviceClient) Artifact(ctx context.Context, id string, version int64, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersion, error) {
	var output agentsdk.ConversationArtifactVersion
	err := value.invoke(ctx, "service.artifact", struct {
		ID        string                         `json:"id"`
		Version   int64                          `json:"version"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, version, authority}, &output)
	return output, err
}
func (value *serviceClient) ArtifactAccess(ctx context.Context, authority agentsdk.ConversationAuthority, key string, input any) (persistence.ConversationArtifactRepository, error) {
	if err := value.invoke(ctx, "service.artifact_access", struct {
		Authority agentsdk.ConversationAuthority `json:"authority"`
		Key       string                         `json:"key"`
		Input     any                            `json:"input"`
	}{authority, key, input}, nil); err != nil {
		return nil, err
	}
	return artifactRepository{client: value.client}, nil
}
func (value *serviceClient) ArtifactBody(ctx context.Context, record persistence.ConversationArtifactRecord, content agentsdk.ConversationArtifactContent, authority agentsdk.ConversationAuthority) (persistence.ConversationArtifactRecord, error) {
	var output persistence.ConversationArtifactRecord
	err := value.invoke(ctx, "service.artifact_body", struct {
		Record    persistence.ConversationArtifactRecord `json:"record"`
		Content   agentsdk.ConversationArtifactContent   `json:"content"`
		Authority agentsdk.ConversationAuthority         `json:"authority"`
	}{record, content, authority}, &output)
	return output, err
}
func (value *serviceClient) ArtifactContent(ctx context.Context, record persistence.ConversationArtifactRecord, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactContent, error) {
	var output agentsdk.ConversationArtifactContent
	err := value.invoke(ctx, "service.artifact_content", struct {
		Record    persistence.ConversationArtifactRecord `json:"record"`
		Authority agentsdk.ConversationAuthority         `json:"authority"`
	}{record, authority}, &output)
	return output, err
}
func (value *serviceClient) ArtifactOrigin(ctx context.Context, conversationID, runID string, authority agentsdk.ConversationAuthority) (*agentsdk.ConversationSources, error) {
	var output *agentsdk.ConversationSources
	err := value.invoke(ctx, "service.artifact_origin", struct {
		ConversationID string                         `json:"conversation_id"`
		RunID          string                         `json:"run_id"`
		Authority      agentsdk.ConversationAuthority `json:"authority"`
	}{conversationID, runID, authority}, &output)
	return output, err
}
func (value *serviceClient) ArtifactVersions(ctx context.Context, id string, before int64, limit int, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersions, error) {
	var output agentsdk.ConversationArtifactVersions
	err := value.invoke(ctx, "service.artifact_versions", struct {
		ID        string                         `json:"id"`
		Before    int64                          `json:"before"`
		Limit     int                            `json:"limit"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, before, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) ArtifactView(ctx context.Context, record persistence.ConversationArtifactRecord, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersion, error) {
	var output agentsdk.ConversationArtifactVersion
	err := value.invoke(ctx, "service.artifact_view", struct {
		Record    persistence.ConversationArtifactRecord `json:"record"`
		Authority agentsdk.ConversationAuthority         `json:"authority"`
	}{record, authority}, &output)
	return output, err
}
func (value *serviceClient) Artifacts(ctx context.Context, input agentsdk.ConversationArtifactQuery, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactPage, error) {
	var output agentsdk.ConversationArtifactPage
	err := value.invoke(ctx, "service.artifacts", struct {
		Input     agentsdk.ConversationArtifactQuery `json:"input"`
		Authority agentsdk.ConversationAuthority     `json:"authority"`
	}{input, authority}, &output)
	return output, err
}
func (value *serviceClient) CreateArtifact(ctx context.Context, input agentsdk.ConversationArtifactCreate, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersion, error) {
	var output agentsdk.ConversationArtifactVersion
	err := value.invoke(ctx, "service.artifact_create", struct {
		Input     agentsdk.ConversationArtifactCreate `json:"input"`
		Authority agentsdk.ConversationAuthority      `json:"authority"`
	}{input, authority}, &output)
	return output, err
}
func (value *serviceClient) EditArtifact(ctx context.Context, id string, input agentsdk.ConversationArtifactEdit, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersion, error) {
	var output agentsdk.ConversationArtifactVersion
	err := value.invoke(ctx, "service.artifact_edit", struct {
		ID        string                            `json:"id"`
		Input     agentsdk.ConversationArtifactEdit `json:"input"`
		Authority agentsdk.ConversationAuthority    `json:"authority"`
	}{id, input, authority}, &output)
	return output, err
}
func (value *serviceClient) ExportArtifact(ctx context.Context, id string, input agentsdk.ConversationArtifactExportRequest, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactExport, error) {
	var output agentsdk.ConversationArtifactExport
	err := value.invoke(ctx, "service.artifact_export", struct {
		ID        string                                     `json:"id"`
		Input     agentsdk.ConversationArtifactExportRequest `json:"input"`
		Authority agentsdk.ConversationAuthority             `json:"authority"`
	}{id, input, authority}, &output)
	return output, err
}
func (value *serviceClient) DownloadArtifact(ctx context.Context, exportID string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactDownload, error) {
	var output agentsdk.ConversationArtifactDownload
	err := value.invoke(ctx, "service.artifact_download", struct {
		ExportID  string                         `json:"export_id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{exportID, authority}, &output)
	return output, err
}

func (value *serviceClient) Attachment(ctx context.Context, conversationID, id string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error) {
	var output agentsdk.ConversationAttachment
	err := value.invoke(ctx, "service.attachment", struct {
		ConversationID string                         `json:"conversation_id"`
		ID             string                         `json:"id"`
		Authority      agentsdk.ConversationAuthority `json:"authority"`
	}{conversationID, id, authority}, &output)
	return output, err
}
func (value *serviceClient) AttachmentAccess(ctx context.Context, action string, authority agentsdk.ConversationAuthority) (persistence.ConversationAttachmentRepository, error) {
	if err := value.invoke(ctx, "service.attachment_access", struct {
		Action    string                         `json:"action"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{action, authority}, nil); err != nil {
		return nil, err
	}
	return attachmentRepository{client: value.client}, nil
}
func (value *serviceClient) AttachmentIndexView(ctx context.Context, record persistence.ConversationAttachmentRecord, authority agentsdk.ConversationAuthority) agentsdk.ConversationAttachment {
	var output agentsdk.ConversationAttachment
	if value.invoke(ctx, "service.attachment_index_view", struct {
		Record    persistence.ConversationAttachmentRecord `json:"record"`
		Authority agentsdk.ConversationAuthority           `json:"authority"`
	}{record, authority}, &output) != nil {
		return agentsdk.ConversationAttachment{}
	}
	return output
}
func (*serviceClient) AttachmentIndexWriteAllowed(context.Context, persistence.ConversationAttachmentRecord) bool {
	return false
}
func (value *serviceClient) AttachmentKnowledge(ctx context.Context, conversation, operation, query, id string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationKnowledgeResult, error) {
	var output agentsdk.ConversationKnowledgeResult
	err := value.invoke(ctx, "service.attachment_knowledge", struct {
		Conversation string                         `json:"conversation"`
		Operation    string                         `json:"operation"`
		Query        string                         `json:"query"`
		ID           string                         `json:"id"`
		Authority    agentsdk.ConversationAuthority `json:"authority"`
	}{conversation, operation, query, id, authority}, &output)
	return output, err
}
func (value *serviceClient) AttachmentKnowledgeAccess(ctx context.Context, conversation string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentKnowledgeScope, map[string]persistence.ConversationAttachmentRecord, error) {
	var output struct {
		Scope   remoteSourceDescriptor                              `json:"scope"`
		Records map[string]persistence.ConversationAttachmentRecord `json:"records"`
	}
	err := value.invoke(ctx, "service.attachment_knowledge_access", struct {
		Conversation string                         `json:"conversation"`
		Authority    agentsdk.ConversationAuthority `json:"authority"`
	}{conversation, authority}, &output)
	if err != nil {
		return agentsdk.ConversationAttachmentKnowledgeScope{}, nil, err
	}
	return agentsdk.ConversationAttachmentKnowledgeScope{Source: managedSource{client: value.client, descriptor: output.Scope}, PermissionID: output.Scope.PermissionID}, output.Records, nil
}
func (value *serviceClient) AttachmentKnowledgeBinding(ctx context.Context, conversation string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentKnowledgeScope, error) {
	var output remoteSourceDescriptor
	err := value.invoke(ctx, "service.attachment_knowledge_binding", struct {
		Conversation string                         `json:"conversation"`
		Authority    agentsdk.ConversationAuthority `json:"authority"`
	}{conversation, authority}, &output)
	if err != nil {
		return agentsdk.ConversationAttachmentKnowledgeScope{}, err
	}
	return agentsdk.ConversationAttachmentKnowledgeScope{Source: managedSource{client: value.client, descriptor: output}, PermissionID: output.PermissionID}, nil
}
func (value *serviceClient) AttachmentRecord(ctx context.Context, _ persistence.ConversationAttachmentRepository, conversationID, id string, authority agentsdk.ConversationAuthority) (persistence.ConversationAttachmentRecord, error) {
	var output persistence.ConversationAttachmentRecord
	err := value.invoke(ctx, "service.attachment_record", struct {
		ConversationID string                         `json:"conversation_id"`
		ID             string                         `json:"id"`
		Authority      agentsdk.ConversationAuthority `json:"authority"`
	}{conversationID, id, authority}, &output)
	return output, err
}
func (*serviceClient) AttachmentView(item agentsdk.ConversationAttachment) agentsdk.ConversationAttachment {
	item.Indexing = nil
	return item
}
func (value *serviceClient) Attachments(ctx context.Context, conversationID, after string, limit int, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentPage, error) {
	var output agentsdk.ConversationAttachmentPage
	err := value.invoke(ctx, "service.attachments", struct {
		ConversationID string                         `json:"conversation_id"`
		After          string                         `json:"after"`
		Limit          int                            `json:"limit"`
		Authority      agentsdk.ConversationAuthority `json:"authority"`
	}{conversationID, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) AuthorizeAttachmentKnowledgeResult(ctx context.Context, input agentsdk.ConversationToolRequest, result agentsdk.ConversationToolResult) error {
	return value.invoke(ctx, "service.attachment_knowledge_authorize", struct {
		Input  agentsdk.ConversationToolRequest `json:"input"`
		Result agentsdk.ConversationToolResult  `json:"result"`
	}{input, result}, nil)
}
func (value *serviceClient) CheckAttachmentIndex(ctx context.Context, conversation, id string, expected int64, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error) {
	var output agentsdk.ConversationAttachment
	err := value.invoke(ctx, "service.attachment_index_check", struct {
		Conversation string                         `json:"conversation"`
		ID           string                         `json:"id"`
		Expected     int64                          `json:"expected"`
		Authority    agentsdk.ConversationAuthority `json:"authority"`
	}{conversation, id, expected, authority}, &output)
	return output, err
}
func (value *serviceClient) DeleteAttachment(ctx context.Context, conversationID, id string, expected int64, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error) {
	var output agentsdk.ConversationAttachment
	err := value.invoke(ctx, "service.attachment_delete", struct {
		ConversationID string                         `json:"conversation_id"`
		ID             string                         `json:"id"`
		Expected       int64                          `json:"expected"`
		Authority      agentsdk.ConversationAuthority `json:"authority"`
	}{conversationID, id, expected, authority}, &output)
	return output, err
}
func (value *serviceClient) DownloadAttachment(ctx context.Context, conversationID, id string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentDownload, error) {
	var output agentsdk.ConversationAttachmentDownload
	err := value.invoke(ctx, "service.attachment_download", struct {
		ConversationID string                         `json:"conversation_id"`
		ID             string                         `json:"id"`
		Authority      agentsdk.ConversationAuthority `json:"authority"`
	}{conversationID, id, authority}, &output)
	return output, err
}
func (value *serviceClient) IndexAttachment(ctx context.Context, conversation, id string, expected int64, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error) {
	var output agentsdk.ConversationAttachment
	err := value.invoke(ctx, "service.attachment_index", struct {
		Conversation string                         `json:"conversation"`
		ID           string                         `json:"id"`
		Expected     int64                          `json:"expected"`
		Authority    agentsdk.ConversationAuthority `json:"authority"`
	}{conversation, id, expected, authority}, &output)
	return output, err
}
func (value *serviceClient) UploadAttachment(ctx context.Context, conversationID string, input agentsdk.ConversationAttachmentUpload, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error) {
	var output agentsdk.ConversationAttachment
	err := value.invoke(ctx, "service.attachment_upload", struct {
		ConversationID string                                `json:"conversation_id"`
		Input          agentsdk.ConversationAttachmentUpload `json:"input"`
		Authority      agentsdk.ConversationAuthority        `json:"authority"`
	}{conversationID, input, authority}, &output)
	return output, err
}

type attachmentRepository struct{ client *client }
type attachmentReserveWire struct {
	ClientID       string `json:"client_id"`
	ConversationID string `json:"conversation_id"`
	Filename       string `json:"filename"`
	ContentType    string `json:"content_type"`
	SHA256         string `json:"sha256"`
	Bytes          int64  `json:"bytes"`
	Content        []byte `json:"content"`
}

func (value attachmentRepository) ReserveAttachment(ctx context.Context, input persistence.ConversationAttachmentReserve, authority agentsdk.ConversationAuthority) (persistence.ConversationAttachmentRecord, error) {
	var output persistence.ConversationAttachmentRecord
	wire := attachmentReserveWire{input.ClientID, input.ConversationID, input.Filename, input.ContentType, input.SHA256, input.Bytes, input.Content}
	err := value.client.invoke(ctx, "attachments.reserve", struct {
		Input     attachmentReserveWire          `json:"input"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{wire, authority}, &output)
	return output, err
}
func (value attachmentRepository) AttachmentRecord(ctx context.Context, id string, authority agentsdk.ConversationAuthority) (persistence.ConversationAttachmentRecord, error) {
	var output persistence.ConversationAttachmentRecord
	err := value.client.invoke(ctx, "attachments.record", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, authority}, &output)
	return output, err
}
func (value attachmentRepository) Attachments(ctx context.Context, conversationID, after string, limit int, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentPage, error) {
	var output agentsdk.ConversationAttachmentPage
	err := value.client.invoke(ctx, "attachments.list", struct {
		ConversationID string                         `json:"conversation_id"`
		After          string                         `json:"after"`
		Limit          int                            `json:"limit"`
		Authority      agentsdk.ConversationAuthority `json:"authority"`
	}{conversationID, after, limit, authority}, &output)
	return output, err
}
func (value attachmentRepository) TransitionAttachment(ctx context.Context, id string, expected int64, input persistence.ConversationAttachmentTransition, authority agentsdk.ConversationAuthority) (persistence.ConversationAttachmentRecord, error) {
	var output persistence.ConversationAttachmentRecord
	err := value.client.invoke(ctx, "attachments.transition", struct {
		ID        string                                       `json:"id"`
		Expected  int64                                        `json:"expected"`
		Input     persistence.ConversationAttachmentTransition `json:"input"`
		Authority agentsdk.ConversationAuthority               `json:"authority"`
	}{id, expected, input, authority}, &output)
	return output, err
}
func (value attachmentRepository) AttachmentContent(ctx context.Context, id string, authority agentsdk.ConversationAuthority) ([]byte, error) {
	var output []byte
	err := value.client.invoke(ctx, "attachments.content", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, authority}, &output)
	return output, err
}
func (value attachmentRepository) DeleteAttachmentContent(ctx context.Context, id string, authority agentsdk.ConversationAuthority) error {
	return value.client.invoke(ctx, "attachments.delete_content", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, authority}, nil)
}

func (value *serviceClient) BindKnowledgeLibrarySource(ctx context.Context, id string, input agentsdk.KnowledgeLibrarySourceWrite, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	var output agentsdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.library_source_bind", struct {
		ID        string                               `json:"id"`
		Input     agentsdk.KnowledgeLibrarySourceWrite `json:"input"`
		Authority agentsdk.ConversationAuthority       `json:"authority"`
	}{id, input, authority}, &output)
	return output, err
}
func (value *serviceClient) CreateKnowledgeLibrary(ctx context.Context, input agentsdk.KnowledgeLibraryCreate, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	var output agentsdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.library_create", struct {
		Input     agentsdk.KnowledgeLibraryCreate `json:"input"`
		Authority agentsdk.ConversationAuthority  `json:"authority"`
	}{input, authority}, &output)
	return output, err
}
func (value *serviceClient) DatasourceAccess(ctx context.Context, id, operation string, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	var output agentsdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.datasource_access", struct {
		ID        string                         `json:"id"`
		Operation string                         `json:"operation"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, operation, authority}, &output)
	return output, err
}
func (value *serviceClient) DatasourceDefinitions(ctx context.Context, library string, authority agentsdk.ConversationAuthority) ([]agentsdk.KnowledgeDatasourceDefinition, error) {
	var output []agentsdk.KnowledgeDatasourceDefinition
	err := value.invoke(ctx, "service.datasource_definitions", struct {
		Library   string                         `json:"library"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{library, authority}, &output)
	return output, err
}
func (value *serviceClient) DeleteKnowledgeDocument(ctx context.Context, library, id string, expected int64, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error) {
	var output agentsdk.KnowledgeDocument
	err := value.invoke(ctx, "service.document_delete", struct {
		Library   string                         `json:"library"`
		ID        string                         `json:"id"`
		Expected  int64                          `json:"expected"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{library, id, expected, authority}, &output)
	return output, err
}
func (value *serviceClient) DocumentAccess(ctx context.Context, library, operation string, authority agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRepository, error) {
	if err := value.invoke(ctx, "service.document_access", struct {
		Library   string                         `json:"library"`
		Operation string                         `json:"operation"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{library, operation, authority}, nil); err != nil {
		return nil, err
	}
	return documentRepository{client: value.client}, nil
}
func (value *serviceClient) DocumentBinding(ctx context.Context, library string, authority agentsdk.ConversationAuthority) (agentsdk.ManagedKnowledgeDocumentSource, error) {
	var output remoteSourceDescriptor
	err := value.invoke(ctx, "service.document_binding", struct {
		Library   string                         `json:"library"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{library, authority}, &output)
	if err != nil {
		return nil, err
	}
	return managedSource{client: value.client, descriptor: output}, nil
}
func (value *serviceClient) DownloadKnowledgeDocument(ctx context.Context, library, id string, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocumentDownload, error) {
	var output agentsdk.KnowledgeDocumentDownload
	err := value.invoke(ctx, "service.document_download", struct {
		Library   string                         `json:"library"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{library, id, authority}, &output)
	return output, err
}
func (value *serviceClient) ImportConversationAttachment(ctx context.Context, library string, input agentsdk.KnowledgeAttachmentImport, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error) {
	var output agentsdk.KnowledgeDocument
	err := value.invoke(ctx, "service.document_import_attachment", struct {
		Library   string                             `json:"library"`
		Input     agentsdk.KnowledgeAttachmentImport `json:"input"`
		Authority agentsdk.ConversationAuthority     `json:"authority"`
	}{library, input, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeDocument(ctx context.Context, library, id string, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error) {
	var output agentsdk.KnowledgeDocument
	err := value.invoke(ctx, "service.document", struct {
		Library   string                         `json:"library"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{library, id, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeDocumentRecord(ctx context.Context, _ persistence.KnowledgeDocumentRepository, library, id string, authority agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRecord, error) {
	var output persistence.KnowledgeDocumentRecord
	err := value.invoke(ctx, "service.document_record", struct {
		Library   string                         `json:"library"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{library, id, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeDocuments(ctx context.Context, library, after string, limit int, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocumentPage, error) {
	var output agentsdk.KnowledgeDocumentPage
	err := value.invoke(ctx, "service.documents", struct {
		Library   string                         `json:"library"`
		After     string                         `json:"after"`
		Limit     int                            `json:"limit"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{library, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibraries(ctx context.Context, after string, limit int, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibraryPage, error) {
	var output agentsdk.KnowledgeLibraryPage
	err := value.invoke(ctx, "service.libraries", struct {
		After     string                         `json:"after"`
		Limit     int                            `json:"limit"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibrary(ctx context.Context, id string, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	var output agentsdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.library", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibraryMembers(ctx context.Context, id, after string, limit int, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibraryMembers, error) {
	var output agentsdk.KnowledgeLibraryMembers
	err := value.invoke(ctx, "service.library_members", struct {
		ID        string                         `json:"id"`
		After     string                         `json:"after"`
		Limit     int                            `json:"limit"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) KnowledgeLibrarySources(ctx context.Context, id, after string, limit int, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrarySources, error) {
	var output agentsdk.KnowledgeLibrarySources
	err := value.invoke(ctx, "service.library_sources", struct {
		ID        string                         `json:"id"`
		After     string                         `json:"after"`
		Limit     int                            `json:"limit"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, after, limit, authority}, &output)
	return output, err
}
func (value *serviceClient) LibraryAccess(ctx context.Context, operation string, authority agentsdk.ConversationAuthority) (persistence.KnowledgeLibraryRepository, error) {
	if err := value.invoke(ctx, "service.library_access", struct {
		Operation string                         `json:"operation"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{operation, authority}, nil); err != nil {
		return nil, err
	}
	return libraryRepository{service: value}, nil
}
func (value *serviceClient) LibraryAuthorize(ctx context.Context, operation string, item agentsdk.KnowledgeLibrary, authority agentsdk.ConversationAuthority) error {
	return value.invoke(ctx, "service.library_authorize", struct {
		Operation string                         `json:"operation"`
		Item      agentsdk.KnowledgeLibrary      `json:"item"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{operation, item, authority}, nil)
}
func (value *serviceClient) LibraryResult(ctx context.Context, item agentsdk.KnowledgeLibrary, inputErr error, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	if inputErr != nil {
		return agentsdk.KnowledgeLibrary{}, inputErr
	}
	var output agentsdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.library_result", struct {
		Item      agentsdk.KnowledgeLibrary      `json:"item"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{item, authority}, &output)
	return output, err
}
func (value *serviceClient) RemoveKnowledgeLibraryMember(ctx context.Context, id, user string, revision int64, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	var output agentsdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.library_member_remove", struct {
		ID        string                         `json:"id"`
		User      string                         `json:"user"`
		Revision  int64                          `json:"revision"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, user, revision, authority}, &output)
	return output, err
}
func (value *serviceClient) SetKnowledgeLibraryMember(ctx context.Context, id, user string, input agentsdk.KnowledgeLibraryMemberWrite, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	var output agentsdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.library_member_set", struct {
		ID        string                               `json:"id"`
		User      string                               `json:"user"`
		Input     agentsdk.KnowledgeLibraryMemberWrite `json:"input"`
		Authority agentsdk.ConversationAuthority       `json:"authority"`
	}{id, user, input, authority}, &output)
	return output, err
}
func (value *serviceClient) TransferKnowledgeDocument(ctx context.Context, library string, input agentsdk.KnowledgeDocumentTransfer, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error) {
	var output agentsdk.KnowledgeDocument
	err := value.invoke(ctx, "service.document_transfer", struct {
		Library   string                             `json:"library"`
		Input     agentsdk.KnowledgeDocumentTransfer `json:"input"`
		Authority agentsdk.ConversationAuthority     `json:"authority"`
	}{library, input, authority}, &output)
	return output, err
}
func (value *serviceClient) UpdateKnowledgeLibrary(ctx context.Context, id string, input agentsdk.KnowledgeLibraryUpdate, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	var output agentsdk.KnowledgeLibrary
	err := value.invoke(ctx, "service.library_update", struct {
		ID        string                          `json:"id"`
		Input     agentsdk.KnowledgeLibraryUpdate `json:"input"`
		Authority agentsdk.ConversationAuthority  `json:"authority"`
	}{id, input, authority}, &output)
	return output, err
}
func (value *serviceClient) UploadDocumentContent(ctx context.Context, library string, input agentsdk.KnowledgeDocumentUpload, origin *persistence.KnowledgeAttachmentOrigin, documentOrigin *persistence.KnowledgeDocumentOrigin, recheck func() error, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error) {
	if recheck != nil {
		if err := recheck(); err != nil {
			return agentsdk.KnowledgeDocument{}, err
		}
	}
	var output agentsdk.KnowledgeDocument
	err := value.invoke(ctx, "service.document_upload_content", struct {
		Library        string                                 `json:"library"`
		Input          agentsdk.KnowledgeDocumentUpload       `json:"input"`
		Origin         *persistence.KnowledgeAttachmentOrigin `json:"origin,omitempty"`
		DocumentOrigin *persistence.KnowledgeDocumentOrigin   `json:"document_origin,omitempty"`
		Authority      agentsdk.ConversationAuthority         `json:"authority"`
	}{library, input, origin, documentOrigin, authority}, &output)
	return output, err
}
func (value *serviceClient) UploadKnowledgeDocument(ctx context.Context, library string, input agentsdk.KnowledgeDocumentUpload, authority agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error) {
	var output agentsdk.KnowledgeDocument
	err := value.invoke(ctx, "service.document_upload", struct {
		Library   string                           `json:"library"`
		Input     agentsdk.KnowledgeDocumentUpload `json:"input"`
		Authority agentsdk.ConversationAuthority   `json:"authority"`
	}{library, input, authority}, &output)
	return output, err
}

type libraryRepository struct{ service *serviceClient }

func (v libraryRepository) CreateKnowledgeLibrary(c context.Context, i agentsdk.KnowledgeLibraryCreate, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	return v.service.CreateKnowledgeLibrary(c, i, a)
}
func (v libraryRepository) KnowledgeLibraries(c context.Context, x string, l int, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibraryPage, error) {
	return v.service.KnowledgeLibraries(c, x, l, a)
}
func (v libraryRepository) KnowledgeLibrary(c context.Context, id string, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	return v.service.KnowledgeLibrary(c, id, a)
}
func (v libraryRepository) UpdateKnowledgeLibrary(c context.Context, id string, i agentsdk.KnowledgeLibraryUpdate, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	return v.service.UpdateKnowledgeLibrary(c, id, i, a)
}
func (v libraryRepository) KnowledgeLibraryMembers(c context.Context, id, x string, l int, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibraryMembers, error) {
	return v.service.KnowledgeLibraryMembers(c, id, x, l, a)
}
func (v libraryRepository) SetKnowledgeLibraryMember(c context.Context, id, u string, i agentsdk.KnowledgeLibraryMemberWrite, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	return v.service.SetKnowledgeLibraryMember(c, id, u, i, a)
}
func (v libraryRepository) RemoveKnowledgeLibraryMember(c context.Context, id, u string, r int64, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error) {
	return v.service.RemoveKnowledgeLibraryMember(c, id, u, r, a)
}

type documentRepository struct{ client *client }

func (v documentRepository) ActivateKnowledgeDocumentSource(c context.Context, s agentsdk.KnowledgeDocumentStorageScope, id string) error {
	return v.client.invoke(c, "documents.activate_source", struct {
		Scope agentsdk.KnowledgeDocumentStorageScope `json:"scope"`
		ID    string                                 `json:"id"`
	}{s, id}, nil)
}
func (v documentRepository) KnowledgeDocumentLibrarySource(c context.Context, id string, a agentsdk.ConversationAuthority) (string, error) {
	var o string
	e := v.client.invoke(c, "documents.library_source", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, a}, &o)
	return o, e
}
func (v documentRepository) KnowledgeSourceManaged(c context.Context, id string) (bool, error) {
	var o bool
	e := v.client.invoke(c, "documents.source_managed", struct {
		ID string `json:"id"`
	}{id}, &o)
	return o, e
}
func (v documentRepository) ReserveKnowledgeDocument(c context.Context, i persistence.KnowledgeDocumentReserve, a agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRecord, error) {
	var o persistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.reserve", struct {
		Input     persistence.KnowledgeDocumentReserve `json:"input"`
		Authority agentsdk.ConversationAuthority       `json:"authority"`
	}{i, a}, &o)
	return o, e
}
func (v documentRepository) CommitKnowledgeDocumentContent(c context.Context, id string, r int64, ref string, a agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRecord, error) {
	var o persistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.commit_content", struct {
		ID        string                         `json:"id"`
		Revision  int64                          `json:"revision"`
		BodyRef   string                         `json:"body_ref"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, r, ref, a}, &o)
	return o, e
}
func (v documentRepository) KnowledgeDocumentRecord(c context.Context, id string, a agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRecord, error) {
	var o persistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.record", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, a}, &o)
	return o, e
}
func (v documentRepository) KnowledgeDocumentByRemoteID(c context.Context, l, s, id string, a agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRecord, error) {
	var o persistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.by_remote", struct {
		Library   string                         `json:"library"`
		Source    string                         `json:"source"`
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{l, s, id, a}, &o)
	return o, e
}
func (v documentRepository) KnowledgeDocuments(c context.Context, l, x string, n int, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocumentPage, error) {
	var o agentsdk.KnowledgeDocumentPage
	e := v.client.invoke(c, "documents.list", struct {
		Library   string                         `json:"library"`
		After     string                         `json:"after"`
		Limit     int                            `json:"limit"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{l, x, n, a}, &o)
	return o, e
}
func (v documentRepository) RequestKnowledgeDocumentDeletion(c context.Context, id string, r int64, a agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRecord, error) {
	var o persistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.request_delete", struct {
		ID        string                         `json:"id"`
		Revision  int64                          `json:"revision"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, r, a}, &o)
	return o, e
}
func (v documentRepository) ClaimKnowledgeDocumentWork(c context.Context, owner, scope string, now time.Time, d time.Duration) (persistence.KnowledgeDocumentLease, bool, error) {
	var o struct {
		Lease   persistence.KnowledgeDocumentLease `json:"lease"`
		Claimed bool                               `json:"claimed"`
	}
	e := v.client.invoke(c, "documents.claim_work", struct {
		Owner    string        `json:"owner"`
		Scope    string        `json:"scope"`
		Now      time.Time     `json:"now"`
		Duration time.Duration `json:"duration"`
	}{owner, scope, now, d}, &o)
	return o.Lease, o.Claimed, e
}
func (v documentRepository) KnowledgeDocumentWorkRecord(c context.Context, l persistence.KnowledgeDocumentLease) (persistence.KnowledgeDocumentRecord, error) {
	var o persistence.KnowledgeDocumentRecord
	e := v.client.invoke(c, "documents.work_record", l, &o)
	return o, e
}
func (v documentRepository) StartKnowledgeDocumentPut(c context.Context, l persistence.KnowledgeDocumentLease) (persistence.KnowledgeDocumentRecord, bool, error) {
	var o struct {
		Record  persistence.KnowledgeDocumentRecord `json:"record"`
		Started bool                                `json:"started"`
	}
	e := v.client.invoke(c, "documents.start_put", l, &o)
	return o.Record, o.Started, e
}
func (v documentRepository) StartKnowledgeDocumentDelete(c context.Context, l persistence.KnowledgeDocumentLease) error {
	return v.client.invoke(c, "documents.start_delete", l, nil)
}
func (v documentRepository) ApplyKnowledgeDocumentProgress(c context.Context, l persistence.KnowledgeDocumentLease, p persistence.KnowledgeDocumentProgress) error {
	return v.client.invoke(c, "documents.apply_progress", struct {
		Lease    persistence.KnowledgeDocumentLease    `json:"lease"`
		Progress persistence.KnowledgeDocumentProgress `json:"progress"`
	}{l, p}, nil)
}

type artifactRepository struct{ client *client }

func (value artifactRepository) Artifacts(ctx context.Context, input agentsdk.ConversationArtifactQuery, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactPage, error) {
	var output agentsdk.ConversationArtifactPage
	err := value.client.invoke(ctx, "artifacts.list", struct {
		Input     agentsdk.ConversationArtifactQuery `json:"input"`
		Authority agentsdk.ConversationAuthority     `json:"authority"`
	}{input, authority}, &output)
	return output, err
}
func (value artifactRepository) ArtifactRecord(ctx context.Context, id string, version int64, authority agentsdk.ConversationAuthority) (persistence.ConversationArtifactRecord, error) {
	return (artifactMutations{client: value.client}).ArtifactRecord(ctx, id, version, authority)
}
func (value artifactRepository) ArtifactVersions(ctx context.Context, id string, before int64, limit int, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersions, error) {
	var output agentsdk.ConversationArtifactVersions
	err := value.client.invoke(ctx, "artifacts.versions", struct {
		ID        string                         `json:"id"`
		Before    int64                          `json:"before"`
		Limit     int                            `json:"limit"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, before, limit, authority}, &output)
	return output, err
}
func (value artifactRepository) SaveArtifact(ctx context.Context, input persistence.ConversationArtifactWrite, authority agentsdk.ConversationAuthority) (persistence.ConversationArtifactRecord, error) {
	return (artifactMutations{client: value.client}).SaveArtifact(ctx, input, authority)
}
func (value artifactRepository) SaveArtifactExport(ctx context.Context, input persistence.ConversationArtifactExportWrite, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactExport, error) {
	return (artifactMutations{client: value.client}).SaveArtifactExport(ctx, input, authority)
}
func (value artifactRepository) ArtifactExport(ctx context.Context, id string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactExport, error) {
	var output agentsdk.ConversationArtifactExport
	err := value.client.invoke(ctx, "artifacts.export_record", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, authority}, &output)
	return output, err
}
func (value artifactRepository) ArtifactExportContent(ctx context.Context, id string, authority agentsdk.ConversationAuthority) ([]byte, error) {
	var output []byte
	err := value.client.invoke(ctx, "artifacts.export_content", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, authority}, &output)
	return output, err
}
func (value artifactRepository) RecordArtifactDownload(ctx context.Context, id string, authority agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactExport, error) {
	var output agentsdk.ConversationArtifactExport
	err := value.client.invoke(ctx, "artifacts.record_download", struct {
		ID        string                         `json:"id"`
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{id, authority}, &output)
	return output, err
}

var _ knowledgecontract.Service = (*serviceClient)(nil)
var _ = json.Valid
