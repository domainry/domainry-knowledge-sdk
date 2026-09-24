package contract

import (
	"context"
	"fmt"
	"slices"
	"strings"

	agentsdk "github.com/domainry/domainry-agent-sdk"
	persistence "github.com/domainry/domainry-agent-sdk/persistence"
)

const (
	SaaSProtocolVersionV1 = "knowledge.saas.v1"
	DeploymentModeSaaS    = "saas"
)

var SaaSCapabilitiesV1 = []string{
	"artifacts.read",
	"artifacts.write",
	"attachments.read",
	"attachments.write",
	"knowledge_documents.read",
	"knowledge_documents.write",
	"knowledge_libraries.read",
	"knowledge_libraries.write",
	"knowledge_sources.read",
	"runtime.knowledge",
	"subjects.lifecycle",
}

// Descriptor is returned by the SaaS discovery endpoint before a client is
// allowed to bind the service to one Runtime audience.
type Descriptor struct {
	ProtocolVersion string   `json:"protocol_version"`
	Mode            string   `json:"mode"`
	Audience        string   `json:"audience"`
	Capabilities    []string `json:"capabilities"`
}

func (descriptor Descriptor) Validate() error {
	if descriptor.ProtocolVersion != SaaSProtocolVersionV1 || descriptor.Mode != DeploymentModeSaaS {
		return fmt.Errorf("Knowledge SaaS protocol descriptor is unsupported")
	}
	if strings.TrimSpace(descriptor.Audience) == "" || descriptor.Audience != strings.TrimSpace(descriptor.Audience) || len(descriptor.Audience) > 255 {
		return fmt.Errorf("Knowledge SaaS protocol audience is invalid")
	}
	capabilities := slices.Clone(descriptor.Capabilities)
	slices.Sort(capabilities)
	expected := slices.Clone(SaaSCapabilitiesV1)
	slices.Sort(expected)
	if !slices.Equal(capabilities, expected) {
		return fmt.Errorf("Knowledge SaaS protocol capabilities are unsupported")
	}
	return nil
}

// ArtifactMutationService is the deployment-neutral idempotent artifact write
// boundary. Module mode may additionally provide a caller-owned SQL transaction
// extension, but SaaS callers only use these Knowledge-owned commands.
type ArtifactMutationService interface {
	ArtifactRecord(context.Context, string, int64, agentsdk.ConversationAuthority) (persistence.ConversationArtifactRecord, error)
	SaveArtifact(context.Context, persistence.ConversationArtifactWrite, agentsdk.ConversationAuthority) (persistence.ConversationArtifactRecord, error)
	SaveArtifactExport(context.Context, persistence.ConversationArtifactExportWrite, agentsdk.ConversationAuthority) (agentsdk.ConversationArtifactExport, error)
}
