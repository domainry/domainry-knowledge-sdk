package knowledgesdk

import (
	"context"
	"time"
)

const (
	KnowledgePageContentPolicyPage         = "page"
	KnowledgePageContentPolicySubtree      = "subtree"
	KnowledgePageContentOperationCopy      = "copy"
	KnowledgePageContentOperationDownload  = "download"
	KnowledgePageContentOperationPrint     = "print"
	KnowledgePageContentOperationExport    = "export"
	KnowledgePageContentOperationDuplicate = "duplicate"
	KnowledgePageContentOperationShare     = "share"

	KnowledgePageClassificationPublic       = "public"
	KnowledgePageClassificationInternal     = "internal"
	KnowledgePageClassificationConfidential = "confidential"
	KnowledgePageClassificationRestricted   = "restricted"

	KnowledgePageWatermarkNone     = "none"
	KnowledgePageWatermarkIdentity = "identity"

	KnowledgePageSensitiveContentAllow         = "allow"
	KnowledgePageSensitiveContentBlockOutbound = "block_outbound"
)

type KnowledgePageDynamicWatermark struct {
	Text        string    `json:"text"`
	GeneratedAt time.Time `json:"generated_at"`
}

// KnowledgePageContentGovernance is derived for the current reader. The
// dynamic watermark is never persisted in page content or policy.
type KnowledgePageContentGovernance struct {
	Classification   string                         `json:"classification"`
	DynamicWatermark *KnowledgePageDynamicWatermark `json:"dynamic_watermark,omitempty"`
}

// KnowledgePageContentPolicy controls product-owned exfiltration entry points.
// A page policy always applies to its own page; subtree scope also applies to
// descendants until a nearer policy overrides it.
type KnowledgePageContentPolicy struct {
	PageID          string    `json:"page_id"`
	LibraryID       string    `json:"library_id"`
	Scope           string    `json:"scope"`
	AllowCopy       bool      `json:"allow_copy"`
	AllowDownload   bool      `json:"allow_download"`
	AllowPrint      bool      `json:"allow_print"`
	AllowExport     bool      `json:"allow_export"`
	AllowDuplicate  bool      `json:"allow_duplicate"`
	Classification  string    `json:"classification"`
	WatermarkMode   string    `json:"watermark_mode"`
	SensitiveAction string    `json:"sensitive_content_action"`
	RequireApproval bool      `json:"require_outbound_approval"`
	Revision        int64     `json:"revision"`
	UpdatedByUserID string    `json:"updated_by_user_id,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

type KnowledgePageContentPolicyWrite struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
	Scope            string `json:"scope"`
	AllowCopy        bool   `json:"allow_copy"`
	AllowDownload    bool   `json:"allow_download"`
	AllowPrint       bool   `json:"allow_print"`
	AllowExport      bool   `json:"allow_export"`
	AllowDuplicate   bool   `json:"allow_duplicate"`
	Classification   string `json:"classification"`
	WatermarkMode    string `json:"watermark_mode"`
	SensitiveAction  string `json:"sensitive_content_action"`
	RequireApproval  bool   `json:"require_outbound_approval"`
}

type KnowledgePageContentDecision struct {
	Operation         string `json:"operation"`
	Allowed           bool   `json:"allowed"`
	SourcePageID      string `json:"source_page_id,omitempty"`
	Scope             string `json:"scope"`
	Revision          int64  `json:"revision"`
	LibraryID         string `json:"library_id"`
	PageID            string `json:"page_id"`
	Classification    string `json:"classification"`
	WatermarkMode     string `json:"watermark_mode"`
	SensitiveAction   string `json:"sensitive_content_action"`
	ApprovalRequired  bool   `json:"approval_required"`
	ApprovalReference string `json:"approval_reference,omitempty"`
}

type KnowledgePageOutboundApprovalRequest struct {
	LibraryID      string `json:"library_id"`
	PageID         string `json:"page_id"`
	Operation      string `json:"operation"`
	Classification string `json:"classification"`
	PolicyPageID   string `json:"policy_page_id,omitempty"`
	PolicyRevision int64  `json:"policy_revision"`
}

type KnowledgePageOutboundApprovalDecision struct {
	Allowed   bool   `json:"allowed"`
	Reference string `json:"reference,omitempty"`
}

// KnowledgePageOutboundApprover is supplied by the host approval system. It
// resolves current approval truth for the current Identity authority;
// Knowledge never accepts a client assertion as approval evidence.
type KnowledgePageOutboundApprover interface {
	AuthorizeKnowledgePageOutbound(context.Context, KnowledgePageOutboundApprovalRequest, Authority) (KnowledgePageOutboundApprovalDecision, error)
}

type KnowledgePageContentPolicyService interface {
	KnowledgePageContentPolicy(context.Context, string, string, Authority) (KnowledgePageContentPolicy, error)
	SetKnowledgePageContentPolicy(context.Context, string, string, KnowledgePageContentPolicyWrite, Authority) (KnowledgePageContentPolicy, error)
	AuthorizeKnowledgePageContentOperation(context.Context, string, string, string, Authority) (KnowledgePageContentDecision, error)
	DownloadKnowledgePageDocument(context.Context, string, string, string, Authority) (KnowledgeDocumentDownload, error)
}

type KnowledgePageContentPolicyRepository interface {
	KnowledgePageContentPolicy(context.Context, string, string, Authority) (KnowledgePageContentPolicy, error)
	SetKnowledgePageContentPolicy(context.Context, string, string, KnowledgePageContentPolicyWrite, Authority) (KnowledgePageContentPolicy, error)
	EvaluateKnowledgePageContentPolicy(context.Context, string, string, string, Authority) (KnowledgePageContentDecision, error)
}
