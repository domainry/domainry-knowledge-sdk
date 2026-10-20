package knowledgesdk

import (
	"context"
	"time"
)

// KnowledgeCollaborationUpdate carries a native Yjs update plus the canonical
// Block document produced from that same Y.Doc. The server owns authorization,
// ordering, schema validation, version boundaries and persistence; Yjs owns
// concurrent tree/text convergence in the editor.
type KnowledgeCollaborationUpdate struct {
	ClientID       string   `json:"client_id"`
	UpdateID       string   `json:"update_id"`
	Title          *string  `json:"title,omitempty"`
	YjsUpdate      string   `json:"yjs_update,omitempty"`
	YjsState       string   `json:"yjs_state,omitempty"`
	Content        string   `json:"content"`
	ContentSHA256  string   `json:"content_sha256"`
	EditedBlockIDs []string `json:"edited_block_ids,omitempty"`
}

type KnowledgeCollaborationSnapshot struct {
	PageID        string `json:"page_id"`
	DraftRevision int64  `json:"draft_revision"`
	ServerClock   int64  `json:"server_clock"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	ContentSHA256 string `json:"content_sha256"`
	YjsState      string `json:"yjs_state,omitempty"`
}

type KnowledgeCollaborationPresence struct {
	SessionID        string    `json:"session_id"`
	UserID           string    `json:"user_id"`
	DisplayName      string    `json:"display_name,omitempty"`
	CursorBlockID    string    `json:"cursor_block_id,omitempty"`
	SelectionBlockID string    `json:"selection_block_id,omitempty"`
	EditedBlocks     int       `json:"edited_blocks"`
	LastSeenAt       time.Time `json:"last_seen_at"`
}

type KnowledgeCollaborationSession struct {
	ID          string                           `json:"id"`
	LibraryID   string                           `json:"library_id"`
	PageID      string                           `json:"page_id"`
	ClientID    string                           `json:"client_id"`
	UserID      string                           `json:"user_id"`
	State       string                           `json:"state"` // connected, revoked, closed, expired
	ServerClock int64                            `json:"server_clock"`
	ExpiresAt   time.Time                        `json:"expires_at"`
	Snapshot    KnowledgeCollaborationSnapshot   `json:"snapshot"`
	Presence    []KnowledgeCollaborationPresence `json:"presence"`
}

type KnowledgeCollaborationConnect struct {
	ClientID string `json:"client_id"`
}

type KnowledgeCollaborationSync struct {
	SessionID        string                        `json:"session_id"`
	AfterClock       int64                         `json:"after_clock"`
	Update           *KnowledgeCollaborationUpdate `json:"update,omitempty"`
	CursorBlockID    string                        `json:"cursor_block_id,omitempty"`
	SelectionBlockID string                        `json:"selection_block_id,omitempty"`
}

type KnowledgeCollaborationUpdateEvent struct {
	Clock    int64                           `json:"clock"`
	UserID   string                          `json:"user_id,omitempty"`
	Update   *KnowledgeCollaborationUpdate   `json:"update,omitempty"`
	Snapshot *KnowledgeCollaborationSnapshot `json:"snapshot,omitempty"`
}

type KnowledgeCollaborationSyncResult struct {
	Session  KnowledgeCollaborationSession       `json:"session"`
	Events   []KnowledgeCollaborationUpdateEvent `json:"events"`
	Applied  bool                                `json:"applied"`
	Complete bool                                `json:"complete"`
}

type KnowledgeCollaborationEventType string

const (
	KnowledgeCollaborationEventUpdate    KnowledgeCollaborationEventType = "update"
	KnowledgeCollaborationEventSnapshot  KnowledgeCollaborationEventType = "snapshot"
	KnowledgeCollaborationEventPresence  KnowledgeCollaborationEventType = "presence"
	KnowledgeCollaborationEventRevoked   KnowledgeCollaborationEventType = "revoked"
	KnowledgeCollaborationEventHeartbeat KnowledgeCollaborationEventType = "heartbeat"
)

// KnowledgeCollaborationServerEvent is the resumable server-to-client event
// contract. Clock is the durable page clock used as the SSE id. Presence does
// not advance that clock and is therefore delivered alongside the current id.
type KnowledgeCollaborationServerEvent struct {
	Type     KnowledgeCollaborationEventType  `json:"type"`
	Clock    int64                            `json:"clock"`
	UserID   string                           `json:"user_id,omitempty"`
	Update   *KnowledgeCollaborationUpdate    `json:"update,omitempty"`
	Snapshot *KnowledgeCollaborationSnapshot  `json:"snapshot,omitempty"`
	Presence []KnowledgeCollaborationPresence `json:"presence,omitempty"`
	Code     string                           `json:"code,omitempty"`
}

// KnowledgeCollaborationEvents asks for events after one durable clock. The
// bounded wait keeps the Service usable through both embedded and SaaS
// transports; the public HTTP adapter turns successive batches into one SSE
// stream. WaitMillis is capped by the implementation.
type KnowledgeCollaborationEvents struct {
	SessionID  string `json:"session_id"`
	AfterClock int64  `json:"after_clock"`
	WaitMillis int64  `json:"wait_millis,omitempty"`
}

type KnowledgeCollaborationEventBatch struct {
	SessionID   string                              `json:"session_id"`
	ServerClock int64                               `json:"server_clock"`
	Events      []KnowledgeCollaborationServerEvent `json:"events"`
	Complete    bool                                `json:"complete"`
}

type KnowledgeCollaborationPublish struct {
	SessionID   string `json:"session_id"`
	ClientID    string `json:"client_id"`
	PublishNote string `json:"publish_note,omitempty"`
}

type KnowledgeRealtimeCollaborationService interface {
	ConnectKnowledgeCollaboration(context.Context, string, string, KnowledgeCollaborationConnect, Authority) (KnowledgeCollaborationSession, error)
	SyncKnowledgeCollaboration(context.Context, string, string, KnowledgeCollaborationSync, Authority) (KnowledgeCollaborationSyncResult, error)
	KnowledgeCollaborationEvents(context.Context, string, string, KnowledgeCollaborationEvents, Authority) (KnowledgeCollaborationEventBatch, error)
	CloseKnowledgeCollaboration(context.Context, string, string, string, Authority) error
	PublishKnowledgeCollaboration(context.Context, string, string, KnowledgeCollaborationPublish, Authority) (KnowledgePageView, error)
}
