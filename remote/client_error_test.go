package remote

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
)

func TestDecodeRemoteErrorRestoresPageConflictType(t *testing.T) {
	want := knowledgesdk.KnowledgePageConflict{
		CurrentDraftRevision: 11, BasePublishedVersion: 5, CurrentContentSHA256: "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc", CurrentUpdatedAt: time.Date(2026, 10, 19, 1, 2, 3, 0, time.UTC),
	}
	details, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	decoded := decodeRemoteError(&knowledgecontract.SaaSError{Class: "conflict", Code: "page_draft_revision_conflict", Details: details})
	var conflict *knowledgesdk.KnowledgePageConflictError
	if !errors.As(decoded, &conflict) || conflict.Conflict != want {
		t.Fatalf("decoded=%T %+v", decoded, decoded)
	}
}

func TestDecodeRemoteErrorRejectsMalformedPageConflictDetails(t *testing.T) {
	decoded := decodeRemoteError(&knowledgecontract.SaaSError{Class: "conflict", Code: "page_draft_revision_conflict", Details: json.RawMessage(`{"current_draft_revision":2}`)})
	var conflict *knowledgesdk.KnowledgePageConflictError
	if errors.As(decoded, &conflict) {
		t.Fatalf("malformed details must not become a typed page conflict: %+v", conflict)
	}
}
