package pgxconfig

import (
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// Connect opens a new connection with pgx connector.
// It implements bunframework.DriverConfig interface.
func (config Driver) Connect(options ...bun.DBOption) (*bun.DB, *sql.DB, error) {
	pgxConfig, err := pgx.ParseConfig(config.DSN)
	if err != nil {
		return nil, nil, err
	}

	if !config.KeepImplicitPreparedStatements {
		pgxConfig.PreferSimpleProtocol = true
	}

	if config.Logger != nil {
		pgxConfig.Logger = config.Logger
	}
	if config.LogLevel != 0 {
		pgxConfig.LogLevel = config.LogLevel
	}
	if config.BuildStatementCache != nil {
		pgxConfig.BuildStatementCache = config.BuildStatementCache
	}

	sqlDB := stdlib.OpenDB(*pgxConfig)
	database := bun.NewDB(sqlDB, pgdialect.New(), options...)
	return database, sqlDB, nil
}
