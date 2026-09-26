package files

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemoteFileClientUsesIsolatedAuthenticatedCapability(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if request.Header.Get("Authorization") != "Bearer file-token" || request.Header.Get(runtimeHeader) != "runtime-a" {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}
		switch request.URL.Path {
		case discoveryPath:
			_ = json.NewEncoder(writer).Encode(remoteDescriptor{ProtocolVersion: protocolV1, Mode: deploymentSaaS, Audience: "runtime-a", Capabilities: []string{"files.read", "files.write"}})
		case invokePath:
			var envelope remoteRequest
			if json.NewDecoder(request.Body).Decode(&envelope) != nil || envelope.Operation != "files.upload" {
				t.Fatalf("envelope=%+v", envelope)
			}
			value := File{ID: "file_0123456789abcdef0123456789abcdef", WorkspaceID: "workspace-a", CreatedBy: "user-a", Filename: "note.txt", ContentType: "text/plain", SizeBytes: 4}
			raw, _ := json.Marshal(value)
			_ = json.NewEncoder(writer).Encode(remoteResponse{Result: raw})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := OpenRemote(t.Context(), RemoteConfig{Endpoint: server.URL, ServiceAccessToken: "file-token"}, "runtime-a")
	if err != nil {
		t.Fatal(err)
	}
	file, err := client.Upload(t.Context(), Authority{RuntimeID: "runtime-a", WorkspaceID: "workspace-a", SubjectID: "user-a"}, Upload{ClientID: "upload-a", Filename: "note.txt", ContentType: "text/plain", Data: []byte("note")})
	if err != nil || file.ID == "" || file.WorkspaceID != "workspace-a" {
		t.Fatalf("file=%+v err=%v", file, err)
	}
}
