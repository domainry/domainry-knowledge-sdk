package knowledgesdk

import (
	"context"
	"time"
)

type KnowledgePageCommentAnchor struct {
	NodeID      string `json:"node_id,omitempty"`
	Quote       string `json:"quote,omitempty"`
	BlockSHA256 string `json:"block_sha256,omitempty"`
	Changed     bool   `json:"changed,omitempty"`
}

type KnowledgePageComment struct {
	ID               string                     `json:"id"`
	PageID           string                     `json:"page_id"`
	LibraryID        string                     `json:"library_id"`
	ParentCommentID  string                     `json:"parent_comment_id,omitempty"`
	PublishedVersion int64                      `json:"published_version"`
	Anchor           KnowledgePageCommentAnchor `json:"anchor"`
	Body             string                     `json:"body"`
	MentionedUserIDs []string                   `json:"mentioned_user_ids"`
	AuthorUserID     string                     `json:"author_user_id"`
	Resolved         bool                       `json:"resolved"`
	ResolvedByUserID string                     `json:"resolved_by_user_id,omitempty"`
	Deleted          bool                       `json:"deleted"`
	Revision         int64                      `json:"revision"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
	Reactions        []KnowledgeCommentReaction `json:"reactions"`
}

type KnowledgeCommentReaction struct {
	Emoji   string `json:"emoji"`
	Count   int64  `json:"count"`
	Reacted bool   `json:"reacted"`
}

type KnowledgePageCommentCreate struct {
	ClientID         string                     `json:"client_id"`
	ParentCommentID  string                     `json:"parent_comment_id,omitempty"`
	PublishedVersion int64                      `json:"published_version"`
	Anchor           KnowledgePageCommentAnchor `json:"anchor"`
	Body             string                     `json:"body"`
	MentionedUserIDs []string                   `json:"mentioned_user_ids"`
}

type KnowledgePageCommentUpdate struct {
	ClientID         string   `json:"client_id"`
	ExpectedRevision int64    `json:"expected_revision"`
	Body             string   `json:"body"`
	MentionedUserIDs []string `json:"mentioned_user_ids"`
}

type KnowledgePageCommentResolve struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
	Resolved         bool   `json:"resolved"`
}

type KnowledgePageCommentDelete struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
}

type KnowledgePageCommentReactionWrite struct {
	ClientID string `json:"client_id"`
	Emoji    string `json:"emoji"`
	Active   bool   `json:"active"`
}

type KnowledgePageCommentPage struct {
	Items     []KnowledgePageComment `json:"items"`
	NextAfter string                 `json:"next_after,omitempty"`
	Complete  bool                   `json:"complete"`
}

type KnowledgePageBlockChange struct {
	NodeID        string `json:"node_id"`
	Kind          string `json:"kind"` // insert, update, delete
	BeforeSHA256  string `json:"before_sha256,omitempty"`
	AfterSHA256   string `json:"after_sha256,omitempty"`
	BeforePreview string `json:"before_preview,omitempty"`
	AfterPreview  string `json:"after_preview,omitempty"`
}

type KnowledgePageEditHistoryEntry struct {
	Clock             int64                      `json:"clock"`
	DraftRevision     int64                      `json:"draft_revision"`
	UserID            string                     `json:"user_id"`
	TitleChanged      bool                       `json:"title_changed"`
	ChangedBlockCount int                        `json:"changed_block_count"`
	Changes           []KnowledgePageBlockChange `json:"changes"`
	Truncated         bool                       `json:"truncated"`
	CreatedAt         time.Time                  `json:"created_at"`
}

type KnowledgePageEditHistoryPage struct {
	Items      []KnowledgePageEditHistoryEntry `json:"items"`
	NextBefore int64                           `json:"next_before,omitempty"`
	Complete   bool                            `json:"complete"`
}

const (
	KnowledgeSuggestionInsert = "insert"
	KnowledgeSuggestionDelete = "delete"

	KnowledgeSuggestionOpen      = "open"
	KnowledgeSuggestionAccepted  = "accepted"
	KnowledgeSuggestionRejected  = "rejected"
	KnowledgeSuggestionWithdrawn = "withdrawn"

	KnowledgeSuggestionDecisionAccept   = "accept"
	KnowledgeSuggestionDecisionReject   = "reject"
	KnowledgeSuggestionDecisionWithdraw = "withdraw"
)

type KnowledgePageSuggestion struct {
	ID                string                     `json:"id"`
	LibraryID         string                     `json:"library_id"`
	PageID            string                     `json:"page_id"`
	Kind              string                     `json:"kind"`
	Status            string                     `json:"status"`
	BaseDraftRevision int64                      `json:"base_draft_revision"`
	Anchor            KnowledgePageCommentAnchor `json:"anchor"`
	ProposedBlockID   string                     `json:"proposed_block_id,omitempty"`
	ProposedText      string                     `json:"proposed_text,omitempty"`
	Comment           string                     `json:"comment,omitempty"`
	MentionedUserIDs  []string                   `json:"mentioned_user_ids"`
	AuthorUserID      string                     `json:"author_user_id"`
	DecidedByUserID   string                     `json:"decided_by_user_id,omitempty"`
	Revision          int64                      `json:"revision"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}

type KnowledgePageSuggestionCreate struct {
	ClientID          string                     `json:"client_id"`
	Kind              string                     `json:"kind"`
	BaseDraftRevision int64                      `json:"base_draft_revision"`
	Anchor            KnowledgePageCommentAnchor `json:"anchor"`
	ProposedText      string                     `json:"proposed_text,omitempty"`
	Comment           string                     `json:"comment,omitempty"`
	MentionedUserIDs  []string                   `json:"mentioned_user_ids"`
}

type KnowledgePageSuggestionDecision struct {
	ClientID              string `json:"client_id"`
	ExpectedRevision      int64  `json:"expected_revision"`
	ExpectedDraftRevision int64  `json:"expected_draft_revision"`
	Decision              string `json:"decision"` // accept, reject, withdraw
}

type KnowledgePageSuggestionPage struct {
	Items []KnowledgePageSuggestion `json:"items"`
}

type KnowledgePageAttachment struct {
	PageID          string    `json:"page_id"`
	LibraryID       string    `json:"library_id"`
	DocumentID      string    `json:"document_id"`
	Title           string    `json:"title"`
	Rank            string    `json:"rank"`
	CreatedByUserID string    `json:"created_by_user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type KnowledgePageAttachmentWrite struct {
	ClientID string `json:"client_id"`
	Title    string `json:"title"`
	Rank     string `json:"rank"`
}

type KnowledgePageAttachmentRemove struct {
	ClientID string `json:"client_id"`
}

type KnowledgePageAttachmentPage struct {
	Items []KnowledgePageAttachment `json:"items"`
}

type CommentService interface {
	KnowledgePageComments(context.Context, string, string, string, int, bool, Authority) (KnowledgePageCommentPage, error)
	CreateKnowledgePageComment(context.Context, string, string, KnowledgePageCommentCreate, Authority) (KnowledgePageComment, error)
	UpdateKnowledgePageComment(context.Context, string, string, string, KnowledgePageCommentUpdate, Authority) (KnowledgePageComment, error)
	ResolveKnowledgePageComment(context.Context, string, string, string, KnowledgePageCommentResolve, Authority) (KnowledgePageComment, error)
	DeleteKnowledgePageComment(context.Context, string, string, string, KnowledgePageCommentDelete, Authority) (KnowledgePageComment, error)
	SetKnowledgePageCommentReaction(context.Context, string, string, string, KnowledgePageCommentReactionWrite, Authority) (KnowledgePageComment, error)
}

type PageEditHistoryService interface {
	KnowledgePageEditHistory(context.Context, string, string, int64, int, Authority) (KnowledgePageEditHistoryPage, error)
}

type PageSuggestionService interface {
	KnowledgePageSuggestions(context.Context, string, string, Authority) (KnowledgePageSuggestionPage, error)
	CreateKnowledgePageSuggestion(context.Context, string, string, KnowledgePageSuggestionCreate, Authority) (KnowledgePageSuggestion, error)
	DecideKnowledgePageSuggestion(context.Context, string, string, string, KnowledgePageSuggestionDecision, Authority) (KnowledgePageSuggestion, error)
}

type PageAttachmentService interface {
	KnowledgePageAttachments(context.Context, string, string, Authority) (KnowledgePageAttachmentPage, error)
	SetKnowledgePageAttachment(context.Context, string, string, string, KnowledgePageAttachmentWrite, Authority) (KnowledgePageAttachment, error)
	RemoveKnowledgePageAttachment(context.Context, string, string, string, KnowledgePageAttachmentRemove, Authority) error
}

type KnowledgePrincipal struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	Status      string `json:"status"`
}

type KnowledgePrincipalPage struct {
	Items     []KnowledgePrincipal `json:"items"`
	NextAfter string               `json:"next_after,omitempty"`
	Complete  bool                 `json:"complete"`
}

type KnowledgePrincipalResolveRequest struct {
	UserIDs []string `json:"user_ids"`
}

type KnowledgeOrganizationPrincipal struct {
	OrganizationID string `json:"organization_id"`
	DisplayName    string `json:"display_name"`
	NodeType       string `json:"node_type,omitempty"`
	Status         string `json:"status"`
}

type KnowledgeOrganizationPrincipalResolveRequest struct {
	OrganizationIDs []string `json:"organization_ids"`
}

type KnowledgePrincipalDirectory interface {
	SearchKnowledgePrincipals(context.Context, string, string, int, Authority) (KnowledgePrincipalPage, error)
	ResolveKnowledgePrincipals(context.Context, []string, Authority) ([]KnowledgePrincipal, error)
	ResolveKnowledgeOrganizationPrincipals(context.Context, []string, Authority) ([]KnowledgeOrganizationPrincipal, error)
}

type KnowledgeAccessLayer struct {
	Allowed bool   `json:"allowed"`
	Code    string `json:"code"`
	Label   string `json:"label"`
	Source  string `json:"source,omitempty"`
}

type KnowledgeAccessExplanation struct {
	UserID        string                `json:"user_id"`
	LibraryID     string                `json:"library_id"`
	PageID        string                `json:"page_id,omitempty"`
	Action        string                `json:"action"`
	Identity      KnowledgeAccessLayer  `json:"identity"`
	Library       KnowledgeAccessLayer  `json:"library"`
	Page          KnowledgeAccessLayer  `json:"page"`
	Share         *KnowledgeAccessLayer `json:"share,omitempty"`
	Final         KnowledgeAccessLayer  `json:"final"`
	ShareID       string                `json:"share_id,omitempty"`
	LibraryRole   string                `json:"library_role,omitempty"`
	InheritedFrom string                `json:"inherited_from,omitempty"`
	ImpactPages   int                   `json:"impact_pages"`
	Proposed      bool                  `json:"proposed,omitempty"`
}

type KnowledgeAccessExplainRequest struct {
	UserID       string `json:"user_id"`
	PageID       string `json:"page_id,omitempty"`
	Action       string `json:"action"`
	ShareID      string `json:"share_id,omitempty"`
	ProposedRole string `json:"proposed_role,omitempty"`
}

type KnowledgeAccessExplainer interface {
	ResolveKnowledgeAccessAuthority(context.Context, string, Authority) (Authority, error)
	ExplainKnowledgeAccess(context.Context, KnowledgeAccessExplainRequest, KnowledgeLibrary, Authority) (KnowledgeAccessExplanation, error)
}

type CollaborationService interface {
	CommentService
	PageEditHistoryService
	PageSuggestionService
	PageAttachmentService
	PageTaskService
	PagePollService
	PageFormService
	KnowledgePageAccessService
	KnowledgeRealtimeCollaborationService
	KnowledgePageLiveEventService
	KnowledgePageNavigationService
	KnowledgePageMetadataService
	KnowledgePageDiscoveryService
	KnowledgePageGovernanceAudit(context.Context, KnowledgeGovernanceAuditQuery, Authority) (KnowledgeGovernanceAuditPage, error)
	SearchKnowledgePrincipals(context.Context, string, string, int, Authority) (KnowledgePrincipalPage, error)
	ResolveKnowledgePrincipals(context.Context, []string, Authority) ([]KnowledgePrincipal, error)
	ResolveKnowledgeOrganizationPrincipals(context.Context, []string, Authority) ([]KnowledgeOrganizationPrincipal, error)
	SearchKnowledgeLibraryMembers(context.Context, string, string, string, int, Authority) (KnowledgePrincipalPage, error)
	ExplainKnowledgeAccess(context.Context, string, KnowledgeAccessExplainRequest, Authority) (KnowledgeAccessExplanation, error)
}
