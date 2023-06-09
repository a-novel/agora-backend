package test_utils

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"github.com/a-novel/agora-backend/framework/bunframework"
	"github.com/a-novel/agora-backend/framework/bunframework/pgconfig"
	"github.com/a-novel/agora-backend/migrations"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"io/fs"
	"os"
	"testing"
	"time"
)

// GetPostgres returns a bun service with a test database.
func GetPostgres(t *testing.T) (*bun.DB, *sql.DB) {
	db, sqlDB, err := bunframework.NewClient(context.TODO(), bunframework.Config{
		Driver:      pgconfig.NewDriverWithDSN(os.Getenv("POSTGRES_URL")),
		Migrations:  &bunframework.MigrateConfig{Files: []fs.FS{migrations.Migrations}},
		ResetOnConn: true,
	})
	require.NoError(t, err)

	return db, sqlDB
}

// GetSecurityGenerateCode returns a mocked function for security.GenerateCode.
func GetSecurityGenerateCode(public, hashed string, err error) func() (string, string, error) {
	return func() (string, string, error) {
		return public, hashed, err
	}
}

// GetSecurityVerifyCode returns a mocked function for security.VerifyCode.
func GetSecurityVerifyCode(ok bool, err error) func(string, string) (bool, error) {
	return func(_ string, _ string) (bool, error) {
		return ok, err
	}
}

// GetBcryptGenerateFromPassword returns a mocked function for bcrypt.GenerateFromPassword.
func GetBcryptGenerateFromPassword(hashed string, err error) func([]byte, int) ([]byte, error) {
	return func(_ []byte, _ int) ([]byte, error) {
		return []byte(hashed), err
	}
}

// GetBcryptCompareHashAndPassword returns a mocked function for bcrypt.CompareHashAndPassword.
func GetBcryptCompareHashAndPassword(err error) func([]byte, []byte) error {
	return func(_, _ []byte) error {
		return err
	}
}

// GetTimeNow returns a mocked function for time.Now.
func GetTimeNow(now time.Time) func() time.Time {
	return func() time.Time {
		return now
	}
}

// GetUUID returns a mocked function for uuid.New.
func GetUUID(id uuid.UUID) func() uuid.UUID {
	return func() uuid.UUID {
		return id
	}
}

// GetJWKKeyGen returns a mocked function for security.JWKKeyGen.
func GetJWKKeyGen(key ed25519.PrivateKey, err error) func() (ed25519.PrivateKey, error) {
	return func() (ed25519.PrivateKey, error) {
		return key, err
	}
}
