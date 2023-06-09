package authentication

import (
	"context"
	"crypto/ed25519"
	"errors"
	"github.com/a-novel/agora-backend/domains/keys/service/jwk"
	"github.com/a-novel/agora-backend/domains/keys/storage/jwk"
	"github.com/a-novel/agora-backend/domains/user/service/credentials"
	"github.com/a-novel/agora-backend/domains/user/service/token"
	"github.com/a-novel/agora-backend/framework/test"
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

func TestAuthenticationProvider_Authenticate(t *testing.T) {
	data := []struct {
		name string

		token           string
		autoRenew       bool
		tokenTTL        time.Duration
		tokenRenewDelta time.Duration

		now  time.Time
		id   uuid.UUID
		keys []ed25519.PrivateKey

		tokenDecodeData  *models.UserToken
		tokenDecodeError error
		tokenEncodeData  string
		tokenEncodeError error

		shouldCallTokenDecodeService     bool
		shouldCallTokenEncodeService     bool
		shouldCallTokenEncodeServiceWith models.UserTokenPayload

		expectedToken string
		expectedError error
	}{
		{
			name:            "Success",
			token:           "foo.bar.qux",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			keys: []ed25519.PrivateKey{
				jwk_storage.MockedKeys[0],
				jwk_storage.MockedKeys[1],
				jwk_storage.MockedKeys[2],
			},
			tokenDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-10 * time.Minute),
					EXP: baseTime.Add(2 * time.Minute),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(1)},
			},
			shouldCallTokenDecodeService: true,
			expectedToken:                "foo.bar.qux",
		},
		{
			name:            "Success/AutoRenewalNoTrigger",
			autoRenew:       true,
			token:           "foo.bar.qux",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			keys: []ed25519.PrivateKey{
				jwk_storage.MockedKeys[0],
				jwk_storage.MockedKeys[1],
				jwk_storage.MockedKeys[2],
			},
			tokenDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-10 * time.Minute),
					EXP: baseTime.Add(10 * time.Minute),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(1)},
			},
			shouldCallTokenDecodeService: true,
			expectedToken:                "foo.bar.qux",
		},
		{
			name:            "Success/AutoRenewal",
			autoRenew:       true,
			token:           "foo.bar.qux",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			keys: []ed25519.PrivateKey{
				jwk_storage.MockedKeys[0],
				jwk_storage.MockedKeys[1],
				jwk_storage.MockedKeys[2],
			},
			tokenDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-10 * time.Minute),
					EXP: baseTime.Add(2 * time.Minute),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(1)},
			},
			shouldCallTokenDecodeService:     true,
			shouldCallTokenEncodeService:     true,
			shouldCallTokenEncodeServiceWith: models.UserTokenPayload{ID: test_utils.NumberUUID(1)},
			tokenEncodeData:                  "qux.bar.foo",
			expectedToken:                    "qux.bar.foo",
		},
		{
			name:            "Error/NoToken",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			expectedError:   validation.ErrInvalidCredentials,
		},
		{
			name:            "Error/TokenEncodeFailure",
			autoRenew:       true,
			token:           "foo.bar.qux",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			keys: []ed25519.PrivateKey{
				jwk_storage.MockedKeys[0],
				jwk_storage.MockedKeys[1],
				jwk_storage.MockedKeys[2],
			},
			tokenDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-10 * time.Minute),
					EXP: baseTime.Add(2 * time.Minute),
					ID:  test_utils.NumberUUID(100),
				},
				Payload: models.UserTokenPayload{ID: test_utils.NumberUUID(1)},
			},
			shouldCallTokenDecodeService:     true,
			shouldCallTokenEncodeService:     true,
			shouldCallTokenEncodeServiceWith: models.UserTokenPayload{ID: test_utils.NumberUUID(1)},
			tokenEncodeError:                 fooErr,
			expectedToken:                    "foo.bar.qux",
			expectedError:                    fooErr,
		},
		{
			name:            "Error/TokenDecodeServiceFailure",
			token:           "foo.bar.qux",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			keys: []ed25519.PrivateKey{
				jwk_storage.MockedKeys[0],
				jwk_storage.MockedKeys[1],
				jwk_storage.MockedKeys[2],
			},
			tokenDecodeError:             fooErr,
			shouldCallTokenDecodeService: true,
			expectedError:                fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)
			now := test_utils.GetTimeNow(d.now)
			id := test_utils.GetUUID(d.id)

			provider := NewProvider(Config{
				TokenService:    tokenService,
				KeysService:     keysService,
				Time:            now,
				ID:              id,
				TokenTTL:        d.tokenTTL,
				TokenRenewDelta: d.tokenRenewDelta,
			})

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenDecodeService {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenDecodeData, d.tokenDecodeError)
			}

			if d.shouldCallTokenEncodeService {
				keysService.
					On("GetPrivate").
					Return(d.keys[0])

				tokenService.
					On("Encode", d.shouldCallTokenEncodeServiceWith, d.tokenTTL, d.keys[0], d.id, d.now).
					Return(d.tokenEncodeData, d.tokenEncodeError)
			}

			token, err := provider.Authenticate(context.TODO(), d.token, d.autoRenew)
			test_utils.RequireError(st, d.expectedError, err)
			require.Equal(st, d.expectedToken, token)

			tokenService.AssertExpectations(st)
			keysService.AssertExpectations(st)
		})
	}
}

func TestAuthenticationProvider_Login(t *testing.T) {
	data := []struct {
		name string

		form models.UserCredentialsLoginForm

		tokenTTL        time.Duration
		tokenRenewDelta time.Duration

		now  time.Time
		id   uuid.UUID
		keys []ed25519.PrivateKey

		tokenEncodeData  string
		tokenEncodeError error
		credentialsData  *models.UserCredentials
		credentialsError error

		shouldCallTokenEncodeService      bool
		shouldCallCredentialsService      bool
		shouldCallUserServicesWithPayload models.UserTokenPayload

		expectedToken string
		expectedError error
	}{
		{
			name:            "Success",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			keys: []ed25519.PrivateKey{
				jwk_storage.MockedKeys[0],
				jwk_storage.MockedKeys[1],
				jwk_storage.MockedKeys[2],
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &baseTime,
				Email:     "user@company.com",
				NewEmail:  "user2@company.com",
				Validated: true,
			},
			tokenEncodeData:              "foo.bar.qux",
			shouldCallCredentialsService: true,
			shouldCallTokenEncodeService: true,
			shouldCallUserServicesWithPayload: models.UserTokenPayload{
				ID: test_utils.NumberUUID(1),
			},
			expectedToken: "foo.bar.qux",
		},
		{
			name:            "Error/TokenEncodeFailure",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			keys: []ed25519.PrivateKey{
				jwk_storage.MockedKeys[0],
				jwk_storage.MockedKeys[1],
				jwk_storage.MockedKeys[2],
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &baseTime,
				Email:     "user@company.com",
				NewEmail:  "user2@company.com",
				Validated: true,
			},
			shouldCallCredentialsService: true,
			shouldCallUserServicesWithPayload: models.UserTokenPayload{
				ID: test_utils.NumberUUID(1),
			},
			shouldCallTokenEncodeService: true,
			tokenEncodeError:             fooErr,
			expectedError:                fooErr,
		},
		{
			name:            "Error/CredentialsServiceFailure",
			tokenTTL:        time.Hour,
			tokenRenewDelta: 5 * time.Minute,
			now:             baseTime,
			id:              test_utils.NumberUUID(12),
			keys: []ed25519.PrivateKey{
				jwk_storage.MockedKeys[0],
				jwk_storage.MockedKeys[1],
				jwk_storage.MockedKeys[2],
			},
			credentialsError:             fooErr,
			shouldCallCredentialsService: true,
			shouldCallUserServicesWithPayload: models.UserTokenPayload{
				ID: test_utils.NumberUUID(1),
			},
			expectedError: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			credentialsService := credentials_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)
			now := test_utils.GetTimeNow(d.now)
			id := test_utils.GetUUID(d.id)

			provider := NewProvider(Config{
				CredentialsService: credentialsService,
				TokenService:       tokenService,
				KeysService:        keysService,
				Time:               now,
				ID:                 id,
				TokenTTL:           d.tokenTTL,
				TokenRenewDelta:    d.tokenRenewDelta,
			})

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenEncodeService {
				keysService.
					On("GetPrivate").
					Return(d.keys[0])

				tokenService.
					On("Encode", d.shouldCallUserServicesWithPayload, d.tokenTTL, d.keys[0], d.id, d.now).
					Return(d.tokenEncodeData, d.tokenEncodeError)
			}

			if d.shouldCallCredentialsService {
				credentialsService.
					On("Authenticate", context.TODO(), &models.UserCredentialsLoginForm{
						Email:    d.form.Email,
						Password: d.form.Password,
					}).
					Return(d.credentialsData, d.credentialsError)
			}

			token, err := provider.Login(context.TODO(), d.form)
			test_utils.RequireError(st, d.expectedError, err)
			require.Equal(st, d.expectedToken, token)

			credentialsService.AssertExpectations(st)
			tokenService.AssertExpectations(st)
			keysService.AssertExpectations(st)
		})
	}
}
