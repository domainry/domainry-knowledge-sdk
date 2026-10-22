package knowledgesdk

import (
	"context"
	"time"
)

const (
	KnowledgeTaskStatusTodo       = "todo"
	KnowledgeTaskStatusInProgress = "in_progress"
	KnowledgeTaskStatusDone       = "done"
	KnowledgeTaskStatusCanceled   = "canceled"

	KnowledgeTaskPriorityNone   = "none"
	KnowledgeTaskPriorityLow    = "low"
	KnowledgeTaskPriorityMedium = "medium"
	KnowledgeTaskPriorityHigh   = "high"
	KnowledgeTaskPriorityUrgent = "urgent"
)

// KnowledgePageTask is live business state referenced by one stable document
// Block. It is intentionally not embedded in KnowledgePageVersion: completing
// or assigning a task must not create or rewrite an immutable document version.
type KnowledgePageTask struct {
	ID              string     `json:"id"`
	LibraryID       string     `json:"library_id"`
	PageID          string     `json:"page_id"`
	SourceBlockID   string     `json:"source_block_id"`
	Title           string     `json:"title"`
	Status          string     `json:"status"`
	AssigneeUserIDs []string   `json:"assignee_user_ids"`
	DueAt           *time.Time `json:"due_at,omitempty"`
	ReminderAt      *time.Time `json:"reminder_at,omitempty"`
	Priority        string     `json:"priority"`
	Revision        int64      `json:"revision"`
	CreatedByUserID string     `json:"created_by_user_id"`
	UpdatedByUserID string     `json:"updated_by_user_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type KnowledgePageTaskCreate struct {
	ClientID        string     `json:"client_id"`
	SourceBlockID   string     `json:"source_block_id"`
	Title           string     `json:"title"`
	AssigneeUserIDs []string   `json:"assignee_user_ids"`
	DueAt           *time.Time `json:"due_at,omitempty"`
	ReminderAt      *time.Time `json:"reminder_at,omitempty"`
	Priority        string     `json:"priority"`
}

type KnowledgePageTaskUpdate struct {
	ClientID         string     `json:"client_id"`
	ExpectedRevision int64      `json:"expected_revision"`
	Title            string     `json:"title"`
	Status           string     `json:"status"`
	AssigneeUserIDs  []string   `json:"assignee_user_ids"`
	DueAt            *time.Time `json:"due_at,omitempty"`
	ReminderAt       *time.Time `json:"reminder_at,omitempty"`
	Priority         string     `json:"priority"`
}

type KnowledgePageTaskPage struct {
	Items []KnowledgePageTask `json:"items"`
}

type PageTaskService interface {
	KnowledgePageTasks(context.Context, string, string, Authority) (KnowledgePageTaskPage, error)
	CreateKnowledgePageTask(context.Context, string, string, KnowledgePageTaskCreate, Authority) (KnowledgePageTask, error)
	UpdateKnowledgePageTask(context.Context, string, string, string, KnowledgePageTaskUpdate, Authority) (KnowledgePageTask, error)
}
