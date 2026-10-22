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
	ClientID        string    `json:"client_id"`
	UpdateID        string    `json:"update_id"`
	ClientCreatedAt time.Time `json:"client_created_at,omitempty"`
	Title           *string   `json:"title,omitempty"`
	YjsUpdate       string    `json:"yjs_update,omitempty"`
	YjsState        string    `json:"yjs_state,omitempty"`
	Content         string    `json:"content"`
	ContentSHA256   string    `json:"content_sha256"`
	EditedBlockIDs  []string  `json:"edited_block_ids,omitempty"`
}

type KnowledgeCollaborationSnapshot struct {
	PageID         string                            `json:"page_id"`
	DraftRevision  int64                             `json:"draft_revision"`
	ServerClock    int64                             `json:"server_clock"`
	Title          string                            `json:"title"`
	Content        string                            `json:"content"`
	ContentSHA256  string                            `json:"content_sha256"`
	YjsState       string                            `json:"yjs_state,omitempty"`
	YjsStateClock  int64                             `json:"yjs_state_clock"`
	YjsUpdateBytes int64                             `json:"yjs_update_bytes_since_checkpoint"`
	YjsUpdates     []KnowledgeCollaborationYjsUpdate `json:"yjs_updates,omitempty"`
}

// KnowledgeCollaborationYjsUpdate is one durable Yjs delta after the latest
// checkpoint. Snapshot responses carry a bounded, clock-ordered replay so a
// joining client can restore current CRDT state without uploading a full state
// on every edit.
type KnowledgeCollaborationYjsUpdate struct {
	Clock     int64  `json:"clock"`
	UpdateID  string `json:"update_id"`
	YjsUpdate string `json:"yjs_update"`
}

type KnowledgeCollaborationPresence struct {
	SessionID        string    `json:"session_id"`
	UserID           string    `json:"user_id"`
	DisplayName      string    `json:"display_name,omitempty"`
	AwarenessUpdate  string    `json:"awareness_update,omitempty"`
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
	AwarenessUpdate  string                        `json:"awareness_update,omitempty"`
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
	SessionID         string `json:"session_id"`
	AfterClock        int64  `json:"after_clock"`
	WaitMillis        int64  `json:"wait_millis,omitempty"`
	ConnectionAttempt int64  `json:"connection_attempt,omitempty"`
}

type KnowledgeCollaborationEventBatch struct {
	SessionID   string                              `json:"session_id"`
	ServerClock int64                               `json:"server_clock"`
	Events      []KnowledgeCollaborationServerEvent `json:"events"`
	Complete    bool                                `json:"complete"`
}

type KnowledgePageLiveEventType string

const (
	KnowledgePageLiveEventChanged   KnowledgePageLiveEventType = "changed"
	KnowledgePageLiveEventHeartbeat KnowledgePageLiveEventType = "heartbeat"
	KnowledgePageLiveEventRevoked   KnowledgePageLiveEventType = "revoked"
)

// KnowledgePageLiveEvent is deliberately content-free. Readers use it only as
// a signal to repeat the normal authorized page read, so this stream can never
// expose draft content or collaboration presence to a caller that lost access.
type KnowledgePageLiveEvent struct {
	Type     KnowledgePageLiveEventType `json:"type"`
	PageID   string                     `json:"page_id"`
	Revision int64                      `json:"revision"`
	Code     string                     `json:"code,omitempty"`
}

type KnowledgePageLiveEvents struct {
	AfterRevision int64 `json:"after_revision"`
	WaitMillis    int64 `json:"wait_millis,omitempty"`
}

type KnowledgePageLiveEventBatch struct {
	PageID   string                   `json:"page_id"`
	Revision int64                    `json:"revision"`
	Events   []KnowledgePageLiveEvent `json:"events"`
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

type KnowledgePageLiveEventService interface {
	KnowledgePageLiveEvents(context.Context, string, string, KnowledgePageLiveEvents, Authority) (KnowledgePageLiveEventBatch, error)
}

// KnowledgeCollaborationNotifier is an optional process-external wake-up bus.
// Messages carry only bounded scope/page keys. Durable collaboration rows remain
// the replay and recovery source of truth when a notification is duplicated or
// lost.
type KnowledgeCollaborationNotifier interface {
	PublishKnowledgeCollaborationWake(context.Context, []string) error
	SubscribeKnowledgeCollaborationWake(context.Context, []string) (<-chan struct{}, func(), error)
}
