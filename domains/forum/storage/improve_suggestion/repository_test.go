package improve_suggestion_storage

import (
	"context"
	"encoding/json"
	"github.com/a-novel/agora-backend/framework"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

func TestImproveSuggestionRepository_Read(t *testing.T) {
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
			name: "Success",
			id:   test_utils.NumberUUID(1001),
			expect: &Model{
				ID:        test_utils.NumberUUID(1001),
				CreatedAt: baseTime,
				SourceID:  test_utils.NumberUUID(1000),
				UserID:    test_utils.NumberUUID(201),
				Validated: true,
				UpVotes:   7,
				DownVotes: 2,
				Core: Core{
					RequestID: test_utils.NumberUUID(1000),
					Title:     "Test",
					Content:   "Smart content.",
				},
			},
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1010),
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

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

func TestImproveSuggestionRepository_Create(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		data     *Core
		userID   uuid.UUID
		sourceID uuid.UUID
		id       uuid.UUID
		now      time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			data: &Core{
				RequestID: test_utils.NumberUUID(1000),
				Title:     "Test",
				Content:   "Intelligent content.",
			},
			userID:   test_utils.NumberUUID(200),
			sourceID: test_utils.NumberUUID(1000),
			id:       test_utils.NumberUUID(1),
			now:      baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				SourceID:  test_utils.NumberUUID(1000),
				UserID:    test_utils.NumberUUID(200),
				Core: Core{
					RequestID: test_utils.NumberUUID(1000),
					Title:     "Test",
					Content:   "Intelligent content.",
				},
			},
		},
		{
			name: "Success/OnRevision",
			data: &Core{
				RequestID: test_utils.NumberUUID(1001),
				Title:     "Test",
				Content:   "Intelligent content.",
			},
			userID:   test_utils.NumberUUID(200),
			sourceID: test_utils.NumberUUID(1000),
			id:       test_utils.NumberUUID(1),
			now:      baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				SourceID:  test_utils.NumberUUID(1000),
				UserID:    test_utils.NumberUUID(200),
				Core: Core{
					RequestID: test_utils.NumberUUID(1001),
					Title:     "Test",
					Content:   "Intelligent content.",
				},
			},
		},
		{
			name: "Error/SourceIsARevision",
			data: &Core{
				RequestID: test_utils.NumberUUID(1001),
				Title:     "Test",
				Content:   "Intelligent content.",
			},
			userID:    test_utils.NumberUUID(200),
			sourceID:  test_utils.NumberUUID(1001),
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrMissingRelation,
		},
		{
			name: "Error/MissingSource",
			data: &Core{
				RequestID: test_utils.NumberUUID(1000),
				Title:     "Test",
				Content:   "Intelligent content.",
			},
			userID:    test_utils.NumberUUID(200),
			sourceID:  test_utils.NumberUUID(100),
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrMissingRelation,
		},
		{
			name: "Error/MissingRequest",
			data: &Core{
				RequestID: test_utils.NumberUUID(100),
				Title:     "Test",
				Content:   "Intelligent content.",
			},
			userID:    test_utils.NumberUUID(200),
			sourceID:  test_utils.NumberUUID(1000),
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrMissingRelation,
		},
		{
			name: "Error/RevisionAndSourceMismatch",
			data: &Core{
				RequestID: test_utils.NumberUUID(2000),
				Title:     "Test",
				Content:   "Intelligent content.",
			},
			userID:    test_utils.NumberUUID(200),
			sourceID:  test_utils.NumberUUID(1000),
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrMissingRelation,
		},
		{
			name: "Error/AlreadyExists",
			data: &Core{
				RequestID: test_utils.NumberUUID(1000),
				Title:     "Test",
				Content:   "Intelligent content.",
			},
			userID:    test_utils.NumberUUID(200),
			sourceID:  test_utils.NumberUUID(1000),
			id:        test_utils.NumberUUID(1000),
			now:       baseTime,
			expectErr: validation.ErrUniqConstraintViolation,
		},
		{
			name: "Error/NoTitle",
			data: &Core{
				RequestID: test_utils.NumberUUID(1000),
				Content:   "Intelligent content.",
			},
			userID:    test_utils.NumberUUID(200),
			sourceID:  test_utils.NumberUUID(1000),
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name: "Error/NoContent",
			data: &Core{
				RequestID: test_utils.NumberUUID(1000),
				Title:     "Test",
			},
			userID:    test_utils.NumberUUID(200),
			sourceID:  test_utils.NumberUUID(1000),
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx, 10)

				res, err := repository.Create(ctx, d.data, d.userID, d.sourceID, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveSuggestionRepository_Update(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		data *Core
		id   uuid.UUID
		now  time.Time

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			data: &Core{
				RequestID: test_utils.NumberUUID(1002),
				Title:     "Test",
				Content:   "Good content.",
			},
			id:  test_utils.NumberUUID(1002),
			now: updateTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(1002),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				SourceID:  test_utils.NumberUUID(1000),
				UserID:    test_utils.NumberUUID(200),
				UpVotes:   21,
				DownVotes: 8,
				Core: Core{
					RequestID: test_utils.NumberUUID(1002),
					Title:     "Test",
					Content:   "Good content.",
				},
			},
		},
		{
			name: "Error/MismatchingSourceAndRevision",
			data: &Core{
				RequestID: test_utils.NumberUUID(5000),
				Title:     "Test",
				Content:   "Good content.",
			},
			id:        test_utils.NumberUUID(1002),
			now:       baseTime,
			expectErr: validation.ErrMissingRelation,
		},
		{
			name: "Error/MissingSource",
			data: &Core{
				RequestID: test_utils.NumberUUID(1010),
				Title:     "Test",
				Content:   "Good content.",
			},
			id:        test_utils.NumberUUID(1002),
			now:       baseTime,
			expectErr: validation.ErrMissingRelation,
		},
		{
			name: "Error/NoTitle",
			data: &Core{
				RequestID: test_utils.NumberUUID(1002),
				Content:   "Good content.",
			},
			id:        test_utils.NumberUUID(1002),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name: "Error/NoContent",
			data: &Core{
				RequestID: test_utils.NumberUUID(1002),
				Title:     "Test",
			},
			id:        test_utils.NumberUUID(1002),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx, 10)

				res, err := repository.Update(ctx, d.data, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveSuggestionRepository_Delete(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id uuid.UUID

		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1002),
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1010),
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx, 10)
				test_utils.RequireError(t, d.expectErr, repository.Delete(ctx, d.id))
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveSuggestionRepository_Validate(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		validated bool
		id        uuid.UUID

		expect    *Model
		expectErr error
	}{
		{
			name:      "Success/Validated",
			validated: true,
			id:        test_utils.NumberUUID(1002),
			expect: &Model{
				ID:        test_utils.NumberUUID(1002),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				SourceID:  test_utils.NumberUUID(1000),
				UserID:    test_utils.NumberUUID(200),
				UpVotes:   21,
				DownVotes: 8,
				Validated: true,
				Core: Core{
					RequestID: test_utils.NumberUUID(1001),
					Title:     "Test",
					Content:   "Simple content.",
				},
			},
		},
		{
			name:      "Success/Unvalidated",
			validated: false,
			id:        test_utils.NumberUUID(1001),
			expect: &Model{
				ID:        test_utils.NumberUUID(1001),
				CreatedAt: baseTime,
				SourceID:  test_utils.NumberUUID(1000),
				UserID:    test_utils.NumberUUID(201),
				UpVotes:   7,
				DownVotes: 2,
				Core: Core{
					RequestID: test_utils.NumberUUID(1000),
					Title:     "Test",
					Content:   "Smart content.",
				},
			},
		},
		{
			name:      "Error/NotFound",
			validated: true,
			id:        test_utils.NumberUUID(1010),
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx, 10)

				res, err := repository.Validate(ctx, d.validated, d.id)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveSuggestionRepository_List(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		query  ListQuery
		limit  int
		offset int

		expect      []*Model
		expectCount int64
		expectErr   error
	}{
		// Validated.
		{
			name: "Success/Validated",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				Validated: framework.ToPTR(true),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			expectCount: 4,
			limit:       10,
			offset:      0,
			// Sorted by score.
			expect: []*Model{
				{ // 10, validated
					ID:        test_utils.NumberUUID(1000),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					Validated: true,
					UpVotes:   10,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test 2",
						Content:   "Dummy cont",
					},
				},
				{ // 5, validated
					ID:        test_utils.NumberUUID(1001),
					CreatedAt: baseTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					Validated: true,
					UpVotes:   7,
					DownVotes: 2,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test",
						Content:   "Smart cont",
					},
				},
				{ // 3, validated
					ID:        test_utils.NumberUUID(1003),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					Validated: true,
					UpVotes:   16,
					DownVotes: 13,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 3",
						Content:   "Simple con",
					},
				},
				{ // 0, validated
					ID:        test_utils.NumberUUID(1007),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					Validated: true,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 7",
						Content:   "Simple con",
					},
				},
			},
		},
		{
			name: "Success/ValidatedWithRequestID",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(1001)),
				Validated: framework.ToPTR(true),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			expectCount: 2,
			limit:       10,
			offset:      0,
			// Sorted by score.
			expect: []*Model{
				{ // 3, validated
					ID:        test_utils.NumberUUID(1003),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					Validated: true,
					UpVotes:   16,
					DownVotes: 13,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 3",
						Content:   "Simple con",
					},
				},
				{ // 0, validated
					ID:        test_utils.NumberUUID(1007),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					Validated: true,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 7",
						Content:   "Simple con",
					},
				},
			},
		},
		{
			name:        "Success/ValidatedWithPagination",
			expectCount: 4,
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				Validated: framework.ToPTR(true),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			limit:  3,
			offset: 2,
			// Sorted by score.
			expect: []*Model{
				{ // 3, validated
					ID:        test_utils.NumberUUID(1003),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					Validated: true,
					UpVotes:   16,
					DownVotes: 13,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 3",
						Content:   "Simple con",
					},
				},
				{ // 0, validated
					ID:        test_utils.NumberUUID(1007),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					Validated: true,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 7",
						Content:   "Simple con",
					},
				},
			},
		},
		{
			name: "Success/ValidatedWithLimit",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				Validated: framework.ToPTR(true),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			expectCount: 4,
			limit:       2,
			offset:      0,
			// Sorted by score.
			expect: []*Model{
				{ // 10, validated
					ID:        test_utils.NumberUUID(1000),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					Validated: true,
					UpVotes:   10,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test 2",
						Content:   "Dummy cont",
					},
				},
				{ // 5, validated
					ID:        test_utils.NumberUUID(1001),
					CreatedAt: baseTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					Validated: true,
					UpVotes:   7,
					DownVotes: 2,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test",
						Content:   "Smart cont",
					},
				},
			},
		},
		// Not validated.
		{
			name: "Success",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				Validated: framework.ToPTR(false),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			expectCount: 4,
			limit:       10,
			offset:      0,
			// Sorted by score.
			expect: []*Model{
				{ // 13
					ID:        test_utils.NumberUUID(1002),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					UpVotes:   21,
					DownVotes: 8,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test",
						Content:   "Simple con",
					},
				},
				{ // 9
					ID:        test_utils.NumberUUID(1006),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					UpVotes:   9,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test 6",
						Content:   "Simple con",
					},
				},
				{ // 8
					ID:        test_utils.NumberUUID(1005),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					UpVotes:   32,
					DownVotes: 24,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test 5",
						Content:   "Simple con",
					},
				},
				{ // -4
					ID:        test_utils.NumberUUID(1004),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					UpVotes:   4,
					DownVotes: 8,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 4",
						Content:   "Simple con",
					},
				},
			},
		},
		{
			name: "Success/WithRequestID",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(1001)),
				Validated: framework.ToPTR(false),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			expectCount: 2,
			limit:       10,
			offset:      0,
			// Sorted by score.
			expect: []*Model{
				{ // 13
					ID:        test_utils.NumberUUID(1002),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					UpVotes:   21,
					DownVotes: 8,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test",
						Content:   "Simple con",
					},
				},
				{ // -4
					ID:        test_utils.NumberUUID(1004),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					UpVotes:   4,
					DownVotes: 8,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 4",
						Content:   "Simple con",
					},
				},
			},
		},
		{
			name: "Success/WithPagination",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				Validated: framework.ToPTR(false),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			expectCount: 4,
			limit:       3,
			offset:      2,
			// Sorted by score.
			expect: []*Model{
				{ // 8
					ID:        test_utils.NumberUUID(1005),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					UpVotes:   32,
					DownVotes: 24,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test 5",
						Content:   "Simple con",
					},
				},
				{ // -4
					ID:        test_utils.NumberUUID(1004),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					UpVotes:   4,
					DownVotes: 8,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 4",
						Content:   "Simple con",
					},
				},
			},
		},
		{
			name: "Success/Limit",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				Validated: framework.ToPTR(false),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			expectCount: 4,
			limit:       2,
			offset:      0,
			// Sorted by score.
			expect: []*Model{
				{ // 13
					ID:        test_utils.NumberUUID(1002),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					UpVotes:   21,
					DownVotes: 8,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test",
						Content:   "Simple con",
					},
				},
				{ // 9
					ID:        test_utils.NumberUUID(1006),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					UpVotes:   9,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test 6",
						Content:   "Simple con",
					},
				},
			},
		},
		// User id.
		{
			name: "Success/UserID",
			query: ListQuery{
				UserID: framework.ToPTR(test_utils.NumberUUID(201)),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			expectCount: 4,
			limit:       10,
			offset:      0,
			// Sorted by score.
			expect: []*Model{
				{ // 5, validated
					ID:        test_utils.NumberUUID(1001),
					CreatedAt: baseTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					Validated: true,
					UpVotes:   7,
					DownVotes: 2,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test",
						Content:   "Smart cont",
					},
				},
				{ // 3
					ID:        test_utils.NumberUUID(2000),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(5000),
					UserID:    test_utils.NumberUUID(201),
					UpVotes:   4,
					DownVotes: 1,
					Core: Core{
						RequestID: test_utils.NumberUUID(5000),
						Title:     "Ipsum Lorem",
						Content:   "Lorem ipsu",
					},
				},
				{ // 0, validated
					ID:        test_utils.NumberUUID(1007),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					Validated: true,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 7",
						Content:   "Simple con",
					},
				},
				{ // -4
					ID:        test_utils.NumberUUID(1004),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					UpVotes:   4,
					DownVotes: 8,
					Core: Core{
						RequestID: test_utils.NumberUUID(1001),
						Title:     "Test 4",
						Content:   "Simple con",
					},
				},
			},
		},
		// No results.
		{
			name: "Success/MissingSource",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1010)),
				Validated: framework.ToPTR(true),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			limit:  10,
			offset: 0,
			// Sorted by score.
			expect: []*Model(nil),
		},
		{
			name: "Success/MissingRequest",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(1000)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(1010)),
				Validated: framework.ToPTR(true),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			limit:  10,
			offset: 0,
			// Sorted by score.
			expect: []*Model(nil),
		},
		{
			name: "Success/NoSuggestions",
			query: ListQuery{
				SourceID:  framework.ToPTR(test_utils.NumberUUID(6000)),
				Validated: framework.ToPTR(true),
				Order: &SearchQueryOrder{
					Score: true,
				},
			},
			limit:  10,
			offset: 0,
			// Sorted by score.
			expect: []*Model(nil),
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, count, err := repository.List(ctx, d.query, d.limit, d.offset)
				test_utils.RequireError(t, d.expectErr, err)

				mrshRes, _ := json.Marshal(res)
				mrshExpect, _ := json.Marshal(d.expect)

				require.Equal(t, d.expect, res, string(mrshRes), string(mrshExpect))
				require.Equal(t, d.expectCount, count)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveSuggestionRepository_IsCreator(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id     uuid.UUID
		userID uuid.UUID

		expect    bool
		expectErr error
	}{
		{
			name:   "Success",
			id:     test_utils.NumberUUID(1001),
			userID: test_utils.NumberUUID(201),
			expect: true,
		},
		{
			name:   "Success/NotFound",
			id:     test_utils.NumberUUID(1001),
			userID: test_utils.NumberUUID(202),
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.IsCreator(ctx, d.userID, d.id)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveSuggestionRepository_GetPreviews(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		ids []uuid.UUID

		expect    []*Model
		expectErr error
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1001),
				test_utils.NumberUUID(1006),
				test_utils.NumberUUID(3000),
			},
			expect: []*Model{
				{
					ID:        test_utils.NumberUUID(1001),
					CreatedAt: baseTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(201),
					Validated: true,
					UpVotes:   7,
					DownVotes: 2,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test",
						Content:   "Smart cont",
					},
				},
				{
					ID:        test_utils.NumberUUID(1006),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(1000),
					UserID:    test_utils.NumberUUID(200),
					UpVotes:   9,
					Core: Core{
						RequestID: test_utils.NumberUUID(1000),
						Title:     "Test 6",
						Content:   "Simple con",
					},
				},
			},
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.GetPreviews(ctx, d.ids)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}
