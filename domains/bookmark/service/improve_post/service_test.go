package improve_post_service

import (
	"context"
	"errors"
	"github.com/a-novel/agora-backend/framework"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	fooErr   = errors.New("it broken")
)

func TestImprovePostService_Bookmark(t *testing.T) {
	data := []struct {
		name string

		userID    uuid.UUID
		requestID uuid.UUID
		target    models.BookmarkTarget
		level     models.BookmarkLevel
		now       time.Time

		shouldCallRepository bool
		bookmarkData         *improve_post_storage.Model
		bookmarkErr          error

		expect    *models.Bookmark
		expectErr error
	}{
		{
			name:                 "Success",
			userID:               test_utils.NumberUUID(10),
			requestID:            test_utils.NumberUUID(20),
			target:               models.BookmarkTargetImproveSuggestion,
			level:                models.BookmarkLevelFavorite,
			now:                  baseTime,
			shouldCallRepository: true,
			bookmarkData: &improve_post_storage.Model{
				UserID:    test_utils.NumberUUID(10),
				RequestID: test_utils.NumberUUID(20),
				CreatedAt: baseTime,
				Target:    improve_post_storage.BookmarkTargetImproveSuggestion,
				Level:     bookmark_storage.LevelFavorite,
			},
			expect: &models.Bookmark{
				UserID:    test_utils.NumberUUID(10),
				RequestID: test_utils.NumberUUID(20),
				CreatedAt: baseTime,
				Target:    models.BookmarkTargetImproveSuggestion,
				Level:     models.BookmarkLevelFavorite,
			},
		},
		{
			name:      "Error/InvalidLevel",
			userID:    test_utils.NumberUUID(10),
			requestID: test_utils.NumberUUID(20),
			target:    models.BookmarkTargetImproveSuggestion,
			level:     models.BookmarkLevel("foo"),
			now:       baseTime,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:      "Error/InvalidTarget",
			userID:    test_utils.NumberUUID(10),
			requestID: test_utils.NumberUUID(20),
			target:    models.BookmarkTarget("foo"),
			level:     models.BookmarkLevelFavorite,
			now:       baseTime,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:                 "Error/RepositoryFailure",
			userID:               test_utils.NumberUUID(10),
			requestID:            test_utils.NumberUUID(20),
			target:               models.BookmarkTargetImproveSuggestion,
			level:                models.BookmarkLevelFavorite,
			now:                  baseTime,
			shouldCallRepository: true,
			bookmarkErr:          fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_post_storage.NewMockRepository(st)

			if d.shouldCallRepository {
				repository.
					On(
						"Bookmark", context.TODO(), d.userID, d.requestID,
						improve_post_storage.BookmarkTarget(d.target),
						bookmark_storage.Level(d.level),
						d.now,
					).
					Return(d.bookmarkData, d.bookmarkErr)
			}

			service := NewService(repository)

			res, err := service.Bookmark(context.TODO(), d.userID, d.requestID, d.target, d.level, d.now)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImprovePostService_UnBookmark(t *testing.T) {
	data := []struct {
		name string

		userID    uuid.UUID
		requestID uuid.UUID
		target    models.BookmarkTarget

		shouldCallRepository bool
		bookmarkErr          error

		expectErr error
	}{
		{
			name:                 "Success",
			userID:               test_utils.NumberUUID(10),
			requestID:            test_utils.NumberUUID(20),
			target:               models.BookmarkTargetImproveSuggestion,
			shouldCallRepository: true,
		},
		{
			name:      "Error/InvalidTarget",
			userID:    test_utils.NumberUUID(10),
			requestID: test_utils.NumberUUID(20),
			target:    models.BookmarkTarget("foo"),
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:                 "Error/RepositoryFailure",
			userID:               test_utils.NumberUUID(10),
			requestID:            test_utils.NumberUUID(20),
			target:               models.BookmarkTargetImproveSuggestion,
			shouldCallRepository: true,
			bookmarkErr:          fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_post_storage.NewMockRepository(st)

			if d.shouldCallRepository {
				repository.
					On("UnBookmark", context.TODO(), d.userID, d.requestID, improve_post_storage.BookmarkTarget(d.target)).
					Return(d.bookmarkErr)
			}

			service := NewService(repository)

			test_utils.RequireError(st, d.expectErr, service.UnBookmark(context.TODO(), d.userID, d.requestID, d.target))
			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImprovePostService_IsBookmarked(t *testing.T) {
	data := []struct {
		name string

		userID    uuid.UUID
		requestID uuid.UUID
		target    models.BookmarkTarget

		shouldCallRepository bool
		bookmarkData         *bookmark_storage.Level
		bookmarkErr          error

		expect    *models.BookmarkLevel
		expectErr error
	}{
		{
			name:                 "Success",
			userID:               test_utils.NumberUUID(10),
			requestID:            test_utils.NumberUUID(20),
			target:               models.BookmarkTargetImproveSuggestion,
			shouldCallRepository: true,
			bookmarkData:         framework.ToPTR(bookmark_storage.LevelFavorite),
			expect:               framework.ToPTR(models.BookmarkLevelFavorite),
		},
		{
			name:                 "Success/NotBookmarked",
			userID:               test_utils.NumberUUID(10),
			requestID:            test_utils.NumberUUID(20),
			target:               models.BookmarkTargetImproveSuggestion,
			shouldCallRepository: true,
		},
		{
			name:      "Error/InvalidTarget",
			userID:    test_utils.NumberUUID(10),
			requestID: test_utils.NumberUUID(20),
			target:    models.BookmarkTarget("foo"),
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:                 "Error/RepositoryFailure",
			userID:               test_utils.NumberUUID(10),
			requestID:            test_utils.NumberUUID(20),
			target:               models.BookmarkTargetImproveSuggestion,
			shouldCallRepository: true,
			bookmarkErr:          fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_post_storage.NewMockRepository(st)

			if d.shouldCallRepository {
				repository.
					On("IsBookmarked", context.TODO(), d.userID, d.requestID, improve_post_storage.BookmarkTarget(d.target)).
					Return(d.bookmarkData, d.bookmarkErr)
			}

			service := NewService(repository)

			res, err := service.IsBookmarked(context.TODO(), d.userID, d.requestID, d.target)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImprovePostService_List(t *testing.T) {
	data := []struct {
		name string

		userID uuid.UUID
		level  models.BookmarkLevel
		target models.BookmarkTarget
		limit  int
		offset int

		shouldCallRepository bool
		bookmarkData         []*improve_post_storage.Model
		bookmarkCount        int64
		bookmarkErr          error

		expect      []*models.Bookmark
		expectCount int64
		expectErr   error
	}{
		{
			name:                 "Success",
			userID:               test_utils.NumberUUID(10),
			level:                models.BookmarkLevelFavorite,
			target:               models.BookmarkTargetImproveRequest,
			limit:                10,
			offset:               20,
			shouldCallRepository: true,
			bookmarkData: []*improve_post_storage.Model{
				{
					UserID:    test_utils.NumberUUID(10),
					RequestID: test_utils.NumberUUID(101),
					CreatedAt: baseTime,
					Level:     bookmark_storage.LevelFavorite,
					Target:    improve_post_storage.BookmarkTargetImproveRequest,
				},
				{
					UserID:    test_utils.NumberUUID(10),
					RequestID: test_utils.NumberUUID(111),
					CreatedAt: baseTime,
					Level:     bookmark_storage.LevelFavorite,
					Target:    improve_post_storage.BookmarkTargetImproveRequest,
				},
			},
			bookmarkCount: 123,
			expect: []*models.Bookmark{
				{
					UserID:    test_utils.NumberUUID(10),
					RequestID: test_utils.NumberUUID(101),
					CreatedAt: baseTime,
					Level:     models.BookmarkLevelFavorite,
					Target:    models.BookmarkTargetImproveRequest,
				},
				{
					UserID:    test_utils.NumberUUID(10),
					RequestID: test_utils.NumberUUID(111),
					CreatedAt: baseTime,
					Level:     models.BookmarkLevelFavorite,
					Target:    models.BookmarkTargetImproveRequest,
				},
			},
			expectCount: 123,
		},
		{
			name:                 "Success/NoResults",
			userID:               test_utils.NumberUUID(10),
			level:                models.BookmarkLevelFavorite,
			target:               models.BookmarkTargetImproveRequest,
			limit:                10,
			offset:               20,
			shouldCallRepository: true,
			bookmarkCount:        123,
			expect:               []*models.Bookmark{},
			expectCount:          123,
		},
		{
			name:      "Error/InvalidLevel",
			userID:    test_utils.NumberUUID(10),
			level:     models.BookmarkLevel("foo"),
			target:    models.BookmarkTargetImproveRequest,
			limit:     10,
			offset:    20,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:      "Error/InvalidTarget",
			userID:    test_utils.NumberUUID(10),
			level:     models.BookmarkLevelFavorite,
			target:    models.BookmarkTarget("foo"),
			limit:     10,
			offset:    20,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:      "Error/NoLevel",
			userID:    test_utils.NumberUUID(10),
			target:    models.BookmarkTargetImproveRequest,
			limit:     10,
			offset:    20,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:      "Error/NoTarget",
			userID:    test_utils.NumberUUID(10),
			level:     models.BookmarkLevelFavorite,
			limit:     10,
			offset:    20,
			expectErr: validation.ErrNotAllowed,
		},
		{
			name:                 "Error/RepositoryFailure",
			userID:               test_utils.NumberUUID(10),
			level:                models.BookmarkLevelFavorite,
			target:               models.BookmarkTargetImproveRequest,
			limit:                10,
			offset:               20,
			shouldCallRepository: true,
			bookmarkErr:          fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_post_storage.NewMockRepository(st)

			if d.shouldCallRepository {
				repository.
					On(
						"List", context.TODO(),
						d.userID,
						bookmark_storage.Level(d.level),
						improve_post_storage.BookmarkTarget(d.target),
						d.limit, d.offset,
					).
					Return(d.bookmarkData, d.bookmarkCount, d.bookmarkErr)
			}

			service := NewService(repository)

			res, count, err := service.List(context.TODO(), d.userID, d.level, d.target, d.limit, d.offset)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expectCount, count)
			require.Equal(st, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}
