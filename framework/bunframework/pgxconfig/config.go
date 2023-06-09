package pgxconfig

import (
	"github.com/jackc/pgx/v4"
)

// Driver for usage with jackc/pgx.
// See https://bun.uptrace.dev/postgres/#pgx.
type Driver struct {
	// DSN is a string used to connect to a psql instance.
	DSN string `json:"dsn" yaml:"dsn"`

	// KeepImplicitPreparedStatements sets the PreferSimpleProtocol option to false on pgx config.
	// It is not the default with bun since it does not benefit from using implicit prepared statements.
	KeepImplicitPreparedStatements bool `json:"keepImplicitPreparedStatements" yaml:"keepImplicitPreparedStatements"`

	Logger              pgx.Logger                  `json:"-" yaml:"-"`
	LogLevel            pgx.LogLevel                `json:"logLevel" yaml:"logLevel"`
	BuildStatementCache pgx.BuildStatementCacheFunc `json:"-" yaml:"-"`
}

// NewDriverWithDSN generates a quick configuration object with a DSN.
//
// For a more granular configuration, use the Driver object directly.
func NewDriverWithDSN(dsn string) Driver {
	return Driver{DSN: dsn}
}
