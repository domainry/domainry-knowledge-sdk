package knowledgesdk

import (
	"context"
	"time"
)

type KnowledgePageProjectionLease struct {
	RuntimeID      string    `json:"runtime_id"`
	WorkspaceID    string    `json:"workspace_id"`
	LibraryID      string    `json:"library_id"`
	PageID         string    `json:"page_id"`
	Version        int64     `json:"version"`
	AccessRevision int64     `json:"access_revision"`
	PermissionID   string    `json:"permission_id,omitempty"`
	Owner          string    `json:"owner"`
	Token          int64     `json:"token"`
	ExpiresAt      time.Time `json:"expires_at"`
	Actor          Authority `json:"actor"`
}

type KnowledgePageProjectionWork struct {
	Lease   KnowledgePageProjectionLease `json:"lease"`
	Page    KnowledgePage                `json:"page"`
	Version KnowledgePageVersion         `json:"version"`
}

type KnowledgePageProjectionRepository interface {
	ClaimKnowledgePageProjection(context.Context, string, string, time.Time, time.Duration) (KnowledgePageProjectionWork, bool, error)
	CompleteKnowledgePageProjection(context.Context, KnowledgePageProjectionLease, string, int64, string, string, string) error
	RetryKnowledgePageProjection(context.Context, KnowledgePageProjectionLease, time.Time) error
}
