package knowledgesdk

import (
	"testing"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

func TestAuthorizationActionsExposeLocalizedPermissionPresentation(t *testing.T) {
	definitions, err := AuthorizationActions()
	if err != nil {
		t.Fatal(err)
	}
	for _, definition := range definitions {
		permission := definition.Permission
		if permission == nil {
			continue
		}
		if permission.ModuleKey != "knowledge" || permission.GroupKey == "" || permission.ScopeKind == "" || permission.DefaultLocale != "en" {
			t.Fatalf("permission %q has incomplete presentation contract: %+v", permission.Key, permission)
		}
		wantGrantMode := actioncontract.PermissionGrantBundled
		if permission.Key == ActionLibrariesCreate {
			wantGrantMode = actioncontract.PermissionGrantDirect
		}
		if permission.GrantMode != wantGrantMode {
			t.Fatalf("permission %q grant mode=%q want=%q", permission.Key, permission.GrantMode, wantGrantMode)
		}
		for _, locale := range []string{"en", "zh-CN"} {
			text := permission.I18n[locale]
			if text.ModuleLabel == "" || text.GroupLabel == "" || text.ResourceLabel == "" || text.OperationLabel == "" || text.Label == "" || text.Description == "" {
				t.Fatalf("permission %q has incomplete %s presentation: %+v", permission.Key, locale, text)
			}
			if locale == "zh-CN" && (text.GroupLabel == "知识库功能" || text.ResourceLabel == "知识库功能") {
				t.Fatalf("permission %q uses fallback Chinese presentation: %+v", permission.Key, text)
			}
		}
	}
}
