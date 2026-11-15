package knowledgesdk

import (
	"context"
	"encoding/hex"
	"strings"
)

// KnowledgeDocumentSnapshot identifies the exact immutable content and index
// generation used by an earlier, authorized managed-document read. It is not
// an authorization grant, a source credential, or permission to disclose data.
type KnowledgeDocumentSnapshot struct {
	DocumentID      string `json:"document_id"`
	ContentVersion  int64  `json:"content_version"`
	IndexGeneration int64  `json:"index_generation"`
	SHA256          string `json:"sha256"`
	PageID          string `json:"page_id,omitempty"`
	PageVersion     int64  `json:"page_version,omitempty"`
}

func (s KnowledgeDocumentSnapshot) Valid() bool {
	if s.DocumentID == "" || s.DocumentID != strings.TrimSpace(s.DocumentID) || len(s.DocumentID) > 256 || s.ContentVersion < 1 || s.IndexGeneration < 1 || len(s.SHA256) != 64 || s.SHA256 != strings.ToLower(s.SHA256) {
		return false
	}
	if _, err := hex.DecodeString(s.SHA256); err != nil {
		return false
	}
	return (s.PageID == "" && s.PageVersion == 0) || (s.PageID != "" && s.PageID == strings.TrimSpace(s.PageID) && len(s.PageID) <= 256 && s.PageVersion > 0)
}

// KnowledgeDocumentSnapshotService performs an Owner-local, metadata-only
// authorization and lifecycle check. It rechecks current library download
// rights, source-record/page access, source binding, scan and index state, and
// exact snapshot identity. It does not read original content or call a remote
// search/fetch Provider. External-source freshness still belongs to Source.
// An application must separately authorize any customer disclosure or effect.
type KnowledgeDocumentSnapshotService interface {
	CheckKnowledgeDocumentSnapshot(context.Context, string, KnowledgeDocumentSnapshot, Authority) (KnowledgeDocument, error)
}
