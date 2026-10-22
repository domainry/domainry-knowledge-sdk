package knowledgesdk

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

const KnowledgeAuditEventFamily = "knowledge.collaboration"

type KnowledgeGovernanceEvent struct {
	ID         string            `json:"id"`
	Family     string            `json:"family"`
	Type       string            `json:"type"`
	Action     string            `json:"action"`
	Actor      Authority         `json:"actor"`
	LibraryID  string            `json:"library_id,omitempty"`
	PageID     string            `json:"page_id,omitempty"`
	Version    int64             `json:"version,omitempty"`
	ResourceID string            `json:"resource_id,omitempty"`
	ClientID   string            `json:"client_id,omitempty"`
	Before     map[string]string `json:"before,omitempty"`
	After      map[string]string `json:"after,omitempty"`
	OccurredAt time.Time         `json:"occurred_at"`
}

type KnowledgeAuditAppender interface {
	AppendKnowledgeAudit(context.Context, KnowledgeGovernanceEvent) error
}

type KnowledgeGovernanceAuditQuery struct {
	LibraryID string `json:"library_id"`
	PageID    string `json:"page_id"`
	Event     string `json:"event,omitempty"`
	ActorID   string `json:"actor_id,omitempty"`
	After     string `json:"after,omitempty"`
	Limit     int    `json:"limit"`
}

type KnowledgeGovernanceAuditPage struct {
	Items     []KnowledgeGovernanceEvent `json:"items"`
	NextAfter string                     `json:"next_after,omitempty"`
	Complete  bool                       `json:"complete"`
}

// KnowledgeAuditReader is the narrow read port supplied by the embedding
// host. Audit remains the immutable event owner; Knowledge only asks for its
// own page-scoped family and never reads an Audit database directly.
type KnowledgeAuditReader interface {
	ReadKnowledgeAudit(context.Context, KnowledgeGovernanceAuditQuery, Authority) (KnowledgeGovernanceAuditPage, error)
}

type KnowledgeCollaborationEvent struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"`
	Actor            Authority `json:"actor"`
	LibraryID        string    `json:"library_id"`
	PageID           string    `json:"page_id"`
	Version          int64     `json:"version"`
	CommentID        string    `json:"comment_id,omitempty"`
	TaskID           string    `json:"task_id,omitempty"`
	SuggestionID     string    `json:"suggestion_id,omitempty"`
	RecipientUserIDs []string  `json:"recipient_user_ids"`
	OccurredAt       time.Time `json:"occurred_at"`
}

type KnowledgeCollaborationPublisher interface {
	PublishKnowledgeCollaboration(context.Context, KnowledgeCollaborationEvent) error
}

func StableKnowledgeEventID(parts ...string) string {
	normalized := make([]string, len(parts))
	for index, part := range parts {
		normalized[index] = strings.TrimSpace(part)
	}
	digest := sha256.Sum256([]byte(strings.Join(normalized, "\x00")))
	return "kevt_" + hex.EncodeToString(digest[:16])
}
