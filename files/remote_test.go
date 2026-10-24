package files

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
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
		case "/knowledge/v1/files/upload":
			reader, err := request.MultipartReader()
			if err != nil {
				t.Fatal(err)
			}
			metadataPart, err := reader.NextPart()
			if err != nil || metadataPart.FormName() != "metadata" {
				t.Fatalf("metadata part err=%v", err)
			}
			var metadata struct {
				Authority Authority    `json:"authority"`
				Input     StreamUpload `json:"input"`
			}
			if json.NewDecoder(metadataPart).Decode(&metadata) != nil || metadata.Input.Content != nil {
				t.Fatalf("metadata=%+v", metadata)
			}
			contentPart, err := reader.NextPart()
			content, readErr := io.ReadAll(contentPart)
			want := []byte("note")
			if metadata.Input.ClientID == "upload-b" {
				want = []byte("streamed")
			}
			if err != nil || readErr != nil || contentPart.FormName() != "content" || !bytes.Equal(content, want) || metadata.Input.SizeBytes != int64(len(want)) {
				t.Fatalf("content=%q err=%v readErr=%v", content, err, readErr)
			}
			value := File{ID: "file_0123456789abcdef0123456789abcdef", WorkspaceID: "workspace-a", CreatedBy: "user-a", Filename: metadata.Input.Filename, ContentType: "text/plain", SizeBytes: metadata.Input.SizeBytes}
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
	raw := []byte("streamed")
	digest := sha256.Sum256(raw)
	file, err = client.UploadStream(t.Context(), Authority{RuntimeID: "runtime-a", WorkspaceID: "workspace-a", SubjectID: "user-a"}, StreamUpload{ClientID: "upload-b", Filename: "stream.txt", ContentType: "text/plain", SizeBytes: int64(len(raw)), SHA256: hex.EncodeToString(digest[:]), Content: bytes.NewReader(raw)})
	if err != nil || file.ID == "" {
		t.Fatalf("stream file=%+v err=%v", file, err)
	}
}
