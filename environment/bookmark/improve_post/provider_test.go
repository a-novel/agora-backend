package improve_post

import (
	"context"
	"crypto/ed25519"
	"errors"
	"github.com/a-novel/agora-backend/domains/bookmark/service/improve_post"
	"github.com/a-novel/agora-backend/domains/keys/service/jwk"
	"github.com/a-novel/agora-backend/domains/keys/storage/jwk"
	"github.com/a-novel/agora-backend/domains/user/service/token"
	user_service "github.com/a-novel/agora-backend/domains/user/service/user"
	"github.com/a-novel/agora-backend/framework"
	"github.com/a-novel/agora-backend/framework/test"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	fooErr   = errors.New("it broken")
)

func TestBookmarkImprovePostProvider_Bookmark(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token     string
		requestID uuid.UUID
		target    models.BookmarkTarget
		level     models.BookmarkLevel

		shouldCallBookmarkService bool
		shouldCallUserService     bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		bookmarkData           *models.Bookmark
		bookmarkErr            error
		hasAuthorization       bool
		hasAuthorizationErr    error

		expect    *models.Bookmark
		expectErr error
	}{
		{
			name:                      "Success",
			now:                       baseTime,
			keys:                      jwk_storage.MockedKeys,
			userID:                    test_utils.NumberUUID(10),
			token:                     "foo.bar.qux",
			requestID:                 test_utils.NumberUUID(100),
			target:                    models.BookmarkTargetImproveRequest,
			level:                     models.BookmarkLevelBookmark,
			shouldCallUserService:     true,
			hasAuthorization:          true,
			shouldCallBookmarkService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			bookmarkData: &models.Bookmark{
				UserID:    test_utils.NumberUUID(10),
				RequestID: test_utils.NumberUUID(100),
				CreatedAt: baseTime,
				Target:    models.BookmarkTargetImproveRequest,
				Level:     models.BookmarkLevelBookmark,
			},
			expect: &models.Bookmark{
				UserID:    test_utils.NumberUUID(10),
				RequestID: test_utils.NumberUUID(100),
				CreatedAt: baseTime,
				Target:    models.BookmarkTargetImproveRequest,
				Level:     models.BookmarkLevelBookmark,
			},
		},
		{
			name:                      "Error/BookmarkServiceFailure",
			now:                       baseTime,
			keys:                      jwk_storage.MockedKeys,
			userID:                    test_utils.NumberUUID(10),
			token:                     "foo.bar.qux",
			requestID:                 test_utils.NumberUUID(100),
			target:                    models.BookmarkTargetImproveRequest,
			level:                     models.BookmarkLevelBookmark,
			shouldCallBookmarkService: true,
			shouldCallUserService:     true,
			hasAuthorization:          true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			bookmarkErr: fooErr,
			expectErr:   fooErr,
		},
		{
			name:                  "Error/UserNotValidated",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			userID:                test_utils.NumberUUID(10),
			token:                 "foo.bar.qux",
			requestID:             test_utils.NumberUUID(100),
			target:                models.BookmarkTargetImproveRequest,
			level:                 models.BookmarkLevelBookmark,
			shouldCallUserService: true,
			hasAuthorization:      false,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			expectErr: validation.ErrUnauthorized,
		},
		{
			name:                  "Error/TokenServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			userID:                test_utils.NumberUUID(10),
			token:                 "foo.bar.qux",
			requestID:             test_utils.NumberUUID(100),
			target:                models.BookmarkTargetImproveRequest,
			level:                 models.BookmarkLevelBookmark,
			tokenServiceDecodeErr: fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			bookmarkService := improve_post_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)
			userService := user_service.NewMockService(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			keysService.
				On("ListPublic").
				Return(publicKeys)

			tokenService.
				On("Decode", d.token, publicKeys, d.now).
				Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)

			if d.shouldCallUserService {
				userService.
					On("HasAuthorizations", context.TODO(), d.tokenServiceDecodeData.Payload.ID, mock.Anything).
					Return(d.hasAuthorization, d.hasAuthorizationErr)
			}

			if d.shouldCallBookmarkService {
				bookmarkService.
					On("Bookmark", context.TODO(), d.userID, d.requestID, d.target, d.level, d.now).
					Return(d.bookmarkData, d.bookmarkErr)
			}

			provider := NewProvider(Config{
				BookmarkService: bookmarkService,
				TokenService:    tokenService,
				KeysService:     keysService,
				UserService:     userService,
				Time:            test_utils.GetTimeNow(d.now),
			})

			res, err := provider.Bookmark(context.TODO(), d.token, d.requestID, d.target, d.level)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			bookmarkService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
			userService.AssertExpectations(t)
		})
	}
}

func TestBookmarkImprovePostProvider_UnBookmark(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token     string
		requestID uuid.UUID
		target    models.BookmarkTarget

		shouldCallBookmarkService bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		bookmarkErr            error

		expectErr error
	}{
		{
			name:                      "Success",
			now:                       baseTime,
			keys:                      jwk_storage.MockedKeys,
			userID:                    test_utils.NumberUUID(10),
			token:                     "foo.bar.qux",
			requestID:                 test_utils.NumberUUID(100),
			target:                    models.BookmarkTargetImproveRequest,
			shouldCallBookmarkService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
		},
		{
			name:                      "Error/BookmarkServiceFailure",
			now:                       baseTime,
			keys:                      jwk_storage.MockedKeys,
			userID:                    test_utils.NumberUUID(10),
			token:                     "foo.bar.qux",
			requestID:                 test_utils.NumberUUID(100),
			target:                    models.BookmarkTargetImproveRequest,
			shouldCallBookmarkService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			bookmarkErr: fooErr,
			expectErr:   fooErr,
		},
		{
			name:                  "Error/TokenServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			userID:                test_utils.NumberUUID(10),
			token:                 "foo.bar.qux",
			requestID:             test_utils.NumberUUID(100),
			target:                models.BookmarkTargetImproveRequest,
			tokenServiceDecodeErr: fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			bookmarkService := improve_post_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			keysService.
				On("ListPublic").
				Return(publicKeys)

			tokenService.
				On("Decode", d.token, publicKeys, d.now).
				Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)

			if d.shouldCallBookmarkService {
				bookmarkService.
					On(
						"UnBookmark", context.TODO(), d.userID, d.requestID,
						models.BookmarkTarget(d.target),
					).
					Return(d.bookmarkErr)
			}

			provider := NewProvider(Config{
				BookmarkService: bookmarkService,
				TokenService:    tokenService,
				KeysService:     keysService,
				Time:            test_utils.GetTimeNow(d.now),
			})

			err := provider.UnBookmark(context.TODO(), d.token, d.requestID, d.target)
			test_utils.RequireError(t, d.expectErr, err)

			bookmarkService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestBookmarkImprovePostProvider_IsBookmarked(t *testing.T) {
	data := []struct {
		name string

		userID    uuid.UUID
		requestID uuid.UUID
		target    models.BookmarkTarget

		bookmarkData *models.BookmarkLevel
		bookmarkErr  error

		expect    *models.BookmarkLevel
		expectErr error
	}{
		{
			name:         "Success",
			userID:       test_utils.NumberUUID(10),
			requestID:    test_utils.NumberUUID(100),
			target:       models.BookmarkTargetImproveRequest,
			bookmarkData: framework.ToPTR(models.BookmarkLevelBookmark),
			expect:       framework.ToPTR(models.BookmarkLevelBookmark),
		},
		{
			name:        "Error/BookmarkServiceFailure",
			userID:      test_utils.NumberUUID(10),
			requestID:   test_utils.NumberUUID(100),
			target:      models.BookmarkTargetImproveRequest,
			bookmarkErr: fooErr,
			expectErr:   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			bookmarkService := improve_post_service.NewMockService(t)

			bookmarkService.
				On(
					"IsBookmarked", context.TODO(), d.userID, d.requestID,
					models.BookmarkTarget(d.target),
				).
				Return(d.bookmarkData, d.bookmarkErr)
			provider := NewProvider(Config{
				BookmarkService: bookmarkService,
			})

			res, err := provider.IsBookmarked(context.TODO(), d.userID, d.requestID, d.target)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			bookmarkService.AssertExpectations(t)
		})
	}
}

func TestBookmarkImprovePostProvider_List(t *testing.T) {
	data := []struct {
		name string

		userID uuid.UUID
		level  models.BookmarkLevel
		target models.BookmarkTarget
		limit  int
		offset int

		bookmarkData  []*models.Bookmark
		bookmarkCount int64
		bookmarkErr   error

		expect      []*models.Bookmark
		expectCount int64
		expectErr   error
	}{
		{
			name:   "Success",
			userID: test_utils.NumberUUID(10),
			target: models.BookmarkTargetImproveRequest,
			bookmarkData: []*models.Bookmark{
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
		},
		{
			name:        "Error/BookmarkServiceFailure",
			userID:      test_utils.NumberUUID(10),
			target:      models.BookmarkTargetImproveRequest,
			bookmarkErr: fooErr,
			expectErr:   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			bookmarkService := improve_post_service.NewMockService(t)

			bookmarkService.
				On(
					"List", context.TODO(), d.userID,
					models.BookmarkLevel(d.level),
					models.BookmarkTarget(d.target),
					d.limit, d.offset,
				).
				Return(d.bookmarkData, d.bookmarkCount, d.bookmarkErr)

			provider := NewProvider(Config{
				BookmarkService: bookmarkService,
			})

			res, count, err := provider.List(context.TODO(), d.userID, d.level, d.target, d.limit, d.offset)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)
			require.Equal(t, d.expectCount, count)

			bookmarkService.AssertExpectations(t)
		})
	}
}
