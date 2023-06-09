package identity_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

var (
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 9, 0, 0, 0, time.UTC)
)

var Fixtures = []*Model{
	// Standard.
	{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: Core{
			FirstName: "Anna",
			LastName:  "Banana",
			Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
			Sex:       models.SexFemale,
		},
	},
}

func TestIdentityRepository_Read(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id uuid.UUID

		expect    *Model
		expectErr error
	}{
		{
			name:   "Success",
			id:     test_utils.NumberUUID(1000),
			expect: Fixtures[0],
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1),
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.Read(ctx, d.id)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestIdentityRepository_Update(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		core *Core
		id   uuid.UUID
		now  time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			core: &Core{
				FirstName: "Robert",
				LastName:  "Linder",
				Birthday:  time.Date(1988, time.November, 16, 19, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			id:  test_utils.NumberUUID(1000),
			now: updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					FirstName: "Robert",
					LastName:  "Linder",
					Birthday:  time.Date(1988, time.November, 16, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
		},
		{
			name: "Success/NonRomanizedName",
			core: &Core{
				FirstName: "ルイズ フランソワーズ ル ブラン",
				LastName:  "ド ラ ヴァリエール",
				Birthday:  time.Date(1988, time.November, 16, 19, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			id:  test_utils.NumberUUID(1000),
			now: updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					FirstName: "ルイズ フランソワーズ ル ブラン",
					LastName:  "ド ラ ヴァリエール",
					Birthday:  time.Date(1988, time.November, 16, 19, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
		},
		{
			name: "Error/NotFound",
			core: &Core{
				FirstName: "Robert",
				LastName:  "Linder",
				Birthday:  time.Date(1988, time.November, 16, 19, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			id:        test_utils.NumberUUID(1),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.Update(ctx, d.core, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}
