package knowledgesdk

import "testing"

func TestRequiresIdentityActionGrantOnlyForSharedLibraryCreation(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		library   KnowledgeLibrary
		want      bool
	}{
		{name: "shared library creation", operation: ActionLibrariesCreate, library: KnowledgeLibrary{Kind: "shared"}, want: true},
		{name: "personal library creation", operation: ActionLibrariesCreate, library: KnowledgeLibrary{Kind: "personal"}},
		{name: "library member edit", operation: ActionPagesUpdate, library: KnowledgeLibrary{ID: "library-a", Kind: "shared", Role: "editor"}},
		{name: "membership filtered list", operation: ActionLibrariesList},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := RequiresIdentityActionGrant(test.operation, test.library); got != test.want {
				t.Fatalf("RequiresIdentityActionGrant(%q, %+v)=%t want=%t", test.operation, test.library, got, test.want)
			}
		})
	}
}
