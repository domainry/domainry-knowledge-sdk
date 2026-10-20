package knowledgesdk

import (
	"context"
	"encoding/json"
)

// Knowledge remains an optional host service. Data preserves the upstream
// response for a configured default source. Managed libraries provide an explicitly mapped,
// locally authorized document projection instead; callers must not infer a
// provider-specific document schema.
// ScopeSHA256 binds a receipt to its source configuration and authenticated owner.
type Result struct {
	LibraryID      string          `json:"library_id,omitempty"`
	Provider       string          `json:"provider"`
	KBID           string          `json:"kb_id"`
	Operation      string          `json:"operation"`
	Query          string          `json:"query,omitempty"`
	DocumentID     string          `json:"doc_id,omitempty"`
	SearchTopK     int             `json:"search_top_k,omitempty"`
	SearchMode     string          `json:"search_mode,omitempty"`
	SearchReranked *bool           `json:"search_reranked,omitempty"`
	SearchDegraded string          `json:"search_degraded_reason,omitempty"`
	ScopeSHA256    string          `json:"scope_sha256"`
	Citations      []Citation      `json:"citations,omitempty"`
	Data           json.RawMessage `json:"data"`
}

// Citations are normalized by the trusted knowledge host from the actual
// response. IDs bind source configuration, owner, response and source position;
// they are not document ACL grants and must be revalidated with the receipt.
type Citation struct {
	LibraryID        string            `json:"library_id,omitempty"`
	ID               string            `json:"id"`
	Source           string            `json:"source,omitempty"`
	Provider         string            `json:"provider"`
	KBID             string            `json:"kb_id"`
	Operation        string            `json:"operation"`
	ObjectKey        string            `json:"object_key,omitempty"`
	RecordID         string            `json:"record_id,omitempty"`
	DocumentID       string            `json:"doc_id"`
	PageID           string            `json:"page_id,omitempty"`
	PageVersion      int64             `json:"page_version,omitempty"`
	ChunkID          string            `json:"chunk_id,omitempty"`
	DocumentVersion  int64             `json:"document_version,omitempty"`
	ContentSHA256    string            `json:"content_sha256,omitempty"`
	SourceRevision   string            `json:"source_revision,omitempty"`
	Score            float64           `json:"score,omitempty"`
	Title            string            `json:"title,omitempty"`
	URL              string            `json:"url,omitempty"`
	Excerpt          string            `json:"excerpt,omitempty"`
	ExcerptTruncated bool              `json:"excerpt_truncated,omitempty"`
	Location         *DocumentLocation `json:"location,omitempty"`
}

type Source interface {
	SearchKnowledge(context.Context, string, Authority) (Result, error)
	ReadKnowledge(context.Context, string, Authority) (Result, error)
	// RevalidateKnowledge must verify current access to ALL the saved data,
	// including document ACL changes. A historical permission receipt is not a grant.
	RevalidateKnowledge(context.Context, Result, Authority) error
}

// ResultReadAuthorizer applies the current read policy for a previously saved
// result. It is intentionally separate from RevalidateKnowledge: replaying a
// frozen execution result and releasing it to another collaboration boundary
// can require a stronger permission such as document download.
type ResultReadAuthorizer interface {
	AuthorizeKnowledgeResultRead(context.Context, Result, Authority) error
}

// ExtractionContentSource exposes normalized passages from a previously
// authorized Knowledge result. It never grants access by itself.
type ExtractionContentSource interface {
	KnowledgeExtractionPassages(context.Context, Result, Authority) ([]KnowledgeDocumentPassage, error)
}
