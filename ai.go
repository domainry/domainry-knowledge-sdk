package knowledgesdk

import (
	"context"
	"time"
)

const (
	KnowledgeAIOperationWrite     = "write"
	KnowledgeAIOperationSummarize = "summarize"
	KnowledgeAIOperationTranslate = "translate"
	KnowledgeAIOperationRewrite   = "rewrite"
	KnowledgeAIOperationAnswer    = "answer"

	KnowledgeAISourcePublished = "published"
	KnowledgeAISourceDraft     = "draft"

	KnowledgeAIExecutionPending   = "pending"
	KnowledgeAIExecutionRunning   = "running"
	KnowledgeAIExecutionCompleted = "completed"
	KnowledgeAIExecutionFailed    = "failed"
	KnowledgeAIExecutionConfirmed = "confirmed"
)

type KnowledgeAISource struct {
	Scope         string                   `json:"scope"`
	PageID        string                   `json:"page_id"`
	Title         string                   `json:"title"`
	Version       int64                    `json:"version,omitempty"`
	DraftRevision int64                    `json:"draft_revision,omitempty"`
	BlockID       string                   `json:"block_id,omitempty"`
	ContentSHA256 string                   `json:"content_sha256"`
	Reference     *KnowledgeAgentReference `json:"reference,omitempty"`
}

type KnowledgeAIExecution struct {
	ID              string            `json:"id"`
	LibraryID       string            `json:"library_id"`
	PageID          string            `json:"page_id"`
	Operation       string            `json:"operation"`
	Prompt          string            `json:"prompt,omitempty"`
	Language        string            `json:"language,omitempty"`
	State           string            `json:"state"`
	Revision        int64             `json:"revision"`
	Attempt         int               `json:"attempt"`
	Source          KnowledgeAISource `json:"source"`
	Result          string            `json:"result,omitempty"`
	Provider        string            `json:"provider,omitempty"`
	Model           string            `json:"model,omitempty"`
	ErrorCode       string            `json:"error_code,omitempty"`
	CreatedByUserID string            `json:"created_by_user_id"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	ConfirmedAt     time.Time         `json:"confirmed_at,omitempty"`
}

type KnowledgeAIExecutionCreate struct {
	ClientID    string `json:"client_id"`
	Operation   string `json:"operation"`
	Prompt      string `json:"prompt,omitempty"`
	Language    string `json:"language,omitempty"`
	SourceScope string `json:"source_scope"`
	BlockID     string `json:"block_id,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
}

type KnowledgeAIExecutionConfirm struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
	SessionID        string `json:"session_id"`
}

type KnowledgeAIGenerationInput struct {
	ExecutionID string            `json:"execution_id"`
	Operation   string            `json:"operation"`
	Prompt      string            `json:"prompt,omitempty"`
	Language    string            `json:"language,omitempty"`
	Source      KnowledgeAISource `json:"source"`
	SourceText  string            `json:"source_text"`
}

type KnowledgeAIGenerationResult struct {
	Text     string `json:"text"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

type KnowledgeAIProvider interface {
	GenerateKnowledgeAI(context.Context, KnowledgeAIGenerationInput) (KnowledgeAIGenerationResult, error)
}

type KnowledgeAIService interface {
	CreateKnowledgeAIExecution(context.Context, string, string, KnowledgeAIExecutionCreate, Authority) (KnowledgeAIExecution, error)
	KnowledgeAIExecution(context.Context, string, string, string, Authority) (KnowledgeAIExecution, error)
	ConfirmKnowledgeAIExecution(context.Context, string, string, string, KnowledgeAIExecutionConfirm, Authority) (KnowledgeAIExecution, error)
}

type KnowledgeAIExecutionWork struct {
	ScopeKey       string               `json:"-"`
	Execution      KnowledgeAIExecution `json:"execution"`
	Authority      Authority            `json:"-"`
	SessionID      string               `json:"-"`
	Fence          int64                `json:"fence"`
	LeaseExpiresAt time.Time            `json:"lease_expires_at"`
}

type KnowledgeAIRepository interface {
	KnowledgeAIService
	ClaimKnowledgeAIExecution(context.Context, string, time.Time, time.Duration) (KnowledgeAIExecutionWork, bool, error)
	ReadKnowledgeAIExecutionSource(context.Context, KnowledgeAIExecutionWork, Authority) (KnowledgeAIGenerationInput, error)
	CompleteKnowledgeAIExecution(context.Context, KnowledgeAIExecutionWork, KnowledgeAIGenerationResult, string) error
}
