package knowledgesdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// PageContentSchemaVersion is the canonical page payload contract. Version 2
// stores text documents as a structured Block tree instead of Markdown. The
// product has not shipped a version-1 page corpus, so callers must write the
// current contract directly; there is intentionally no dual-format fallback.
const PageContentSchemaVersion = 2

const (
	KnowledgePageTypeDocument              = "document"
	KnowledgePageTypeMultidimensionalTable = "multidimensional_table"
	KnowledgePageTypeSpreadsheet           = "spreadsheet"
	KnowledgePageTypePresentation          = "presentation"
	KnowledgePageTypeForm                  = "form"
	KnowledgePageTypeMindNote              = "mind_note"
	KnowledgePageTypeWhiteboard            = "whiteboard"
	KnowledgePageTypeMindMap               = "mind_map"
	KnowledgePageTypeFlowchart             = "flowchart"
	KnowledgePageTypeFile                  = "file"
)

// KnowledgePageTypeCapability is a stable, additive capability bitset shared
// by SDK clients. It describes product behavior; it never grants authority.
// Every operation must still pass the server-owned action and resource checks.
type KnowledgePageTypeCapability uint64

const (
	KnowledgePageCapabilityCreate KnowledgePageTypeCapability = 1 << iota
	KnowledgePageCapabilityEdit
	KnowledgePageCapabilityRead
	KnowledgePageCapabilityPreview
	KnowledgePageCapabilityImport
	KnowledgePageCapabilityRealtime
	KnowledgePageCapabilitySearch
	KnowledgePageCapabilityComment
	KnowledgePageCapabilityVersion
)

type KnowledgePageContentRootKind string

const (
	KnowledgePageContentRootObject KnowledgePageContentRootKind = "object"
	KnowledgePageContentRootArray  KnowledgePageContentRootKind = "array"
	KnowledgePageContentRootScalar KnowledgePageContentRootKind = "scalar"
)

// KnowledgePageTypeContract is the canonical cross-process description of a
// trusted page type. RootField/RootKind describe the stable content envelope;
// deeper invariants remain owned by the server validator registered for ID.
type KnowledgePageTypeContract struct {
	ID           string                       `json:"id"`
	RootField    string                       `json:"root_field"`
	RootKind     KnowledgePageContentRootKind `json:"root_kind"`
	Capabilities KnowledgePageTypeCapability  `json:"capabilities"`
}

var knowledgePageTypeContracts = []KnowledgePageTypeContract{
	{ID: KnowledgePageTypeDocument, RootField: "blocks", RootKind: KnowledgePageContentRootArray, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilityRealtime | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypeMultidimensionalTable, RootField: "fields", RootKind: KnowledgePageContentRootArray, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypeSpreadsheet, RootField: "workbook", RootKind: KnowledgePageContentRootObject, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypePresentation, RootField: "slides", RootKind: KnowledgePageContentRootArray, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypeForm, RootField: "components", RootKind: KnowledgePageContentRootArray, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypeMindNote, RootField: "nodes", RootKind: KnowledgePageContentRootArray, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypeWhiteboard, RootField: "elements", RootKind: KnowledgePageContentRootArray, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypeMindMap, RootField: "root", RootKind: KnowledgePageContentRootObject, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypeFlowchart, RootField: "nodes", RootKind: KnowledgePageContentRootArray, Capabilities: KnowledgePageCapabilityCreate | KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
	{ID: KnowledgePageTypeFile, RootField: "document_id", RootKind: KnowledgePageContentRootScalar, Capabilities: KnowledgePageCapabilityEdit | KnowledgePageCapabilityRead | KnowledgePageCapabilityPreview | KnowledgePageCapabilityImport | KnowledgePageCapabilitySearch | KnowledgePageCapabilityComment | KnowledgePageCapabilityVersion},
}

func KnowledgePageTypeContracts() []KnowledgePageTypeContract {
	return append([]KnowledgePageTypeContract(nil), knowledgePageTypeContracts...)
}

func KnowledgePageTypeContractFor(id string) (KnowledgePageTypeContract, bool) {
	for _, contract := range knowledgePageTypeContracts {
		if contract.ID == id {
			return contract, true
		}
	}
	return KnowledgePageTypeContract{}, false
}

func ValidateKnowledgePageTypeContracts() error {
	seen := make(map[string]bool, len(knowledgePageTypeContracts))
	for _, contract := range knowledgePageTypeContracts {
		if strings.TrimSpace(contract.ID) == "" || seen[contract.ID] || strings.TrimSpace(contract.RootField) == "" || contract.Capabilities == 0 {
			return fmt.Errorf("invalid knowledge page type contract %q", contract.ID)
		}
		if contract.RootKind != KnowledgePageContentRootObject && contract.RootKind != KnowledgePageContentRootArray && contract.RootKind != KnowledgePageContentRootScalar {
			return fmt.Errorf("invalid root kind for knowledge page type %q", contract.ID)
		}
		seen[contract.ID] = true
	}
	return nil
}

func init() {
	if err := ValidateKnowledgePageTypeContracts(); err != nil {
		panic(err)
	}
}

// KnowledgeFilePageContent is the canonical content for a file page. The
// immutable original stays in KnowledgeDocument; pages persist only an
// authorized stable reference and presentation metadata.
type KnowledgeFilePageContent struct {
	DocumentID  string `json:"document_id"`
	DisplayMode string `json:"display_mode"` // preview, card
	Title       string `json:"title"`
}

// The following envelopes are the SDK-owned canonical roots for native page
// types. RawMessage is limited to renderer state that the owning Knowledge
// service validates before persistence; product features use explicit types.
type KnowledgeMultidimensionalTableContent struct {
	Fields  []json.RawMessage `json:"fields"`
	Records []json.RawMessage `json:"records"`
	Views   []json.RawMessage `json:"views"`
}

type KnowledgeSpreadsheetContent struct {
	SchemaVersion   int                                  `json:"schema_version"`
	Workbook        KnowledgeSpreadsheetWorkbook         `json:"workbook"`
	Charts          []KnowledgeSpreadsheetChart          `json:"charts"`
	Pivots          []KnowledgeSpreadsheetPivot          `json:"pivots"`
	ProtectedRanges []KnowledgeSpreadsheetProtectedRange `json:"protected_ranges"`
}

type KnowledgeSpreadsheetWorkbook struct {
	ID                  string                     `json:"id"`
	Revision            int64                      `json:"rev,omitempty"`
	Name                string                     `json:"name"`
	AppVersion          string                     `json:"appVersion"`
	Locale              string                     `json:"locale"`
	DateSystem          string                     `json:"dateSystem,omitempty"`
	CalculationSettings json.RawMessage            `json:"calculationSettings,omitempty"`
	Styles              map[string]json.RawMessage `json:"styles"`
	SheetOrder          []string                   `json:"sheetOrder"`
	Sheets              map[string]json.RawMessage `json:"sheets"`
	DefaultStyle        json.RawMessage            `json:"defaultStyle,omitempty"`
	Resources           []json.RawMessage          `json:"resources,omitempty"`
}

type KnowledgeSpreadsheetRange struct {
	StartRow    int `json:"start_row"`
	EndRow      int `json:"end_row"`
	StartColumn int `json:"start_column"`
	EndColumn   int `json:"end_column"`
}

type KnowledgeSpreadsheetChart struct {
	ID      string                    `json:"id"`
	Name    string                    `json:"name"`
	SheetID string                    `json:"sheet_id"`
	Type    string                    `json:"type"`
	Range   KnowledgeSpreadsheetRange `json:"range"`
}

type KnowledgeSpreadsheetPivot struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	SheetID     string                    `json:"sheet_id"`
	SourceRange KnowledgeSpreadsheetRange `json:"source_range"`
	Aggregate   string                    `json:"aggregate"`
}

type KnowledgeSpreadsheetProtectedRange struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	SheetID     string                    `json:"sheet_id"`
	Range       KnowledgeSpreadsheetRange `json:"range"`
	MinimumRole string                    `json:"minimum_role"`
}

type KnowledgePresentationContent struct {
	Slides []json.RawMessage `json:"slides"`
}

type KnowledgeFormContent struct {
	SchemaVersion int                      `json:"schema_version"`
	Display       string                   `json:"display"`
	Components    []KnowledgeFormComponent `json:"components"`
	Settings      KnowledgeFormSettings    `json:"settings"`
}

type KnowledgeFormComponent struct {
	Type        string                    `json:"type"`
	Key         string                    `json:"key"`
	Label       string                    `json:"label"`
	Description string                    `json:"description,omitempty"`
	Placeholder string                    `json:"placeholder,omitempty"`
	Input       bool                      `json:"input"`
	Validate    KnowledgeFormValidation   `json:"validate"`
	Values      []KnowledgeFormOption     `json:"values,omitempty"`
	Data        *KnowledgeFormSelectData  `json:"data,omitempty"`
	Multiple    bool                      `json:"multiple,omitempty"`
	Conditional *KnowledgeFormConditional `json:"conditional,omitempty"`
}

type KnowledgeFormValidation struct {
	Required  bool     `json:"required"`
	MinLength int      `json:"minLength,omitempty"`
	MaxLength int      `json:"maxLength,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
}

type KnowledgeFormOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type KnowledgeFormSelectData struct {
	Values []KnowledgeFormOption `json:"values"`
}

type KnowledgeFormConditional struct {
	Show bool            `json:"show"`
	When string          `json:"when"`
	Eq   json.RawMessage `json:"eq"`
}

type KnowledgeFormSettings struct {
	SubmissionRole           string `json:"submission_role"`
	AllowMultipleSubmissions bool   `json:"allow_multiple_submissions"`
	OpensAt                  string `json:"opens_at"`
	ClosesAt                 string `json:"closes_at"`
	ConfirmationMessage      string `json:"confirmation_message"`
}

type KnowledgeMindNoteContent struct {
	Nodes []json.RawMessage `json:"nodes"`
}

type KnowledgeWhiteboardPoint [2]float64

type KnowledgeWhiteboardBinding struct {
	ElementID  string                    `json:"element_id"`
	Focus      float64                   `json:"focus"`
	Gap        float64                   `json:"gap"`
	FixedPoint *KnowledgeWhiteboardPoint `json:"fixed_point,omitempty"`
}

type KnowledgeWhiteboardBoundElement struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type KnowledgeWhiteboardCrop struct {
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Width         float64 `json:"width"`
	Height        float64 `json:"height"`
	NaturalWidth  float64 `json:"natural_width"`
	NaturalHeight float64 `json:"natural_height"`
}

// KnowledgeWhiteboardElement is the Knowledge-owned durable drawing model.
// Excalidraw-specific revision, seed, cache, selection, and deletion fields are
// deliberately excluded; they are reconstructed by the client adapter.
type KnowledgeWhiteboardElement struct {
	ID               string                            `json:"id"`
	Type             string                            `json:"type"`
	X                float64                           `json:"x"`
	Y                float64                           `json:"y"`
	Width            float64                           `json:"width"`
	Height           float64                           `json:"height"`
	Angle            float64                           `json:"angle"`
	StrokeColor      string                            `json:"stroke_color"`
	BackgroundColor  string                            `json:"background_color"`
	FillStyle        string                            `json:"fill_style"`
	StrokeWidth      float64                           `json:"stroke_width"`
	StrokeStyle      string                            `json:"stroke_style"`
	Roughness        float64                           `json:"roughness"`
	Opacity          float64                           `json:"opacity"`
	GroupIDs         []string                          `json:"group_ids,omitempty"`
	FrameID          string                            `json:"frame_id,omitempty"`
	Roundness        *float64                          `json:"roundness,omitempty"`
	BoundElements    []KnowledgeWhiteboardBoundElement `json:"bound_elements,omitempty"`
	Link             string                            `json:"link,omitempty"`
	Locked           bool                              `json:"locked,omitempty"`
	CustomKind       string                            `json:"custom_kind,omitempty"`
	Text             string                            `json:"text,omitempty"`
	OriginalText     string                            `json:"original_text,omitempty"`
	FontSize         float64                           `json:"font_size,omitempty"`
	FontFamily       float64                           `json:"font_family,omitempty"`
	TextAlign        string                            `json:"text_align,omitempty"`
	VerticalAlign    string                            `json:"vertical_align,omitempty"`
	ContainerID      string                            `json:"container_id,omitempty"`
	LineHeight       float64                           `json:"line_height,omitempty"`
	AutoResize       *bool                             `json:"auto_resize,omitempty"`
	Points           []KnowledgeWhiteboardPoint        `json:"points,omitempty"`
	Pressures        []float64                         `json:"pressures,omitempty"`
	SimulatePressure *bool                             `json:"simulate_pressure,omitempty"`
	StartBinding     *KnowledgeWhiteboardBinding       `json:"start_binding,omitempty"`
	EndBinding       *KnowledgeWhiteboardBinding       `json:"end_binding,omitempty"`
	StartArrowhead   string                            `json:"start_arrowhead,omitempty"`
	EndArrowhead     string                            `json:"end_arrowhead,omitempty"`
	Elbowed          bool                              `json:"elbowed,omitempty"`
	FileID           string                            `json:"file_id,omitempty"`
	Scale            *KnowledgeWhiteboardPoint         `json:"scale,omitempty"`
	Crop             *KnowledgeWhiteboardCrop          `json:"crop,omitempty"`
	Name             string                            `json:"name,omitempty"`
}

type KnowledgeWhiteboardAsset struct {
	ID        string `json:"id"`
	MIMEType  string `json:"mime_type"`
	DataURL   string `json:"data_url"`
	CreatedAt int64  `json:"created_at"`
}

type KnowledgeWhiteboardSettings struct {
	BackgroundColor string  `json:"background_color"`
	GridSize        float64 `json:"grid_size"`
	GridEnabled     bool    `json:"grid_enabled"`
}

type KnowledgeWhiteboardContent struct {
	SchemaVersion int                                 `json:"schema_version"`
	Elements      []KnowledgeWhiteboardElement        `json:"elements"`
	Assets        map[string]KnowledgeWhiteboardAsset `json:"assets"`
	Settings      KnowledgeWhiteboardSettings         `json:"settings"`
}

type KnowledgeMindMapContent struct {
	Root json.RawMessage `json:"root"`
}

type KnowledgeFlowchartContent struct {
	Nodes []json.RawMessage `json:"nodes"`
	Edges []json.RawMessage `json:"edges"`
}

type PageContent struct {
	ContentSchemaVersion int    `json:"content_schema_version"`
	Content              string `json:"content"`
	ContentSHA256        string `json:"content_sha256"`
	PlainText            string `json:"plain_text"`
	Summary              string `json:"summary,omitempty"`
}

// KnowledgeDocumentContent is the canonical payload for document pages. Block
// shape validation remains server-owned; RawMessage keeps the SDK independent
// from a particular React editor implementation while preserving every native
// and custom Block property losslessly.
type KnowledgeDocumentContent struct {
	Blocks []json.RawMessage `json:"blocks"`
}

type KnowledgePage struct {
	ID               string    `json:"id"`
	LibraryID        string    `json:"library_id"`
	ParentPageID     string    `json:"parent_page_id,omitempty"`
	HasChildren      bool      `json:"has_children"`
	Rank             string    `json:"rank"`
	Title            string    `json:"title"`
	ContentType      string    `json:"content_type"`
	State            string    `json:"state"` // active, archived
	PublishedVersion int64     `json:"published_version,omitempty"`
	ProjectionStatus string    `json:"projection_status,omitempty"` // none, queued, indexing, ready, failed
	AccessMode       string    `json:"access_mode"`                 // inherit, restricted
	AccessRevision   int64     `json:"access_revision"`
	AccessStatus     string    `json:"access_status"` // ready, queued, applying, failed
	Revision         int64     `json:"revision"`
	CreatedByUserID  string    `json:"created_by_user_id"`
	UpdatedByUserID  string    `json:"updated_by_user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type KnowledgePageDraft struct {
	PageID               string      `json:"page_id"`
	LibraryID            string      `json:"library_id"`
	BasePublishedVersion int64       `json:"base_published_version"`
	DraftRevision        int64       `json:"draft_revision"`
	Title                string      `json:"title"`
	ContentType          string      `json:"content_type"`
	ParentPageID         string      `json:"parent_page_id,omitempty"`
	Rank                 string      `json:"rank"`
	Body                 PageContent `json:"body"`
	UpdatedByUserID      string      `json:"updated_by_user_id"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

type KnowledgePageVersion struct {
	PageID            string      `json:"page_id"`
	LibraryID         string      `json:"library_id"`
	Version           int64       `json:"version"`
	Title             string      `json:"title"`
	ContentType       string      `json:"content_type"`
	ParentPageID      string      `json:"parent_page_id,omitempty"`
	Rank              string      `json:"rank"`
	Body              PageContent `json:"body"`
	PublishedByUserID string      `json:"published_by_user_id"`
	PublishedAt       time.Time   `json:"published_at"`
	PublishNote       string      `json:"publish_note,omitempty"`
}

type KnowledgePageView struct {
	Page        KnowledgePage                  `json:"page"`
	Published   *KnowledgePageVersion          `json:"published,omitempty"`
	Draft       *KnowledgePageDraft            `json:"draft,omitempty"`
	Breadcrumbs []KnowledgePageBreadcrumb      `json:"breadcrumbs"`
	GridAccess  *KnowledgeGridAccess           `json:"grid_access,omitempty"`
	Governance  KnowledgePageContentGovernance `json:"governance"`
}

// KnowledgeGridAccess is derived by the server for the redacted table payload
// in this response. Clients may use it for UI state but never as authority.
type KnowledgeGridAccess struct {
	Read        bool            `json:"read"`
	Write       bool            `json:"write"`
	Manage      bool            `json:"manage"`
	ViewWrite   map[string]bool `json:"view_write,omitempty"`
	FieldWrite  map[string]bool `json:"field_write,omitempty"`
	RecordWrite map[string]bool `json:"record_write,omitempty"`
}

type KnowledgePageBreadcrumb struct {
	PageID string `json:"page_id"`
	Title  string `json:"title"`
}

type KnowledgePageCreate struct {
	ClientID     string `json:"client_id"`
	ParentPageID string `json:"parent_page_id,omitempty"`
	Rank         string `json:"rank,omitempty"`
	Title        string `json:"title"`
	ContentType  string `json:"content_type"`
}

// KnowledgePageReferenceCreate creates one child page and replaces the source
// Block in its parent document with a governed reference in the same database
// transaction. ReferenceBlockID is the stable BlockNote Block ID at the slash
// insertion point; clients never manufacture a target page ID.
type KnowledgePageReferenceCreate struct {
	ClientID                    string `json:"client_id"`
	ExpectedParentDraftRevision int64  `json:"expected_parent_draft_revision"`
	ExpectedParentPageRevision  int64  `json:"expected_parent_page_revision"`
	ParentContentSchemaVersion  int    `json:"parent_content_schema_version"`
	ParentContent               string `json:"parent_content"`
	ParentContentSHA256         string `json:"parent_content_sha256"`
	ReferenceBlockID            string `json:"reference_block_id"`
	Title                       string `json:"title"`
	ContentType                 string `json:"content_type"`
	Display                     string `json:"display"`
}

type KnowledgePageReferenceCreateResult struct {
	Parent KnowledgePageView `json:"parent"`
	Target KnowledgePageView `json:"target"`
}

type KnowledgePageReferenceDirection string

const (
	KnowledgePageReferenceOutgoing KnowledgePageReferenceDirection = "outgoing"
	KnowledgePageReferenceIncoming KnowledgePageReferenceDirection = "incoming"
)

// KnowledgePageReferenceLink is a live, permission-filtered projection of one
// stored reference edge. Titles are resolved from the current authorized page
// projection and are never copied from the source block snapshot.
type KnowledgePageReferenceLink struct {
	SourcePageID  string    `json:"source_page_id"`
	SourceTitle   string    `json:"source_title,omitempty"`
	SourceType    string    `json:"source_content_type,omitempty"`
	TargetPageID  string    `json:"target_page_id"`
	TargetTitle   string    `json:"target_title,omitempty"`
	TargetType    string    `json:"target_content_type,omitempty"`
	BlockID       string    `json:"block_id"`
	ReferenceMode string    `json:"reference_mode"`
	TargetVersion int64     `json:"target_version,omitempty"`
	Display       string    `json:"display"`
	Available     bool      `json:"available"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type KnowledgePageReferencePage struct {
	Direction KnowledgePageReferenceDirection `json:"direction"`
	Items     []KnowledgePageReferenceLink    `json:"items"`
	NextAfter string                          `json:"next_after,omitempty"`
	Complete  bool                            `json:"complete"`
}

type KnowledgePageDraftWrite struct {
	ClientID              string `json:"client_id"`
	ExpectedDraftRevision int64  `json:"expected_draft_revision"`
	ExpectedPageRevision  int64  `json:"expected_page_revision"`
	Title                 string `json:"title"`
	ParentPageID          string `json:"parent_page_id,omitempty"`
	Rank                  string `json:"rank"`
	ContentSchemaVersion  int    `json:"content_schema_version"`
	Content               string `json:"content"`
	ContentSHA256         string `json:"content_sha256"`
}

type KnowledgePageMove struct {
	ClientID          string `json:"client_id"`
	ExpectedRevision  int64  `json:"expected_revision"`
	ParentPageID      string `json:"parent_page_id,omitempty"`
	Rank              string `json:"rank"`
	AccessImpactToken string `json:"access_impact_token,omitempty"`
}

type KnowledgePageOrder struct {
	PageID           string `json:"page_id"`
	ExpectedRevision int64  `json:"expected_revision"`
}

// KnowledgePageReorder is an all-siblings compare-and-set. Items are in the
// desired order and must contain every active child of ParentPageID exactly
// once, preventing a paginated or stale client from silently dropping pages.
type KnowledgePageReorder struct {
	ClientID     string               `json:"client_id"`
	ParentPageID string               `json:"parent_page_id,omitempty"`
	Items        []KnowledgePageOrder `json:"items"`
}

type KnowledgePageCopy struct {
	ClientID     string `json:"client_id"`
	ParentPageID string `json:"parent_page_id,omitempty"`
	Rank         string `json:"rank,omitempty"`
	Title        string `json:"title,omitempty"`
}

type KnowledgePageStateWrite struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
}

type KnowledgePagePublish struct {
	ClientID              string `json:"client_id"`
	ExpectedPageRevision  int64  `json:"expected_page_revision"`
	ExpectedDraftRevision int64  `json:"expected_draft_revision"`
	PublishNote           string `json:"publish_note,omitempty"`
}

type KnowledgePageRestoreVersion struct {
	ClientID              string `json:"client_id"`
	ExpectedPageRevision  int64  `json:"expected_page_revision"`
	ExpectedDraftRevision int64  `json:"expected_draft_revision"`
}

type KnowledgePagePage struct {
	Items     []KnowledgePage `json:"items"`
	NextAfter string          `json:"next_after,omitempty"`
	Complete  bool            `json:"complete"`
}

type KnowledgePageVersionPage struct {
	Items      []KnowledgePageVersion `json:"items"`
	NextBefore int64                  `json:"next_before,omitempty"`
	Complete   bool                   `json:"complete"`
}

type KnowledgePageVersionDiffInput struct {
	FromVersion int64 `json:"from_version"`
	ToVersion   int64 `json:"to_version"`
}

type KnowledgePageVersionDiff struct {
	FromVersion int64             `json:"from_version"`
	ToVersion   int64             `json:"to_version"`
	Title       []KnowledgeDiffOp `json:"title"`
	Parent      []KnowledgeDiffOp `json:"parent"`
	Content     []KnowledgeDiffOp `json:"content"`
}

type KnowledgeDiffOp struct {
	Kind string `json:"kind"` // equal, insert, delete
	Text string `json:"text"`
}

type KnowledgePageConflict struct {
	CurrentDraftRevision int64     `json:"current_draft_revision"`
	BasePublishedVersion int64     `json:"base_published_version"`
	CurrentContentSHA256 string    `json:"current_content_sha256"`
	CurrentUpdatedAt     time.Time `json:"current_updated_at"`
}

// KnowledgePageConflictError carries the bounded, non-content snapshot a
// caller needs to recover from an optimistic-lock conflict. Err remains the
// canonical classified error so errors.As continues to work for existing
// authorization and HTTP status handling.
type KnowledgePageConflictError struct {
	Err      *Error
	Conflict KnowledgePageConflict
}

func (e *KnowledgePageConflictError) Error() string {
	if e == nil || e.Err == nil {
		return "knowledge page conflict"
	}
	return e.Err.Error()
}

func (e *KnowledgePageConflictError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *KnowledgePageConflictError) KnowledgePageConflict() KnowledgePageConflict {
	if e == nil {
		return KnowledgePageConflict{}
	}
	return e.Conflict
}

func NewKnowledgePageConflictError(code string, conflict KnowledgePageConflict) error {
	code = strings.TrimSpace(code)
	if code == "" {
		code = "page_draft_revision_conflict"
	}
	return &KnowledgePageConflictError{
		Err:      &Error{Class: "conflict", Code: code},
		Conflict: conflict,
	}
}

// KnowledgePageCoreService contains the page operations shared by the
// application service and its persistence repository. Application-only
// operations such as content-policy authorization and document download do
// not belong here.
type KnowledgePageCoreService interface {
	KnowledgePages(context.Context, string, string, string, int, bool, Authority) (KnowledgePagePage, error)
	KnowledgePage(context.Context, string, string, bool, Authority) (KnowledgePageView, error)
	CreateKnowledgePage(context.Context, string, KnowledgePageCreate, Authority) (KnowledgePageView, error)
	KnowledgePageReferences(context.Context, string, string, KnowledgePageReferenceDirection, string, int, Authority) (KnowledgePageReferencePage, error)
	CreateKnowledgePageReference(context.Context, string, string, KnowledgePageReferenceCreate, Authority) (KnowledgePageReferenceCreateResult, error)
	SaveKnowledgePageDraft(context.Context, string, string, KnowledgePageDraftWrite, Authority) (KnowledgePageView, error)
	MoveKnowledgePage(context.Context, string, string, KnowledgePageMove, Authority) (KnowledgePage, error)
	PreviewKnowledgePageMove(context.Context, string, string, KnowledgePageMovePreview, Authority) (KnowledgePageMoveImpact, error)
	ReorderKnowledgePages(context.Context, string, KnowledgePageReorder, Authority) (KnowledgePagePage, error)
	CopyKnowledgePage(context.Context, string, string, KnowledgePageCopy, Authority) (KnowledgePageView, error)
	SetKnowledgePageArchived(context.Context, string, string, bool, KnowledgePageStateWrite, Authority) (KnowledgePage, error)
	PublishKnowledgePage(context.Context, string, string, KnowledgePagePublish, Authority) (KnowledgePageView, error)
	KnowledgePageVersions(context.Context, string, string, int64, int, Authority) (KnowledgePageVersionPage, error)
	KnowledgePageVersion(context.Context, string, string, int64, Authority) (KnowledgePageVersion, error)
	DiffKnowledgePageVersions(context.Context, string, string, KnowledgePageVersionDiffInput, Authority) (KnowledgePageVersionDiff, error)
	RestoreKnowledgePageVersion(context.Context, string, string, int64, KnowledgePageRestoreVersion, Authority) (KnowledgePageView, error)
	PageNativeOperationService
	PageGridAutomationService
}

// KnowledgePageRepository is the persistence boundary consumed by the
// application service. It evaluates policy but does not perform the
// application layer's final authorization or document delivery.
type KnowledgePageRepository interface {
	KnowledgePageCoreService
	KnowledgePageContentPolicyRepository
}

type PageService interface {
	KnowledgePageCoreService
	KnowledgePageContentPolicyService
}

const (
	KnowledgePageAccessInherit    = "inherit"
	KnowledgePageAccessRestricted = "restricted"
	KnowledgePageGrantPage        = "page"
	KnowledgePageGrantSubtree     = "subtree"
)

// KnowledgePageAccessGrant never replaces library membership. The effective
// role is the lesser of the current library role and this grant's role.
type KnowledgePageAccessGrant struct {
	Subject         KnowledgeAccessSubject `json:"subject"`
	Role            string                 `json:"role"`  // reader, editor, manager
	Scope           string                 `json:"scope"` // page, subtree
	CreatedByUserID string                 `json:"created_by_user_id"`
	CreatedAt       time.Time              `json:"created_at"`
}

type KnowledgePageAccessPolicy struct {
	PageID           string                     `json:"page_id"`
	LibraryID        string                     `json:"library_id"`
	Mode             string                     `json:"mode"`
	Revision         int64                      `json:"revision"`
	ProjectionStatus string                     `json:"projection_status"`
	Grants           []KnowledgePageAccessGrant `json:"grants"`
	UpdatedByUserID  string                     `json:"updated_by_user_id"`
	UpdatedAt        time.Time                  `json:"updated_at"`
}

type KnowledgePageEffectiveAccess struct {
	Allowed     bool   `json:"allowed"`
	Role        string `json:"role,omitempty"`
	Mode        string `json:"mode"`
	Source      string `json:"source"` // library, direct_grant, ancestor_grant, manager
	GrantPageID string `json:"grant_page_id,omitempty"`
	GrantScope  string `json:"grant_scope,omitempty"`
}

type KnowledgePageAccessWrite struct {
	ClientID         string                     `json:"client_id"`
	ExpectedRevision int64                      `json:"expected_revision"`
	Mode             string                     `json:"mode"`
	Grants           []KnowledgePageAccessGrant `json:"grants"`
}

type KnowledgePageMovePreview struct {
	ParentPageID string `json:"parent_page_id,omitempty"`
	Rank         string `json:"rank"`
}

type KnowledgePageMoveImpact struct {
	PageID          string                   `json:"page_id"`
	ParentPageID    string                   `json:"parent_page_id,omitempty"`
	AffectedPages   int                      `json:"affected_pages"`
	ChangedSubjects []KnowledgeAccessSubject `json:"changed_subjects"`
	RequiresConfirm bool                     `json:"requires_confirm"`
	Token           string                   `json:"token"`
}

type KnowledgePageAccessService interface {
	KnowledgePageAccess(context.Context, string, string, Authority) (KnowledgePageAccessPolicy, error)
	SetKnowledgePageAccess(context.Context, string, string, KnowledgePageAccessWrite, Authority) (KnowledgePageAccessPolicy, error)
}

// KnowledgePageAccessRepository is the trusted server-side policy boundary
// used by application authorization and retrieval. Permission tokens are
// derived by Knowledge and are never accepted from browser input.
type KnowledgePageAccessRepository interface {
	KnowledgePageAccessService
	EvaluateKnowledgePageAccess(context.Context, string, string, Authority) (KnowledgePageEffectiveAccess, error)
	KnowledgePagePermissionTokens(context.Context, string, Authority) ([]string, error)
	KnowledgePageProjectionPermission(context.Context, string, string, Authority) (string, int64, error)
}

// JSONClone provides a small contract-test helper without defining a second
// wire representation for page values.
func JSONClone[T any](value T) (T, error) {
	var out T
	raw, err := json.Marshal(value)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(raw, &out)
	return out, err
}
