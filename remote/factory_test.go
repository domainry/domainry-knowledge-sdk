package remote

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	agentsdk "github.com/domainry/domainry-agent-sdk"
	"github.com/domainry/domainry-agent-sdk/persistence"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	"github.com/domainry/domainry-knowledge-sdk/saashost"
)

type sourceAuthorizerStub struct {
	called    bool
	source    agentsdk.KnowledgeDocumentSourceAccess
	authority agentsdk.ConversationAuthority
}

type attachmentPermissionResolverStub struct {
	authority agentsdk.ConversationAuthority
}

func (*attachmentPermissionResolverStub) AuthorizeConversationAttachment(context.Context, string, agentsdk.ConversationAuthority) error {
	return nil
}
func (s *attachmentPermissionResolverStub) ResolveConversationAttachmentPermissions(_ context.Context, authority agentsdk.ConversationAuthority) (agentsdk.ConversationAttachmentPermissionScope, error) {
	s.authority = authority
	return agentsdk.ConversationAttachmentPermissionScope{OrganizationIDs: []string{"org-a", "region-a"}}, nil
}

func (s *sourceAuthorizerStub) AuthorizeKnowledgeDocumentSource(_ context.Context, source agentsdk.KnowledgeDocumentSourceAccess, authority agentsdk.ConversationAuthority) error {
	s.called, s.source, s.authority = true, source, authority
	return nil
}

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

func TestSourceBackedDocumentUploadAnswersHostAuthorizationChallenge(t *testing.T) {
	authority := agentsdk.ConversationAuthority{Known: true, RuntimeID: "runtime-a", WorkspaceID: "workspace-a", UserID: "user-a"}
	source := agentsdk.KnowledgeDocumentSourceAccess{Namespace: agentsdk.KnowledgeDocumentSourceNamespaceRuntimeRecord, ResourceType: "meeting", ResourceID: "meeting-a"}
	authorizer := &sourceAuthorizerStub{}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case knowledgecontract.SaaSDiscoveryPath:
			_ = json.NewEncoder(writer).Encode(knowledgecontract.Descriptor{ProtocolVersion: knowledgecontract.SaaSProtocolVersionV1, Mode: knowledgecontract.DeploymentModeSaaS, Audience: "runtime-a", Capabilities: append([]string(nil), knowledgecontract.SaaSCapabilitiesV1...)})
		case knowledgecontract.SaaSInvokePath:
			var envelope knowledgecontract.SaaSRequest
			if json.NewDecoder(request.Body).Decode(&envelope) != nil || envelope.Operation != "service.document_upload_source" {
				t.Fatalf("request=%+v", envelope)
			}
			if len(envelope.Grants) == 0 {
				input, _ := json.Marshal(struct {
					Source    agentsdk.KnowledgeDocumentSourceAccess `json:"source"`
					Authority agentsdk.ConversationAuthority         `json:"authority"`
				}{source, authority})
				_ = json.NewEncoder(writer).Encode(knowledgecontract.SaaSResponse{Challenge: &knowledgecontract.SaaSChallenge{Token: "source-grant", Kind: "authorize.document_source", Input: input}})
				return
			}
			if len(envelope.Grants) != 1 || envelope.Grants[0].Token != "source-grant" || string(envelope.Grants[0].Result) != `{}` {
				t.Fatalf("grants=%+v", envelope.Grants)
			}
			result, _ := json.Marshal(agentsdk.KnowledgeDocument{ID: "document-a", LibraryID: "library-a", Revision: 1})
			_ = json.NewEncoder(writer).Encode(knowledgecontract.SaaSResponse{Result: result})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	binding, err := NewFactory(Config{Endpoint: server.URL}).OpenSaaS(t.Context(), knowledgecontract.ApplicationRef{RuntimeID: "runtime-a"})
	if err != nil {
		t.Fatal(err)
	}
	service := binding.Runtime().NewService("runtime-a", knowledgecontract.Options{SourceAuthorizer: authorizer})
	document, err := service.UploadKnowledgeDocumentForSource(t.Context(), "library-a", agentsdk.KnowledgeDocumentUpload{ClientID: "upload-a", Filename: "meeting.txt", Data: []byte("content")}, source, authority)
	if err != nil || document.ID != "document-a" || !authorizer.called || authorizer.source != source || authorizer.authority.UserID != authority.UserID {
		t.Fatalf("document=%+v source=%+v authority=%+v called=%t err=%v", document, authorizer.source, authorizer.authority, authorizer.called, err)
	}
}

func TestAttachmentPermissionChallengeUsesLiveHostResolver(t *testing.T) {
	authority := agentsdk.ConversationAuthority{Known: true, RuntimeID: "runtime-a", WorkspaceID: "workspace-a", UserID: "user-a"}
	input, _ := json.Marshal(struct {
		Authority agentsdk.ConversationAuthority `json:"authority"`
	}{authority})
	resolver := &attachmentPermissionResolverStub{}
	grant, err := answerChallenge(t.Context(), knowledgecontract.SaaSChallenge{Token: "permission-grant", Kind: "resolve.attachment_permissions", Input: input}, knowledgecontract.Options{AttachmentAuthorizer: resolver})
	if err != nil || resolver.authority != authority || grant.Token != "permission-grant" {
		t.Fatalf("grant=%+v authority=%+v err=%v", grant, resolver.authority, err)
	}
	var scope agentsdk.ConversationAttachmentPermissionScope
	if json.Unmarshal(grant.Result, &scope) != nil || len(scope.OrganizationIDs) != 2 || scope.OrganizationIDs[0] != "org-a" {
		t.Fatalf("scope=%+v", scope)
	}
}
