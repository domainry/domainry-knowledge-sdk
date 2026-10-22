package knowledgesdk

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	KnowledgeApplicationCardCalendar = "calendar_event"
	KnowledgeApplicationCardMeeting  = "meeting"
	KnowledgeApplicationCardGroup    = "group"
	KnowledgeApplicationCardTask     = "task"
	KnowledgeApplicationCardJira     = "jira"
	KnowledgeApplicationCardWeb      = "web"
)

var knowledgeApplicationCardKinds = map[string]struct{}{
	KnowledgeApplicationCardCalendar: {}, KnowledgeApplicationCardMeeting: {}, KnowledgeApplicationCardGroup: {},
	KnowledgeApplicationCardTask: {}, KnowledgeApplicationCardJira: {}, KnowledgeApplicationCardWeb: {},
}

type KnowledgeApplicationCardResolve struct {
	Provider   string `json:"provider"`
	Kind       string `json:"kind"`
	ResourceID string `json:"resource_id"`
}

type KnowledgeApplicationCard struct {
	Provider   string    `json:"provider"`
	Kind       string    `json:"kind"`
	ResourceID string    `json:"resource_id"`
	Title      string    `json:"title"`
	Subtitle   string    `json:"subtitle,omitempty"`
	Status     string    `json:"status,omitempty"`
	URL        string    `json:"url,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type KnowledgeApplicationCardProviderInput struct {
	LibraryID  string
	PageID     string
	Kind       string
	ResourceID string
	Authority  Authority
}

// KnowledgeApplicationCardProvider resolves one current, user-authorized
// third-party resource. Implementations own credentials and provider ACLs;
// Knowledge never stores provider payloads or credentials in page content.
type KnowledgeApplicationCardProvider interface {
	Key() string
	Kinds() []string
	ResolveKnowledgeApplicationCard(context.Context, KnowledgeApplicationCardProviderInput) (KnowledgeApplicationCard, error)
}

type KnowledgeApplicationCardRegistry struct {
	providers map[string]KnowledgeApplicationCardProvider
	kinds     map[string][]string
}

func NewKnowledgeApplicationCardRegistry(providers ...KnowledgeApplicationCardProvider) (*KnowledgeApplicationCardRegistry, error) {
	registry := &KnowledgeApplicationCardRegistry{providers: make(map[string]KnowledgeApplicationCardProvider, len(providers)), kinds: make(map[string][]string, len(providers))}
	for _, provider := range providers {
		if provider == nil {
			return nil, fmt.Errorf("Knowledge application card provider is nil")
		}
		key := strings.TrimSpace(provider.Key())
		if !validEmbedKey(key) || registry.providers[key] != nil {
			return nil, fmt.Errorf("Knowledge application card provider key %q is invalid or duplicated", key)
		}
		kinds := append([]string(nil), provider.Kinds()...)
		if len(kinds) == 0 || len(kinds) > len(knowledgeApplicationCardKinds) {
			return nil, fmt.Errorf("Knowledge application card provider %q has invalid kinds", key)
		}
		sort.Strings(kinds)
		for index, kind := range kinds {
			if _, valid := knowledgeApplicationCardKinds[kind]; !valid || index > 0 && kinds[index-1] == kind {
				return nil, fmt.Errorf("Knowledge application card provider %q has invalid kind %q", key, kind)
			}
		}
		registry.providers[key], registry.kinds[key] = provider, kinds
	}
	return registry, nil
}

func (registry *KnowledgeApplicationCardRegistry) Resolve(ctx context.Context, input KnowledgeApplicationCardResolve, libraryID, pageID string, authority Authority) (KnowledgeApplicationCard, error) {
	input.Provider, input.Kind, input.ResourceID = strings.TrimSpace(input.Provider), strings.TrimSpace(input.Kind), strings.TrimSpace(input.ResourceID)
	provider := registry.providers[input.Provider]
	kinds := registry.kinds[input.Provider]
	index := sort.SearchStrings(kinds, input.Kind)
	if provider == nil || index >= len(kinds) || kinds[index] != input.Kind || input.ResourceID == "" || len(input.ResourceID) > 2048 || !utf8.ValidString(input.ResourceID) || strings.ContainsRune(input.ResourceID, 0) {
		return KnowledgeApplicationCard{}, fmt.Errorf("Knowledge application card reference is invalid")
	}
	card, err := provider.ResolveKnowledgeApplicationCard(ctx, KnowledgeApplicationCardProviderInput{LibraryID: libraryID, PageID: pageID, Kind: input.Kind, ResourceID: input.ResourceID, Authority: authority})
	if err != nil {
		return KnowledgeApplicationCard{}, err
	}
	card.Provider, card.Kind, card.ResourceID = input.Provider, input.Kind, input.ResourceID
	if !boundedApplicationCardText(card.Title, 256, true) || !boundedApplicationCardText(card.Subtitle, 512, false) || !boundedApplicationCardText(card.Status, 64, false) {
		return KnowledgeApplicationCard{}, fmt.Errorf("Knowledge application card result is invalid")
	}
	if card.URL != "" {
		parsed, parseErr := url.Parse(card.URL)
		if parseErr != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil {
			return KnowledgeApplicationCard{}, fmt.Errorf("Knowledge application card URL is invalid")
		}
	}
	return card, nil
}

func (registry *KnowledgeApplicationCardRegistry) Providers() map[string][]string {
	result := make(map[string][]string, len(registry.kinds))
	for key, kinds := range registry.kinds {
		result[key] = append([]string(nil), kinds...)
	}
	return result
}

func boundedApplicationCardText(value string, maximum int, required bool) bool {
	return utf8.ValidString(value) && !strings.ContainsRune(value, 0) && len(value) <= maximum && (!required || strings.TrimSpace(value) != "")
}

type KnowledgeApplicationCardService interface {
	ResolveKnowledgeApplicationCard(context.Context, string, string, KnowledgeApplicationCardResolve, Authority) (KnowledgeApplicationCard, error)
}
