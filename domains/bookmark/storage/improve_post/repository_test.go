package improve_post_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
)

var Fixtures = []*Model{
	{
		UserID:    test_utils.NumberUUID(1),
		RequestID: test_utils.NumberUUID(2),
		CreatedAt: baseTime,
		Target:    BookmarkTargetImproveRequest,
		Level:     bookmark_storage.LevelBookmark,
	},
}

func TestImprovePostRepository_Bookmark(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		userID    uuid.UUID
		requestID uuid.UUID
		target    BookmarkTarget
		level     bookmark_storage.Level
		now       time.Time

		expect    *Model
		expectErr error
	}{
		{
			name:      "Success",
			userID:    test_utils.NumberUUID(10),
			requestID: test_utils.NumberUUID(20),
			target:    BookmarkTargetImproveSuggestion,
			level:     bookmark_storage.LevelFavorite,
			now:       baseTime,
			expect: &Model{
				BaseModel: bun.BaseModel{},
				UserID:    test_utils.NumberUUID(10),
				RequestID: test_utils.NumberUUID(20),
				CreatedAt: baseTime,
				Target:    BookmarkTargetImproveSuggestion,
				Level:     bookmark_storage.LevelFavorite,
			},
		},
		{
			name:      "Success/Update",
			userID:    test_utils.NumberUUID(1),
			requestID: test_utils.NumberUUID(2),
			target:    BookmarkTargetImproveRequest,
			level:     bookmark_storage.LevelFavorite,
			now:       baseTime.Add(time.Hour),
			expect: &Model{
				BaseModel: bun.BaseModel{},
				UserID:    test_utils.NumberUUID(1),
				RequestID: test_utils.NumberUUID(2),
				CreatedAt: baseTime.Add(time.Hour),
				Target:    BookmarkTargetImproveRequest,
				Level:     bookmark_storage.LevelFavorite,
			},
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx)

				res, err := repository.Bookmark(ctx, d.userID, d.requestID, d.target, d.level, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImprovePostRepository_UnBookmark(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		userID    uuid.UUID
		requestID uuid.UUID
		target    BookmarkTarget

		expectErr error
	}{
		{
			name:      "Success",
			userID:    test_utils.NumberUUID(1),
			requestID: test_utils.NumberUUID(2),
			target:    BookmarkTargetImproveRequest,
		},
		{
			name:      "Error/NotFound",
			userID:    test_utils.NumberUUID(1),
			requestID: test_utils.NumberUUID(2),
			target:    BookmarkTargetImproveSuggestion,
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx)

				test_utils.RequireError(t, d.expectErr, repository.UnBookmark(ctx, d.userID, d.requestID, d.target))
			})
		}
	})
	require.NoError(t, err)
}

func TestImprovePostRepository_IsBookmarked(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		userID    uuid.UUID
		requestID uuid.UUID
		target    BookmarkTarget

		expect    *bookmark_storage.Level
		expectErr error
	}{
		{
			name:      "Success/True",
			userID:    test_utils.NumberUUID(1),
			requestID: test_utils.NumberUUID(2),
			target:    BookmarkTargetImproveRequest,
			expect:    framework.ToPTR(bookmark_storage.LevelBookmark),
		},
		{
			name:      "Success/False",
			userID:    test_utils.NumberUUID(1),
			requestID: test_utils.NumberUUID(2),
			target:    BookmarkTargetImproveSuggestion,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.IsBookmarked(ctx, d.userID, d.requestID, d.target)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImprovePostRepository_List(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*Model{
		// User 1 - Bookmarks
		{
			UserID:    test_utils.NumberUUID(1),
			RequestID: test_utils.NumberUUID(10),
			CreatedAt: baseTime,
			Target:    BookmarkTargetImproveRequest,
			Level:     bookmark_storage.LevelBookmark,
		},
		{
			UserID:    test_utils.NumberUUID(1),
			RequestID: test_utils.NumberUUID(11),
			CreatedAt: baseTime.Add(time.Hour),
			Target:    BookmarkTargetImproveRequest,
			Level:     bookmark_storage.LevelBookmark,
		},
		{
			UserID:    test_utils.NumberUUID(1),
			RequestID: test_utils.NumberUUID(12),
			CreatedAt: baseTime.Add(30 * time.Minute),
			Target:    BookmarkTargetImproveRequest,
			Level:     bookmark_storage.LevelBookmark,
		},
		{
			UserID:    test_utils.NumberUUID(1),
			RequestID: test_utils.NumberUUID(10),
			CreatedAt: baseTime.Add(35 * time.Minute),
			Target:    BookmarkTargetImproveSuggestion,
			Level:     bookmark_storage.LevelBookmark,
		},
		// User 1 - Favorites
		{
			UserID:    test_utils.NumberUUID(1),
			RequestID: test_utils.NumberUUID(20),
			CreatedAt: baseTime,
			Target:    BookmarkTargetImproveRequest,
			Level:     bookmark_storage.LevelFavorite,
		},
		{
			UserID:    test_utils.NumberUUID(1),
			RequestID: test_utils.NumberUUID(21),
			CreatedAt: baseTime.Add(time.Hour),
			Target:    BookmarkTargetImproveRequest,
			Level:     bookmark_storage.LevelFavorite,
		},
		{
			UserID:    test_utils.NumberUUID(1),
			RequestID: test_utils.NumberUUID(22),
			CreatedAt: baseTime.Add(30 * time.Minute),
			Target:    BookmarkTargetImproveRequest,
			Level:     bookmark_storage.LevelFavorite,
		},
		{
			UserID:    test_utils.NumberUUID(1),
			RequestID: test_utils.NumberUUID(20),
			CreatedAt: baseTime.Add(35 * time.Minute),
			Target:    BookmarkTargetImproveSuggestion,
			Level:     bookmark_storage.LevelFavorite,
		},
		// User 2
		{
			UserID:    test_utils.NumberUUID(2),
			RequestID: test_utils.NumberUUID(10),
			CreatedAt: baseTime,
			Target:    BookmarkTargetImproveRequest,
			Level:     bookmark_storage.LevelBookmark,
		},
	}

	data := []struct {
		name string

		userID uuid.UUID
		level  bookmark_storage.Level
		target BookmarkTarget
		limit  int
		offset int

		expect      []*Model
		expectCount int64
		expectErr   error
	}{
		{
			name:   "Success/Bookmarks",
			userID: test_utils.NumberUUID(1),
			level:  bookmark_storage.LevelBookmark,
			target: BookmarkTargetImproveRequest,
			limit:  10,
			offset: 1,
			expect: []*Model{
				fixtures[2],
				fixtures[0],
			},
			expectCount: 3,
		},
		{
			name:   "Success/Favorites",
			userID: test_utils.NumberUUID(1),
			level:  bookmark_storage.LevelFavorite,
			target: BookmarkTargetImproveRequest,
			limit:  10,
			offset: 1,
			expect: []*Model{
				fixtures[6],
				fixtures[4],
			},
			expectCount: 3,
		},
		{
			name:   "Success/ImproveSuggestions",
			userID: test_utils.NumberUUID(1),
			level:  bookmark_storage.LevelBookmark,
			target: BookmarkTargetImproveSuggestion,
			limit:  10,
			offset: 0,
			expect: []*Model{
				fixtures[3],
			},
			expectCount: 1,
		},
		{
			name:        "Success/Empty",
			userID:      test_utils.NumberUUID(4),
			level:       bookmark_storage.LevelFavorite,
			target:      BookmarkTargetImproveRequest,
			limit:       10,
			offset:      1,
			expect:      []*Model(nil),
			expectCount: 0,
		},
	}

	err := test_utils.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, count, err := repository.List(ctx, d.userID, d.level, d.target, d.limit, d.offset)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
				require.Equal(t, d.expectCount, count)
			})
		}
	})
	require.NoError(t, err)
}
