package knowledgesdk

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// EmbedProvider is a trusted, code-registered external-card boundary. Page
// content can select a provider key and submit a URL, but it can never submit
// CSP origins or executable embed policy.
type EmbedProvider interface {
	Key() string
	FrameOrigins() []string
	CanonicalizeEmbedURL(string) (string, error)
}

type EmbedProviderRegistry struct {
	providers map[string]EmbedProvider
	origins   []string
}

func NewEmbedProviderRegistry(providers ...EmbedProvider) (*EmbedProviderRegistry, error) {
	registry := &EmbedProviderRegistry{providers: make(map[string]EmbedProvider, len(providers))}
	origins := map[string]struct{}{}
	for _, provider := range providers {
		if provider == nil {
			return nil, fmt.Errorf("Knowledge EmbedProvider is nil")
		}
		key := strings.TrimSpace(provider.Key())
		if key == "" || len(key) > 64 || !validEmbedKey(key) || registry.providers[key] != nil {
			return nil, fmt.Errorf("Knowledge EmbedProvider key %q is invalid or duplicated", key)
		}
		providerOrigins := provider.FrameOrigins()
		if len(providerOrigins) == 0 {
			return nil, fmt.Errorf("Knowledge EmbedProvider %q has no frame origin", key)
		}
		for _, raw := range providerOrigins {
			parsed, err := url.Parse(strings.TrimSpace(raw))
			if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
				return nil, fmt.Errorf("Knowledge EmbedProvider %q frame origin %q is invalid", key, raw)
			}
			if _, seen := origins[parsed.String()]; !seen {
				origins[parsed.String()] = struct{}{}
				registry.origins = append(registry.origins, parsed.String())
			}
		}
		registry.providers[key] = provider
	}
	return registry, nil
}

func (registry *EmbedProviderRegistry) Canonicalize(provider, rawURL string) (string, error) {
	if registry == nil {
		return "", fmt.Errorf("Knowledge EmbedProvider registry is unavailable")
	}
	value := registry.providers[strings.TrimSpace(provider)]
	if value == nil {
		return "", fmt.Errorf("Knowledge EmbedProvider %q is not registered", provider)
	}
	canonical, err := value.CanonicalizeEmbedURL(strings.TrimSpace(rawURL))
	if err != nil || strings.TrimSpace(canonical) == "" {
		return "", fmt.Errorf("Knowledge EmbedProvider %q rejected URL", provider)
	}
	return canonical, nil
}

func (registry *EmbedProviderRegistry) FrameOrigins() []string {
	if registry == nil {
		return nil
	}
	return append([]string(nil), registry.origins...)
}

func (registry *EmbedProviderRegistry) Keys() []string {
	if registry == nil {
		return nil
	}
	keys := make([]string, 0, len(registry.providers))
	for key := range registry.providers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func validEmbedKey(value string) bool {
	for _, character := range value {
		if !(character >= 'a' && character <= 'z' || character >= '0' && character <= '9' || character == '_' || character == '-') {
			return false
		}
	}
	return true
}
