package knowledgesdk

import (
	"context"
	"time"
)

const (
	KnowledgeWebhookPagePublished     = "knowledge.page.published"
	KnowledgeWebhookCommentMentioned  = "knowledge.page_comment.mentioned"
	KnowledgeWebhookTaskChanged       = "knowledge.page_task.changed"
	KnowledgeWebhookPermissionChanged = "knowledge.page_permission.changed"
	KnowledgeWebhookIndexCompleted    = "knowledge.page_index.completed"

	KnowledgeWebhookSubscriptionActive  = "active"
	KnowledgeWebhookSubscriptionRevoked = "revoked"
	KnowledgeWebhookDeliveryPending     = "pending"
	KnowledgeWebhookDeliveryDelivering  = "delivering"
	KnowledgeWebhookDeliveryDelivered   = "delivered"
	KnowledgeWebhookDeliveryRevoked     = "revoked"
	KnowledgeWebhookDeliveryFailed      = "failed"
)

type KnowledgeWebhookSubscription struct {
	ID              string    `json:"id"`
	LibraryID       string    `json:"library_id"`
	ApplicationKey  string    `json:"application_key"`
	EndpointURL     string    `json:"endpoint_url"`
	EventTypes      []string  `json:"event_types"`
	State           string    `json:"state"`
	Revision        int64     `json:"revision"`
	CreatedByUserID string    `json:"created_by_user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type KnowledgeWebhookSubscriptionCreate struct {
	ClientID       string   `json:"client_id"`
	ApplicationKey string   `json:"application_key"`
	EndpointURL    string   `json:"endpoint_url"`
	EventTypes     []string `json:"event_types"`
}

type KnowledgeWebhookSubscriptionUpdate struct {
	ClientID         string   `json:"client_id"`
	ExpectedRevision int64    `json:"expected_revision"`
	EndpointURL      string   `json:"endpoint_url"`
	EventTypes       []string `json:"event_types"`
	State            string   `json:"state"`
}

type KnowledgeWebhookSubscriptionPage struct {
	Items []KnowledgeWebhookSubscription `json:"items"`
}

type KnowledgeWebhookEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Action      string    `json:"action"`
	WorkspaceID string    `json:"workspace_id"`
	LibraryID   string    `json:"library_id"`
	PageID      string    `json:"page_id,omitempty"`
	ResourceID  string    `json:"resource_id,omitempty"`
	Version     int64     `json:"version,omitempty"`
	ActorUserID string    `json:"actor_user_id,omitempty"`
	OccurredAt  time.Time `json:"occurred_at"`
	Authority   Authority `json:"-"`
}

type KnowledgeWebhookDelivery struct {
	EventID        string    `json:"event_id"`
	SubscriptionID string    `json:"subscription_id"`
	State          string    `json:"state"`
	Attempt        int       `json:"attempt"`
	NextAttemptAt  time.Time `json:"next_attempt_at"`
	LastStatusCode int       `json:"last_status_code,omitempty"`
	LastErrorCode  string    `json:"last_error_code,omitempty"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type KnowledgeWebhookEventRecord struct {
	Event      KnowledgeWebhookEvent      `json:"event"`
	Deliveries []KnowledgeWebhookDelivery `json:"deliveries"`
}

type KnowledgeWebhookEventPage struct {
	Items []KnowledgeWebhookEventRecord `json:"items"`
}

type KnowledgeWebhookWork struct {
	ScopeKey       string                       `json:"-"`
	Event          KnowledgeWebhookEvent        `json:"event"`
	Subscription   KnowledgeWebhookSubscription `json:"subscription"`
	Delivery       KnowledgeWebhookDelivery     `json:"delivery"`
	Fence          int64                        `json:"fence"`
	LeaseExpiresAt time.Time                    `json:"lease_expires_at"`
}

type KnowledgeWebhookDeliveryResult struct {
	Delivered  bool
	StatusCode int
	ErrorCode  string
	RetryAt    time.Time
}

type KnowledgeWebhookDeliveryRequest struct {
	ApplicationKey string
	EndpointURL    string
	Event          KnowledgeWebhookEvent
	IdempotencyKey string
}

type KnowledgeWebhookDeliverer interface {
	DeliverKnowledgeWebhook(context.Context, KnowledgeWebhookDeliveryRequest) (KnowledgeWebhookDeliveryResult, error)
}

type KnowledgeTrustedApplicationResolver interface {
	ResolveKnowledgeTrustedApplication(context.Context, string, string) (Authority, error)
}

type WebhookService interface {
	KnowledgeWebhookSubscriptions(context.Context, string, Authority) (KnowledgeWebhookSubscriptionPage, error)
	CreateKnowledgeWebhookSubscription(context.Context, string, KnowledgeWebhookSubscriptionCreate, Authority) (KnowledgeWebhookSubscription, error)
	UpdateKnowledgeWebhookSubscription(context.Context, string, string, KnowledgeWebhookSubscriptionUpdate, Authority) (KnowledgeWebhookSubscription, error)
	KnowledgeWebhookEvents(context.Context, string, Authority) (KnowledgeWebhookEventPage, error)
	ReplayKnowledgeWebhookEvent(context.Context, string, string, Authority) (KnowledgeWebhookEventRecord, error)
}

type KnowledgeWebhookRepository interface {
	WebhookService
	AppendKnowledgeWebhookEvent(context.Context, KnowledgeWebhookEvent) error
	ClaimKnowledgeWebhookDelivery(context.Context, string, time.Time, time.Duration) (KnowledgeWebhookWork, bool, error)
	CompleteKnowledgeWebhookDelivery(context.Context, KnowledgeWebhookWork, KnowledgeWebhookDeliveryResult) error
}
