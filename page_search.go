package knowledgesdk

import (
	"context"
	"time"
)

const (
	KnowledgePagePublicationPublished   = "published"
	KnowledgePagePublicationUnpublished = "unpublished"

	KnowledgePageSearchKindPage      = "page"
	KnowledgePageSearchKindFile      = "file"
	KnowledgePageSearchKindNativeApp = "native_app"

	KnowledgePageSearchMatchPage  = "page"
	KnowledgePageSearchMatchBlock = "block"
	KnowledgePageSearchMatchCell  = "cell"
	KnowledgePageSearchMatchNode  = "node"

	KnowledgePageSearchEventChanged   = "changed"
	KnowledgePageSearchEventHeartbeat = "heartbeat"
	KnowledgePageSearchEventRevoked   = "revoked"
)

// KnowledgePageSearchInput is the Knowledge-owned, structured search query.
// LibraryID stays in the route so authorization is evaluated before any page
// metadata or excerpt is released. Empty filter slices mean "all".
type KnowledgePageSearchInput struct {
	Query               string     `json:"query"`
	ContentTypes        []string   `json:"content_types"`
	AuthorUserIDs       []string   `json:"author_user_ids"`
	UpdatedAfter        *time.Time `json:"updated_after,omitempty"`
	UpdatedBefore       *time.Time `json:"updated_before,omitempty"`
	Tags                []string   `json:"tags"`
	PathQuery           string     `json:"path_query,omitempty"`
	AncestorPageID      string     `json:"ancestor_page_id,omitempty"`
	PublicationStatuses []string   `json:"publication_statuses"`
	Limit               int        `json:"limit"`
}

// KnowledgePageSearchLocation identifies the canonical editor entity that
// matched. Only the fields meaningful to Kind are populated.
type KnowledgePageSearchLocation struct {
	Kind     string `json:"kind"`
	BlockID  string `json:"block_id,omitempty"`
	SheetID  string `json:"sheet_id,omitempty"`
	Cell     string `json:"cell,omitempty"`
	RecordID string `json:"record_id,omitempty"`
	FieldID  string `json:"field_id,omitempty"`
	NodeID   string `json:"node_id,omitempty"`
	JSONPath string `json:"json_path,omitempty"`
}

type KnowledgePageSearchHit struct {
	LibraryID         string                      `json:"library_id"`
	PageID            string                      `json:"page_id"`
	Kind              string                      `json:"kind"`
	ContentType       string                      `json:"content_type"`
	Title             string                      `json:"title"`
	Excerpt           string                      `json:"excerpt,omitempty"`
	Path              []KnowledgePageBreadcrumb   `json:"path"`
	Tags              []string                    `json:"tags"`
	AuthorUserID      string                      `json:"author_user_id"`
	PublishedVersion  int64                       `json:"published_version,omitempty"`
	PublicationStatus string                      `json:"publication_status"`
	IndexStatus       string                      `json:"index_status"`
	Location          KnowledgePageSearchLocation `json:"location"`
	Score             int                         `json:"score"`
	UpdatedAt         time.Time                   `json:"updated_at"`
	PublishedAt       *time.Time                  `json:"published_at,omitempty"`
}

type KnowledgePageSearchPage struct {
	Items     []KnowledgePageSearchHit `json:"items"`
	Truncated bool                     `json:"truncated"`
}

// KnowledgePageSearchEvent is deliberately signal-only. A changed event never
// carries a cached title, excerpt or path; clients must repeat the authorized
// search request before rendering content again.
type KnowledgePageSearchEvent struct {
	Type     string `json:"type"`
	Revision int64  `json:"revision"`
	Code     string `json:"code,omitempty"`
}

type KnowledgePageSearchEvents struct {
	AfterRevision int64 `json:"after_revision"`
	WaitMillis    int64 `json:"wait_millis,omitempty"`
}

type KnowledgePageSearchEventBatch struct {
	Revision int64                      `json:"revision"`
	Events   []KnowledgePageSearchEvent `json:"events"`
}

type KnowledgePageSearchService interface {
	SearchKnowledgePages(context.Context, string, KnowledgePageSearchInput, Authority) (KnowledgePageSearchPage, error)
	KnowledgePageSearchEvents(context.Context, string, KnowledgePageSearchEvents, Authority) (KnowledgePageSearchEventBatch, error)
}
