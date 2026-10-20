// Package provider defines the provider-neutral construction boundary used by
// product composition roots. Implementations stay in domainry-knowledge.
package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	agentsdk "github.com/domainry/domainry-agent-sdk"
	connector "github.com/domainry/domainry-connector-sdk"
	"github.com/domainry/domainry-knowledge-sdk/contract"
)

// CitationMapping selects trusted citation fields from a provider response.
// Pointers use RFC 6901 syntax and are validated by the implementation.
type CitationMapping struct {
	MetadataObject  *string `json:"metadata_object,omitempty"`
	Items           string  `json:"items"`
	Many            bool    `json:"many,omitempty"`
	DocumentID      string  `json:"doc_id,omitempty"`
	Title           string  `json:"title,omitempty"`
	URL             string  `json:"url,omitempty"`
	Excerpt         string  `json:"excerpt,omitempty"`
	ChunkID         string  `json:"chunk_id,omitempty"`
	DocumentVersion string  `json:"document_version,omitempty"`
	ContentSHA256   string  `json:"content_sha256,omitempty"`
	SourceRevision  string  `json:"source_revision,omitempty"`
	Page            string  `json:"page,omitempty"`
	Section         string  `json:"section,omitempty"`
	AnchorRef       string  `json:"anchor_ref,omitempty"`
	RecordKind      string  `json:"record_kind,omitempty"`
	Score           string  `json:"score,omitempty"`
}

type ResponseMapping struct {
	Search *CitationMapping `json:"search,omitempty"`
	Fetch  *CitationMapping `json:"fetch,omitempty"`
}

// Config is trusted startup configuration. It contains no persistence handle.
// Transport is an optional host-supplied Connector transport used by tests and
// specialized products; normal deployments leave it nil.
type Config struct {
	DocumentManagement                         bool
	DocumentPermissionIDs                      []string
	ReadPermissionIDs                          []string
	AnalysisDocumentIDs                        []string
	ResponseMapping                            *ResponseMapping
	InvalidResponseMapping                     bool
	InvalidAnalysisDocumentIDs                 bool
	InvalidSearchConfig                        bool
	BaseURL, APIKey, TeamID, KBID, WorkspaceID string
	RuntimeID                                  string
	TopK                                       int
	SearchMode                                 string
	SearchRerank                               *bool
	SearchRecordKinds                          []string
	Client                                     *http.Client
	Transport                                  connector.Transport
	PermissionIDs                              func(context.Context, agentsdk.ConversationAuthority) ([]string, error)
	AuthorizeWorkspace                         func(context.Context, agentsdk.ConversationAuthority) error
}

func (config Config) Configured() bool {
	return strings.TrimSpace(config.BaseURL+config.TeamID+config.KBID+config.WorkspaceID) != "" ||
		config.TopK != 0 || config.ResponseMapping != nil || config.InvalidResponseMapping ||
		config.InvalidAnalysisDocumentIDs || config.InvalidSearchConfig || config.DocumentManagement ||
		config.SearchMode != "" || config.SearchRerank != nil || config.SearchRecordKinds != nil ||
		config.DocumentPermissionIDs != nil || config.ReadPermissionIDs != nil || config.AnalysisDocumentIDs != nil
}

// ConfigFromEnvironment parses Agent's default Knowledge source configuration
// without selecting an implementation. The outer composition root supplies the
// Factory separately.
func ConfigFromEnvironment() Config {
	config := Config{
		BaseURL:     os.Getenv("AGENT_KNOWLEDGE_BASE_URL"),
		APIKey:      os.Getenv("AGENT_KNOWLEDGE_API_KEY"),
		TeamID:      os.Getenv("AGENT_KNOWLEDGE_TEAM_ID"),
		KBID:        os.Getenv("AGENT_KNOWLEDGE_KB_ID"),
		WorkspaceID: os.Getenv("AGENT_KNOWLEDGE_WORKSPACE_ID"),
	}
	if strings.TrimSpace(config.APIKey) == "" {
		config.APIKey = os.Getenv("AGENT_PROVIDER_API_KEY")
	}
	if raw := strings.TrimSpace(os.Getenv("AGENT_KNOWLEDGE_TOP_K")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 {
			value = -1
		}
		config.TopK = value
	}
	config.SearchMode = strings.TrimSpace(os.Getenv("AGENT_KNOWLEDGE_SEARCH_MODE"))
	if raw := strings.TrimSpace(os.Getenv("AGENT_KNOWLEDGE_RERANK")); raw != "" {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			config.InvalidSearchConfig = true
		} else {
			config.SearchRerank = &value
		}
	}
	if raw := strings.TrimSpace(os.Getenv("AGENT_KNOWLEDGE_RECORD_KINDS")); raw != "" {
		if json.Unmarshal([]byte(raw), &config.SearchRecordKinds) != nil {
			config.InvalidSearchConfig = true
		}
	}
	if raw := strings.TrimSpace(os.Getenv("AGENT_KNOWLEDGE_RESPONSE_MAPPING")); raw != "" {
		config.ResponseMapping, config.InvalidResponseMapping = parseResponseMapping(raw)
	}
	if raw := strings.TrimSpace(os.Getenv("AGENT_KNOWLEDGE_ANALYSIS_DOCUMENT_IDS")); raw != "" {
		config.InvalidAnalysisDocumentIDs = json.Unmarshal([]byte(raw), &config.AnalysisDocumentIDs) != nil
	}
	return config
}

func parseResponseMapping(raw string) (*ResponseMapping, bool) {
	var mapping ResponseMapping
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&mapping); err != nil {
		return nil, true
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		return nil, true
	}
	return &mapping, false
}

// Source supports the three public views required by Agent's default,
// library and datasource configurations.
type Source interface {
	contract.ConversationKnowledge
	agentsdk.KnowledgeDatasourceSource
}

// AdapterFactory is supplied by the outer product composition root. Knowledge
// owns provider behavior, while the concrete Connector HTTP adapter remains in
// the Connectors implementation repository.
type AdapterFactory func(connector.Transport) (connector.Adapter, error)

type Factory interface {
	NewSource(Config) (Source, error)
	NewAttachmentSource(Config, string) (agentsdk.ConversationAttachmentKnowledge, error)
}

func RequireFactory(factory Factory, configured bool) error {
	if configured && factory == nil {
		return fmt.Errorf("Knowledge provider factory is required by configured sources")
	}
	return nil
}
