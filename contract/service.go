package contract

import (
	"context"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
	knowledgepersistence "github.com/domainry/domainry-knowledge-sdk/persistence"
)

// Service is the complete standalone Knowledge application boundary.
type Service interface {
	knowledgesdk.PageService
	knowledgesdk.CollaborationService
	knowledgesdk.PageTemplateService
	knowledgesdk.PageShareService
	knowledgesdk.WebhookService
	knowledgesdk.KnowledgeAgentReferenceService
	knowledgesdk.KnowledgeAIService
	knowledgesdk.KnowledgeApplicationCardService

	CreateKnowledgeLibrary(context.Context, knowledgesdk.KnowledgeLibraryCreate, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)
	KnowledgeLibraries(context.Context, string, int, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryPage, error)
	KnowledgeLibrary(context.Context, string, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)
	UpdateKnowledgeLibrary(context.Context, string, knowledgesdk.KnowledgeLibraryUpdate, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)
	KnowledgeLibraryMembers(context.Context, string, string, int, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryMembers, error)
	SetKnowledgeLibraryMember(context.Context, string, string, knowledgesdk.KnowledgeLibraryMemberWrite, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)
	RemoveKnowledgeLibraryMember(context.Context, string, string, int64, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)
	KnowledgeLibraryOrganizationMembers(context.Context, string, string, int, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibraryOrganizationMembers, error)
	SetKnowledgeLibraryOrganizationMember(context.Context, string, string, knowledgesdk.KnowledgeLibraryMemberWrite, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)
	RemoveKnowledgeLibraryOrganizationMember(context.Context, string, string, int64, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)

	UploadKnowledgeDocument(context.Context, string, knowledgesdk.KnowledgeDocumentUpload, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error)
	UploadKnowledgeDocumentChunk(context.Context, string, knowledgesdk.KnowledgeDocumentUploadChunk, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentUploadProgress, error)
	UploadKnowledgeDocumentForSource(context.Context, string, knowledgesdk.KnowledgeDocumentUpload, knowledgesdk.KnowledgeDocumentSourceAccess, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error)
	KnowledgeDocuments(context.Context, string, string, int, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentPage, error)
	KnowledgeDocument(context.Context, string, string, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error)
	DownloadKnowledgeDocument(context.Context, string, string, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentDownload, error)
	DeleteKnowledgeDocument(context.Context, string, string, int64, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error)
	RetryKnowledgeDocument(context.Context, string, string, int64, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error)
	ReindexKnowledgeDocument(context.Context, string, string, int64, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error)
	TransferKnowledgeDocument(context.Context, string, knowledgesdk.KnowledgeDocumentTransfer, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocument, error)

	KnowledgeLibrarySources(context.Context, string, string, int, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrarySources, error)
	BindKnowledgeLibrarySource(context.Context, string, knowledgesdk.KnowledgeLibrarySourceWrite, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)

	LibraryAccess(context.Context, string, knowledgesdk.Authority) (knowledgepersistence.KnowledgeLibraryRepository, error)
	DocumentAccess(context.Context, string, string, knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRepository, error)
	DocumentBinding(context.Context, string, knowledgesdk.Authority) (knowledgesdk.ManagedKnowledgeDocumentSource, error)
	DatasourceDefinitions(context.Context, string, knowledgesdk.Authority) ([]knowledgesdk.KnowledgeDatasourceDefinition, error)
	DatasourceAccess(context.Context, string, string, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)
	LibraryAuthorize(context.Context, string, knowledgesdk.KnowledgeLibrary, knowledgesdk.Authority) error
	LibraryResult(context.Context, knowledgesdk.KnowledgeLibrary, error, knowledgesdk.Authority) (knowledgesdk.KnowledgeLibrary, error)
	KnowledgeDocumentRecord(context.Context, knowledgepersistence.KnowledgeDocumentRepository, string, string, knowledgesdk.Authority) (knowledgepersistence.KnowledgeDocumentRecord, error)

	Start(context.Context)
	WakeKnowledgeDocuments()
	Close()
}
