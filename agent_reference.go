package knowledgesdk

import "context"

// KnowledgeAgentReference is an immutable citation address. ContentSHA256
// belongs to the exact published page version or document block, never to a
// mutable draft or a search excerpt.
type KnowledgeAgentReference struct {
	LibraryID     string `json:"library_id"`
	PageID        string `json:"page_id"`
	Version       int64  `json:"version"`
	BlockID       string `json:"block_id,omitempty"`
	ContentSHA256 string `json:"content_sha256"`
}

type KnowledgeAgentResolvedReference struct {
	Reference   KnowledgeAgentReference `json:"reference"`
	Title       string                  `json:"title"`
	ContentType string                  `json:"content_type"`
	Text        string                  `json:"text"`
}

type KnowledgeAgentReferenceService interface {
	ResolveKnowledgeAgentReference(context.Context, KnowledgeAgentReference, Authority) (KnowledgeAgentResolvedReference, error)
}

type KnowledgeAgentReferenceRepository interface {
	KnowledgeAgentReferenceService
}
