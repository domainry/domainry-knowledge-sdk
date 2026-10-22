package knowledgesdk

import (
	"context"
	"time"
)

type KnowledgePageCollectionKind string

const (
	KnowledgePageCollectionFavorites   KnowledgePageCollectionKind = "favorites"
	KnowledgePageCollectionRecentViews KnowledgePageCollectionKind = "recent_views"
	KnowledgePageCollectionRecentEdits KnowledgePageCollectionKind = "recent_edits"
	KnowledgePageCollectionSubscribed  KnowledgePageCollectionKind = "subscribed"
	KnowledgePageCollectionShortcuts   KnowledgePageCollectionKind = "shortcuts"
	KnowledgePageCollectionTrash       KnowledgePageCollectionKind = "trash"
)

// KnowledgePagePersonalState is private to one caller. Page ACL is still
// evaluated independently on every read; saving a favorite or shortcut never
// creates access to the page.
type KnowledgePagePersonalState struct {
	LibraryID    string    `json:"library_id"`
	PageID       string    `json:"page_id"`
	UserID       string    `json:"user_id"`
	Favorite     bool      `json:"favorite"`
	FavoriteAt   time.Time `json:"favorite_at,omitempty"`
	Subscribed   bool      `json:"subscribed"`
	SubscribedAt time.Time `json:"subscribed_at,omitempty"`
	ShortcutRank string    `json:"shortcut_rank,omitempty"`
	LastViewedAt time.Time `json:"last_viewed_at,omitempty"`
	LastEditedAt time.Time `json:"last_edited_at,omitempty"`
	Revision     int64     `json:"revision"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type KnowledgePagePersonalStateWrite struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
	Favorite         bool   `json:"favorite"`
	Subscribed       bool   `json:"subscribed"`
	ShortcutRank     string `json:"shortcut_rank,omitempty"`
}

type KnowledgePageViewRecord struct {
	ClientID string `json:"client_id"`
}

type KnowledgePageCollectionItem struct {
	Page  KnowledgePage              `json:"page"`
	State KnowledgePagePersonalState `json:"state"`
}

type KnowledgePageCollectionPage struct {
	Items     []KnowledgePageCollectionItem `json:"items"`
	NextAfter string                        `json:"next_after,omitempty"`
	Complete  bool                          `json:"complete"`
}

type KnowledgePageNavigationService interface {
	KnowledgePagePersonalState(context.Context, string, string, Authority) (KnowledgePagePersonalState, error)
	SetKnowledgePagePersonalState(context.Context, string, string, KnowledgePagePersonalStateWrite, Authority) (KnowledgePagePersonalState, error)
	RecordKnowledgePageView(context.Context, string, string, KnowledgePageViewRecord, Authority) (KnowledgePagePersonalState, error)
	KnowledgePageCollection(context.Context, string, KnowledgePageCollectionKind, string, int, Authority) (KnowledgePageCollectionPage, error)
}
