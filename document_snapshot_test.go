package knowledgesdk

import (
	"strings"
	"testing"
)

func TestKnowledgeDocumentSnapshotRequiresExactVersionDigestAndPairedPageIdentity(t *testing.T) {
	valid := KnowledgeDocumentSnapshot{DocumentID: "document", ContentVersion: 2, IndexGeneration: 3, SHA256: strings.Repeat("a", 64)}
	if !valid.Valid() {
		t.Fatal("valid immutable snapshot rejected")
	}
	for name, change := range map[string]func(*KnowledgeDocumentSnapshot){
		"missing document":        func(s *KnowledgeDocumentSnapshot) { s.DocumentID = "" },
		"untrimmed document":      func(s *KnowledgeDocumentSnapshot) { s.DocumentID = " document" },
		"missing content version": func(s *KnowledgeDocumentSnapshot) { s.ContentVersion = 0 },
		"missing generation":      func(s *KnowledgeDocumentSnapshot) { s.IndexGeneration = 0 },
		"unknown hash":            func(s *KnowledgeDocumentSnapshot) { s.SHA256 = "" },
		"invalid hash":            func(s *KnowledgeDocumentSnapshot) { s.SHA256 = strings.Repeat("g", 64) },
		"noncanonical hash":       func(s *KnowledgeDocumentSnapshot) { s.SHA256 = strings.Repeat("A", 64) },
		"page without version":    func(s *KnowledgeDocumentSnapshot) { s.PageID = "page" },
		"version without page":    func(s *KnowledgeDocumentSnapshot) { s.PageVersion = 1 },
	} {
		t.Run(name, func(t *testing.T) {
			s := valid
			change(&s)
			if s.Valid() {
				t.Fatal("invalid snapshot accepted")
			}
		})
	}
	valid.PageID, valid.PageVersion = "page", 4
	if !valid.Valid() {
		t.Fatal("valid page snapshot rejected")
	}
}
