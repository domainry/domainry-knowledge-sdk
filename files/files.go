// Package files defines the Agent-independent file capability exposed by
// Domainry Knowledge. Product modules own their resource authorization and
// keep only the returned opaque file ID.
package files

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/domainry/domainry-foundation/apperror"
)

const MaxContentBytes int64 = 25 << 20

type Authority struct {
	RuntimeID   string `json:"runtime_id"`
	WorkspaceID string `json:"workspace_id"`
	SubjectID   string `json:"subject_id"`
}

type Upload struct {
	ClientID    string `json:"client_id"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data"`
}

// StreamUpload carries immutable file metadata separately from its bytes so
// callers and the standalone Knowledge service do not buffer an entire file.
type StreamUpload struct {
	ClientID    string    `json:"client_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	SHA256      string    `json:"sha256"`
	Content     io.Reader `json:"-"`
}

type File struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	CreatedBy   string    `json:"created_by"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	SHA256      string    `json:"sha256"`
	CreatedAt   time.Time `json:"created_at"`
}

type Download struct {
	File File   `json:"file"`
	Data []byte `json:"data"`
}

// Binding records which durable product resource owns a file. It does not
// delegate that product's read permission to Knowledge.
type Binding struct {
	Owner        string `json:"owner"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	FieldKey     string `json:"field_key,omitempty"`
}

type Service interface {
	Upload(context.Context, Authority, Upload) (File, error)
	UploadStream(context.Context, Authority, StreamUpload) (File, error)
	Get(context.Context, Authority, string) (File, error)
	Download(context.Context, Authority, string) (Download, error)
	Bind(context.Context, Authority, string, Binding) error
	Delete(context.Context, Authority, string) error
}

type ErrorClass = apperror.ErrorKind

const (
	ErrorClassBadRequest  = apperror.KindBadRequest
	ErrorClassForbidden   = apperror.KindForbidden
	ErrorClassNotFound    = apperror.KindNotFound
	ErrorClassConflict    = apperror.KindConflict
	ErrorClassRateLimited = apperror.KindRateLimited
	ErrorClassUnavailable = apperror.KindUnavailable
	ErrorClassInternal    = apperror.KindInternal
)

type Error struct {
	Class     ErrorClass
	Code      string
	Message   string
	Retryable bool
	Cause     error
}

func (value *Error) Error() string {
	if value == nil {
		return ""
	}
	message := strings.TrimSpace(value.Message)
	if message == "" {
		message = strings.TrimSpace(value.Code)
	}
	if message == "" {
		message = "knowledge file operation failed"
	}
	if value.Cause != nil {
		return fmt.Sprintf("%s: %v", message, value.Cause)
	}
	return message
}

func (value *Error) Unwrap() error {
	if value == nil {
		return nil
	}
	return value.Cause
}

func (value *Error) ApplicationErrorKind() string {
	if value == nil {
		return ""
	}
	if value.Class == "" && value.Retryable {
		return string(apperror.KindUnavailable)
	}
	return string(apperror.NormalizeKind(string(value.Class)))
}

func (value *Error) ErrorCode() string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(value.Code)
}

func (*Error) ErrorParams() map[string]string { return nil }
