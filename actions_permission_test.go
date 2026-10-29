package knowledgesdk

import (
	"reflect"
	"sort"
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

func TestSharedLibraryCreationPermissionExplainsItsGlobalScope(t *testing.T) {
	action, found, err := AuthorizationAction(ActionLibrariesCreate)
	if err != nil || !found || action.Permission == nil {
		t.Fatalf("shared library creation permission unavailable: found=%t err=%v", found, err)
	}
	english := action.Permission.I18n["en"]
	chinese := action.Permission.I18n["zh-CN"]
	if english.Label != "Shared knowledge libraries · Create" || english.Description != "Allows creating shared knowledge libraries. Personal libraries do not require this permission." {
		t.Fatalf("unexpected English shared-library presentation: %+v", english)
	}
	if chinese.Label != "共享知识库 · 创建" || chinese.Description != "允许创建共享知识库；个人知识库不需要此权限。" {
		t.Fatalf("unexpected Chinese shared-library presentation: %+v", chinese)
	}
	wantLocales := []string{"de", "en", "es", "fr", "it", "ja", "ko", "pt", "tr", "zh-CN", "zh-TW"}
	gotLocales := make([]string, 0, len(action.Permission.I18n))
	for locale, text := range action.Permission.I18n {
		gotLocales = append(gotLocales, locale)
		if text.ModuleLabel == "" || text.GroupLabel == "" || text.ResourceLabel == "" || text.OperationLabel == "" || text.Label == "" || text.Description == "" {
			t.Fatalf("shared-library presentation is incomplete for %s: %+v", locale, text)
		}
	}
	sort.Strings(gotLocales)
	if !reflect.DeepEqual(gotLocales, wantLocales) {
		t.Fatalf("shared-library locale coverage=%v want=%v", gotLocales, wantLocales)
	}
}
