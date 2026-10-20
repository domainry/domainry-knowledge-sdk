package persistence

import knowledgesdk "github.com/domainry/domainry-knowledge-sdk"

// This optional repository enforces workspace isolation and current library
// membership. Each mutation rechecks membership and expected library revision
// inside the same serializable transaction. Global Identity remains the host's
// responsibility; callers must not expose this repository directly over HTTP.
type KnowledgeLibraryRepository interface {
	knowledgesdk.KnowledgeLibraryService
}
