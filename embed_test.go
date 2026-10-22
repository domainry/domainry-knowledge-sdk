package knowledgesdk

import "testing"

type embedProviderTestStub struct{ key string }

func (stub embedProviderTestStub) Key() string                                { return stub.key }
func (embedProviderTestStub) FrameOrigins() []string                          { return []string{"https://embed.example.test"} }
func (embedProviderTestStub) CanonicalizeEmbedURL(raw string) (string, error) { return raw, nil }

func TestEmbedProviderRegistryOwnsKeysAndCSPOrigins(t *testing.T) {
	registry, err := NewEmbedProviderRegistry(embedProviderTestStub{key: "example"})
	if err != nil {
		t.Fatal(err)
	}
	if canonical, canonicalErr := registry.Canonicalize("example", "https://embed.example.test/item/1"); canonicalErr != nil || canonical != "https://embed.example.test/item/1" {
		t.Fatalf("canonical=%q err=%v", canonical, canonicalErr)
	}
	if origins := registry.FrameOrigins(); len(origins) != 1 || origins[0] != "https://embed.example.test" {
		t.Fatalf("origins=%v", origins)
	}
	if _, err = NewEmbedProviderRegistry(embedProviderTestStub{key: "example"}, embedProviderTestStub{key: "example"}); err == nil {
		t.Fatal("duplicate provider must be rejected")
	}
}
