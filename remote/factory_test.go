package remote

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	agentsdk "github.com/domainry/domainry-agent-sdk"
	"github.com/domainry/domainry-agent-sdk/persistence"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	"github.com/domainry/domainry-knowledge-sdk/saashost"
)

func TestFactoryValidatesAudienceAndUsesAuthenticatedRuntimeBoundary(t *testing.T) {
	var invoked bool
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer service-token" || request.Header.Get(saashost.RuntimeIDHeader) != "runtime-a" {
			t.Fatalf("headers=%v", request.Header)
		}
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case knowledgecontract.SaaSDiscoveryPath:
			_ = json.NewEncoder(writer).Encode(knowledgecontract.Descriptor{ProtocolVersion: knowledgecontract.SaaSProtocolVersionV1, Mode: knowledgecontract.DeploymentModeSaaS, Audience: "runtime-a", Capabilities: append([]string(nil), knowledgecontract.SaaSCapabilitiesV1...)})
		case knowledgecontract.SaaSInvokePath:
			invoked = true
			var envelope knowledgecontract.SaaSRequest
			if json.NewDecoder(request.Body).Decode(&envelope) != nil || envelope.Operation != "artifacts.save" {
				t.Fatalf("request=%+v", envelope)
			}
			record := persistence.ConversationArtifactRecord{Artifact: agentsdk.ConversationArtifact{ID: "artifact-a", Version: 1}}
			raw, _ := json.Marshal(record)
			_ = json.NewEncoder(writer).Encode(knowledgecontract.SaaSResponse{Result: raw})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	binding, err := NewFactory(Config{Endpoint: server.URL, ServiceAccessToken: "service-token"}).OpenSaaS(t.Context(), knowledgecontract.ApplicationRef{RuntimeID: "runtime-a"})
	if err != nil {
		t.Fatal(err)
	}
	authority := agentsdk.ConversationAuthority{Known: true, RuntimeID: "runtime-a", WorkspaceID: "workspace-a", UserID: "user-a"}
	record, err := binding.ArtifactMutations().SaveArtifact(t.Context(), persistence.ConversationArtifactWrite{ClientID: "request-a"}, authority)
	if err != nil || record.Artifact.ID != "artifact-a" || !invoked {
		t.Fatalf("record=%+v invoked=%v err=%v", record, invoked, err)
	}
}
