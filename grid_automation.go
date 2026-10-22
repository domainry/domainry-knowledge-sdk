package knowledgesdk

import "context"

const (
	KnowledgeGridAutomationTriggerRecordCreated = "record_created"
	KnowledgeGridAutomationTriggerRecordUpdated = "record_updated"
	KnowledgeGridAutomationActionSetField       = "set_field"
)

// KnowledgeGridAutomationExecute is an explicitly authorized, idempotent
// execution request. ExecutionID identifies one trigger delivery and must be
// reused when the caller retries that delivery.
type KnowledgeGridAutomationExecute struct {
	ExecutionID           string `json:"execution_id"`
	AutomationID          string `json:"automation_id"`
	Trigger               string `json:"trigger"`
	RecordID              string `json:"record_id"`
	ExpectedDraftRevision int64  `json:"expected_draft_revision"`
}

type KnowledgeGridAutomationResult struct {
	View         KnowledgePageView `json:"view"`
	AutomationID string            `json:"automation_id"`
	RecordID     string            `json:"record_id"`
	FieldID      string            `json:"field_id"`
	Applied      bool              `json:"applied"`
}

type PageGridAutomationService interface {
	ExecuteKnowledgeGridAutomation(context.Context, string, string, KnowledgeGridAutomationExecute, Authority) (KnowledgeGridAutomationResult, error)
}
