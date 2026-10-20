package knowledgesdk

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	agentsdk "github.com/domainry/domainry-agent-sdk"
)

// PageContentSchemaVersion is the canonical page payload contract. Version 2
// stores text documents as a structured Block tree instead of Markdown. The
// product has not shipped a version-1 page corpus, so callers must write the
// current contract directly; there is intentionally no dual-format fallback.
const PageContentSchemaVersion = 2

const (
	KnowledgePageTypeDocument              = "document"
	KnowledgePageTypeMultidimensionalTable = "multidimensional_table"
	KnowledgePageTypeSpreadsheet           = "spreadsheet"
	KnowledgePageTypePresentation          = "presentation"
	KnowledgePageTypeForm                  = "form"
	KnowledgePageTypeMindNote              = "mind_note"
	KnowledgePageTypeWhiteboard            = "whiteboard"
	KnowledgePageTypeMindMap               = "mind_map"
	KnowledgePageTypeFlowchart             = "flowchart"
)

type PageContent struct {
	ContentSchemaVersion int    `json:"content_schema_version"`
	Content              string `json:"content"`
	ContentSHA256        string `json:"content_sha256"`
	PlainText            string `json:"plain_text"`
	Summary              string `json:"summary,omitempty"`
}

// KnowledgeDocumentContent is the canonical payload for document pages. Block
// shape validation remains server-owned; RawMessage keeps the SDK independent
// from a particular React editor implementation while preserving every native
// and custom Block property losslessly.
type KnowledgeDocumentContent struct {
	Blocks []json.RawMessage `json:"blocks"`
}

type KnowledgePage struct {
	ID               string    `json:"id"`
	LibraryID        string    `json:"library_id"`
	ParentPageID     string    `json:"parent_page_id,omitempty"`
	Rank             string    `json:"rank"`
	Title            string    `json:"title"`
	ContentType      string    `json:"content_type"`
	State            string    `json:"state"` // active, archived
	PublishedVersion int64     `json:"published_version,omitempty"`
	ProjectionStatus string    `json:"projection_status,omitempty"` // none, queued, indexing, ready, failed
	AccessMode       string    `json:"access_mode"`                 // inherit, restricted
	AccessRevision   int64     `json:"access_revision"`
	AccessStatus     string    `json:"access_status"` // ready, queued, applying, failed
	Revision         int64     `json:"revision"`
	CreatedByUserID  string    `json:"created_by_user_id"`
	UpdatedByUserID  string    `json:"updated_by_user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type KnowledgePageDraft struct {
	PageID               string      `json:"page_id"`
	LibraryID            string      `json:"library_id"`
	BasePublishedVersion int64       `json:"base_published_version"`
	DraftRevision        int64       `json:"draft_revision"`
	Title                string      `json:"title"`
	ContentType          string      `json:"content_type"`
	ParentPageID         string      `json:"parent_page_id,omitempty"`
	Rank                 string      `json:"rank"`
	Body                 PageContent `json:"body"`
	UpdatedByUserID      string      `json:"updated_by_user_id"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

type KnowledgePageVersion struct {
	PageID            string      `json:"page_id"`
	LibraryID         string      `json:"library_id"`
	Version           int64       `json:"version"`
	Title             string      `json:"title"`
	ContentType       string      `json:"content_type"`
	ParentPageID      string      `json:"parent_page_id,omitempty"`
	Rank              string      `json:"rank"`
	Body              PageContent `json:"body"`
	PublishedByUserID string      `json:"published_by_user_id"`
	PublishedAt       time.Time   `json:"published_at"`
	PublishNote       string      `json:"publish_note,omitempty"`
}

type KnowledgePageView struct {
	Page        KnowledgePage             `json:"page"`
	Published   *KnowledgePageVersion     `json:"published,omitempty"`
	Draft       *KnowledgePageDraft       `json:"draft,omitempty"`
	Breadcrumbs []KnowledgePageBreadcrumb `json:"breadcrumbs"`
}

type KnowledgePageBreadcrumb struct {
	PageID string `json:"page_id"`
	Title  string `json:"title"`
}

type KnowledgePageCreate struct {
	ClientID     string `json:"client_id"`
	ParentPageID string `json:"parent_page_id,omitempty"`
	Rank         string `json:"rank,omitempty"`
	Title        string `json:"title"`
	ContentType  string `json:"content_type"`
}

type KnowledgePageDraftWrite struct {
	ClientID              string `json:"client_id"`
	ExpectedDraftRevision int64  `json:"expected_draft_revision"`
	ExpectedPageRevision  int64  `json:"expected_page_revision"`
	Title                 string `json:"title"`
	ParentPageID          string `json:"parent_page_id,omitempty"`
	Rank                  string `json:"rank"`
	ContentSchemaVersion  int    `json:"content_schema_version"`
	Content               string `json:"content"`
	ContentSHA256         string `json:"content_sha256"`
}

type KnowledgePageMove struct {
	ClientID          string `json:"client_id"`
	ExpectedRevision  int64  `json:"expected_revision"`
	ParentPageID      string `json:"parent_page_id,omitempty"`
	Rank              string `json:"rank"`
	AccessImpactToken string `json:"access_impact_token,omitempty"`
}

type KnowledgePageOrder struct {
	PageID           string `json:"page_id"`
	ExpectedRevision int64  `json:"expected_revision"`
}

// KnowledgePageReorder is an all-siblings compare-and-set. Items are in the
// desired order and must contain every active child of ParentPageID exactly
// once, preventing a paginated or stale client from silently dropping pages.
type KnowledgePageReorder struct {
	ClientID     string               `json:"client_id"`
	ParentPageID string               `json:"parent_page_id,omitempty"`
	Items        []KnowledgePageOrder `json:"items"`
}

type KnowledgePageCopy struct {
	ClientID     string `json:"client_id"`
	ParentPageID string `json:"parent_page_id,omitempty"`
	Rank         string `json:"rank,omitempty"`
	Title        string `json:"title,omitempty"`
}

type KnowledgePageStateWrite struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
}

type KnowledgePagePublish struct {
	ClientID              string `json:"client_id"`
	ExpectedPageRevision  int64  `json:"expected_page_revision"`
	ExpectedDraftRevision int64  `json:"expected_draft_revision"`
	PublishNote           string `json:"publish_note,omitempty"`
}

type KnowledgePageRestoreVersion struct {
	ClientID              string `json:"client_id"`
	ExpectedPageRevision  int64  `json:"expected_page_revision"`
	ExpectedDraftRevision int64  `json:"expected_draft_revision"`
}

type KnowledgePagePage struct {
	Items     []KnowledgePage `json:"items"`
	NextAfter string          `json:"next_after,omitempty"`
	Complete  bool            `json:"complete"`
}

type KnowledgePageVersionPage struct {
	Items      []KnowledgePageVersion `json:"items"`
	NextBefore int64                  `json:"next_before,omitempty"`
	Complete   bool                   `json:"complete"`
}

type KnowledgePageVersionDiffInput struct {
	FromVersion int64 `json:"from_version"`
	ToVersion   int64 `json:"to_version"`
}

type KnowledgePageVersionDiff struct {
	FromVersion int64             `json:"from_version"`
	ToVersion   int64             `json:"to_version"`
	Title       []KnowledgeDiffOp `json:"title"`
	Parent      []KnowledgeDiffOp `json:"parent"`
	Content     []KnowledgeDiffOp `json:"content"`
}

type KnowledgeDiffOp struct {
	Kind string `json:"kind"` // equal, insert, delete
	Text string `json:"text"`
}

type KnowledgePageConflict struct {
	CurrentDraftRevision int64     `json:"current_draft_revision"`
	BasePublishedVersion int64     `json:"base_published_version"`
	CurrentContentSHA256 string    `json:"current_content_sha256"`
	CurrentUpdatedAt     time.Time `json:"current_updated_at"`
}

// KnowledgePageConflictError carries the bounded, non-content snapshot a
// caller needs to recover from an optimistic-lock conflict. Err remains the
// canonical classified error so errors.As continues to work for existing
// authorization and HTTP status handling.
type KnowledgePageConflictError struct {
	Err      *agentsdk.Error
	Conflict KnowledgePageConflict
}

func (e *KnowledgePageConflictError) Error() string {
	if e == nil || e.Err == nil {
		return "knowledge page conflict"
	}
	return e.Err.Error()
}

func (e *KnowledgePageConflictError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *KnowledgePageConflictError) KnowledgePageConflict() KnowledgePageConflict {
	if e == nil {
		return KnowledgePageConflict{}
	}
	return e.Conflict
}

func NewKnowledgePageConflictError(code string, conflict KnowledgePageConflict) error {
	code = strings.TrimSpace(code)
	if code == "" {
		code = "page_draft_revision_conflict"
	}
	return &KnowledgePageConflictError{
		Err:      &agentsdk.Error{Class: "conflict", Code: code},
		Conflict: conflict,
	}
}

type PageService interface {
	KnowledgePages(context.Context, string, string, string, int, bool, agentsdk.ConversationAuthority) (KnowledgePagePage, error)
	KnowledgePage(context.Context, string, string, bool, agentsdk.ConversationAuthority) (KnowledgePageView, error)
	CreateKnowledgePage(context.Context, string, KnowledgePageCreate, agentsdk.ConversationAuthority) (KnowledgePageView, error)
	SaveKnowledgePageDraft(context.Context, string, string, KnowledgePageDraftWrite, agentsdk.ConversationAuthority) (KnowledgePageView, error)
	MoveKnowledgePage(context.Context, string, string, KnowledgePageMove, agentsdk.ConversationAuthority) (KnowledgePage, error)
	PreviewKnowledgePageMove(context.Context, string, string, KnowledgePageMovePreview, agentsdk.ConversationAuthority) (KnowledgePageMoveImpact, error)
	ReorderKnowledgePages(context.Context, string, KnowledgePageReorder, agentsdk.ConversationAuthority) (KnowledgePagePage, error)
	CopyKnowledgePage(context.Context, string, string, KnowledgePageCopy, agentsdk.ConversationAuthority) (KnowledgePageView, error)
	SetKnowledgePageArchived(context.Context, string, string, bool, KnowledgePageStateWrite, agentsdk.ConversationAuthority) (KnowledgePage, error)
	PublishKnowledgePage(context.Context, string, string, KnowledgePagePublish, agentsdk.ConversationAuthority) (KnowledgePageView, error)
	KnowledgePageVersions(context.Context, string, string, int64, int, agentsdk.ConversationAuthority) (KnowledgePageVersionPage, error)
	KnowledgePageVersion(context.Context, string, string, int64, agentsdk.ConversationAuthority) (KnowledgePageVersion, error)
	DiffKnowledgePageVersions(context.Context, string, string, KnowledgePageVersionDiffInput, agentsdk.ConversationAuthority) (KnowledgePageVersionDiff, error)
	RestoreKnowledgePageVersion(context.Context, string, string, int64, KnowledgePageRestoreVersion, agentsdk.ConversationAuthority) (KnowledgePageView, error)
}

const (
	KnowledgePageAccessInherit    = "inherit"
	KnowledgePageAccessRestricted = "restricted"
	KnowledgePageGrantPage        = "page"
	KnowledgePageGrantSubtree     = "subtree"
)

// KnowledgePageAccessGrant never replaces library membership. The effective
// role is the lesser of the current library role and this grant's role.
type KnowledgePageAccessGrant struct {
	UserID          string    `json:"user_id"`
	Role            string    `json:"role"`  // reader, editor, manager
	Scope           string    `json:"scope"` // page, subtree
	CreatedByUserID string    `json:"created_by_user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type KnowledgePageAccessPolicy struct {
	PageID           string                     `json:"page_id"`
	LibraryID        string                     `json:"library_id"`
	Mode             string                     `json:"mode"`
	Revision         int64                      `json:"revision"`
	ProjectionStatus string                     `json:"projection_status"`
	Grants           []KnowledgePageAccessGrant `json:"grants"`
	UpdatedByUserID  string                     `json:"updated_by_user_id"`
	UpdatedAt        time.Time                  `json:"updated_at"`
}

type KnowledgePageEffectiveAccess struct {
	Allowed     bool   `json:"allowed"`
	Role        string `json:"role,omitempty"`
	Mode        string `json:"mode"`
	Source      string `json:"source"` // library, direct_grant, ancestor_grant, manager
	GrantPageID string `json:"grant_page_id,omitempty"`
	GrantScope  string `json:"grant_scope,omitempty"`
}

type KnowledgePageAccessWrite struct {
	ClientID         string                     `json:"client_id"`
	ExpectedRevision int64                      `json:"expected_revision"`
	Mode             string                     `json:"mode"`
	Grants           []KnowledgePageAccessGrant `json:"grants"`
}

type KnowledgePageMovePreview struct {
	ParentPageID string `json:"parent_page_id,omitempty"`
	Rank         string `json:"rank"`
}

type KnowledgePageMoveImpact struct {
	PageID          string   `json:"page_id"`
	ParentPageID    string   `json:"parent_page_id,omitempty"`
	AffectedPages   int      `json:"affected_pages"`
	ChangedUserIDs  []string `json:"changed_user_ids"`
	RequiresConfirm bool     `json:"requires_confirm"`
	Token           string   `json:"token"`
}

type KnowledgePageAccessService interface {
	KnowledgePageAccess(context.Context, string, string, agentsdk.ConversationAuthority) (KnowledgePageAccessPolicy, error)
	SetKnowledgePageAccess(context.Context, string, string, KnowledgePageAccessWrite, agentsdk.ConversationAuthority) (KnowledgePageAccessPolicy, error)
}

// KnowledgePageAccessRepository is the trusted server-side policy boundary
// used by application authorization and retrieval. Permission tokens are
// derived by Knowledge and are never accepted from browser input.
type KnowledgePageAccessRepository interface {
	KnowledgePageAccessService
	EvaluateKnowledgePageAccess(context.Context, string, string, agentsdk.ConversationAuthority) (KnowledgePageEffectiveAccess, error)
	KnowledgePagePermissionTokens(context.Context, string, agentsdk.ConversationAuthority) ([]string, error)
	KnowledgePageProjectionPermission(context.Context, string, string, agentsdk.ConversationAuthority) (string, int64, error)
}

// JSONClone provides a small contract-test helper without defining a second
// wire representation for page values.
func JSONClone[T any](value T) (T, error) {
	var out T
	raw, err := json.Marshal(value)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(raw, &out)
	return out, err
}
