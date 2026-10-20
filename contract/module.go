package contract

import (
	"context"
	"fmt"
	"strings"
	"time"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
)

// Knowledge is the retrieval capability owned by the Knowledge module.
type Knowledge interface {
	SearchKnowledge(context.Context, string, knowledgesdk.Authority) (knowledgesdk.Result, error)
	ReadKnowledge(context.Context, string, knowledgesdk.Authority) (knowledgesdk.Result, error)
	RevalidateKnowledge(context.Context, knowledgesdk.Result, knowledgesdk.Authority) error
}

type Options struct {
	DocumentStorage           knowledgesdk.KnowledgeDocumentStorage
	DocumentSecurity          DocumentSecurityScanner
	DocumentPoll              time.Duration
	DocumentWorkTimeout       time.Duration
	DocumentUploadIdleTimeout time.Duration
	LibraryKnowledge          []LibraryKnowledgeBinding
	KnowledgeDatasources      knowledgesdk.KnowledgeDatasourceCatalog
	LibraryAuthorizer         knowledgesdk.KnowledgeLibraryAuthorizer
	SourceAuthorizer          knowledgesdk.KnowledgeDocumentSourceAuthorizer
	Knowledge                 Knowledge
	PrincipalDirectory        knowledgesdk.KnowledgePrincipalDirectory
	AccessExplainer           knowledgesdk.KnowledgeAccessExplainer
	AuditAppender             knowledgesdk.KnowledgeAuditAppender
	CollaborationPublisher    knowledgesdk.KnowledgeCollaborationPublisher
}

type LibraryKnowledgeBinding struct {
	WorkspaceID     string
	LibraryID       string
	Source          knowledgesdk.Source
	ManageDocuments bool
}

type Runtime interface {
	Validate(*Options) error
	Prepare(string, Options) (Knowledge, error)
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
