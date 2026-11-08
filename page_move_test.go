package knowledgesdk

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestKnowledgePageMoveCarriesCrossLibraryBoundary(t *testing.T) {
	input := KnowledgePageMove{
		ClientID:          "move-request",
		ExpectedRevision:  7,
		TargetLibraryID:   "klib_target",
		ParentPageID:      "kpage_parent",
		Rank:              "0001",
		AccessImpactToken: "impact-token",
	}
	raw, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	encoded := string(raw)
	for _, expected := range []string{`"target_library_id":"klib_target"`, `"parent_page_id":"kpage_parent"`, `"access_impact_token":"impact-token"`} {
		if !strings.Contains(encoded, expected) {
			t.Fatalf("cross-library move lost %s in %s", expected, encoded)
		}
	}

	impact := KnowledgePageMoveImpact{PageID: "kpage_source", SourceLibraryID: "klib_source", TargetLibraryID: "klib_target", AffectedPages: 3, RequiresConfirm: true, Token: "impact-token"}
	raw, err = json.Marshal(impact)
	if err != nil {
		t.Fatal(err)
	}
	encoded = string(raw)
	for _, expected := range []string{`"source_library_id":"klib_source"`, `"target_library_id":"klib_target"`, `"affected_pages":3`} {
		if !strings.Contains(encoded, expected) {
			t.Fatalf("move impact lost %s in %s", expected, encoded)
		}
	}
}
