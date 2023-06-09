package improve_post

import (
	"context"
	"crypto/ed25519"
	"errors"
	"github.com/a-novel/agora-backend/domains/forum/service/improve_request"
	"github.com/a-novel/agora-backend/domains/forum/service/improve_suggestion"
	"github.com/a-novel/agora-backend/domains/forum/service/votes"
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
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 9, 0, 0, 0, time.UTC)
	fooErr     = errors.New("it broken")
)

func TestImprovePostProvider_ReadImproveRequest(t *testing.T) {
	data := []struct {
		name string

		id uuid.UUID

		serviceData []*models.ImproveRequest
		serviceErr  error

		expect    []*models.ImproveRequest
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1),
			serviceData: []*models.ImproveRequest{
				{
					ID:        test_utils.NumberUUID(1),
					Source:    test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UserID:    test_utils.NumberUUID(10),
					Title:     "Dummy request.",
					Content:   "Foo bar qux.",
					UpVotes:   17,
					DownVotes: 3,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Source:    test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UserID:    test_utils.NumberUUID(10),
					Title:     "Smart request.",
					Content:   "Qux bar foo.",
					UpVotes:   12,
				},
			},
			expect: []*models.ImproveRequest{
				{
					ID:        test_utils.NumberUUID(1),
					Source:    test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UserID:    test_utils.NumberUUID(10),
					Title:     "Dummy request.",
					Content:   "Foo bar qux.",
					UpVotes:   17,
					DownVotes: 3,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Source:    test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UserID:    test_utils.NumberUUID(10),
					Title:     "Smart request.",
					Content:   "Qux bar foo.",
					UpVotes:   12,
				},
			},
		},
		{
			name:       "Error/ServiceFailure",
			id:         test_utils.NumberUUID(1),
			serviceErr: fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveRequestService := improve_request_service.NewMockService(t)

			improveRequestService.
				On("ReadRevisions", context.TODO(), d.id).
				Return(d.serviceData, d.serviceErr)

			provider := NewProvider(Config{
				ImproveRequestService: improveRequestService,
			})

			res, err := provider.ReadImproveRequest(context.TODO(), d.id)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			improveRequestService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_CreateImproveRequest(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		id     uuid.UUID
		userID uuid.UUID

		token   string
		title   string
		content string

		shouldCallImproveRequestService bool
		shouldCallUserService           bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		improveRequestData     *models.ImproveRequest
		improveRequestErr      error
		hasAuthorization       bool
		hasAuthorizationErr    error

		expect    *models.ImproveRequest
		expectErr error
	}{
		{
			name:                            "Success",
			now:                             baseTime,
			keys:                            jwk_storage.MockedKeys,
			id:                              test_utils.NumberUUID(1),
			userID:                          test_utils.NumberUUID(10),
			token:                           "foo.bar.qux",
			title:                           "Dummy request",
			content:                         "Foo bar qux.",
			shouldCallUserService:           true,
			hasAuthorization:                true,
			shouldCallImproveRequestService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveRequestData: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Dummy request",
				Content:   "Foo bar qux.",
				UpVotes:   0,
				DownVotes: 0,
			},
			expect: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Dummy request",
				Content:   "Foo bar qux.",
				UpVotes:   0,
				DownVotes: 0,
			},
		},
		{
			name:                            "Error/ImproveRequestServiceFailure",
			now:                             baseTime,
			keys:                            jwk_storage.MockedKeys,
			id:                              test_utils.NumberUUID(1),
			userID:                          test_utils.NumberUUID(10),
			token:                           "foo.bar.qux",
			title:                           "Dummy request",
			content:                         "Foo bar qux.",
			shouldCallImproveRequestService: true,
			shouldCallUserService:           true,
			hasAuthorization:                true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveRequestErr: fooErr,
			expectErr:         fooErr,
		},
		{
			name:                  "Error/UserNotValidated",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			id:                    test_utils.NumberUUID(1),
			userID:                test_utils.NumberUUID(10),
			token:                 "foo.bar.qux",
			title:                 "Dummy request",
			content:               "Foo bar qux.",
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
			id:                    test_utils.NumberUUID(1),
			userID:                test_utils.NumberUUID(10),
			token:                 "foo.bar.qux",
			title:                 "Dummy request",
			content:               "Foo bar qux.",
			tokenServiceDecodeErr: fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveRequestService := improve_request_service.NewMockService(t)

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

			if d.shouldCallImproveRequestService {
				improveRequestService.
					On("Create", context.TODO(), d.userID, d.title, d.content, d.id, d.now).
					Return(d.improveRequestData, d.improveRequestErr)
			}

			provider := NewProvider(Config{
				ImproveRequestService: improveRequestService,
				TokenService:          tokenService,
				KeysService:           keysService,
				UserService:           userService,
				Time:                  test_utils.GetTimeNow(d.now),
				ID:                    test_utils.GetUUID(d.id),
			})

			res, err := provider.CreateImproveRequest(context.TODO(), d.token, d.title, d.content)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			improveRequestService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_CreateImproveRequestRevision(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		id     uuid.UUID
		userID uuid.UUID

		token    string
		title    string
		content  string
		sourceID uuid.UUID

		shouldCallImproveRequestGetService    bool
		shouldCallImproveRequestCreateService bool
		shouldCallUserService                 bool

		tokenServiceDecodeData   *models.UserToken
		tokenServiceDecodeErr    error
		improveRequestGetData    *models.ImproveRequest
		improveRequestGetErr     error
		improveRequestCreateData *models.ImproveRequest
		improveRequestCreateErr  error
		hasAuthorization         bool
		hasAuthorizationErr      error

		expect    *models.ImproveRequest
		expectErr error
	}{
		{
			name:                                  "Success",
			now:                                   baseTime,
			keys:                                  jwk_storage.MockedKeys,
			id:                                    test_utils.NumberUUID(2),
			userID:                                test_utils.NumberUUID(10),
			token:                                 "foo.bar.qux",
			title:                                 "Smart request",
			content:                               "Qux bar foo.",
			sourceID:                              test_utils.NumberUUID(1),
			shouldCallImproveRequestGetService:    true,
			shouldCallImproveRequestCreateService: true,
			shouldCallUserService:                 true,
			hasAuthorization:                      true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveRequestGetData: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Dummy request",
				Content:   "Foo bar qux.",
				UpVotes:   10,
				DownVotes: 2,
			},
			improveRequestCreateData: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Smart request",
				Content:   "Qux bar foo.",
			},
			expect: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Smart request",
				Content:   "Qux bar foo.",
			},
		},
		{
			name:                               "Error/UnauthorizedUser",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			id:                                 test_utils.NumberUUID(2),
			userID:                             test_utils.NumberUUID(11),
			token:                              "foo.bar.qux",
			title:                              "Smart request",
			content:                            "Qux bar foo.",
			sourceID:                           test_utils.NumberUUID(1),
			shouldCallImproveRequestGetService: true,
			shouldCallUserService:              true,
			hasAuthorization:                   true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(11)},
			},
			improveRequestGetData: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Dummy request",
				Content:   "Foo bar qux.",
				UpVotes:   10,
				DownVotes: 2,
			},
			expectErr: validation.ErrInvalidCredentials,
		},
		{
			name:                                  "Error/ImproveRequestServiceCreateFailure",
			now:                                   baseTime,
			keys:                                  jwk_storage.MockedKeys,
			id:                                    test_utils.NumberUUID(2),
			userID:                                test_utils.NumberUUID(10),
			token:                                 "foo.bar.qux",
			title:                                 "Smart request",
			content:                               "Qux bar foo.",
			sourceID:                              test_utils.NumberUUID(1),
			shouldCallImproveRequestGetService:    true,
			shouldCallImproveRequestCreateService: true,
			shouldCallUserService:                 true,
			hasAuthorization:                      true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveRequestGetData: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Dummy request",
				Content:   "Foo bar qux.",
				UpVotes:   10,
				DownVotes: 2,
			},
			improveRequestCreateErr: fooErr,
			expectErr:               fooErr,
		},
		{
			name:                               "Error/ImproveRequestServiceGetFailure",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			id:                                 test_utils.NumberUUID(2),
			userID:                             test_utils.NumberUUID(10),
			token:                              "foo.bar.qux",
			title:                              "Smart request",
			content:                            "Qux bar foo.",
			sourceID:                           test_utils.NumberUUID(1),
			shouldCallImproveRequestGetService: true,
			shouldCallUserService:              true,
			hasAuthorization:                   true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveRequestGetErr: fooErr,
			expectErr:            fooErr,
		},
		{
			name:                  "Error/UserNotValidated",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			id:                    test_utils.NumberUUID(2),
			userID:                test_utils.NumberUUID(10),
			token:                 "foo.bar.qux",
			title:                 "Smart request",
			content:               "Qux bar foo.",
			sourceID:              test_utils.NumberUUID(1),
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
			id:                    test_utils.NumberUUID(2),
			userID:                test_utils.NumberUUID(10),
			token:                 "foo.bar.qux",
			title:                 "Smart request",
			content:               "Qux bar foo.",
			sourceID:              test_utils.NumberUUID(1),
			tokenServiceDecodeErr: fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveRequestService := improve_request_service.NewMockService(t)

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

			if d.shouldCallImproveRequestGetService {
				improveRequestService.
					On("Read", context.TODO(), d.sourceID).
					Return(d.improveRequestGetData, d.improveRequestGetErr)
			}

			if d.shouldCallImproveRequestCreateService {
				improveRequestService.
					On("CreateRevision", context.TODO(), d.userID, d.sourceID, d.title, d.content, d.id, d.now).
					Return(d.improveRequestCreateData, d.improveRequestCreateErr)
			}

			provider := NewProvider(Config{
				ImproveRequestService: improveRequestService,
				TokenService:          tokenService,
				KeysService:           keysService,
				UserService:           userService,
				Time:                  test_utils.GetTimeNow(d.now),
				ID:                    test_utils.GetUUID(d.id),
			})

			res, err := provider.CreateImproveRequestRevision(context.TODO(), d.token, d.sourceID, d.title, d.content)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			improveRequestService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_DeleteImproveRequest(t *testing.T) {
	data := []struct {
		name string

		now  time.Time
		keys []ed25519.PrivateKey

		token     string
		title     string
		content   string
		requestID uuid.UUID

		shouldCallImproveRequestGetService    bool
		shouldCallImproveRequestDeleteService bool

		tokenServiceDecodeData  *models.UserToken
		tokenServiceDecodeErr   error
		improveRequestGetData   *models.ImproveRequest
		improveRequestGetErr    error
		improveRequestDeleteErr error

		expectErr error
	}{
		{
			name:                                  "Success",
			now:                                   baseTime,
			keys:                                  jwk_storage.MockedKeys,
			token:                                 "foo.bar.qux",
			title:                                 "Smart request",
			content:                               "Qux bar foo.",
			requestID:                             test_utils.NumberUUID(1),
			shouldCallImproveRequestGetService:    true,
			shouldCallImproveRequestDeleteService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveRequestGetData: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Dummy request",
				Content:   "Foo bar qux.",
				UpVotes:   10,
				DownVotes: 2,
			},
		},
		{
			name:                               "Error/UnauthorizedUser",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			token:                              "foo.bar.qux",
			title:                              "Smart request",
			content:                            "Qux bar foo.",
			requestID:                          test_utils.NumberUUID(1),
			shouldCallImproveRequestGetService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(11)},
			},
			improveRequestGetData: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Dummy request",
				Content:   "Foo bar qux.",
				UpVotes:   10,
				DownVotes: 2,
			},
			expectErr: validation.ErrInvalidCredentials,
		},
		{
			name:                                  "Error/ImproveRequestServiceDeleteFailure",
			now:                                   baseTime,
			keys:                                  jwk_storage.MockedKeys,
			token:                                 "foo.bar.qux",
			title:                                 "Smart request",
			content:                               "Qux bar foo.",
			requestID:                             test_utils.NumberUUID(1),
			shouldCallImproveRequestGetService:    true,
			shouldCallImproveRequestDeleteService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveRequestGetData: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(10),
				Title:     "Dummy request",
				Content:   "Foo bar qux.",
				UpVotes:   10,
				DownVotes: 2,
			},
			improveRequestDeleteErr: fooErr,
			expectErr:               fooErr,
		},
		{
			name:                               "Error/ImproveRequestServiceGetFailure",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			token:                              "foo.bar.qux",
			title:                              "Smart request",
			content:                            "Qux bar foo.",
			requestID:                          test_utils.NumberUUID(1),
			shouldCallImproveRequestGetService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveRequestGetErr: fooErr,
			expectErr:            fooErr,
		},
		{
			name:                  "Error/TokenServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			token:                 "foo.bar.qux",
			title:                 "Smart request",
			content:               "Qux bar foo.",
			requestID:             test_utils.NumberUUID(1),
			tokenServiceDecodeErr: fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveRequestService := improve_request_service.NewMockService(t)

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

			if d.shouldCallImproveRequestGetService {
				improveRequestService.
					On("Read", context.TODO(), d.requestID).
					Return(d.improveRequestGetData, d.improveRequestGetErr)
			}

			if d.shouldCallImproveRequestDeleteService {
				improveRequestService.
					On("Delete", context.TODO(), d.requestID).
					Return(d.improveRequestDeleteErr)
			}

			provider := NewProvider(Config{
				ImproveRequestService: improveRequestService,
				TokenService:          tokenService,
				KeysService:           keysService,
				Time:                  test_utils.GetTimeNow(d.now),
			})

			err := provider.DeleteImproveRequest(context.TODO(), d.token, d.requestID)
			test_utils.RequireError(t, d.expectErr, err)

			improveRequestService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_SearchImproveRequests(t *testing.T) {
	data := []struct {
		name string

		query  models.ImproveRequestSearch
		limit  int
		offset int

		serviceData  []*models.ImproveRequestPreview
		serviceTotal int64
		serviceErr   error

		expect      []*models.ImproveRequestPreview
		expectTotal int64
		expectErr   error
	}{
		{
			name: "Success",
			query: models.ImproveRequestSearch{
				UserID: framework.ToPTR(test_utils.NumberUUID(1)),
				Query:  "foo",
			},
			limit:  10,
			offset: 20,
			serviceData: []*models.ImproveRequestPreview{
				{
					ID:                  test_utils.NumberUUID(1),
					Source:              test_utils.NumberUUID(1),
					UserID:              test_utils.NumberUUID(1),
					CreatedAt:           baseTime,
					Title:               "Dummy request",
					Content:             "Foo bar qux.",
					UpVotes:             10,
					DownVotes:           2,
					RevisionCount:       3,
					MoreRecentRevisions: 1,
				},
				{
					ID:            test_utils.NumberUUID(3),
					Source:        test_utils.NumberUUID(2),
					UserID:        test_utils.NumberUUID(1),
					CreatedAt:     baseTime,
					Title:         "Smart request",
					Content:       "Qux bar foo.",
					UpVotes:       17,
					DownVotes:     1,
					RevisionCount: 1,
				},
			},
			serviceTotal: 20,
			expect: []*models.ImproveRequestPreview{
				{
					ID:                  test_utils.NumberUUID(1),
					Source:              test_utils.NumberUUID(1),
					UserID:              test_utils.NumberUUID(1),
					CreatedAt:           baseTime,
					Title:               "Dummy request",
					Content:             "Foo bar qux.",
					UpVotes:             10,
					DownVotes:           2,
					RevisionCount:       3,
					MoreRecentRevisions: 1,
				},
				{
					ID:            test_utils.NumberUUID(3),
					Source:        test_utils.NumberUUID(2),
					UserID:        test_utils.NumberUUID(1),
					CreatedAt:     baseTime,
					Title:         "Smart request",
					Content:       "Qux bar foo.",
					UpVotes:       17,
					DownVotes:     1,
					RevisionCount: 1,
				},
			},
			expectTotal: 20,
		},
		{
			name: "Error/ServiceFailure",
			query: models.ImproveRequestSearch{
				UserID: framework.ToPTR(test_utils.NumberUUID(1)),
				Query:  "foo",
			},
			limit:      10,
			offset:     20,
			serviceErr: fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveRequestService := improve_request_service.NewMockService(t)

			improveRequestService.
				On("Search", context.TODO(), models.ImproveRequestSearch{
					UserID: d.query.UserID,
					Query:  d.query.Query,
				}, d.limit, d.offset).
				Return(d.serviceData, d.serviceTotal, d.serviceErr)

			provider := NewProvider(Config{
				ImproveRequestService: improveRequestService,
			})

			res, total, err := provider.SearchImproveRequests(context.TODO(), d.query, d.limit, d.offset)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)
			require.Equal(t, d.expectTotal, total)

			improveRequestService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_GetImproveRequestPreviews(t *testing.T) {
	data := []struct {
		name string

		ids []uuid.UUID

		serviceData []*models.ImproveRequestPreview
		serviceErr  error

		expect    []*models.ImproveRequestPreview
		expectErr error
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			serviceData: []*models.ImproveRequestPreview{
				{
					ID:                  test_utils.NumberUUID(1),
					Source:              test_utils.NumberUUID(2),
					CreatedAt:           baseTime,
					UserID:              test_utils.NumberUUID(1),
					Title:               "Dummy post",
					Content:             "Foo bar qux.",
					UpVotes:             10,
					DownVotes:           5,
					MoreRecentRevisions: 1,
					RevisionCount:       10,
				},
				{
					ID:                  test_utils.NumberUUID(42),
					Source:              test_utils.NumberUUID(666),
					CreatedAt:           baseTime,
					UserID:              test_utils.NumberUUID(12),
					Title:               "Smart post",
					Content:             "Cats taking a nap.",
					UpVotes:             4,
					DownVotes:           1,
					MoreRecentRevisions: 3,
					RevisionCount:       4,
				},
			},
			expect: []*models.ImproveRequestPreview{
				{
					ID:                  test_utils.NumberUUID(1),
					Source:              test_utils.NumberUUID(2),
					CreatedAt:           baseTime,
					UserID:              test_utils.NumberUUID(1),
					Title:               "Dummy post",
					Content:             "Foo bar qux.",
					UpVotes:             10,
					DownVotes:           5,
					MoreRecentRevisions: 1,
					RevisionCount:       10,
				},
				{
					ID:                  test_utils.NumberUUID(42),
					Source:              test_utils.NumberUUID(666),
					CreatedAt:           baseTime,
					UserID:              test_utils.NumberUUID(12),
					Title:               "Smart post",
					Content:             "Cats taking a nap.",
					UpVotes:             4,
					DownVotes:           1,
					MoreRecentRevisions: 3,
					RevisionCount:       4,
				},
			},
		},
		{
			name: "Error/ServiceFailure",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			serviceErr: fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveRequestService := improve_request_service.NewMockService(t)

			improveRequestService.
				On("GetPreviews", context.TODO(), d.ids).
				Return(d.serviceData, d.serviceErr)

			provider := NewProvider(Config{
				ImproveRequestService: improveRequestService,
			})

			res, err := provider.GetImproveRequestPreviews(context.TODO(), d.ids)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			improveRequestService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_ReadImproveSuggestion(t *testing.T) {
	data := []struct {
		name string

		id uuid.UUID

		serviceData *models.ImproveSuggestion
		serviceErr  error

		expect    *models.ImproveSuggestion
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1),
			serviceData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
			expect: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
		},
		{
			name:       "Error/ServiceFailure",
			id:         test_utils.NumberUUID(1),
			serviceErr: fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveSuggestionService := improve_suggestion_service.NewMockService(t)

			improveSuggestionService.
				On("Read", context.TODO(), d.id).
				Return(d.serviceData, d.serviceErr)

			provider := NewProvider(Config{
				ImproveSuggestionService: improveSuggestionService,
			})

			res, err := provider.ReadImproveSuggestion(context.TODO(), d.id)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			improveSuggestionService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_CreateImproveSuggestion(t *testing.T) {
	data := []struct {
		name string

		now       time.Time
		keys      []ed25519.PrivateKey
		id        uuid.UUID
		userID    uuid.UUID
		requestID uuid.UUID
		sourceID  uuid.UUID

		token   string
		title   string
		content string

		shouldCallImproveSuggestionService bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		improveSuggestionData  *models.ImproveSuggestion
		improveSuggestionErr   error

		expect    *models.ImproveSuggestion
		expectErr error
	}{
		{
			name:                               "Success",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			id:                                 test_utils.NumberUUID(1),
			userID:                             test_utils.NumberUUID(10),
			token:                              "foo.bar.qux",
			title:                              "Dummy request",
			content:                            "Foo bar qux.",
			shouldCallImproveSuggestionService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveSuggestionData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
			expect: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
		},
		{
			name:                               "Error/ImproveSuggestionServiceFailure",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			id:                                 test_utils.NumberUUID(1),
			userID:                             test_utils.NumberUUID(10),
			token:                              "foo.bar.qux",
			title:                              "Dummy request",
			content:                            "Foo bar qux.",
			shouldCallImproveSuggestionService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveSuggestionErr: fooErr,
			expectErr:            fooErr,
		},
		{
			name:                  "Error/TokenServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			id:                    test_utils.NumberUUID(1),
			userID:                test_utils.NumberUUID(10),
			token:                 "foo.bar.qux",
			title:                 "Dummy request",
			content:               "Foo bar qux.",
			tokenServiceDecodeErr: fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveSuggestionService := improve_suggestion_service.NewMockService(t)

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

			if d.shouldCallImproveSuggestionService {
				improveSuggestionService.
					On("Create", context.TODO(), &models.ImproveSuggestionUpsert{
						RequestID: d.requestID,
						Title:     d.title,
						Content:   d.content,
					}, d.userID, d.sourceID, d.id, d.now).
					Return(d.improveSuggestionData, d.improveSuggestionErr)
			}

			provider := NewProvider(Config{
				ImproveSuggestionService: improveSuggestionService,
				TokenService:             tokenService,
				KeysService:              keysService,
				Time:                     test_utils.GetTimeNow(d.now),
				ID:                       test_utils.GetUUID(d.id),
			})

			res, err := provider.CreateImproveSuggestion(context.TODO(), d.token, d.requestID, d.sourceID, d.title, d.content)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			improveSuggestionService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_UpdateImproveSuggestion(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token     string
		title     string
		content   string
		postID    uuid.UUID
		requestID uuid.UUID

		shouldCallImproveSuggestionGetService    bool
		shouldCallImproveSuggestionUpdateService bool

		tokenServiceDecodeData      *models.UserToken
		tokenServiceDecodeErr       error
		improveSuggestionGetData    *models.ImproveSuggestion
		improveSuggestionGetErr     error
		improveSuggestionUpdateData *models.ImproveSuggestion
		improveSuggestionUpdateErr  error

		expect    *models.ImproveSuggestion
		expectErr error
	}{
		{
			name:                                     "Success",
			now:                                      baseTime,
			keys:                                     jwk_storage.MockedKeys,
			postID:                                   test_utils.NumberUUID(1),
			requestID:                                test_utils.NumberUUID(10),
			userID:                                   test_utils.NumberUUID(100),
			token:                                    "foo.bar.qux",
			title:                                    "Smart request",
			content:                                  "Qux bar foo.",
			shouldCallImproveSuggestionGetService:    true,
			shouldCallImproveSuggestionUpdateService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveSuggestionGetData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
			improveSuggestionUpdateData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Smart request",
				Content:   "Qux bar foo.",
			},
			expect: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Smart request",
				Content:   "Qux bar foo.",
			},
		},
		{
			name:                                  "Error/UnauthorizedUser",
			now:                                   baseTime,
			keys:                                  jwk_storage.MockedKeys,
			token:                                 "foo.bar.qux",
			title:                                 "Smart request",
			content:                               "Qux bar foo.",
			postID:                                test_utils.NumberUUID(1),
			requestID:                             test_utils.NumberUUID(10),
			userID:                                test_utils.NumberUUID(101),
			shouldCallImproveSuggestionGetService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(101)},
			},
			improveSuggestionGetData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
			expectErr: validation.ErrInvalidCredentials,
		},
		{
			name:                                     "Error/ImproveRequestServiceCreateFailure",
			now:                                      baseTime,
			keys:                                     jwk_storage.MockedKeys,
			token:                                    "foo.bar.qux",
			title:                                    "Smart request",
			content:                                  "Qux bar foo.",
			postID:                                   test_utils.NumberUUID(1),
			requestID:                                test_utils.NumberUUID(10),
			userID:                                   test_utils.NumberUUID(100),
			shouldCallImproveSuggestionGetService:    true,
			shouldCallImproveSuggestionUpdateService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveSuggestionGetData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
			improveSuggestionUpdateErr: fooErr,
			expectErr:                  fooErr,
		},
		{
			name:                                  "Error/ImproveRequestServiceGetFailure",
			now:                                   baseTime,
			keys:                                  jwk_storage.MockedKeys,
			token:                                 "foo.bar.qux",
			title:                                 "Smart request",
			content:                               "Qux bar foo.",
			postID:                                test_utils.NumberUUID(1),
			requestID:                             test_utils.NumberUUID(10),
			userID:                                test_utils.NumberUUID(100),
			shouldCallImproveSuggestionGetService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveSuggestionGetErr: fooErr,
			expectErr:               fooErr,
		},
		{
			name:                  "Error/TokenServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			token:                 "foo.bar.qux",
			title:                 "Smart request",
			content:               "Qux bar foo.",
			postID:                test_utils.NumberUUID(1),
			requestID:             test_utils.NumberUUID(10),
			userID:                test_utils.NumberUUID(100),
			tokenServiceDecodeErr: fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveSuggestionService := improve_suggestion_service.NewMockService(t)

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

			if d.shouldCallImproveSuggestionGetService {
				improveSuggestionService.
					On("Read", context.TODO(), d.postID).
					Return(d.improveSuggestionGetData, d.improveSuggestionGetErr)
			}

			if d.shouldCallImproveSuggestionUpdateService {
				improveSuggestionService.
					On("Update", context.TODO(), &models.ImproveSuggestionUpsert{
						RequestID: d.requestID,
						Title:     d.title,
						Content:   d.content,
					}, d.postID, d.now).
					Return(d.improveSuggestionUpdateData, d.improveSuggestionUpdateErr)
			}

			provider := NewProvider(Config{
				ImproveSuggestionService: improveSuggestionService,
				TokenService:             tokenService,
				KeysService:              keysService,
				Time:                     test_utils.GetTimeNow(d.now),
			})

			res, err := provider.UpdateImproveSuggestion(context.TODO(), d.token, d.postID, d.requestID, d.title, d.content)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			improveSuggestionService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_DeleteImproveSuggestion(t *testing.T) {
	data := []struct {
		name string

		now  time.Time
		keys []ed25519.PrivateKey

		token     string
		title     string
		content   string
		requestID uuid.UUID

		shouldCallImproveSuggestionGetService    bool
		shouldCallImproveSuggestionDeleteService bool

		tokenServiceDecodeData     *models.UserToken
		tokenServiceDecodeErr      error
		improveSuggestionGetData   *models.ImproveSuggestion
		improveSuggestionGetErr    error
		improveSuggestionDeleteErr error

		expectErr error
	}{
		{
			name:                                     "Success",
			now:                                      baseTime,
			keys:                                     jwk_storage.MockedKeys,
			token:                                    "foo.bar.qux",
			title:                                    "Smart request",
			content:                                  "Qux bar foo.",
			requestID:                                test_utils.NumberUUID(1),
			shouldCallImproveSuggestionGetService:    true,
			shouldCallImproveSuggestionDeleteService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveSuggestionGetData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
		},
		{
			name:                                  "Error/UnauthorizedUser",
			now:                                   baseTime,
			keys:                                  jwk_storage.MockedKeys,
			token:                                 "foo.bar.qux",
			title:                                 "Smart request",
			content:                               "Qux bar foo.",
			requestID:                             test_utils.NumberUUID(1),
			shouldCallImproveSuggestionGetService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(101)},
			},
			improveSuggestionGetData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
			expectErr: validation.ErrInvalidCredentials,
		},
		{
			name:                                     "Error/ImproveSuggestionServiceDeleteFailure",
			now:                                      baseTime,
			keys:                                     jwk_storage.MockedKeys,
			token:                                    "foo.bar.qux",
			title:                                    "Smart request",
			content:                                  "Qux bar foo.",
			requestID:                                test_utils.NumberUUID(1),
			shouldCallImproveSuggestionGetService:    true,
			shouldCallImproveSuggestionDeleteService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveSuggestionGetData: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   32,
				DownVotes: 2,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy suggestion",
				Content:   "Foo bar qux.",
			},
			improveSuggestionDeleteErr: fooErr,
			expectErr:                  fooErr,
		},
		{
			name:                                  "Error/ImproveSuggestionServiceGetFailure",
			now:                                   baseTime,
			keys:                                  jwk_storage.MockedKeys,
			token:                                 "foo.bar.qux",
			title:                                 "Smart request",
			content:                               "Qux bar foo.",
			requestID:                             test_utils.NumberUUID(1),
			shouldCallImproveSuggestionGetService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(10)},
			},
			improveSuggestionGetErr: fooErr,
			expectErr:               fooErr,
		},
		{
			name:                  "Error/TokenServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			token:                 "foo.bar.qux",
			title:                 "Smart request",
			content:               "Qux bar foo.",
			requestID:             test_utils.NumberUUID(1),
			tokenServiceDecodeErr: fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveSuggestionService := improve_suggestion_service.NewMockService(t)

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

			if d.shouldCallImproveSuggestionGetService {
				improveSuggestionService.
					On("Read", context.TODO(), d.requestID).
					Return(d.improveSuggestionGetData, d.improveSuggestionGetErr)
			}

			if d.shouldCallImproveSuggestionDeleteService {
				improveSuggestionService.
					On("Delete", context.TODO(), d.requestID).
					Return(d.improveSuggestionDeleteErr)
			}

			provider := NewProvider(Config{
				ImproveSuggestionService: improveSuggestionService,
				TokenService:             tokenService,
				KeysService:              keysService,
				Time:                     test_utils.GetTimeNow(d.now),
			})

			err := provider.DeleteImproveSuggestion(context.TODO(), d.token, d.requestID)
			test_utils.RequireError(t, d.expectErr, err)

			improveSuggestionService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_ListImproveSuggestions(t *testing.T) {
	data := []struct {
		name string

		query  models.ImproveSuggestionsList
		limit  int
		offset int

		improveSuggestionData  []*models.ImproveSuggestion
		improveSuggestionTotal int64
		improveSuggestionErr   error

		expectData  []*models.ImproveSuggestion
		expectTotal int64
		expectErr   error
	}{
		{
			name: "Success",
			query: models.ImproveSuggestionsList{
				UserID:    framework.ToPTR(test_utils.NumberUUID(100)),
				SourceID:  framework.ToPTR(test_utils.NumberUUID(10)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(1)),
				Validated: framework.ToPTR(true),
			},
			limit:  10,
			offset: 20,
			improveSuggestionData: []*models.ImproveSuggestion{
				{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
					SourceID:  test_utils.NumberUUID(10),
					UserID:    test_utils.NumberUUID(100),
					Validated: true,
					UpVotes:   32,
					DownVotes: 2,
					RequestID: test_utils.NumberUUID(11),
					Title:     "Dummy suggestion",
					Content:   "Foo bar qux.",
				},
				{
					ID:        test_utils.NumberUUID(2),
					CreatedAt: baseTime,
					UpdatedAt: framework.ToPTR(baseTime.Add(time.Minute)),
					SourceID:  test_utils.NumberUUID(11),
					UserID:    test_utils.NumberUUID(101),
					Validated: true,
					UpVotes:   10,
					DownVotes: 8,
					RequestID: test_utils.NumberUUID(12),
					Title:     "Smart suggestion",
					Content:   "Qux bar foo.",
				},
			},
			improveSuggestionTotal: 200,
			expectData: []*models.ImproveSuggestion{
				{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
					SourceID:  test_utils.NumberUUID(10),
					UserID:    test_utils.NumberUUID(100),
					Validated: true,
					UpVotes:   32,
					DownVotes: 2,
					RequestID: test_utils.NumberUUID(11),
					Title:     "Dummy suggestion",
					Content:   "Foo bar qux.",
				},
				{
					ID:        test_utils.NumberUUID(2),
					CreatedAt: baseTime,
					UpdatedAt: framework.ToPTR(baseTime.Add(time.Minute)),
					SourceID:  test_utils.NumberUUID(11),
					UserID:    test_utils.NumberUUID(101),
					Validated: true,
					UpVotes:   10,
					DownVotes: 8,
					RequestID: test_utils.NumberUUID(12),
					Title:     "Smart suggestion",
					Content:   "Qux bar foo.",
				},
			},
			expectTotal: 200,
		},
		{
			name: "Error/ServiceFailure",
			query: models.ImproveSuggestionsList{
				UserID:    framework.ToPTR(test_utils.NumberUUID(100)),
				SourceID:  framework.ToPTR(test_utils.NumberUUID(10)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(1)),
				Validated: framework.ToPTR(true),
			},
			limit:                10,
			offset:               20,
			improveSuggestionErr: fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveSuggestionService := improve_suggestion_service.NewMockService(t)

			improveSuggestionService.
				On("List", context.TODO(), models.ImproveSuggestionsList{
					SourceID:  d.query.SourceID,
					UserID:    d.query.UserID,
					RequestID: d.query.RequestID,
					Validated: d.query.Validated,
				}, d.limit, d.offset).
				Return(d.improveSuggestionData, d.improveSuggestionTotal, d.improveSuggestionErr)

			provider := NewProvider(Config{
				ImproveSuggestionService: improveSuggestionService,
			})

			data, total, err := provider.ListImproveSuggestions(context.TODO(), d.query, d.limit, d.offset)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expectData, data)
			require.Equal(t, d.expectTotal, total)

			improveSuggestionService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_GetImproveSuggestionPreviews(t *testing.T) {
	data := []struct {
		name string

		ids []uuid.UUID

		serviceData []*models.ImproveSuggestion
		serviceErr  error

		expect    []*models.ImproveSuggestion
		expectErr error
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			serviceData: []*models.ImproveSuggestion{
				{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(10),
					UserID:    test_utils.NumberUUID(100),
					Validated: false,
					UpVotes:   17,
					DownVotes: 3,
					RequestID: test_utils.NumberUUID(11),
					Title:     "Dummy post",
					Content:   "Foo bar qux.",
				},
				{
					ID:        test_utils.NumberUUID(2),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(12),
					UserID:    test_utils.NumberUUID(101),
					Validated: true,
					UpVotes:   8,
					DownVotes: 1,
					RequestID: test_utils.NumberUUID(14),
					Title:     "Smart post",
					Content:   "Cats on a nap.",
				},
			},
			expect: []*models.ImproveSuggestion{
				{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(10),
					UserID:    test_utils.NumberUUID(100),
					Validated: false,
					UpVotes:   17,
					DownVotes: 3,
					RequestID: test_utils.NumberUUID(11),
					Title:     "Dummy post",
					Content:   "Foo bar qux.",
				},
				{
					ID:        test_utils.NumberUUID(2),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(12),
					UserID:    test_utils.NumberUUID(101),
					Validated: true,
					UpVotes:   8,
					DownVotes: 1,
					RequestID: test_utils.NumberUUID(14),
					Title:     "Smart post",
					Content:   "Cats on a nap.",
				},
			},
		},
		{
			name: "Error/ServiceFailure",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			serviceErr: fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveSuggestionService := improve_suggestion_service.NewMockService(t)

			improveSuggestionService.
				On("GetPreviews", context.TODO(), d.ids).
				Return(d.serviceData, d.serviceErr)

			provider := NewProvider(Config{
				ImproveSuggestionService: improveSuggestionService,
			})

			res, err := provider.GetImproveSuggestionPreviews(context.TODO(), d.ids)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			improveSuggestionService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_Vote(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token  string
		postID uuid.UUID
		target models.VoteTarget
		vote   models.VoteValue

		shouldCallImproveRequestService    bool
		shouldCallImproveSuggestionService bool
		shouldCallVoteService              bool

		tokenServiceDecodeData       *models.UserToken
		tokenServiceDecodeErr        error
		improveRequestServiceData    bool
		improveRequestServiceErr     error
		improveSuggestionServiceData bool
		improveSuggestionServiceErr  error
		voteServiceData              models.VoteValue
		voteServiceErr               error

		expect    models.VoteValue
		expectErr error
	}{
		{
			name:                            "Success/ImproveRequest",
			now:                             baseTime,
			keys:                            jwk_storage.MockedKeys,
			userID:                          test_utils.NumberUUID(100),
			token:                           "foo.bar.qux",
			postID:                          test_utils.NumberUUID(10),
			target:                          models.VoteTargetImproveRequest,
			vote:                            models.VoteUp,
			shouldCallImproveRequestService: true,
			shouldCallVoteService:           true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			voteServiceData: models.VoteUp,
			expect:          models.VoteUp,
		},
		{
			name:                               "Success/ImproveSuggestion",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			userID:                             test_utils.NumberUUID(100),
			token:                              "foo.bar.qux",
			postID:                             test_utils.NumberUUID(10),
			target:                             models.VoteTargetImproveSuggestion,
			vote:                               models.VoteUp,
			shouldCallImproveSuggestionService: true,
			shouldCallVoteService:              true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			voteServiceData: models.VoteUp,
			expect:          models.VoteUp,
		},
		{
			name:                               "Error/VoteServiceFailure",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			userID:                             test_utils.NumberUUID(100),
			token:                              "foo.bar.qux",
			postID:                             test_utils.NumberUUID(10),
			target:                             models.VoteTargetImproveSuggestion,
			vote:                               models.VoteUp,
			shouldCallImproveSuggestionService: true,
			shouldCallVoteService:              true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			voteServiceErr: fooErr,
			expect:         models.NoVote,
			expectErr:      fooErr,
		},
		{
			name:                            "Error/ImproveRequestSelfVote",
			now:                             baseTime,
			keys:                            jwk_storage.MockedKeys,
			userID:                          test_utils.NumberUUID(100),
			token:                           "foo.bar.qux",
			postID:                          test_utils.NumberUUID(10),
			target:                          models.VoteTargetImproveRequest,
			vote:                            models.VoteUp,
			shouldCallImproveRequestService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveRequestServiceData: true,
			expect:                    models.NoVote,
			expectErr:                 validation.ErrInvalidEntity,
		},
		{
			name:                               "Error/ImproveSuggestionSelfVote",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			userID:                             test_utils.NumberUUID(100),
			token:                              "foo.bar.qux",
			postID:                             test_utils.NumberUUID(10),
			target:                             models.VoteTargetImproveSuggestion,
			vote:                               models.VoteUp,
			shouldCallImproveSuggestionService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveSuggestionServiceData: true,
			expect:                       models.NoVote,
			expectErr:                    validation.ErrInvalidEntity,
		},
		{
			name:                            "Error/ImproveRequestServiceFailure",
			now:                             baseTime,
			keys:                            jwk_storage.MockedKeys,
			userID:                          test_utils.NumberUUID(100),
			token:                           "foo.bar.qux",
			postID:                          test_utils.NumberUUID(10),
			target:                          models.VoteTargetImproveRequest,
			vote:                            models.VoteUp,
			shouldCallImproveRequestService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveRequestServiceErr: fooErr,
			expect:                   models.NoVote,
			expectErr:                fooErr,
		},
		{
			name:                               "Error/ImproveSuggestionServiceFailure",
			now:                                baseTime,
			keys:                               jwk_storage.MockedKeys,
			userID:                             test_utils.NumberUUID(100),
			token:                              "foo.bar.qux",
			postID:                             test_utils.NumberUUID(10),
			target:                             models.VoteTargetImproveSuggestion,
			vote:                               models.VoteUp,
			shouldCallImproveSuggestionService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			improveSuggestionServiceErr: fooErr,
			expect:                      models.NoVote,
			expectErr:                   fooErr,
		},
		{
			name:                  "Error/TokenServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			userID:                test_utils.NumberUUID(100),
			token:                 "foo.bar.qux",
			postID:                test_utils.NumberUUID(10),
			target:                models.VoteTargetImproveRequest,
			vote:                  models.VoteUp,
			tokenServiceDecodeErr: fooErr,
			expect:                models.NoVote,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			improveRequestService := improve_request_service.NewMockService(t)
			improveSuggestionService := improve_suggestion_service.NewMockService(t)
			voteService := votes_service.NewMockService(t)
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

			if d.shouldCallImproveRequestService {
				improveRequestService.
					On("IsCreator", context.TODO(), d.userID, d.postID, false).
					Return(d.improveRequestServiceData, d.improveRequestServiceErr)
			}

			if d.shouldCallImproveSuggestionService {
				improveSuggestionService.
					On("IsCreator", context.TODO(), d.userID, d.postID).
					Return(d.improveSuggestionServiceData, d.improveSuggestionServiceErr)
			}

			if d.shouldCallVoteService {
				voteService.
					On(
						"Vote", context.TODO(),
						d.postID, d.userID, models.VoteTarget(d.target), models.VoteValue(d.vote), d.now,
					).
					Return(d.voteServiceData, d.voteServiceErr)
			}

			provider := NewProvider(Config{
				ImproveRequestService:    improveRequestService,
				ImproveSuggestionService: improveSuggestionService,
				VotesService:             voteService,
				TokenService:             tokenService,
				KeysService:              keysService,
				Time:                     test_utils.GetTimeNow(d.now),
			})

			vote, err := provider.Vote(context.TODO(), d.token, d.postID, d.target, d.vote)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, vote)

			improveRequestService.AssertExpectations(t)
			improveSuggestionService.AssertExpectations(t)
			voteService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_HasVoted(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token  string
		postID uuid.UUID
		target models.VoteTarget

		shouldCallVoteService bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		voteServiceData        models.VoteValue
		voteServiceErr         error

		expect    models.VoteValue
		expectErr error
	}{
		{
			name:                  "Success",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			userID:                test_utils.NumberUUID(100),
			token:                 "foo.bar.qux",
			postID:                test_utils.NumberUUID(10),
			target:                models.VoteTargetImproveRequest,
			shouldCallVoteService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			voteServiceData: models.VoteUp,
			expect:          models.VoteUp,
		},
		{
			name:                  "Error/VoteServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			userID:                test_utils.NumberUUID(100),
			token:                 "foo.bar.qux",
			postID:                test_utils.NumberUUID(10),
			target:                models.VoteTargetImproveRequest,
			shouldCallVoteService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(100)},
			},
			voteServiceErr: fooErr,
			expect:         models.NoVote,
			expectErr:      fooErr,
		},
		{
			name:                  "Error/TokenServiceFailure",
			now:                   baseTime,
			keys:                  jwk_storage.MockedKeys,
			userID:                test_utils.NumberUUID(100),
			token:                 "foo.bar.qux",
			postID:                test_utils.NumberUUID(10),
			target:                models.VoteTargetImproveRequest,
			tokenServiceDecodeErr: fooErr,
			expect:                models.NoVote,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			voteService := votes_service.NewMockService(t)
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

			if d.shouldCallVoteService {
				voteService.
					On("HasVoted", context.TODO(), d.postID, d.userID, models.VoteTarget(d.target)).
					Return(d.voteServiceData, d.voteServiceErr)
			}

			provider := NewProvider(Config{
				VotesService: voteService,
				TokenService: tokenService,
				KeysService:  keysService,
				Time:         test_utils.GetTimeNow(d.now),
			})

			vote, err := provider.HasVoted(context.TODO(), d.token, d.postID, d.target)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, vote)

			voteService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestImprovePostProvider_GetVotedPosts(t *testing.T) {
	data := []struct {
		name   string
		userID uuid.UUID
		target models.VoteTarget
		limit  int
		offset int

		voteServiceData  []*models.VotedPost
		voteServiceCount int64
		voteServiceErr   error

		expect      []*models.VotedPost
		expectCount int64
		expectErr   error
	}{
		{
			name:   "Success",
			userID: test_utils.NumberUUID(100),
			target: models.VoteTargetImproveRequest,
			limit:  10,
			offset: 20,
			voteServiceData: []*models.VotedPost{
				{
					PostID:    test_utils.NumberUUID(10),
					Vote:      models.VoteUp,
					UpdatedAt: baseTime.Add(-time.Hour),
				},
				{
					PostID:    test_utils.NumberUUID(11),
					Vote:      models.VoteDown,
					UpdatedAt: baseTime.Add(30 * time.Minute),
				},
			},
			voteServiceCount: 200,
			expect: []*models.VotedPost{
				{
					PostID:    test_utils.NumberUUID(10),
					Vote:      models.VoteUp,
					UpdatedAt: baseTime.Add(-time.Hour),
				},
				{
					PostID:    test_utils.NumberUUID(11),
					Vote:      models.VoteDown,
					UpdatedAt: baseTime.Add(30 * time.Minute),
				},
			},
			expectCount: 200,
		},
		{
			name:           "Error/ServiceFailure",
			userID:         test_utils.NumberUUID(100),
			target:         models.VoteTargetImproveRequest,
			limit:          10,
			offset:         20,
			voteServiceErr: fooErr,
			expectErr:      fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			voteService := votes_service.NewMockService(t)

			voteService.
				On("GetVotedPosts", context.TODO(), d.userID, models.VoteTarget(d.target), d.limit, d.offset).
				Return(d.voteServiceData, d.voteServiceCount, d.voteServiceErr)

			provider := NewProvider(Config{
				VotesService: voteService,
			})

			posts, count, err := provider.GetVotedPosts(context.TODO(), d.userID, d.target, d.limit, d.offset)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, posts)
			require.Equal(t, d.expectCount, count)

			voteService.AssertExpectations(t)
		})
	}
}
