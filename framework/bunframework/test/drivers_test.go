package test

import (
	"github.com/a-novel/agora-backend/framework/bunframework"
	"github.com/a-novel/agora-backend/framework/bunframework/pgconfig"
	"github.com/a-novel/agora-backend/framework/bunframework/pgxconfig"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/driver/pgdriver"
	"testing"
)

func TestPGConfig(t *testing.T) {
	t.Run("NewDriverWithDSN", func(t *testing.T) {
		db, sqlDB, err := bunframework.NewClientWithDriver(pgconfig.NewDriverWithDSN(dsn))
		require.NoError(t, err)
		require.NoError(t, db.Ping())
		require.NoError(t, db.Close())
		require.NoError(t, sqlDB.Close())
	})

	t.Run("NewDriver", func(t *testing.T) {
		db, sqlDB, err := bunframework.NewClientWithDriver(pgconfig.NewDriver(pgdriver.WithDSN(dsn)))
		require.NoError(t, err)
		require.NoError(t, db.Ping())
		require.NoError(t, db.Close())
		require.NoError(t, sqlDB.Close())
	})

	t.Run("Driver", func(t *testing.T) {
		db, sqlDB, err := bunframework.NewClientWithDriver(pgconfig.Driver{DSN: dsn})
		require.NoError(t, err)
		require.NoError(t, db.Ping())
		require.NoError(t, db.Close())
		require.NoError(t, sqlDB.Close())
	})
}

func TestPGXConfig(t *testing.T) {
	t.Run("NewDriverWithDSN", func(t *testing.T) {
		db, sqlDB, err := bunframework.NewClientWithDriver(pgxconfig.NewDriverWithDSN(dsn))
		require.NoError(t, err)
		require.NoError(t, db.Ping())
		require.NoError(t, db.Close())
		require.NoError(t, sqlDB.Close())
	})

	t.Run("Driver", func(t *testing.T) {
		db, sqlDB, err := bunframework.NewClientWithDriver(pgxconfig.Driver{DSN: dsn, KeepImplicitPreparedStatements: true})
		require.NoError(t, err)
		require.NoError(t, db.Ping())
		require.NoError(t, db.Close())
		require.NoError(t, sqlDB.Close())
	})
}
