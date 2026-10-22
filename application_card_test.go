package knowledgesdk

import (
	"context"
	"testing"
)

type applicationCardProviderStub struct{ key string }

func (provider applicationCardProviderStub) Key() string { return provider.key }
func (applicationCardProviderStub) Kinds() []string {
	return []string{KnowledgeApplicationCardCalendar, KnowledgeApplicationCardJira}
}
func (applicationCardProviderStub) ResolveKnowledgeApplicationCard(_ context.Context, input KnowledgeApplicationCardProviderInput) (KnowledgeApplicationCard, error) {
	return KnowledgeApplicationCard{Title: input.ResourceID, Status: "current", URL: "https://apps.example.org/resources/" + input.ResourceID}, nil
}

func TestKnowledgeApplicationCardRegistryOwnsKindsAndSanitizesResults(t *testing.T) {
	registry, err := NewKnowledgeApplicationCardRegistry(applicationCardProviderStub{key: "work-apps"})
	if err != nil {
		t.Fatal(err)
	}
	card, err := registry.Resolve(t.Context(), KnowledgeApplicationCardResolve{Provider: "work-apps", Kind: KnowledgeApplicationCardJira, ResourceID: "PROJ-42"}, "library", "page", Authority{Known: true})
	if err != nil || card.Provider != "work-apps" || card.Kind != KnowledgeApplicationCardJira || card.Title != "PROJ-42" {
		t.Fatalf("card=%+v err=%v", card, err)
	}
	if _, err = registry.Resolve(t.Context(), KnowledgeApplicationCardResolve{Provider: "work-apps", Kind: KnowledgeApplicationCardWeb, ResourceID: "https://example.org"}, "library", "page", Authority{Known: true}); err == nil {
		t.Fatal("unregistered kind must fail closed")
	}
	if _, err = NewKnowledgeApplicationCardRegistry(applicationCardProviderStub{key: "work-apps"}, applicationCardProviderStub{key: "work-apps"}); err == nil {
		t.Fatal("duplicate provider must fail")
	}
}
