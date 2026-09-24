package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	agentsdk "github.com/domainry/domainry-agent-sdk"
)

type ConversationKnowledge interface {
	Search(context.Context, string, agentsdk.ConversationAuthority) (json.RawMessage, error)
}
type SourcePolicy interface {
	CheckSources(context.Context, agentsdk.ConversationAuthority, string, *agentsdk.ConversationSources) ([]agentsdk.ConversationRunReference, error)
	CheckRun(context.Context, agentsdk.ConversationAuthority, string, agentsdk.ConversationRunReference) ([]agentsdk.ConversationRunReference, error)
}

// ContextReader is the narrow Agent-owned read port required by Knowledge to
// validate a referenced conversation or completed run. Remote clients answer
// server challenges through this port; no repository or Store crosses the
// network boundary.
type ContextReader interface {
	Conversation(context.Context, string, agentsdk.ConversationAuthority) (agentsdk.Conversation, error)
	Run(context.Context, string, string, agentsdk.ConversationAuthority) (agentsdk.ConversationRun, error)
}
type Options struct {
	DocumentStorage      agentsdk.KnowledgeDocumentStorage
	DocumentPoll         time.Duration
	LibraryKnowledge     []LibraryKnowledgeBinding
	KnowledgeDatasources agentsdk.KnowledgeDatasourceCatalog
	LibraryAuthorizer    agentsdk.KnowledgeLibraryAuthorizer
	AttachmentAuthorizer agentsdk.ConversationAttachmentAuthorizer
	AttachmentKnowledge  []agentsdk.ConversationAttachmentKnowledgeBinding
	ArtifactStorage      agentsdk.ConversationArtifactStorage
	ArtifactExportTTL    time.Duration
	PersonalAuthorizer   agentsdk.ConversationToolAuthorizer
	Knowledge            ConversationKnowledge
	Sources              SourcePolicy
	Context              ContextReader
}
type LibraryKnowledgeBinding struct {
	WorkspaceID     string
	LibraryID       string
	Source          agentsdk.ConversationKnowledgeSource
	ManageDocuments bool // Trusted host opt-in; never a client supplied grant.
}

// Runtime is a process-local Knowledge capability whose private repository was
// already opened by the Knowledge module. Agent never supplies a Store here.
type Runtime interface {
	Validate(*Options) error
	Prepare(string, Options) (ConversationKnowledge, error)
	Activate(string, Options) error
	NewService(string, Options) Service
}

type ApplicationRef struct{ RuntimeID string }

func (ref ApplicationRef) Validate() error {
	if strings.TrimSpace(ref.RuntimeID) == "" || len(strings.TrimSpace(ref.RuntimeID)) > 255 {
		return fmt.Errorf("Knowledge Runtime identity is required")
	}
	return nil
}
