package knowledgesdk

import (
	"context"
	"time"
)

const (
	KnowledgePageStatusInProgress = "in_progress"
	KnowledgePageStatusInReview   = "in_review"
	KnowledgePageStatusVerified   = "verified"
	KnowledgePageStatusDeprecated = "deprecated"

	KnowledgePagePropertyText     = "text"
	KnowledgePagePropertyNumber   = "number"
	KnowledgePagePropertyDate     = "date"
	KnowledgePagePropertySelect   = "select"
	KnowledgePagePropertyCheckbox = "checkbox"
	KnowledgePagePropertyPerson   = "person"

	KnowledgePageDefaultArchiveRetentionDays = 365
	KnowledgePageDefaultReviewIntervalDays   = 180
)

// KnowledgePageCustomProperty is a bounded, typed page-level value. Value is
// canonical text interpreted by Type: RFC3339 date/time for date, a finite
// decimal for number, true/false for checkbox, and a current library user ID
// for person. Page content and immutable versions do not own this metadata.
type KnowledgePageCustomProperty struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// KnowledgePageMetadata is shared page metadata, unlike the caller-private
// navigation state. Owner is descriptive responsibility and never grants page
// access. Archived state remains owned by KnowledgePage.State.
type KnowledgePageMetadata struct {
	LibraryID            string                        `json:"library_id"`
	PageID               string                        `json:"page_id"`
	Aliases              []string                      `json:"aliases"`
	RelatedPageIDs       []string                      `json:"related_page_ids"`
	Tags                 []string                      `json:"tags"`
	OwnerUserID          string                        `json:"owner_user_id,omitempty"`
	ContentAdminUserIDs  []string                      `json:"content_admin_user_ids"`
	Status               string                        `json:"status,omitempty"`
	Properties           []KnowledgePageCustomProperty `json:"properties"`
	NextReviewAt         *time.Time                    `json:"next_review_at,omitempty"`
	ReviewIntervalDays   int                           `json:"review_interval_days"`
	LastReviewedAt       *time.Time                    `json:"last_reviewed_at,omitempty"`
	LastReviewedByUserID string                        `json:"last_reviewed_by_user_id,omitempty"`
	ArchiveRetentionDays int                           `json:"archive_retention_days"`
	ArchivedAt           *time.Time                    `json:"archived_at,omitempty"`
	RetainUntil          *time.Time                    `json:"retain_until,omitempty"`
	Revision             int64                         `json:"revision"`
	UpdatedByUserID      string                        `json:"updated_by_user_id,omitempty"`
	UpdatedAt            *time.Time                    `json:"updated_at,omitempty"`
}

type KnowledgePageMetadataWrite struct {
	ClientID         string                        `json:"client_id"`
	ExpectedRevision int64                         `json:"expected_revision"`
	Aliases          []string                      `json:"aliases"`
	RelatedPageIDs   []string                      `json:"related_page_ids"`
	Tags             []string                      `json:"tags"`
	Status           string                        `json:"status,omitempty"`
	Properties       []KnowledgePageCustomProperty `json:"properties"`
}

// KnowledgePageGovernanceWrite owns responsibility and lifecycle policy. It
// deliberately has a separate manager-only Action from ordinary metadata so
// an editor who can tag a page cannot appoint administrators or change how
// long archived content is retained.
type KnowledgePageGovernanceWrite struct {
	ClientID             string    `json:"client_id"`
	ExpectedRevision     int64     `json:"expected_revision"`
	OwnerUserID          string    `json:"owner_user_id,omitempty"`
	ContentAdminUserIDs  []string  `json:"content_admin_user_ids"`
	ArchiveRetentionDays int       `json:"archive_retention_days"`
	ReviewIntervalDays   int       `json:"review_interval_days"`
	NextReviewAt         time.Time `json:"next_review_at,omitempty"`
}

type KnowledgePageGovernanceReview struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
}

// KnowledgePageIdentity is an authorized, live page name projection used by
// pickers and mentions. PageID is the only identity; Title, Aliases and Path
// are mutable labels provided only to disambiguate what the user can see.
type KnowledgePageIdentity struct {
	PageID      string                    `json:"page_id"`
	Title       string                    `json:"title"`
	Aliases     []string                  `json:"aliases"`
	ContentType string                    `json:"content_type"`
	Path        []KnowledgePageBreadcrumb `json:"path"`
}

type KnowledgePageIdentityPage struct {
	Items []KnowledgePageIdentity `json:"items"`
}

const (
	KnowledgePageRelationshipHierarchy = "hierarchy"
	KnowledgePageRelationshipReference = "reference"
	KnowledgePageRelationshipRelated   = "related"
)

type KnowledgePageRelationshipNode struct {
	KnowledgePageIdentity
	Center bool `json:"center,omitempty"`
}

// KnowledgePageRelationshipEdge always stores stable page IDs. Directed is
// true for hierarchy/reference edges and false for an explicit related-page
// association rendered symmetrically in the graph.
type KnowledgePageRelationshipEdge struct {
	SourcePageID string `json:"source_page_id"`
	TargetPageID string `json:"target_page_id"`
	Kind         string `json:"kind"`
	Directed     bool   `json:"directed"`
}

type KnowledgePageRelationshipGraph struct {
	CenterPageID string                          `json:"center_page_id"`
	Nodes        []KnowledgePageRelationshipNode `json:"nodes"`
	Edges        []KnowledgePageRelationshipEdge `json:"edges"`
	Truncated    bool                            `json:"truncated"`
}

type KnowledgePageDiscoveryService interface {
	KnowledgePageSearchService
	SearchKnowledgePageIdentities(context.Context, string, string, int, Authority) (KnowledgePageIdentityPage, error)
	KnowledgePageRelationshipGraph(context.Context, string, string, int, int, Authority) (KnowledgePageRelationshipGraph, error)
}

type KnowledgePageMetadataService interface {
	KnowledgePageMetadata(context.Context, string, string, Authority) (KnowledgePageMetadata, error)
	SetKnowledgePageMetadata(context.Context, string, string, KnowledgePageMetadataWrite, Authority) (KnowledgePageMetadata, error)
	SetKnowledgePageGovernance(context.Context, string, string, KnowledgePageGovernanceWrite, Authority) (KnowledgePageMetadata, error)
	ReviewKnowledgePageGovernance(context.Context, string, string, KnowledgePageGovernanceReview, Authority) (KnowledgePageMetadata, error)
}
