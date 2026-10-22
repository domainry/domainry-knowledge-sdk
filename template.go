package knowledgesdk

import (
	"context"
	"time"
)

// KnowledgePageTemplate is a library-scoped, versioned page-subtree recipe.
// Instances never retain a live link to the template, so later template
// versions cannot overwrite page content that users have already created.
type KnowledgePageTemplate struct {
	ID              string    `json:"id"`
	LibraryID       string    `json:"library_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	RootContentType string    `json:"root_content_type"`
	NodeCount       int       `json:"node_count"`
	LatestVersion   int64     `json:"latest_version"`
	Revision        int64     `json:"revision"`
	CreatedByUserID string    `json:"created_by_user_id"`
	UpdatedByUserID string    `json:"updated_by_user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// KnowledgePageTemplateNode is a frozen page draft in a template version.
// NodeKey and ParentNodeKey describe only the template tree; SourcePageID is
// retained so internal page and Block references can be remapped atomically
// when an independent instance is created.
type KnowledgePageTemplateNode struct {
	NodeKey       string                            `json:"node_key"`
	ParentNodeKey string                            `json:"parent_node_key,omitempty"`
	SourcePageID  string                            `json:"source_page_id"`
	Rank          string                            `json:"rank"`
	Title         string                            `json:"title"`
	ContentType   string                            `json:"content_type"`
	Body          PageContent                       `json:"body"`
	Attachments   []KnowledgePageTemplateAttachment `json:"attachments,omitempty"`
}

// KnowledgePageTemplateAttachment freezes the page-to-document relation, not
// the immutable document itself. Instances reuse the authorized source file
// while owning an independent page relation.
type KnowledgePageTemplateAttachment struct {
	DocumentID string `json:"document_id"`
	Title      string `json:"title,omitempty"`
	Rank       string `json:"rank"`
}

type KnowledgePageTemplateVersion struct {
	TemplateID      string                      `json:"template_id"`
	LibraryID       string                      `json:"library_id"`
	Version         int64                       `json:"version"`
	Title           string                      `json:"title"`
	Description     string                      `json:"description,omitempty"`
	Nodes           []KnowledgePageTemplateNode `json:"nodes"`
	CreatedByUserID string                      `json:"created_by_user_id"`
	CreatedAt       time.Time                   `json:"created_at"`
}

type KnowledgePageTemplateVersionSummary struct {
	TemplateID      string    `json:"template_id"`
	Version         int64     `json:"version"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	NodeCount       int       `json:"node_count"`
	CreatedByUserID string    `json:"created_by_user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type KnowledgePageTemplatePage struct {
	Items     []KnowledgePageTemplate `json:"items"`
	NextAfter string                  `json:"next_after,omitempty"`
	Complete  bool                    `json:"complete"`
}

type KnowledgePageTemplateVersionPage struct {
	Items      []KnowledgePageTemplateVersionSummary `json:"items"`
	NextBefore int64                                 `json:"next_before,omitempty"`
	Complete   bool                                  `json:"complete"`
}

type KnowledgePageTemplateCreate struct {
	ClientID     string `json:"client_id"`
	SourcePageID string `json:"source_page_id"`
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
}

type KnowledgePageTemplateUpdate struct {
	ClientID         string `json:"client_id"`
	ExpectedRevision int64  `json:"expected_revision"`
	SourcePageID     string `json:"source_page_id"`
	Title            string `json:"title"`
	Description      string `json:"description,omitempty"`
}

type KnowledgePageTemplateInstantiate struct {
	ClientID     string `json:"client_id"`
	Version      int64  `json:"version,omitempty"`
	ParentPageID string `json:"parent_page_id,omitempty"`
	Rank         string `json:"rank,omitempty"`
	RootTitle    string `json:"root_title,omitempty"`
}

type KnowledgePageTemplateInstance struct {
	TemplateID     string            `json:"template_id"`
	Version        int64             `json:"version"`
	Root           KnowledgePageView `json:"root"`
	CreatedPageIDs []string          `json:"created_page_ids"`
}

type PageTemplateService interface {
	KnowledgePageTemplates(context.Context, string, string, int, Authority) (KnowledgePageTemplatePage, error)
	CreateKnowledgePageTemplate(context.Context, string, KnowledgePageTemplateCreate, Authority) (KnowledgePageTemplate, error)
	UpdateKnowledgePageTemplate(context.Context, string, string, KnowledgePageTemplateUpdate, Authority) (KnowledgePageTemplate, error)
	KnowledgePageTemplateVersions(context.Context, string, string, int64, int, Authority) (KnowledgePageTemplateVersionPage, error)
	KnowledgePageTemplateVersion(context.Context, string, string, int64, Authority) (KnowledgePageTemplateVersion, error)
	InstantiateKnowledgePageTemplate(context.Context, string, string, KnowledgePageTemplateInstantiate, Authority) (KnowledgePageTemplateInstance, error)
}
