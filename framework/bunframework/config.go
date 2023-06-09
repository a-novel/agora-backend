package bunframework

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/a-novel/agora-backend/framework/communication"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"io/fs"
	"os"
)

// DriverConfig provides a generic representation of a driver implementation.
//
// Available drivers are:
//
// - pgconfig.Driver
//
// - pgxconfig.Driver
//
// For example:
//
//	driver := pgconfig.NewDriverWithDSN(os.Getenv("POSTGRES_URL"))
type DriverConfig interface {
	// Connect to the driver and return a bun.DB object.
	Connect(options ...bun.DBOption) (*bun.DB, *sql.DB, error)
}

// GoMigration represents a set of migration functions for bun.
type GoMigration struct {
	// Up applies a migration.
	Up migrate.MigrationFunc
	// Down rollbacks the changes made by a migration, if applied.
	Down migrate.MigrationFunc
}

// MigrateConfig configures migrations for a given *bun.DB instance.
// See https://bun.uptrace.dev/guide/migrations.html.
//
// migrations.go
//
//	import "embed"
//
//	//go:embed *.sql
//	var Migrations embed.FS
//
// main.go
//
//	import "mypackage/path/to/migrations"
//
//	func main() {
//	  myMigrations := &bunframework.MigrateConfig{
//	    Files: []fs.FS{migrations.Migrations},
//	  }
//	}
type MigrateConfig struct {
	// Files represents a bunch of filesystems to look for SQL migration files.
	// See https://bun.uptrace.dev/guide/migrations.html#sql-based-migrations.
	Files []fs.FS
	// Go represents a bunch of GoMigration to execute when loading the driver.
	// See https://bun.uptrace.dev/guide/migrations.html#go-based-migrations.
	Go []GoMigration

	// Cache migration results.
	migrations *migrate.MigrationGroup
}

// Config is the main configuration object for a *bun.DB instance.
//
//	config := bunframework.Config{
//	  Driver: pgconfig.Driver{
//	    DSN:     os.Getenv("POSTGRES_URL"),
//	    AppName: "My Application",
//	  },
//	  Migrations: &bunframework.MigrateConfig{
//	    Files: []fs.FS{migrations.Migrations},
//	  },
//	}
type Config struct {
	// Driver used to communicate with the instance.
	Driver DriverConfig
	// Migrations is an optional value to run migrations automatically when loading the driver.
	// See https://bun.uptrace.dev/guide/migrations.html.
	Migrations *MigrateConfig

	// Production optimization for bun. It is recommended to set this to true for production builds.
	// See https://bun.uptrace.dev/guide/running-bun-in-production.html#bun-withdiscardunknowncolumns.
	DiscardUnknownColumns bool
	// ResetOnConn resets the whole database content when opening a new connection. Only use this under test
	// environments.
	ResetOnConn bool

	// Options is a fallback/security, to still allow to pass options in a conventional way. Also, it
	// allows Config to accept new options that have not or cannot (for any reason) be configured within
	// the object.
	Options []bun.DBOption
}

// NewClient creates a new *bun.DB instance from a Config object.
//
//	client, sqlClient, err := bunframework.NewClient(context.Background(), config)
//
// sqlClient is used to close all connections. You are responsible for freeing them once your application exits.
//
//	defer client.Close()
//	defer sqlClient.Close()
func NewClient(ctx context.Context, config Config) (*bun.DB, *sql.DB, error) {
	database, sqlDB, err := config.Driver.Connect(config.computeInternalOptions()...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database client: %w", err)
	}

	// Make sure database is up before running configuration. It avoids race conditions when the database instance
	// is updated/started during the deployment.
	if err := communication.WaitForPingAuto(PingDatabase(ctx, database)); err != nil {
		_ = database.Close()
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("failed to reach postgres database: %w", err)
	}

	if config.ResetOnConn {
		// Prevents accidental settings.
		if os.Getenv("ENV") != "test" {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, fmt.Errorf("ResetOnConn flag is only available in test environments")
		}

		// the below line alone erases the database, but makes it, so it may not be reusable normally
		if _, err := database.Exec("DROP SCHEMA IF EXISTS test CASCADE;"); err != nil {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, err
		}

		// ensures the database can be reused normally after the operation
		if _, err = database.Exec("CREATE SCHEMA IF NOT EXISTS test;"); err != nil {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, err
		}
		if _, err = database.Exec("GRANT ALL ON SCHEMA test TO test;"); err != nil {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, err
		}
		if _, err = database.Exec("GRANT ALL ON SCHEMA test TO test;"); err != nil {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, err
		}
	}

	// Apply migrations.
	if config.Migrations != nil {
		if err := config.Migrations.Execute(ctx, database); err != nil {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, err
		}
	}

	return database, sqlDB, nil
}

// NewClientWithDriver is a quick way to create a *bun.DB object with minimal configuration.
//
//	client, sqlClient, err := bunframework.NewClientWithDriver(pgconfig.NewDriverWithDSN(dsn))
//
// sqlClient is used to close all connections. You are responsible for freeing them once your application exits.
//
//	defer client.Close()
//	defer sqlClient.Close()
func NewClientWithDriver(driver DriverConfig) (*bun.DB, *sql.DB, error) {
	return driver.Connect()
}
