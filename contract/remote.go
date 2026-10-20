package contract

import (
	"fmt"
	"slices"
	"strings"
)

const (
	SaaSProtocolVersionV1 = "knowledge.saas.v1"
	DeploymentModeSaaS    = "saas"
)

var SaaSCapabilitiesV1 = []string{
	"files.read",
	"files.write",
	"knowledge_documents.read",
	"knowledge_documents.write",
	"knowledge_libraries.read",
	"knowledge_libraries.write",
	"knowledge_sources.read",
	"runtime.knowledge",
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
