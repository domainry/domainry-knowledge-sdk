package contract

import (
	"context"
	"io"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
)

// DocumentSecurityScanner runs while the immutable original is quarantined.
// A rejected result is terminal until an editor explicitly retries after the
// scanner policy changes; technical errors remain retryable worker failures.
type DocumentSecurityScanner interface {
	ScanKnowledgeDocument(context.Context, KnowledgeDocumentSecurityInput) (KnowledgeDocumentSecurityResult, error)
}

type KnowledgeDocumentSecurityInput struct {
	Document knowledgesdk.KnowledgeDocument
	Content  io.ReadSeeker
}

type KnowledgeDocumentSecurityResult struct {
	Clean        bool
	Scanner      string
	ReportSHA256 string
	ErrorCode    string
}
