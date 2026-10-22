package contract

import (
	"encoding/json"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
)

const (
	SaaSDiscoveryPath         = "/knowledge/v1/discovery"
	SaaSInvokePath            = "/knowledge/v1/invoke"
	SaaSDocumentUploadPath    = "/knowledge/v1/documents/upload"
	SaaSDocumentDownloadPath  = "/knowledge/v1/documents/download"
	SaaSSourceDocumentPutPath = "/knowledge/v1/sources/documents"

	SaaSDocumentUploadPrepareOperation   = "service.document_upload_prepare"
	SaaSDocumentDownloadPrepareOperation = "service.document_download_prepare"

	SaaSDocumentUploadOperation       = "service.document_upload"
	SaaSDocumentUploadChunkOperation  = "service.document_upload_chunk"
	SaaSDocumentUploadSourceOperation = "service.document_upload_source"
	SaaSTransferHeader                = "X-Domainry-Knowledge-Transfer"
)

// SaaSDocumentUploadRequest is the bounded metadata part that accompanies a
// streamed multipart document body. Content is deliberately excluded from the
// JSON contract by KnowledgeDocumentUpload.
type SaaSDocumentUploadRequest struct {
	Operation string                                      `json:"operation"`
	Library   string                                      `json:"library"`
	Input     knowledgesdk.KnowledgeDocumentUpload        `json:"input"`
	Source    *knowledgesdk.KnowledgeDocumentSourceAccess `json:"source,omitempty"`
	Authority knowledgesdk.Authority                      `json:"authority"`
	Grants    []SaaSGrant                                 `json:"grants,omitempty"`
}

type SaaSDocumentChunkUploadRequest struct {
	Operation string                                    `json:"operation"`
	Library   string                                    `json:"library"`
	Input     knowledgesdk.KnowledgeDocumentUploadChunk `json:"input"`
	Authority knowledgesdk.Authority                    `json:"authority"`
	Grants    []SaaSGrant                               `json:"grants,omitempty"`
}

// SaaSDocumentDownloadRequest is encoded into a bounded private request header
// so the download remains a standard GET with Range support.
type SaaSDocumentDownloadRequest struct {
	Library   string                 `json:"library"`
	ID        string                 `json:"id"`
	PageID    string                 `json:"page_id,omitempty"`
	Authority knowledgesdk.Authority `json:"authority"`
	Grants    []SaaSGrant            `json:"grants,omitempty"`
}

type SaaSSourceDocumentPutRequest struct {
	Handle    string                                     `json:"handle"`
	Input     knowledgesdk.KnowledgeDocumentIndexContent `json:"input"`
	Authority knowledgesdk.Authority                     `json:"authority"`
}

// SaaSRequest is the bounded transport envelope used by the official remote
// client. Operation names are owned by this SDK; Input never contains a Store,
// database handle, SQL transaction or deployment credential.
type SaaSRequest struct {
	Operation string          `json:"operation"`
	Input     json.RawMessage `json:"input,omitempty"`
	Grants    []SaaSGrant     `json:"grants,omitempty"`
}

// SaaSResponse keeps implementation errors and response bodies behind a
// stable SDK-owned boundary. Implementations must not place provider, SQL or
// credential details in Message.
type SaaSResponse struct {
	Result    json.RawMessage `json:"result,omitempty"`
	Error     *SaaSError      `json:"error,omitempty"`
	Challenge *SaaSChallenge  `json:"challenge,omitempty"`
}

type SaaSError struct {
	Class     string          `json:"class"`
	Code      string          `json:"code"`
	Message   string          `json:"message,omitempty"`
	Retryable bool            `json:"retryable,omitempty"`
	Details   json.RawMessage `json:"details,omitempty"`
}

type SaaSChallenge struct {
	Token string          `json:"token"`
	Kind  string          `json:"kind"`
	Input json.RawMessage `json:"input"`
}

type SaaSGrant struct {
	Token  string          `json:"token"`
	Result json.RawMessage `json:"result,omitempty"`
}
