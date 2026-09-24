// Package modulehost defines infrastructure and narrow Agent source ports used
// to embed Knowledge without exposing either module's Store.
package modulehost

import (
	"context"
	"database/sql"

	agentsdk "github.com/domainry/domainry-agent-sdk"
	agentmodulehost "github.com/domainry/domainry-agent-sdk/modulehost"
	agentpersistence "github.com/domainry/domainry-agent-sdk/persistence"
	sharedartifact "github.com/domainry/domainry-foundation/artifact"
	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	lifecyclecontract "github.com/domainry/domainry-lifecycle-sdk/contract"
	ormdriver "github.com/domainry/domainry-orm/driver"
)

type DB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type SourceReader interface {
	Conversation(context.Context, DB, string, agentsdk.ConversationAuthority) (agentsdk.Conversation, error)
	Run(context.Context, DB, string, string, agentsdk.ConversationAuthority) (agentsdk.ConversationRun, error)
}

type Host interface {
	Database() agentmodulehost.Database
	Dialect() agentmodulehost.Dialect
	Migrations() agentmodulehost.MigrationRegistrar
	Profile() ormdriver.Profile
	SourceReader
}

type ArtifactHost interface {
	ArtifactStore() sharedartifact.ManagedStore
	ArtifactContentStore() sharedartifact.ContentStore
	ArtifactContentWriter() sharedartifact.ContentWriter
}

// ArtifactTransactions is the narrow cross-module capability needed to keep
// an Agent execution receipt and a Knowledge-owned artifact mutation atomic.
// It exposes no Knowledge Store and no connection lifecycle.
type ArtifactTransactions interface {
	ArtifactRecordInTransaction(context.Context, *sql.Tx, string, int64, agentsdk.ConversationAuthority) (agentpersistence.ConversationArtifactRecord, error)
	SaveArtifactInTransaction(context.Context, *sql.Tx, agentpersistence.ConversationArtifactWrite, agentsdk.ConversationAuthority) (agentpersistence.ConversationArtifactRecord, error)
	SaveArtifactExportInTransaction(context.Context, *sql.Tx, agentpersistence.ConversationArtifactExportWrite, agentsdk.ConversationAuthority) (any, error)
}

type ModuleBinding interface {
	Runtime() knowledgecontract.Runtime
	ArtifactTransactions() ArtifactTransactions
	SubjectLifecycle(knowledgecontract.Options) lifecyclecontract.SubjectExecutionHandler
	Close(context.Context) error
}

type Factory interface {
	OpenModule(context.Context, knowledgecontract.ApplicationRef, Host) (ModuleBinding, error)
}
