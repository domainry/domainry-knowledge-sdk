package knowledgesdk

import (
	"context"
	"encoding/json"
	"time"
)

const (
	KnowledgeNativeOperationUpsert = "upsert"
	KnowledgeNativeOperationDelete = "delete"

	KnowledgeNativeTargetField             = "field"
	KnowledgeNativeTargetRecord            = "record"
	KnowledgeNativeTargetView              = "view"
	KnowledgeNativeTargetCell              = "cell"
	KnowledgeNativeTargetTableStructure    = "table_structure"
	KnowledgeNativeTargetSheet             = "sheet"
	KnowledgeNativeTargetWorkbookStructure = "workbook_structure"
	KnowledgeNativeTargetSheetStructure    = "sheet_structure"
	KnowledgeNativeTargetFormula           = "formula"
	KnowledgeNativeTargetCellStyle         = "cell_style"
	KnowledgeNativeTargetStyle             = "style"
	KnowledgeNativeTargetMerge             = "merge"
	KnowledgeNativeTargetChart             = "chart"
	KnowledgeNativeTargetPivot             = "pivot"
	KnowledgeNativeTargetProtectedRange    = "protected_range"
	KnowledgeNativeTargetNode              = "node"
	KnowledgeNativeTargetEdge              = "edge"
	KnowledgeNativeTargetGraphStructure    = "graph_structure"
	KnowledgeNativeTargetSlide             = "slide"
	KnowledgeNativeTargetElement           = "element"
	KnowledgeNativeTargetAsset             = "asset"
	KnowledgeNativeTargetPresentation      = "presentation_structure"
	KnowledgeNativeTargetComponent         = "component"
	KnowledgeNativeTargetFormStructure     = "form_structure"
)

// KnowledgePageNativeOperation is a semantic mutation against one stable
// editor entity. ExpectedHash is the SHA-256 of the current canonical target;
// an empty hash means the target must not exist. This makes stale batches
// merge when they touch different entities without accepting same-entity lost
// updates.
type KnowledgePageNativeOperation struct {
	ID           string          `json:"id"`
	Target       string          `json:"target"`
	TargetID     string          `json:"target_id"`
	ParentID     string          `json:"parent_id,omitempty"`
	Action       string          `json:"action"`
	ExpectedHash string          `json:"expected_hash"`
	Value        json.RawMessage `json:"value,omitempty"`
}

type KnowledgePageNativeOperationBatch struct {
	ClientID          string                         `json:"client_id"`
	BaseDraftRevision int64                          `json:"base_draft_revision"`
	Operations        []KnowledgePageNativeOperation `json:"operations"`
}

type KnowledgePageNativeOperationRecord struct {
	Operation     KnowledgePageNativeOperation `json:"operation"`
	LibraryID     string                       `json:"library_id"`
	PageID        string                       `json:"page_id"`
	ContentType   string                       `json:"content_type"`
	DraftRevision int64                        `json:"draft_revision"`
	UserID        string                       `json:"user_id"`
	CreatedAt     time.Time                    `json:"created_at"`
}

type KnowledgePageNativeOperationResult struct {
	View       KnowledgePageView                    `json:"view"`
	Operations []KnowledgePageNativeOperationRecord `json:"operations"`
}

type PageNativeOperationService interface {
	ApplyKnowledgePageNativeOperations(context.Context, string, string, KnowledgePageNativeOperationBatch, Authority) (KnowledgePageNativeOperationResult, error)
}
