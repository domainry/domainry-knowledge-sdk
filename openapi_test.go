package knowledgesdk

import "testing"

func TestKnowledgeOpenAPIDocumentProjectsOnlyAuthorizedIntegrationActions(t *testing.T) {
	document, err := KnowledgeOpenAPIDocument()
	if err != nil {
		t.Fatal(err)
	}
	if document["openapi"] != "3.1.0" {
		t.Fatalf("openapi=%v", document["openapi"])
	}
	paths, ok := document["paths"].(map[string]any)
	if !ok || len(paths) == 0 {
		t.Fatalf("paths=%T %#v", document["paths"], document["paths"])
	}
	for key := range openAPIActionSet() {
		action, found, actionErr := AuthorizationAction(key)
		if actionErr != nil || !found || action.HTTP == nil {
			t.Fatalf("resolve action %q: found=%t err=%v action=%+v", key, found, actionErr, action)
		}
		pathItem, ok := paths[action.HTTP.RouteTemplate].(map[string]any)
		if !ok {
			t.Fatalf("path %q missing for action %q", action.HTTP.RouteTemplate, key)
		}
		operation, ok := pathItem[lowerASCII(action.HTTP.Method)].(map[string]any)
		if !ok || operation["x-domainry-action"] != key {
			t.Fatalf("operation mismatch for %q: %#v", key, pathItem)
		}
		security, ok := operation["security"].([]any)
		if !ok || len(security) != 2 {
			t.Fatalf("security missing for %q: %#v", key, operation["security"])
		}
	}
	if _, found := paths["/knowledge/libraries/{libraryID}/members"]; found {
		t.Fatal("member administration must not leak into the integration OpenAPI surface")
	}
	for _, path := range []string{
		"/knowledge/libraries/{libraryID}/pages/{pageID}/agent-references/resolve",
		"/knowledge/libraries/{libraryID}/pages/{pageID}/ai-executions",
		"/knowledge/libraries/{libraryID}/pages/{pageID}/ai-executions/{executionID}",
		"/knowledge/libraries/{libraryID}/pages/{pageID}/ai-executions/{executionID}/confirm",
		"/knowledge/libraries/{libraryID}/pages/{pageID}/application-cards/resolve",
	} {
		if _, found := paths[path]; !found {
			t.Fatalf("integration path %q missing", path)
		}
	}
}

func lowerASCII(value string) string {
	bytes := []byte(value)
	for index, current := range bytes {
		if current >= 'A' && current <= 'Z' {
			bytes[index] = current + ('a' - 'A')
		}
	}
	return string(bytes)
}
