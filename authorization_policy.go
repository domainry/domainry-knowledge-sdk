package knowledgesdk

import "strings"

// RequiresIdentityActionGrant reports whether a Knowledge operation needs an
// explicit product-wide Identity grant in addition to an active principal.
//
// Library membership and page access are the authority for resource-bound
// operations. Personal libraries are owned by their creator. Creating a
// shared library is currently the only operation that has no existing
// resource on which Knowledge can make that decision.
func RequiresIdentityActionGrant(operation string, library KnowledgeLibrary) bool {
	return strings.TrimSpace(operation) == ActionLibrariesCreate && strings.TrimSpace(library.Kind) == "shared"
}
