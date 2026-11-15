package remote

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	knowledgesdk "github.com/domainry/domainry-knowledge-sdk"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	"github.com/domainry/domainry-knowledge-sdk/saashost"
)

func TestRemoteLibraryRetrievalKeepsScopeAndDoesNotFallBack(t *testing.T) {
	authority := knowledgesdk.Authority{Known: true, RuntimeID: "runtime-a", WorkspaceID: "workspace-a", UserID: "reader-a", RoleKey: "support"}
	var operations []string
	var operationsMu sync.Mutex
	var revoked, wrongScope atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get(saashost.RuntimeIDHeader) != "runtime-a" {
			t.Error("service credential or Runtime audience was not retained")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Path == knowledgecontract.SaaSDiscoveryPath {
			_ = json.NewEncoder(w).Encode(knowledgecontract.Descriptor{ProtocolVersion: knowledgecontract.SaaSProtocolVersionV1, Mode: knowledgecontract.DeploymentModeSaaS, Audience: "runtime-a", Capabilities: knowledgecontract.SaaSCapabilitiesV1})
			return
		}
		var envelope knowledgecontract.SaaSRequest
		if r.Method != http.MethodPost || r.URL.Path != knowledgecontract.SaaSInvokePath || json.NewDecoder(r.Body).Decode(&envelope) != nil {
			t.Error("unexpected retrieval request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		operationsMu.Lock()
		operations = append(operations, envelope.Operation)
		operationsMu.Unlock()
		var result knowledgesdk.Result
		switch envelope.Operation {
		case knowledgecontract.SaaSLibrarySearchOperation:
			var input knowledgecontract.SaaSLibrarySearchRequest
			if json.Unmarshal(envelope.Input, &input) != nil || input.LibraryID != "library-a" || input.Query != "PDF export" || !reflect.DeepEqual(input.Authority, authority) {
				t.Errorf("search lost scope: %+v", input)
			}
			result = knowledgesdk.Result{LibraryID: input.LibraryID, Operation: "search", Query: input.Query, ScopeSHA256: "scope-a"}
		case knowledgecontract.SaaSLibraryReadOperation:
			var input knowledgecontract.SaaSLibraryReadRequest
			if json.Unmarshal(envelope.Input, &input) != nil || input.LibraryID != "library-a" || input.DocumentID != "document-a" || !reflect.DeepEqual(input.Authority, authority) {
				t.Errorf("read lost scope: %+v", input)
			}
			result = knowledgesdk.Result{LibraryID: input.LibraryID, Operation: "fetch", DocumentID: input.DocumentID, ScopeSHA256: "scope-a"}
		default:
			t.Errorf("library request fell back to %q", envelope.Operation)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if revoked.Load() {
			_ = json.NewEncoder(w).Encode(knowledgecontract.SaaSResponse{Error: &knowledgecontract.SaaSError{Class: "forbidden", Code: "knowledge.library_access_denied"}})
			return
		}
		if wrongScope.Load() {
			result.LibraryID = "another-library"
		}
		raw, _ := json.Marshal(result)
		_ = json.NewEncoder(w).Encode(knowledgecontract.SaaSResponse{Result: raw})
	}))
	defer server.Close()
	binding, err := NewFactory(Config{Endpoint: server.URL, ServiceAccessToken: "test-token"}).OpenSaaS(t.Context(), knowledgecontract.ApplicationRef{RuntimeID: "runtime-a"})
	if err != nil {
		t.Fatal(err)
	}
	defer binding.Close(t.Context())
	source, err := binding.Runtime().Prepare("runtime-a", knowledgecontract.Options{})
	if err != nil {
		t.Fatal(err)
	}
	library, ok := source.(knowledgecontract.LibraryKnowledge)
	if !ok {
		t.Fatal("remote retrieval does not publish the scoped library port")
	}
	search, err := library.SearchLibraryKnowledge(t.Context(), "library-a", "PDF export", authority)
	if err != nil || search.LibraryID != "library-a" || search.Operation != "search" {
		t.Fatalf("search=%+v error=%v", search, err)
	}
	read, err := library.ReadLibraryKnowledge(t.Context(), "library-a", "document-a", authority)
	if err != nil || read.DocumentID != "document-a" || read.Operation != "fetch" {
		t.Fatalf("read=%+v error=%v", read, err)
	}
	revoked.Store(true)
	_, err = library.ReadLibraryKnowledge(t.Context(), "library-a", "document-a", authority)
	var denied *knowledgesdk.Error
	if !errors.As(err, &denied) || denied.Class != "forbidden" || denied.Code != "knowledge.library_access_denied" {
		t.Fatalf("revocation was lost: %v", err)
	}
	revoked.Store(false)
	wrongScope.Store(true)
	invalid, err := library.ReadLibraryKnowledge(t.Context(), "library-a", "document-a", authority)
	var malformed *knowledgesdk.Error
	if !errors.As(err, &malformed) || malformed.Code != "knowledge.library_result_invalid" || invalid.LibraryID != "" {
		t.Fatalf("another library's result was exposed: %+v error=%v", invalid, err)
	}
	operationsMu.Lock()
	defer operationsMu.Unlock()
	if !reflect.DeepEqual(operations, []string{knowledgecontract.SaaSLibrarySearchOperation, knowledgecontract.SaaSLibraryReadOperation, knowledgecontract.SaaSLibraryReadOperation, knowledgecontract.SaaSLibraryReadOperation}) {
		t.Fatal("unexpected fallback or duplicate requests", operations)
	}
}
