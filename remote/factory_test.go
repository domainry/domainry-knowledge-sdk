package remote

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	agentsdk "github.com/domainry/domainry-agent-sdk"
	"github.com/domainry/domainry-agent-sdk/persistence"
	"github.com/domainry/domainry-foundation/requestcontext"
	"github.com/domainry/domainry-foundation/telemetry"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	"github.com/domainry/domainry-knowledge-sdk/saashost"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	apiTrace "go.opentelemetry.io/otel/trace"
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

func TestFactoryPropagatesTraceAndCorrelationWithoutRequestIdentity(t *testing.T) {
	provider := trace.NewTracerProvider()
	previousProvider := otel.GetTracerProvider()
	otel.SetTracerProvider(provider)
	t.Cleanup(func() {
		_ = provider.Shutdown(context.Background())
		otel.SetTracerProvider(previousProvider)
	})

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		traceparent := request.Header.Get("traceparent")
		if len(traceparent) != 55 || traceparent[3:35] != "4bf92f3577b34da6a3ce929d0e0e4736" {
			t.Fatalf("traceparent=%q", traceparent)
		}
		if got := request.Header.Get(telemetry.CorrelationIDHeader); got != "correlation-a" {
			t.Fatalf("correlation id=%q", got)
		}
		if got := request.Header.Get(telemetry.RequestIDHeader); got != "" {
			t.Fatalf("request id crossed service boundary: %q", got)
		}
		_ = json.NewEncoder(writer).Encode(knowledgecontract.Descriptor{ProtocolVersion: knowledgecontract.SaaSProtocolVersionV1, Mode: knowledgecontract.DeploymentModeSaaS, Audience: "runtime-a", Capabilities: append([]string(nil), knowledgecontract.SaaSCapabilitiesV1...)})
	}))
	t.Cleanup(server.Close)

	traceID, _ := apiTrace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := apiTrace.SpanIDFromHex("00f067aa0ba902b7")
	ctx := apiTrace.ContextWithSpanContext(t.Context(), apiTrace.NewSpanContext(apiTrace.SpanContextConfig{TraceID: traceID, SpanID: spanID, TraceFlags: apiTrace.FlagsSampled}))
	ctx = requestcontext.WithCorrelationID(ctx, "correlation-a")
	if _, err := NewFactory(Config{Endpoint: server.URL, HTTPClient: server.Client()}).OpenSaaS(ctx, knowledgecontract.ApplicationRef{RuntimeID: "runtime-a"}); err != nil {
		t.Fatal(err)
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
			if json.NewDecoder(request.Body).Decode(&envelope) != nil || envelope.Operation != knowledgecontract.SaaSDocumentUploadPrepareOperation {
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
			result, _ := json.Marshal(struct{}{})
			_ = json.NewEncoder(writer).Encode(knowledgecontract.SaaSResponse{Result: result})
		case knowledgecontract.SaaSDocumentUploadPath:
			mediaType, parameters, parseErr := mime.ParseMediaType(request.Header.Get("Content-Type"))
			if parseErr != nil || mediaType != "multipart/form-data" {
				t.Fatalf("content-type=%q err=%v", request.Header.Get("Content-Type"), parseErr)
			}
			reader := multipart.NewReader(request.Body, parameters["boundary"])
			metadataPart, partErr := reader.NextPart()
			if partErr != nil || metadataPart.FormName() != "metadata" {
				t.Fatalf("metadata part err=%v", partErr)
			}
			var metadata knowledgecontract.SaaSDocumentUploadRequest
			if json.NewDecoder(metadataPart).Decode(&metadata) != nil || metadata.Operation != knowledgecontract.SaaSDocumentUploadSourceOperation || metadata.Source == nil || *metadata.Source != source || len(metadata.Grants) != 1 {
				t.Fatalf("metadata=%+v", metadata)
			}
			contentPart, partErr := reader.NextPart()
			content, readErr := io.ReadAll(contentPart)
			if partErr != nil || readErr != nil || !bytes.Equal(content, []byte("content")) {
				t.Fatalf("content=%q partErr=%v readErr=%v", content, partErr, readErr)
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
	raw := []byte("content")
	document, err := service.UploadKnowledgeDocumentForSource(t.Context(), "library-a", agentsdk.KnowledgeDocumentUpload{ClientID: "upload-a", Filename: "meeting.txt", Bytes: int64(len(raw)), SHA256: "ed7002b439e9ac845f22357d822bac144473b68318c28233ffb2446e8f300212", Content: bytes.NewReader(raw)}, source, authority)
	if err != nil || document.ID != "document-a" || !authorizer.called || authorizer.source != source || authorizer.authority.UserID != authority.UserID {
		t.Fatalf("document=%+v source=%+v authority=%+v called=%t err=%v", document, authorizer.source, authorizer.authority, authorizer.called, err)
	}
}

func TestResumableDocumentChunkUsesBoundedMultipartTransfer(t *testing.T) {
	authority := agentsdk.ConversationAuthority{Known: true, RuntimeID: "runtime-a", WorkspaceID: "workspace-a", UserID: "user-a"}
	raw := []byte("content")
	digest := "ed7002b439e9ac845f22357d822bac144473b68318c28233ffb2446e8f300212"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case knowledgecontract.SaaSDiscoveryPath:
			_ = json.NewEncoder(writer).Encode(knowledgecontract.Descriptor{ProtocolVersion: knowledgecontract.SaaSProtocolVersionV1, Mode: knowledgecontract.DeploymentModeSaaS, Audience: "runtime-a", Capabilities: append([]string(nil), knowledgecontract.SaaSCapabilitiesV1...)})
		case knowledgecontract.SaaSInvokePath:
			var envelope knowledgecontract.SaaSRequest
			if json.NewDecoder(request.Body).Decode(&envelope) != nil || envelope.Operation != knowledgecontract.SaaSDocumentUploadPrepareOperation {
				t.Fatalf("request=%+v", envelope)
			}
			result, _ := json.Marshal(struct{}{})
			_ = json.NewEncoder(writer).Encode(knowledgecontract.SaaSResponse{Result: result})
		case knowledgecontract.SaaSDocumentUploadPath:
			mediaType, parameters, parseErr := mime.ParseMediaType(request.Header.Get("Content-Type"))
			if parseErr != nil || mediaType != "multipart/form-data" {
				t.Fatalf("content-type=%q err=%v", request.Header.Get("Content-Type"), parseErr)
			}
			reader := multipart.NewReader(request.Body, parameters["boundary"])
			metadataPart, partErr := reader.NextPart()
			var metadata knowledgecontract.SaaSDocumentChunkUploadRequest
			if partErr != nil || json.NewDecoder(metadataPart).Decode(&metadata) != nil || metadata.Operation != knowledgecontract.SaaSDocumentUploadChunkOperation || metadata.Input.Offset != 0 || metadata.Input.ChunkBytes != int64(len(raw)) || metadata.Input.ChunkSHA256 != digest {
				t.Fatalf("metadata=%+v err=%v", metadata, partErr)
			}
			contentPart, partErr := reader.NextPart()
			content, readErr := io.ReadAll(contentPart)
			if partErr != nil || readErr != nil || !bytes.Equal(content, raw) {
				t.Fatalf("content=%q partErr=%v readErr=%v", content, partErr, readErr)
			}
			document := agentsdk.KnowledgeDocument{ID: "document-a", LibraryID: "library-a", Revision: 1}
			result, _ := json.Marshal(agentsdk.KnowledgeDocumentUploadProgress{NextOffset: int64(len(raw)), Complete: true, Document: &document})
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
	service := binding.Runtime().NewService("runtime-a", knowledgecontract.Options{})
	progress, err := service.UploadKnowledgeDocumentChunk(t.Context(), "library-a", agentsdk.KnowledgeDocumentUploadChunk{ClientID: "upload-a", Filename: "notes.txt", Bytes: int64(len(raw)), SHA256: digest, Offset: 0, ChunkBytes: int64(len(raw)), ChunkSHA256: digest, Content: bytes.NewReader(raw)}, authority)
	if err != nil || !progress.Complete || progress.NextOffset != int64(len(raw)) || progress.Document == nil || progress.Document.ID != "document-a" {
		t.Fatalf("progress=%+v err=%v", progress, err)
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

func TestDocumentDownloadUsesAuthorizedRangeReader(t *testing.T) {
	content := []byte("abcdef")
	authority := agentsdk.ConversationAuthority{Known: true, RuntimeID: "runtime-a", WorkspaceID: "workspace-a", UserID: "user-a"}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case knowledgecontract.SaaSDiscoveryPath:
			_ = json.NewEncoder(writer).Encode(knowledgecontract.Descriptor{ProtocolVersion: knowledgecontract.SaaSProtocolVersionV1, Mode: knowledgecontract.DeploymentModeSaaS, Audience: "runtime-a", Capabilities: append([]string(nil), knowledgecontract.SaaSCapabilitiesV1...)})
		case knowledgecontract.SaaSInvokePath:
			var envelope knowledgecontract.SaaSRequest
			if json.NewDecoder(request.Body).Decode(&envelope) != nil || envelope.Operation != knowledgecontract.SaaSDocumentDownloadPrepareOperation {
				t.Fatalf("request=%+v", envelope)
			}
			result, _ := json.Marshal(agentsdk.KnowledgeDocument{ID: "document-a", LibraryID: "library-a", Filename: "sample.txt", Bytes: int64(len(content)), SHA256: "bef57ec7f53a6d40beb640a780a639c83bc29ac8a9816f1c5c6dcd93c4721f33", Revision: 1})
			_ = json.NewEncoder(writer).Encode(knowledgecontract.SaaSResponse{Result: result})
		case knowledgecontract.SaaSDocumentDownloadPath:
			raw, decodeErr := base64.RawURLEncoding.DecodeString(request.Header.Get(knowledgecontract.SaaSTransferHeader))
			var transfer knowledgecontract.SaaSDocumentDownloadRequest
			if decodeErr != nil || json.Unmarshal(raw, &transfer) != nil || transfer.Library != "library-a" || transfer.ID != "document-a" || transfer.Authority != authority {
				t.Fatalf("transfer=%+v err=%v", transfer, decodeErr)
			}
			value := strings.TrimPrefix(request.Header.Get("Range"), "bytes=")
			bounds := strings.Split(value, "-")
			if len(bounds) != 2 {
				t.Fatalf("range=%q", value)
			}
			start, startErr := strconv.Atoi(bounds[0])
			end, endErr := strconv.Atoi(bounds[1])
			if startErr != nil || endErr != nil || start < 0 || end < start || end >= len(content) {
				t.Fatalf("range=%q", value)
			}
			writer.Header().Set("Content-Range", "bytes "+strconv.Itoa(start)+"-"+strconv.Itoa(end)+"/"+strconv.Itoa(len(content)))
			writer.Header().Set("Content-Length", strconv.Itoa(end-start+1))
			writer.WriteHeader(http.StatusPartialContent)
			_, _ = writer.Write(content[start : end+1])
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	binding, err := NewFactory(Config{Endpoint: server.URL}).OpenSaaS(t.Context(), knowledgecontract.ApplicationRef{RuntimeID: "runtime-a"})
	if err != nil {
		t.Fatal(err)
	}
	download, err := binding.Runtime().NewService("runtime-a", knowledgecontract.Options{}).DownloadKnowledgeDocument(t.Context(), "library-a", "document-a", authority)
	if err != nil {
		t.Fatal(err)
	}
	defer download.Content.Close()
	first := make([]byte, 2)
	if n, readErr := download.Content.Read(first); readErr != nil || n != 2 || string(first) != "ab" {
		t.Fatalf("first=%q n=%d err=%v", first, n, readErr)
	}
	if offset, seekErr := download.Content.Seek(-3, io.SeekEnd); seekErr != nil || offset != 3 {
		t.Fatalf("offset=%d err=%v", offset, seekErr)
	}
	rest, readErr := io.ReadAll(download.Content)
	if readErr != nil || string(rest) != "def" {
		t.Fatalf("rest=%q err=%v", rest, readErr)
	}
}
