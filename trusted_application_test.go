package knowledgesdk

import "testing"

func TestTrustedApplicationRegistryRequiresExactKnownActions(t *testing.T) {
	registry, err := NewTrustedApplicationRegistry(TrustedApplicationRegistration{
		ApplicationKey: "calendar-app", ServiceUserID: "calendar-bot", AllowedActions: []string{ActionPagesGet, ActionPageTasksCreate},
	})
	if err != nil {
		t.Fatal(err)
	}
	registration, allowed := registry.Authorize("calendar-app", ActionPagesGet)
	if !allowed || registration.ServiceUserID != "calendar-bot" {
		t.Fatalf("registration=%+v allowed=%t", registration, allowed)
	}
	if _, allowed = registry.Authorize("calendar-app", ActionPagesUpdate); allowed {
		t.Fatal("undeclared Action was allowed")
	}
	if _, err = NewTrustedApplicationRegistry(TrustedApplicationRegistration{ApplicationKey: "bad", ServiceUserID: "bot", AllowedActions: []string{"knowledge.unknown"}}); err == nil {
		t.Fatal("unknown Action must be rejected")
	}
}
