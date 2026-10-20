package persistence

import (
	"context"
	"time"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
)

type KnowledgeDocumentOrigin struct {
	LibraryID  string `json:"library_id"`
	DocumentID string `json:"document_id"`
	Revision   int64  `json:"revision"`
	Mode       string `json:"mode"`
}

type KnowledgePageProjectionOrigin struct {
	PageID         string `json:"page_id"`
	Version        int64  `json:"version"`
	AccessRevision int64  `json:"access_revision"`
	PermissionID   string `json:"permission_id,omitempty"`
}

type KnowledgeDocumentReplacement struct {
	PageProjection KnowledgePageProjectionOrigin `json:"page_projection"`
	Filename       string                        `json:"filename"`
	Bytes          int64                         `json:"bytes"`
	SHA256         string                        `json:"sha256"`
	ContentVersion int64                         `json:"content_version"`
	BodyRef        string                        `json:"body_ref"`
}
type KnowledgeDocumentTransferRepository interface {
	FindKnowledgeDocumentTransfer(context.Context, string, string, KnowledgeDocumentOrigin, knowledgesdk.Authority) (KnowledgeDocumentRecord, bool, error)
}
type KnowledgeDocumentRecord struct {
	DocumentOrigin            *KnowledgeDocumentOrigin                    `json:"document_origin,omitempty"`
	PageProjection            *KnowledgePageProjectionOrigin              `json:"page_projection,omitempty"`
	Replacement               *KnowledgeDocumentReplacement               `json:"replacement,omitempty"`
	Document                  knowledgesdk.KnowledgeDocument              `json:"document"`
	RequestSHA256             string                                      `json:"request_sha256"`
	SourceID                  string                                      `json:"source_id"`
	AccessPolicySHA256        string                                      `json:"access_policy_sha256,omitempty"`
	SourceAccess              *knowledgesdk.KnowledgeDocumentSourceAccess `json:"source_access,omitempty"`
	PutRequestID              string                                      `json:"put_request_id,omitempty"`
	RemoteID                  string                                      `json:"remote_id"`
	BodyRef                   string                                      `json:"body_ref,omitempty"`
	Actor                     knowledgesdk.Authority                      `json:"actor"`
	PutStarted                bool                                        `json:"put_started"`
	SecurityScanStarted       bool                                        `json:"security_scan_started"`
	SecurityScanCompleted     bool                                        `json:"security_scan_completed"`
	PutAcknowledged           bool                                        `json:"put_acknowledged"`
	IndexObserved             bool                                        `json:"index_observed"`
	DeleteStarted             bool                                        `json:"delete_started"`
	DeleteAcknowledged        bool                                        `json:"delete_acknowledged,omitempty"`
	Reindex                   *KnowledgeDocumentReindex                   `json:"reindex,omitempty"`
	RetiredRemoteID           string                                      `json:"retired_remote_id,omitempty"`
	RetiredDeleteStarted      bool                                        `json:"retired_delete_started,omitempty"`
	RetiredDeleteAcknowledged bool                                        `json:"retired_delete_acknowledged,omitempty"`
}

type KnowledgeDocumentReindex struct {
	Generation      int64  `json:"generation"`
	RemoteID        string `json:"remote_id"`
	PutRequestID    string `json:"put_request_id"`
	PutStarted      bool   `json:"put_started"`
	PutAcknowledged bool   `json:"put_acknowledged"`
}
type KnowledgeDocumentReserve struct {
	DocumentOrigin                                               *KnowledgeDocumentOrigin `json:",omitempty"`
	LibraryID, ClientID, Filename, ContentType, SHA256, SourceID string
	Bytes                                                        int64
	AccessPolicySHA256                                           string                                      `json:",omitempty"`
	SourceAccess                                                 *knowledgesdk.KnowledgeDocumentSourceAccess `json:",omitempty"`
}
type KnowledgeDocumentLease struct {
	RuntimeID   string    `json:"runtime_id"`
	WorkspaceID string    `json:"workspace_id"`
	LibraryID   string    `json:"library_id"`
	DocumentID  string    `json:"document_id"`
	Owner       string    `json:"owner"`
	Token       int64     `json:"token"`
	ExpiresAt   time.Time `json:"expires_at"`
}
type KnowledgeDocumentProgress struct {
	Event, IndexStatus, ErrorCode         string
	SecurityScanner, SecurityReportSHA256 string
	RetryAt                               time.Time
}

// Management markers and tombstones are durable. A previously managed KB may
// not silently become an unrestricted legacy source when configuration changes.
// Worker methods are trusted runtime maintenance, never user HTTP operations.
type KnowledgeDocumentRepository interface {
	ActivateKnowledgeDocumentSource(context.Context, knowledgesdk.KnowledgeDocumentStorageScope, string) error
	KnowledgeDocumentLibrarySource(context.Context, string, knowledgesdk.Authority) (string, error)
	KnowledgeSourceManaged(context.Context, string) (bool, error)
	ReserveKnowledgeDocument(context.Context, KnowledgeDocumentReserve, knowledgesdk.Authority) (KnowledgeDocumentRecord, error)
	CommitKnowledgeDocumentContent(context.Context, string, int64, string, knowledgesdk.Authority) (KnowledgeDocumentRecord, error)
	KnowledgeDocumentRecord(context.Context, string, knowledgesdk.Authority) (KnowledgeDocumentRecord, error)
	KnowledgeDocumentByRemoteID(context.Context, string, string, string, knowledgesdk.Authority) (KnowledgeDocumentRecord, error)
	KnowledgeDocuments(context.Context, string, string, int, knowledgesdk.Authority) (knowledgesdk.KnowledgeDocumentPage, error)
	RequestKnowledgeDocumentDeletion(context.Context, string, int64, knowledgesdk.Authority) (KnowledgeDocumentRecord, error)
	ClaimKnowledgeDocumentWork(context.Context, string, string, time.Time, time.Duration) (KnowledgeDocumentLease, bool, error)
	KnowledgeDocumentWorkRecord(context.Context, KnowledgeDocumentLease) (KnowledgeDocumentRecord, error)
	StartKnowledgeDocumentScan(context.Context, KnowledgeDocumentLease) (KnowledgeDocumentRecord, bool, error)
	StartKnowledgeDocumentPut(context.Context, KnowledgeDocumentLease) (KnowledgeDocumentRecord, bool, error)
	StartKnowledgeDocumentReindexPut(context.Context, KnowledgeDocumentLease) (KnowledgeDocumentRecord, bool, error)
	StartKnowledgeDocumentRetiredDelete(context.Context, KnowledgeDocumentLease) (KnowledgeDocumentRecord, bool, error)
	StartKnowledgeDocumentDelete(context.Context, KnowledgeDocumentLease) error
	ApplyKnowledgeDocumentProgress(context.Context, KnowledgeDocumentLease, KnowledgeDocumentProgress) error
}

// KnowledgeDocumentRetryRepository is the user-triggered scheduler boundary.
// The implementation must preserve provider write/delete markers and fencing.
type KnowledgeDocumentRetryRepository interface {
	RequestKnowledgeDocumentRetry(context.Context, string, int64, knowledgesdk.Authority) (KnowledgeDocumentRecord, error)
}

type KnowledgeDocumentReindexRepository interface {
	RequestKnowledgeDocumentReindex(context.Context, string, int64, knowledgesdk.Authority) (KnowledgeDocumentRecord, error)
}
