// Package contract defines the optional Knowledge capabilities consumed by hosts.
package contract

import (
	context "context"
	agentsdk "github.com/domainry/domainry-agent-sdk"
	persistence "github.com/domainry/domainry-agent-sdk/persistence"
	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
)

// Service is a host port; storage and business implementations are private to Knowledge.
type Service interface {
	knowledgesdk.PageService
	knowledgesdk.CollaborationService
	Artifact(ctx context.Context, id string, version int64, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersion, error)
	ArtifactAccess(ctx context.Context, a agentsdk.ConversationAuthority, key string, input any) (persistence.ConversationArtifactRepository, error)
	ArtifactBody(ctx context.Context, record persistence.ConversationArtifactRecord, content agentsdk.ConversationArtifactContent, a agentsdk.ConversationAuthority) (persistence.ConversationArtifactRecord, error)
	ArtifactContent(ctx context.Context, record persistence.ConversationArtifactRecord, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactContent, error)
	ArtifactOrigin(ctx context.Context, conversationID, runID string, a agentsdk.ConversationAuthority) (*agentsdk.ConversationSources, error)
	ArtifactVersions(ctx context.Context, id string, before int64, limit int, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersions, error)
	ArtifactView(ctx context.Context, record persistence.ConversationArtifactRecord, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersion, error)
	Artifacts(ctx context.Context, in agentsdk.ConversationArtifactQuery, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactPage, error)
	Attachment(ctx context.Context, conversationID, id string, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error)
	AttachmentAccess(ctx context.Context, action string, a agentsdk.ConversationAuthority) (persistence.ConversationAttachmentRepository, error)
	AttachmentIndexView(ctx context.Context, r persistence.ConversationAttachmentRecord, a agentsdk.ConversationAuthority) agentsdk.ConversationAttachment
	AttachmentIndexWriteAllowed(ctx context.Context, r persistence.ConversationAttachmentRecord) bool
	AttachmentKnowledge(ctx context.Context, conversation, op, q, id string, a agentsdk.ConversationAuthority) (agentsdk.ConversationKnowledgeResult, error)
	AttachmentKnowledgeAccess(ctx context.Context, conversation string, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentKnowledgeScope, map[string]persistence.ConversationAttachmentRecord, error)
	AttachmentKnowledgeBinding(ctx context.Context, conversation string, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentKnowledgeScope, error)
	AttachmentRecord(ctx context.Context, repo persistence.ConversationAttachmentRepository, conversationID, id string, a agentsdk.ConversationAuthority) (persistence.ConversationAttachmentRecord, error)
	AttachmentView(item agentsdk.ConversationAttachment) agentsdk.ConversationAttachment
	Attachments(ctx context.Context, conversationID, after string, limit int, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentPage, error)
	AuthorizeAttachmentKnowledgeResult(ctx context.Context, in agentsdk.ConversationToolRequest, result agentsdk.ConversationToolResult) error
	BindKnowledgeLibrarySource(ctx context.Context, id string, in agentsdk.KnowledgeLibrarySourceWrite, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error)
	CheckAttachmentIndex(ctx context.Context, conversation, id string, expected int64, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error)
	Close()
	CreateArtifact(ctx context.Context, in agentsdk.ConversationArtifactCreate, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersion, error)
	CreateKnowledgeLibrary(ctx context.Context, in agentsdk.KnowledgeLibraryCreate, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error)
	DatasourceAccess(ctx context.Context, id, op string, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error)
	DatasourceDefinitions(ctx context.Context, library string, a agentsdk.ConversationAuthority) ([]agentsdk.KnowledgeDatasourceDefinition, error)
	DeleteAttachment(ctx context.Context, conversationID, id string, expected int64, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error)
	DeleteKnowledgeDocument(ctx context.Context, library, id string, expected int64, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error)
	RetryKnowledgeDocument(ctx context.Context, library, id string, expected int64, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error)
	ReindexKnowledgeDocument(ctx context.Context, library, id string, expected int64, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error)
	DocumentAccess(ctx context.Context, library, op string, a agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRepository, error)
	DocumentBinding(ctx context.Context, library string, a agentsdk.ConversationAuthority) (agentsdk.ManagedKnowledgeDocumentSource, error)
	DownloadArtifact(ctx context.Context, exportID string, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactDownload, error)
	DownloadAttachment(ctx context.Context, conversationID, id string, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentDownload, error)
	DownloadKnowledgeDocument(ctx context.Context, library, id string, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocumentDownload, error)
	EditArtifact(ctx context.Context, id string, in agentsdk.ConversationArtifactEdit, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactVersion, error)
	ExportArtifact(ctx context.Context, id string, in agentsdk.ConversationArtifactExportRequest, a agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactExport, error)
	ImportConversationAttachment(ctx context.Context, library string, in agentsdk.KnowledgeAttachmentImport, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error)
	IndexAttachment(ctx context.Context, conversation, id string, expected int64, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error)
	KnowledgeDocument(ctx context.Context, library, id string, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error)
	KnowledgeDocumentRecord(ctx context.Context, repo persistence.KnowledgeDocumentRepository, library, id string, a agentsdk.ConversationAuthority) (persistence.KnowledgeDocumentRecord, error)
	KnowledgeDocuments(ctx context.Context, library, after string, limit int, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocumentPage, error)
	KnowledgeLibraries(ctx context.Context, after string, limit int, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibraryPage, error)
	KnowledgeLibrary(ctx context.Context, id string, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error)
	KnowledgeLibraryMembers(ctx context.Context, id, after string, limit int, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibraryMembers, error)
	KnowledgeLibrarySources(ctx context.Context, id, after string, limit int, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrarySources, error)
	LibraryAccess(ctx context.Context, op string, a agentsdk.ConversationAuthority) (persistence.KnowledgeLibraryRepository, error)
	LibraryAuthorize(ctx context.Context, op string, item agentsdk.KnowledgeLibrary, a agentsdk.ConversationAuthority) error
	LibraryResult(ctx context.Context, item agentsdk.KnowledgeLibrary, err error, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error)
	RemoveKnowledgeLibraryMember(ctx context.Context, id, user string, revision int64, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error)
	SetKnowledgeLibraryMember(ctx context.Context, id, user string, in agentsdk.KnowledgeLibraryMemberWrite, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error)
	Start(parent context.Context)
	TransferKnowledgeDocument(ctx context.Context, library string, in agentsdk.KnowledgeDocumentTransfer, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error)
	UpdateKnowledgeLibrary(ctx context.Context, id string, in agentsdk.KnowledgeLibraryUpdate, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeLibrary, error)
	UploadAttachment(ctx context.Context, conversationID string, in agentsdk.ConversationAttachmentUpload, a agentsdk.ConversationAuthority) (agentsdk.ConversationAttachment, error)
	UploadKnowledgeDocumentChunk(ctx context.Context, library string, in agentsdk.KnowledgeDocumentUploadChunk, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocumentUploadProgress, error)
	UploadKnowledgeDocument(ctx context.Context, library string, in agentsdk.KnowledgeDocumentUpload, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error)
	UploadKnowledgeDocumentForSource(ctx context.Context, library string, in agentsdk.KnowledgeDocumentUpload, source agentsdk.KnowledgeDocumentSourceAccess, a agentsdk.ConversationAuthority) (agentsdk.KnowledgeDocument, error)
	WakeAttachmentIndex()
	WakeKnowledgeDocuments()
}
