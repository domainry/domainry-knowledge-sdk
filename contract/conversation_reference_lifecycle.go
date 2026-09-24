package contract

import (
	"context"
	"encoding/json"

	agentsdk "github.com/domainry/domainry-agent-sdk"
)

// ConversationReferenceLifecycle is the public owner port used when another
// domain removes a conversation. Knowledge keeps the cleanup reference and an
// idempotent receipt in its own transaction.
type ConversationReferenceLifecycle interface {
	DeleteConversationReferencesForRequest(context.Context, string, string, agentsdk.ConversationAuthority) (json.RawMessage, error)
}
