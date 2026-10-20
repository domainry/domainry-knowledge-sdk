// Package modulehost defines the infrastructure required to embed Knowledge.
package modulehost

import (
	"context"
	"database/sql"

	knowledgecontract "github.com/domainry/domainry-knowledge-sdk/contract"
	ormdriver "github.com/domainry/domainry-orm/driver"
	ormmigration "github.com/domainry/domainry-orm/migration"
	"github.com/domainry/domainry-orm/sqlhost"
)

type DB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Executor = sqlhost.Executor
type Queryer = sqlhost.Queryer
type Database = sqlhost.Database

type Dialect interface {
	Identifier(string) string
	Table(string) string
	Placeholder(int) string
}

type SchemaMigration = ormmigration.Migration
type SchemaBaseline = ormmigration.Baseline
type SchemaTable = ormmigration.Table
type SchemaColumn = ormmigration.Column
type SchemaIndex = ormmigration.Index

type MigrationRegistrar interface {
	Driver() string
	Schema() string
	ApplyOwnedMigrations(context.Context, string, []SchemaMigration) error
}

type Host interface {
	Database() Database
	Dialect() Dialect
	Migrations() MigrationRegistrar
	Profile() ormdriver.Profile
}

type ModuleBinding interface {
	Runtime() knowledgecontract.Runtime
	Close(context.Context) error
}

type Factory interface {
	OpenModule(context.Context, knowledgecontract.ApplicationRef, Host) (ModuleBinding, error)
}
