package knowledgesdk

import (
	"context"
	"time"
)

type KnowledgePageCommentAnchor struct {
	NodeID  string `json:"node_id,omitempty"`
	Quote   string `json:"quote,omitempty"`
	Changed bool   `json:"changed,omitempty"`
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

type KnowledgePageCommentPage struct {
	Items     []KnowledgePageComment `json:"items"`
	NextAfter string                 `json:"next_after,omitempty"`
	Complete  bool                   `json:"complete"`
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

type KnowledgePrincipalDirectory interface {
	SearchKnowledgePrincipals(context.Context, string, string, int, Authority) (KnowledgePrincipalPage, error)
	ResolveKnowledgePrincipals(context.Context, []string, Authority) ([]KnowledgePrincipal, error)
}

type KnowledgeAccessLayer struct {
	Allowed bool   `json:"allowed"`
	Code    string `json:"code"`
	Label   string `json:"label"`
	Source  string `json:"source,omitempty"`
}

type KnowledgeAccessExplanation struct {
	UserID        string               `json:"user_id"`
	LibraryID     string               `json:"library_id"`
	PageID        string               `json:"page_id,omitempty"`
	Action        string               `json:"action"`
	Identity      KnowledgeAccessLayer `json:"identity"`
	Library       KnowledgeAccessLayer `json:"library"`
	Page          KnowledgeAccessLayer `json:"page"`
	Final         KnowledgeAccessLayer `json:"final"`
	LibraryRole   string               `json:"library_role,omitempty"`
	InheritedFrom string               `json:"inherited_from,omitempty"`
	ImpactPages   int                  `json:"impact_pages"`
	Proposed      bool                 `json:"proposed,omitempty"`
}

type KnowledgeAccessExplainRequest struct {
	UserID       string `json:"user_id"`
	PageID       string `json:"page_id,omitempty"`
	Action       string `json:"action"`
	ProposedRole string `json:"proposed_role,omitempty"`
}

type KnowledgeAccessExplainer interface {
	ExplainKnowledgeAccess(context.Context, KnowledgeAccessExplainRequest, KnowledgeLibrary, Authority) (KnowledgeAccessExplanation, error)
}

type CollaborationService interface {
	CommentService
	PageAttachmentService
	KnowledgePageAccessService
	KnowledgeRealtimeCollaborationService
	SearchKnowledgePrincipals(context.Context, string, string, int, Authority) (KnowledgePrincipalPage, error)
	ResolveKnowledgePrincipals(context.Context, []string, Authority) ([]KnowledgePrincipal, error)
	SearchKnowledgeLibraryMembers(context.Context, string, string, string, int, Authority) (KnowledgePrincipalPage, error)
	ExplainKnowledgeAccess(context.Context, string, KnowledgeAccessExplainRequest, Authority) (KnowledgeAccessExplanation, error)
}
