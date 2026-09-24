// Package saashost defines the remote Knowledge deployment boundary. It is
// separate from modulehost so a network client cannot claim participation in a
// caller-owned SQL transaction or receive a caller-owned Store.
package saashost

import (
	"context"

	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	lifecyclecontract "github.com/domainry/domainry-lifecycle-sdk/contract"
)

const RuntimeIDHeader = "X-Domainry-Runtime-ID"

type Binding interface {
	Descriptor() knowledgecontract.Descriptor
	Runtime() knowledgecontract.Runtime
	ArtifactMutations() knowledgecontract.ArtifactMutationService
	SubjectLifecycle() lifecyclecontract.SubjectExecutionHandler
	Close(context.Context) error
}

type Factory interface {
	OpenSaaS(context.Context, knowledgecontract.ApplicationRef) (Binding, error)
}
