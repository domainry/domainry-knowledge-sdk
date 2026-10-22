package knowledgesdk

import (
	"context"
	"encoding/json"
	"time"
)

// KnowledgeFormSubmission is an immutable answer set against one immutable
// published form version. Draft form content is never accepted as a submission
// schema.
type KnowledgeFormSubmission struct {
	ID                string                     `json:"id"`
	LibraryID         string                     `json:"library_id"`
	PageID            string                     `json:"page_id"`
	PublishedVersion  int64                      `json:"published_version"`
	ContentSHA256     string                     `json:"content_sha256"`
	Answers           map[string]json.RawMessage `json:"answers"`
	SubmittedByUserID string                     `json:"submitted_by_user_id"`
	SubmittedAt       time.Time                  `json:"submitted_at"`
}

type KnowledgeFormSubmit struct {
	ClientID string                     `json:"client_id"`
	Answers  map[string]json.RawMessage `json:"answers"`
}

type KnowledgeFormSubmissionPage struct {
	Items      []KnowledgeFormSubmission `json:"items"`
	NextBefore string                    `json:"next_before,omitempty"`
	Complete   bool                      `json:"complete"`
}

type KnowledgeFormChoiceCount struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type KnowledgeFormQuestionAnalytics struct {
	Key           string                     `json:"key"`
	Label         string                     `json:"label"`
	Type          string                     `json:"type"`
	Answered      int64                      `json:"answered"`
	Skipped       int64                      `json:"skipped"`
	ChoiceCounts  []KnowledgeFormChoiceCount `json:"choice_counts"`
	NumberMin     *float64                   `json:"number_min,omitempty"`
	NumberMax     *float64                   `json:"number_max,omitempty"`
	NumberAverage *float64                   `json:"number_average,omitempty"`
}

type KnowledgeFormAnalytics struct {
	PublishedVersion int64                            `json:"published_version"`
	TotalSubmissions int64                            `json:"total_submissions"`
	Questions        []KnowledgeFormQuestionAnalytics `json:"questions"`
}

type KnowledgeFormExport struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type PageFormService interface {
	SubmitKnowledgePageForm(context.Context, string, string, KnowledgeFormSubmit, Authority) (KnowledgeFormSubmission, error)
	KnowledgePageFormSubmissions(context.Context, string, string, string, int, Authority) (KnowledgeFormSubmissionPage, error)
	KnowledgePageFormAnalytics(context.Context, string, string, Authority) (KnowledgeFormAnalytics, error)
	ExportKnowledgePageFormSubmissions(context.Context, string, string, Authority) (KnowledgeFormExport, error)
}
