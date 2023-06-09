package bunframework

import (
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"io/fs"
)

var (
	ErrNoDB = errors.New("missing db object")
)

// RegisterSQLMigrations adds new sql-based migrations to the current configuration.
// See https://bun.uptrace.dev/guide/migrations.html#sql-based-migrations.
func (config *MigrateConfig) RegisterSQLMigrations(migrations ...fs.FS) {
	config.Files = append(config.Files, migrations...)
}

// RegisterGoMigrations adds new Go-based migrations to the current configuration.
// See https://bun.uptrace.dev/guide/migrations.html#go-based-migrations.
func (config *MigrateConfig) RegisterGoMigrations(migrations ...GoMigration) {
	config.Go = append(config.Go, migrations...)
}

func (config *MigrateConfig) getMigrations() (*migrate.Migrations, error) {
	migrations := migrate.NewMigrations()

	for i, migration := range config.Files {
		if err := migrations.Discover(migration); err != nil {
			return nil, fmt.Errorf("failed to discover migrations on filesystem %v: %w", i, err)
		}
	}

	for i, migration := range config.Go {
		if err := migrations.Register(migration.Up, migration.Down); err != nil {
			return nil, fmt.Errorf("failed to discover migrations on go migration %v: %w", i, err)
		}
	}

	return migrations, nil
}

// Execute runs all registered migrations. You can call this method multiple times, as bun knows
// to skip already executed migrations.
func (config *MigrateConfig) Execute(ctx context.Context, db *bun.DB) error {
	if db == nil {
		return ErrNoDB
	}

	migrations, err := config.getMigrations()
	if err != nil {
		return err
	}
	migrator := migrate.NewMigrator(db, migrations)

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if groups, err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	} else {
		config.migrations = groups
	}

	return nil
}

// Rollback previously executed migrations.
func (config *MigrateConfig) Rollback(ctx context.Context, db *bun.DB, opts ...migrate.MigrationOption) error {
	if db == nil {
		return ErrNoDB
	}

	migrations, err := config.getMigrations()
	if err != nil {
		return err
	}
	migrator := migrate.NewMigrator(db, migrations)

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if group, err := migrator.Rollback(ctx, opts...); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	} else {
		config.migrations = group
	}
	return nil
}

// Report returns the status of migrations after an Execute statement. It returns nil if Execute has not been called
// or has failed.
func (config *MigrateConfig) Report() *migrate.MigrationGroup {
	return config.migrations
}
