// Package contract defines structured document values and storage ports.
package contract

import (
	"context"
	"encoding/json"
	sdk "github.com/domainry/domainry-tools-sdk"
	"time"
)

type Record struct {
	ID        string          `json:"id"`
	Kind      string          `json:"kind"`
	Title     string          `json:"title"`
	Status    string          `json:"status"`
	Revision  int64           `json:"revision"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
type Write struct {
	ClientID         string          `json:"client_id"`
	ID               string          `json:"id,omitempty"`
	ExpectedRevision int64           `json:"expected_revision"`
	Title            string          `json:"title"`
	Status           string          `json:"status"`
	Data             json.RawMessage `json:"data"`
}
type Page struct {
	Items      []Record `json:"items"`
	NextCursor string   `json:"next_cursor,omitempty"`
	Complete   bool     `json:"complete"`
}
type Validator func(previous *Record, next *Record) error

// Repository owns scoped records, immutable revisions and command receipts.
// Callers supply product rules; implementations keep validation and commit atomic.
type Repository interface {
	Get(context.Context, string, string, int64, sdk.Authority) (Record, error)
	List(context.Context, string, string, string, int, sdk.Authority) (Page, error)
	Save(context.Context, string, Write, sdk.Authority, Validator) (Record, error)
	Receipt(context.Context, string, Write, sdk.Authority) (Record, bool, error)
}
