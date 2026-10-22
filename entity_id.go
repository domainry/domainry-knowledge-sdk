package knowledgesdk

import "github.com/domainry/domainry-foundation/entityid"

// NewKnowledgeEntityID generates a canonical public ID for a Knowledge entity.
func NewKnowledgeEntityID(prefix string) (string, error) {
	return entityid.New(prefix)
}

// DeriveKnowledgeEntityID deterministically generates a canonical public ID
// for idempotent Knowledge entity creation.
func DeriveKnowledgeEntityID(prefix, namespace string, parts ...string) (string, error) {
	return entityid.Derive(prefix, namespace, parts...)
}

// ValidateKnowledgeEntityID validates a Knowledge entity ID against its
// registered short prefix.
func ValidateKnowledgeEntityID(prefix, value string) error {
	return entityid.Validate(prefix, value)
}
