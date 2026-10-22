package knowledgesdk

import (
	"context"
	"time"
)

const (
	KnowledgePageShareOrganization = "organization"
	KnowledgePageShareUser         = "user"
	KnowledgePageShareExternal     = "external"
	KnowledgePageShareScopePage    = "page"
	KnowledgePageShareScopeSubtree = "subtree"
	KnowledgePageShareActive       = "active"
	KnowledgePageShareRevoked      = "revoked"
)

// KnowledgePageShare is a non-secret locator and a server-side access policy.
// A share never contains a reusable bearer credential: redemption is bound to
// the authenticated Identity user and every read re-evaluates Identity plus
// the current share state and expiry.
type KnowledgePageShare struct {
	ID                string                 `json:"id"`
	LibraryID         string                 `json:"library_id"`
	PageID            string                 `json:"page_id"`
	Kind              string                 `json:"kind"`
	Target            KnowledgeAccessSubject `json:"target,omitempty"`
	Scope             string                 `json:"scope"`
	Role              string                 `json:"role"`
	PasswordProtected bool                   `json:"password_protected"`
	State             string                 `json:"state"`
	ExpiresAt         time.Time              `json:"expires_at"`
	Revision          int64                  `json:"revision"`
	CreatedByUserID   string                 `json:"created_by_user_id"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type KnowledgePageShareCreate struct {
	ClientID  string                 `json:"client_id"`
	Kind      string                 `json:"kind"`
	Target    KnowledgeAccessSubject `json:"target,omitempty"`
	Scope     string                 `json:"scope"`
	ExpiresAt time.Time              `json:"expires_at"`
	Password  string                 `json:"password,omitempty"`
}

type KnowledgePageShareRedeem struct {
	Password string `json:"password,omitempty"`
}

type KnowledgePageShareRevoke struct {
	ExpectedRevision int64 `json:"expected_revision"`
}

type KnowledgePageSharePage struct {
	Items     []KnowledgePageShare `json:"items"`
	NextAfter string               `json:"next_after,omitempty"`
	Complete  bool                 `json:"complete"`
}

type KnowledgePageSharedView struct {
	Share KnowledgePageShare `json:"share"`
	View  KnowledgePageView  `json:"view"`
}

type PageShareService interface {
	KnowledgePageShares(context.Context, string, string, string, int, Authority) (KnowledgePageSharePage, error)
	CreateKnowledgePageShare(context.Context, string, string, KnowledgePageShareCreate, Authority) (KnowledgePageShare, error)
	RevokeKnowledgePageShare(context.Context, string, string, string, int64, Authority) (KnowledgePageShare, error)
	RedeemKnowledgePageShare(context.Context, string, KnowledgePageShareRedeem, Authority) (KnowledgePageShare, error)
	ReadKnowledgePageShare(context.Context, string, string, Authority) (KnowledgePageSharedView, error)
	AuthorizeKnowledgePageShareContentOperation(context.Context, string, string, string, Authority) (KnowledgePageContentDecision, error)
}
