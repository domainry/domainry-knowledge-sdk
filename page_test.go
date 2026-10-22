package knowledgesdk

import (
	"encoding/json"
	"testing"
)

func TestKnowledgeFilePageContentContract(t *testing.T) {
	if KnowledgePageTypeFile != "file" {
		t.Fatalf("unstable file page type id: %q", KnowledgePageTypeFile)
	}
	want := KnowledgeFilePageContent{DocumentID: "kdoc_0123456789abcdef0123456789abcdef", DisplayMode: "preview", Title: "产品说明.pdf"}
	raw, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got KnowledgeFilePageContent
	if err = json.Unmarshal(raw, &got); err != nil || got != want {
		t.Fatalf("file page content contract did not round-trip: got=%+v err=%v", got, err)
	}
}

func TestKnowledgePageTypeContractsAreStableAndComplete(t *testing.T) {
	if KnowledgePageCapabilityCreate != 1 || KnowledgePageCapabilityEdit != 2 || KnowledgePageCapabilityRead != 4 || KnowledgePageCapabilityPreview != 8 || KnowledgePageCapabilityImport != 16 || KnowledgePageCapabilityRealtime != 32 || KnowledgePageCapabilitySearch != 64 || KnowledgePageCapabilityComment != 128 || KnowledgePageCapabilityVersion != 256 {
		t.Fatal("knowledge page capability bits changed")
	}
	if err := ValidateKnowledgePageTypeContracts(); err != nil {
		t.Fatal(err)
	}
	contracts := KnowledgePageTypeContracts()
	if len(contracts) != 10 {
		t.Fatalf("unexpected page type contract count: %d", len(contracts))
	}
	contracts[0].ID = "mutated"
	document, found := KnowledgePageTypeContractFor(KnowledgePageTypeDocument)
	if !found || document.ID != KnowledgePageTypeDocument || document.RootField != "blocks" || document.RootKind != KnowledgePageContentRootArray {
		t.Fatalf("document contract is incomplete: %+v found=%v", document, found)
	}
	if document.Capabilities&KnowledgePageCapabilityRealtime == 0 || document.Capabilities&KnowledgePageCapabilityEdit == 0 {
		t.Fatalf("document capabilities are incomplete: %b", document.Capabilities)
	}
	file, found := KnowledgePageTypeContractFor(KnowledgePageTypeFile)
	if !found || file.RootField != "document_id" || file.Capabilities&KnowledgePageCapabilityImport == 0 || file.Capabilities&KnowledgePageCapabilityCreate != 0 {
		t.Fatalf("file contract is incomplete: %+v found=%v", file, found)
	}
}
