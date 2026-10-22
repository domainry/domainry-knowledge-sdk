package knowledgesdk

import (
	"fmt"
	"sort"
	"strings"
)

// TrustedApplicationRegistration binds one Identity application to a current
// service-account user and an explicit subset of Knowledge Action keys.
// Credentials remain exclusively in Identity.
type TrustedApplicationRegistration struct {
	ApplicationKey string   `json:"application_key"`
	ServiceUserID  string   `json:"service_user_id"`
	AllowedActions []string `json:"allowed_actions"`
}

type TrustedApplicationRegistry struct {
	registrations map[string]TrustedApplicationRegistration
}

func NewTrustedApplicationRegistry(values ...TrustedApplicationRegistration) (*TrustedApplicationRegistry, error) {
	known := map[string]struct{}{}
	actions, err := AuthorizationActions()
	if err != nil {
		return nil, err
	}
	for _, action := range actions {
		if action.HTTP != nil {
			known[action.Key] = struct{}{}
		}
	}
	registry := &TrustedApplicationRegistry{registrations: make(map[string]TrustedApplicationRegistration, len(values))}
	for _, value := range values {
		value.ApplicationKey, value.ServiceUserID = strings.TrimSpace(value.ApplicationKey), strings.TrimSpace(value.ServiceUserID)
		if value.ApplicationKey == "" || value.ServiceUserID == "" || len(value.ApplicationKey) > 255 || len(value.ServiceUserID) > 255 || registry.registrations[value.ApplicationKey].ApplicationKey != "" || len(value.AllowedActions) == 0 {
			return nil, fmt.Errorf("Knowledge trusted application registration is invalid")
		}
		seen := map[string]struct{}{}
		value.AllowedActions = append([]string(nil), value.AllowedActions...)
		for index := range value.AllowedActions {
			value.AllowedActions[index] = strings.TrimSpace(value.AllowedActions[index])
			if _, ok := known[value.AllowedActions[index]]; !ok {
				return nil, fmt.Errorf("Knowledge trusted application %q uses unknown Action %q", value.ApplicationKey, value.AllowedActions[index])
			}
			if _, duplicate := seen[value.AllowedActions[index]]; duplicate {
				return nil, fmt.Errorf("Knowledge trusted application %q duplicates Action %q", value.ApplicationKey, value.AllowedActions[index])
			}
			seen[value.AllowedActions[index]] = struct{}{}
		}
		sort.Strings(value.AllowedActions)
		registry.registrations[value.ApplicationKey] = value
	}
	return registry, nil
}

func (registry *TrustedApplicationRegistry) Authorize(applicationKey, action string) (TrustedApplicationRegistration, bool) {
	if registry == nil {
		return TrustedApplicationRegistration{}, false
	}
	registration, ok := registry.registrations[strings.TrimSpace(applicationKey)]
	if !ok {
		return TrustedApplicationRegistration{}, false
	}
	index := sort.SearchStrings(registration.AllowedActions, strings.TrimSpace(action))
	return registration, index < len(registration.AllowedActions) && registration.AllowedActions[index] == strings.TrimSpace(action)
}

func (registry *TrustedApplicationRegistry) Registrations() []TrustedApplicationRegistration {
	if registry == nil {
		return nil
	}
	keys := make([]string, 0, len(registry.registrations))
	for key := range registry.registrations {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := make([]TrustedApplicationRegistration, 0, len(keys))
	for _, key := range keys {
		value := registry.registrations[key]
		value.AllowedActions = append([]string(nil), value.AllowedActions...)
		values = append(values, value)
	}
	return values
}
