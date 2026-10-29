package knowledgesdk

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestKnowledgeReplaceSuggestionContractCarriesConflictAndAgentProvenance(t *testing.T) {
	input := KnowledgePageSuggestionCreate{
		ClientID:            "knowledge-propose:1",
		Source:              KnowledgeSuggestionSourceVerdentAgent,
		ToolCallID:          "tool-call-1",
		Kind:                KnowledgeSuggestionReplace,
		BaseDraftRevision:   7,
		Anchor:              KnowledgePageCommentAnchor{NodeID: "block-1", Quote: "before"},
		ExpectedBlockSHA256: strings.Repeat("a", 64),
		ProposedText:        "after",
	}
	raw, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{
		`"source":"verdent_agent"`,
		`"tool_call_id":"tool-call-1"`,
		`"kind":"replace"`,
		`"base_draft_revision":7`,
		`"node_id":"block-1"`,
		`"expected_block_sha256":"` + strings.Repeat("a", 64) + `"`,
	} {
		if !strings.Contains(string(raw), fragment) {
			t.Fatalf("replace suggestion contract lost %s: %s", fragment, raw)
		}
	}
}

func TestKnowledgeSuggestionDecisionCarriesIndependentAgentInvocation(t *testing.T) {
	input := KnowledgePageSuggestionDecision{
		ClientID:              "knowledge-apply:1",
		Source:                KnowledgeSuggestionSourceVerdentAgent,
		ToolCallID:            "tool-call-2",
		ExpectedRevision:      2,
		ExpectedDraftRevision: 7,
		Decision:              KnowledgeSuggestionDecisionAccept,
	}
	raw, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"source":"verdent_agent"`) || !strings.Contains(string(raw), `"tool_call_id":"tool-call-2"`) {
		t.Fatalf("decision contract lost agent invocation: %s", raw)
	}
}
