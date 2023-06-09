package profile_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
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
			Username: "AnnaLaBanana",
			Slug:     "fruit-basket",
		},
	},
	// Without username.
	{
		ID:        test_utils.NumberUUID(1001),
		CreatedAt: baseTime,
		UpdatedAt: &baseTime,
		Core: Core{
			Slug: "square-the-circle",
		},
	},
}

func TestProfileRepository_Read(t *testing.T) {
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

func TestProfileRepository_ReadSlug(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		slug string

		expect    *Model
		expectErr error
	}{
		{
			name:   "Success",
			slug:   "fruit-basket",
			expect: Fixtures[0],
		},
		{
			name:      "Error/NotFound",
			slug:      "poivresel",
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.ReadSlug(ctx, d.slug)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestProfileRepository_SlugExists(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		slug string

		expect    bool
		expectErr error
	}{
		{
			name:   "Success",
			slug:   "fruit-basket",
			expect: true,
		},
		{
			name:   "Success/NotExist",
			slug:   "poivresel",
			expect: false,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.SlugExists(ctx, d.slug)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestProfileRepository_Update(t *testing.T) {
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
				Username: "BanAnna",
				Slug:     "coco-nuts",
			},
			id:  test_utils.NumberUUID(1000),
			now: updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Username: "BanAnna",
					Slug:     "coco-nuts",
				},
			},
		},
		{
			name: "Success/RemoveUsername",
			core: &Core{
				Slug: "coco-nuts",
			},
			id:  test_utils.NumberUUID(1000),
			now: updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Slug: "coco-nuts",
				},
			},
		},
		{
			name: "Success/AddUsername",
			core: &Core{
				Username: "Dynasty of Penguins",
				Slug:     "square-the-circle",
			},
			id:  test_utils.NumberUUID(1001),
			now: updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1001),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: Core{
					Username: "Dynasty of Penguins",
					Slug:     "square-the-circle",
				},
			},
		},
		{
			name: "Error/NotFound",
			core: &Core{
				Username: "BanAnna",
				Slug:     "coco-nuts",
			},
			id:        test_utils.NumberUUID(1),
			now:       updateTime,
			expectErr: validation.ErrNotFound,
		},
		{
			name: "Error/RemoveSlug",
			core: &Core{
				Username: "BanAnna",
			},
			id:        test_utils.NumberUUID(1000),
			now:       updateTime,
			expectErr: validation.ErrConstraintViolation,
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
