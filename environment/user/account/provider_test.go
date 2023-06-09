package account

import (
	"context"
	"crypto/ed25519"
	"errors"
	"github.com/a-novel/agora-backend/domains/generics"
	"github.com/a-novel/agora-backend/framework"
	"github.com/a-novel/agora-backend/framework/mailer"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	fooErr   = errors.New("it broken")
)

func TestAccountProvider_Register(t *testing.T) {
	data := []struct {
		name string

		now time.Time
		id  uuid.UUID

		tokenTTL            time.Duration
		keys                []ed25519.PrivateKey
		emailValidationLink generics.URL
		emailTemplate       string

		form models.UserCreateForm

		shouldCallTokenServiceEncode bool
		shouldCallUserService        bool
		shouldCallMailer             bool
		shouldReturnDeferred         bool

		shouldCallMailerWithEmail    *mail.Email
		shouldCallMailerWithData     map[string]interface{}
		shouldCallMailerWithTemplate string

		tokenServiceEncodeData string
		tokenServiceEncodeErr  error
		userServiceData        *models.User
		userServicePRModel     *models.UserPostRegistration
		userServiceErr         error
		mailerErr              error

		expect         *models.UserFlat
		expectToken    string
		expectErr      error
		expectDeferErr error
	}{
		{
			name:     "Success",
			now:      baseTime,
			id:       test_utils.NumberUUID(1),
			tokenTTL: 5 * time.Hour,
			keys:     jwk_storage.MockedKeys,
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailTemplate: "foo_template",
			form: models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "sylvester@worldcompany.com",
					Password: "secret",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Sylvestre",
					LastName:  "Stalle",
					Birthday:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "CashCash",
					Slug:     "world-company-international",
				},
			},
			shouldCallTokenServiceEncode: true,
			shouldCallUserService:        true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "sylvester@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			tokenServiceEncodeData:       "foo.bar.qux",
			userServiceData: &models.User{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &baseTime,
				Credentials: models.UserCredentials{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &baseTime,
					Email:     "sylvester@worldcompany.com",
					Validated: false,
				},
				Identity: models.UserIdentity{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &baseTime,
					FirstName: "Sylvestre",
					LastName:  "Stalle",
					Birthday:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: models.UserProfile{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &baseTime,
					Username:  "CashCash",
					Slug:      "world-company-international",
				},
			},
			userServicePRModel: &models.UserPostRegistration{
				EmailValidationCode: "super_validation_code_9000",
			},
			expect: &models.UserFlat{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &baseTime,
				Email:     "sylvester@worldcompany.com",
				Validated: false,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
				Username:  "CashCash",
				Slug:      "world-company-international",
			},
			expectToken: "foo.bar.qux",
		},
		{
			name:     "Success/DeferredFailure",
			now:      baseTime,
			id:       test_utils.NumberUUID(1),
			tokenTTL: 5 * time.Hour,
			keys:     jwk_storage.MockedKeys,
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailTemplate: "foo_template",
			form: models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "sylvester@worldcompany.com",
					Password: "secret",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Sylvestre",
					LastName:  "Stalle",
					Birthday:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "CashCash",
					Slug:     "world-company-international",
				},
			},
			shouldCallTokenServiceEncode: true,
			shouldCallUserService:        true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "sylvester@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			tokenServiceEncodeData:       "foo.bar.qux",
			userServiceData: &models.User{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &baseTime,
				Credentials: models.UserCredentials{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &baseTime,
					Email:     "sylvester@worldcompany.com",
					Validated: false,
				},
				Identity: models.UserIdentity{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &baseTime,
					FirstName: "Sylvestre",
					LastName:  "Stalle",
					Birthday:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: models.UserProfile{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &baseTime,
					Username:  "CashCash",
					Slug:      "world-company-international",
				},
			},
			userServicePRModel: &models.UserPostRegistration{
				EmailValidationCode: "super_validation_code_9000",
			},
			mailerErr: fooErr,
			expect: &models.UserFlat{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &baseTime,
				Email:     "sylvester@worldcompany.com",
				Validated: false,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
				Username:  "CashCash",
				Slug:      "world-company-international",
			},
			expectToken:    "foo.bar.qux",
			expectDeferErr: fooErr,
		},
		{
			name:     "Error/UserServiceFailure",
			now:      baseTime,
			id:       test_utils.NumberUUID(1),
			tokenTTL: 5 * time.Hour,
			keys:     jwk_storage.MockedKeys,
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailTemplate: "foo_template",
			form: models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "sylvester@worldcompany.com",
					Password: "secret",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Sylvestre",
					LastName:  "Stalle",
					Birthday:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "CashCash",
					Slug:     "world-company-international",
				},
			},
			shouldCallTokenServiceEncode: true,
			shouldCallUserService:        true,
			tokenServiceEncodeData:       "foo.bar.qux",
			userServiceErr:               fooErr,
			expectErr:                    fooErr,
		},
		{
			name:     "Error/TokenServiceFailure",
			now:      baseTime,
			id:       test_utils.NumberUUID(1),
			tokenTTL: 5 * time.Hour,
			keys:     jwk_storage.MockedKeys,
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailTemplate: "foo_template",
			form: models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "sylvester@worldcompany.com",
					Password: "secret",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Sylvestre",
					LastName:  "Stalle",
					Birthday:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "CashCash",
					Slug:     "world-company-international",
				},
			},
			shouldCallTokenServiceEncode: true,
			tokenServiceEncodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			userService := user_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)
			tokenService := token_service.NewMockService(t)
			mailerService := mailer.NewMockMailer(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceEncode {
				keysService.
					On("GetPrivate").
					Return(d.keys[0])

				tokenService.
					On("Encode", models.UserTokenPayload{ID: d.id}, d.tokenTTL, d.keys[0], d.id, d.now).
					Return(d.tokenServiceEncodeData, d.tokenServiceEncodeErr)
			}

			if d.shouldCallUserService {
				userService.
					On("Create", context.TODO(), &models.UserCreateForm{
						Credentials: models.UserCredentialsLoginForm{
							Email:    d.form.Credentials.Email,
							Password: d.form.Credentials.Password,
						},
						Identity: models.UserIdentityUpdateForm{
							FirstName: d.form.Identity.FirstName,
							LastName:  d.form.Identity.LastName,
							Birthday:  d.form.Identity.Birthday,
							Sex:       d.form.Identity.Sex,
						},
						Profile: models.UserProfileUpdateForm{
							Username: d.form.Profile.Username,
							Slug:     d.form.Profile.Slug,
						},
					}, d.id, d.now).
					Return(d.userServiceData, d.userServicePRModel, d.userServiceErr)
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", d.shouldCallMailerWithEmail, d.shouldCallMailerWithTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			provider := NewProvider(Config{
				UserService:                userService,
				TokenService:               tokenService,
				KeysService:                keysService,
				Mailer:                     mailerService,
				Time:                       test_utils.GetTimeNow(d.now),
				ID:                         test_utils.GetUUID(d.id),
				TokenTTL:                   d.tokenTTL,
				EmailValidationLink:        d.emailValidationLink,
				EmailValidationTemplate:    d.emailTemplate,
				NewEMailValidationTemplate: "NEW_EMAIL_TEMPLATE",
				PasswordResetTemplate:      "RESET_PASSWORD_TEMPLATE",
			})

			user, token, deferred, err := provider.Register(context.TODO(), d.form)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, user)
			require.Equal(t, d.expectToken, token)

			if d.shouldReturnDeferred {
				require.NotNil(t, deferred)
				test_utils.RequireError(t, d.expectDeferErr, deferred())
			}

			userService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
			mailerService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_GetAccountInfo(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token string

		shouldCallTokenServiceDecode bool
		shouldCallCredentialsService bool
		shouldCallIdentityService    bool
		shouldCallProfileService     bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		credentialsData        *models.UserCredentials
		credentialsErr         error
		identityData           *models.UserIdentity
		identityErr            error
		profileData            *models.UserProfile
		profileErr             error

		expect    *models.UserInfo
		expectErr error
	}{
		{
			name:                         "Success",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallProfileService:     true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				Email:     "user@company.com",
				NewEmail:  "new_user@company.com",
				Validated: true,
			},
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(2 * time.Hour)),
				FirstName: "Foo",
				LastName:  "Bar",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			profileData: &models.UserProfile{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Username:  "foo_bar",
				Slug:      "foo-bar",
			},
			expect: &models.UserInfo{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(2 * time.Hour)),
				Email:     "user@company.com",
				NewEmail:  "new_user@company.com",
				Identity: models.UserInfoIdentity{
					FirstName: "Foo",
					LastName:  "Bar",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
				Profile: models.UserInfoProfile{
					Username: "foo_bar",
					Slug:     "foo-bar",
				},
			},
		},
		{
			name:                         "Error/ProfileServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallProfileService:     true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				Email:     "user@company.com",
				NewEmail:  "new_user@company.com",
				Validated: true,
			},
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(2 * time.Hour)),
				FirstName: "Foo",
				LastName:  "Bar",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			profileErr: fooErr,
			expectErr:  fooErr,
		},
		{
			name:                         "Error/IdentityServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: framework.ToPTR(baseTime.Add(time.Hour)),
				Email:     "user@company.com",
				NewEmail:  "new_user@company.com",
				Validated: true,
			},
			identityErr: fooErr,
			expectErr:   fooErr,
		},
		{
			name:                         "Error/CredentialsServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsErr: fooErr,
			expectErr:      fooErr,
		},
		{
			name:                         "Error/TokenServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)
			identityService := identity_service.NewMockService(t)
			profileService := profile_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallCredentialsService {
				credentialsService.
					On("Read", context.TODO(), d.userID).
					Return(d.credentialsData, d.credentialsErr)
			}

			if d.shouldCallIdentityService {
				identityService.
					On("Read", context.TODO(), d.userID).
					Return(d.identityData, d.identityErr)
			}

			if d.shouldCallProfileService {
				profileService.
					On("Read", context.TODO(), d.userID).
					Return(d.profileData, d.profileErr)
			}

			provider := NewProvider(Config{
				CredentialsService: credentialsService,
				IdentityService:    identityService,
				ProfileService:     profileService,
				TokenService:       tokenService,
				KeysService:        keysService,
				Time:               test_utils.GetTimeNow(d.now),
			})

			info, err := provider.GetAccountInfo(context.TODO(), d.token)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, info)

			credentialsService.AssertExpectations(t)
			identityService.AssertExpectations(t)
			profileService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_GetAccountPreview(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token string

		shouldCallTokenServiceDecode bool
		shouldCallUserService        bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		userData               *models.UserPreview
		userErr                error

		expect    *models.UserPreview
		expectErr error
	}{
		{
			name:                         "Success",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallUserService:        true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			userData: &models.UserPreview{
				ID:        test_utils.NumberUUID(1),
				Username:  "qwerty",
				FirstName: "Foo",
				LastName:  "Bar",
				Email:     "user@company.com",
			},
			expect: &models.UserPreview{
				ID:        test_utils.NumberUUID(1),
				Email:     "user@company.com",
				Username:  "qwerty",
				FirstName: "Foo",
				LastName:  "Bar",
			},
		},
		{
			name:                         "Error/UserServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallUserService:        true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			userErr:   fooErr,
			expectErr: fooErr,
		},
		{
			name:                         "Error/TokenServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			userService := user_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallUserService {
				userService.
					On("GetPreview", context.TODO(), d.userID).
					Return(d.userData, d.userErr)
			}

			provider := NewProvider(Config{
				UserService:  userService,
				TokenService: tokenService,
				KeysService:  keysService,
				Time:         test_utils.GetTimeNow(d.now),
			})

			res, err := provider.GetAccountPreview(context.TODO(), d.token)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			userService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_GetEmailValidationStatus(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token string

		shouldCallTokenServiceDecode bool
		shouldCallCredentialsService bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		credentialsData        *models.UserCredentials
		credentialsErr         error

		expect    *models.UserEmailValidationStatus
		expectErr error
	}{
		{
			name:                         "Success",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "user@company.com",
				NewEmail:  "new_user@company.com",
				Validated: true,
			},
			expect: &models.UserEmailValidationStatus{
				Email:     "user@company.com",
				NewEmail:  "new_user@company.com",
				Validated: true,
			},
		},
		{
			name:                         "Error/CredentialsServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsErr: fooErr,
			expectErr:      fooErr,
		},
		{
			name:                         "Error/TokenServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallCredentialsService {
				credentialsService.
					On("Read", context.TODO(), d.userID).
					Return(d.credentialsData, d.credentialsErr)
			}

			provider := NewProvider(Config{
				CredentialsService: credentialsService,
				TokenService:       tokenService,
				KeysService:        keysService,
				Time:               test_utils.GetTimeNow(d.now),
			})

			res, err := provider.GetEmailValidationStatus(context.TODO(), d.token)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			credentialsService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_UpdateIdentity(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token string
		form  models.UserIdentityUpdateForm

		shouldCallTokenServiceDecode bool
		shouldCallIdentityService    bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		identityData           *models.UserIdentity
		identityErr            error

		expect    *models.UserInfoIdentity
		expectErr error
	}{
		{
			name: "Success",
			now:  baseTime,
			keys: jwk_storage.MockedKeys,
			form: models.UserIdentityUpdateForm{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallIdentityService:    true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			expect: &models.UserInfoIdentity{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		{
			name: "Error/IdentityServiceFailure",
			now:  baseTime,
			keys: jwk_storage.MockedKeys,
			form: models.UserIdentityUpdateForm{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallIdentityService:    true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			identityErr: fooErr,
			expectErr:   fooErr,
		},
		{
			name: "Error/TokenServiceFailure",
			now:  baseTime,
			keys: jwk_storage.MockedKeys,
			form: models.UserIdentityUpdateForm{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			identityService := identity_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallIdentityService {
				identityService.
					On("Update", context.TODO(), &models.UserIdentityUpdateForm{
						FirstName: d.form.FirstName,
						LastName:  d.form.LastName,
						Birthday:  d.form.Birthday,
						Sex:       d.form.Sex,
					}, d.userID, d.now).
					Return(d.identityData, d.identityErr)
			}

			provider := NewProvider(Config{
				IdentityService: identityService,
				TokenService:    tokenService,
				KeysService:     keysService,
				Time:            test_utils.GetTimeNow(d.now),
			})

			res, err := provider.UpdateIdentity(context.TODO(), d.token, d.form)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			identityService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_UpdateProfile(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token string
		form  models.UserProfileUpdateForm

		shouldCallTokenServiceDecode bool
		shouldCallProfileService     bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		profileData            *models.UserProfile
		profileErr             error

		expect    *models.UserInfoProfile
		expectErr error
	}{
		{
			name: "Success",
			now:  baseTime,
			keys: jwk_storage.MockedKeys,
			form: models.UserProfileUpdateForm{
				Username: "banananana",
				Slug:     "fruit-basket",
			},
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallProfileService:     true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			profileData: &models.UserProfile{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Username:  "banananana",
				Slug:      "fruit-basket",
			},
			expect: &models.UserInfoProfile{
				Username: "banananana",
				Slug:     "fruit-basket",
			},
		},
		{
			name: "Error/ProfileServiceFailure",
			now:  baseTime,
			keys: jwk_storage.MockedKeys,
			form: models.UserProfileUpdateForm{
				Username: "banananana",
				Slug:     "fruit-basket",
			},
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallProfileService:     true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			profileErr: fooErr,
			expectErr:  fooErr,
		},
		{
			name: "Error/TokenServiceFailure",
			now:  baseTime,
			keys: jwk_storage.MockedKeys,
			form: models.UserProfileUpdateForm{
				Username: "banananana",
				Slug:     "fruit-basket",
			},
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			profileService := profile_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallProfileService {
				profileService.
					On("Update", context.TODO(), &models.UserProfileUpdateForm{
						Username: d.form.Username,
						Slug:     d.form.Slug,
					}, d.userID, d.now).
					Return(d.profileData, d.profileErr)
			}

			provider := NewProvider(Config{
				ProfileService: profileService,
				TokenService:   tokenService,
				KeysService:    keysService,
				Time:           test_utils.GetTimeNow(d.now),
			})

			res, err := provider.UpdateProfile(context.TODO(), d.token, d.form)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			profileService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_UpdatePassword(t *testing.T) {
	data := []struct {
		name string

		now  time.Time
		form models.UserPasswordUpdateForm

		credentialsErr error
		expectErr      error
	}{
		{
			name: "Success",
			now:  baseTime,
			form: models.UserPasswordUpdateForm{
				ID:          test_utils.NumberUUID(1),
				Password:    "123456",
				OldPassword: "abcdef",
			},
		},
		{
			name: "Error/CredentialsServiceFailure",
			now:  baseTime,
			form: models.UserPasswordUpdateForm{
				ID:          test_utils.NumberUUID(1),
				Password:    "123456",
				OldPassword: "abcdef",
			},
			credentialsErr: fooErr,
			expectErr:      fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)

			credentialsService.
				On("UpdatePassword", context.TODO(), d.form.OldPassword, d.form.Password, d.form.ID, d.now).
				Return(nil, d.credentialsErr)

			provider := NewProvider(Config{
				CredentialsService: credentialsService,
				Time:               test_utils.GetTimeNow(d.now),
			})

			err := provider.UpdatePassword(context.TODO(), d.form)
			test_utils.RequireError(t, d.expectErr, err)

			credentialsService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_UpdateEmail(t *testing.T) {
	data := []struct {
		name string

		now                    time.Time
		keys                   []ed25519.PrivateKey
		userID                 uuid.UUID
		newEmailValidationLink generics.URL
		newEmailTemplate       string

		token string
		form  models.UserEmailUpdateForm

		shouldCallTokenServiceDecode bool
		shouldCallCredentialsService bool
		shouldCallIdentityService    bool
		shouldCallMailer             bool
		shouldReturnDeferred         bool

		shouldCallMailerWithEmail    *mail.Email
		shouldCallMailerWithData     map[string]interface{}
		shouldCallMailerWithTemplate string

		tokenServiceDecodeData    *models.UserToken
		tokenServiceDecodeErr     error
		credentialsData           *models.UserCredentials
		credentialsValidationCode string
		credentialsErr            error
		identityData              *models.UserIdentity
		identityErr               error
		mailerErr                 error

		expect         *models.UserEmailValidationStatus
		expectErr      error
		expectDeferErr error
	}{
		{
			name:   "Success",
			now:    baseTime,
			keys:   jwk_storage.MockedKeys,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailTemplate:             "foo_template",
			token:                        "foo.bar.qux",
			form:                         models.UserEmailUpdateForm{Email: "boss@worldcompany.com"},
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "boss@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsValidationCode: "super_validation_code_9000",
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			expect: &models.UserEmailValidationStatus{
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
		},
		{
			name:   "Error/MailerFailure",
			now:    baseTime,
			keys:   jwk_storage.MockedKeys,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailTemplate:             "foo_template",
			token:                        "foo.bar.qux",
			form:                         models.UserEmailUpdateForm{Email: "boss@worldcompany.com"},
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "boss@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsValidationCode: "super_validation_code_9000",
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			mailerErr: fooErr,
			expect: &models.UserEmailValidationStatus{
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			expectDeferErr: fooErr,
		},
		{
			name:   "Error/IdentityServiceFailure",
			now:    baseTime,
			keys:   jwk_storage.MockedKeys,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailTemplate:             "foo_template",
			token:                        "foo.bar.qux",
			form:                         models.UserEmailUpdateForm{Email: "boss@worldcompany.com"},
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsValidationCode: "super_validation_code_9000",
			identityErr:               fooErr,
			expectErr:                 fooErr,
		},
		{
			name:   "Error/CredentialsServiceFailure",
			now:    baseTime,
			keys:   jwk_storage.MockedKeys,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailTemplate:             "foo_template",
			token:                        "foo.bar.qux",
			form:                         models.UserEmailUpdateForm{Email: "boss@worldcompany.com"},
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsErr: fooErr,
			expectErr:      fooErr,
		},
		{
			name:   "Error/TokenServiceFailure",
			now:    baseTime,
			keys:   jwk_storage.MockedKeys,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailTemplate:             "foo_template",
			token:                        "foo.bar.qux",
			form:                         models.UserEmailUpdateForm{Email: "boss@worldcompany.com"},
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)
			identityService := identity_service.NewMockService(t)
			profileService := profile_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)
			tokenService := token_service.NewMockService(t)
			mailerService := mailer.NewMockMailer(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallCredentialsService {
				credentialsService.
					On("UpdateEmail", context.TODO(), d.form.Email, d.userID, d.now).
					Return(d.credentialsData, d.credentialsValidationCode, d.credentialsErr)
			}

			if d.shouldCallIdentityService {
				identityService.
					On("Read", context.TODO(), d.userID).
					Return(d.identityData, d.identityErr)
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", d.shouldCallMailerWithEmail, d.shouldCallMailerWithTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			provider := NewProvider(Config{
				CredentialsService:         credentialsService,
				IdentityService:            identityService,
				ProfileService:             profileService,
				TokenService:               tokenService,
				KeysService:                keysService,
				Mailer:                     mailerService,
				Time:                       test_utils.GetTimeNow(d.now),
				NewEmailValidationLink:     d.newEmailValidationLink,
				NewEMailValidationTemplate: d.newEmailTemplate,
			})

			user, deferred, err := provider.UpdateEmail(context.TODO(), d.token, d.form)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, user)

			if d.shouldReturnDeferred {
				require.NotNil(t, deferred)
				test_utils.RequireError(t, d.expectDeferErr, deferred())
			}

			credentialsService.AssertExpectations(t)
			identityService.AssertExpectations(t)
			profileService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
			mailerService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_CancelNewEmail(t *testing.T) {
	data := []struct {
		name string

		now    time.Time
		keys   []ed25519.PrivateKey
		userID uuid.UUID

		token string

		shouldCallTokenServiceDecode bool
		shouldCallCredentialsService bool

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		credentialsErr         error

		expectErr error
	}{
		{
			name:                         "Success",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
		},
		{
			name:                         "Error/UserServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			shouldCallCredentialsService: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			credentialsErr: fooErr,
			expectErr:      fooErr,
		},
		{
			name:                         "Error/TokenServiceFailure",
			now:                          baseTime,
			keys:                         jwk_storage.MockedKeys,
			userID:                       test_utils.NumberUUID(1),
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallCredentialsService {
				credentialsService.
					On("CancelNewEmail", context.TODO(), d.userID, d.now).
					Return(nil, d.credentialsErr)
			}

			provider := NewProvider(Config{
				CredentialsService: credentialsService,
				TokenService:       tokenService,
				KeysService:        keysService,
				Time:               test_utils.GetTimeNow(d.now),
			})

			err := provider.CancelNewEmail(context.TODO(), d.token)
			test_utils.RequireError(t, d.expectErr, err)

			credentialsService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_ResetPassword(t *testing.T) {
	data := []struct {
		name string

		now                   time.Time
		userID                uuid.UUID
		passwordResetLink     generics.URL
		passwordResetTemplate string

		form models.UserPasswordResetForm

		shouldCallCredentialsService bool
		shouldCallIdentityService    bool
		shouldCallMailer             bool
		shouldReturnDeferred         bool

		shouldCallMailerWithEmail    *mail.Email
		shouldCallMailerWithData     map[string]interface{}
		shouldCallMailerWithTemplate string

		credentialsData      *models.UserCredentials
		credentialsResetCode string
		credentialsErr       error
		identityData         *models.UserIdentity
		identityErr          error
		mailerErr            error

		expectErr      error
		expectDeferErr error
	}{
		{
			name:   "Success",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			passwordResetLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			passwordResetTemplate:        "foo_template",
			form:                         models.UserPasswordResetForm{Email: "sylvester@worldcompany.com"},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "sylvester@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		{
			name:   "Error/MailerFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			passwordResetLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			passwordResetTemplate:        "foo_template",
			form:                         models.UserPasswordResetForm{Email: "sylvester@worldcompany.com"},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "sylvester@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			mailerErr:      fooErr,
			expectDeferErr: fooErr,
		},
		{
			name:   "Error/IdentityServiceFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			passwordResetLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			passwordResetTemplate:        "foo_template",
			form:                         models.UserPasswordResetForm{Email: "sylvester@worldcompany.com"},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityErr:          fooErr,
			expectErr:            fooErr,
		},
		{
			name:   "Error/CredentialsServiceFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			passwordResetLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			passwordResetTemplate:        "foo_template",
			form:                         models.UserPasswordResetForm{Email: "sylvester@worldcompany.com"},
			shouldCallCredentialsService: true,
			credentialsErr:               fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)
			identityService := identity_service.NewMockService(t)
			profileService := profile_service.NewMockService(t)
			mailerService := mailer.NewMockMailer(t)

			if d.shouldCallCredentialsService {
				credentialsService.
					On("ResetPassword", context.TODO(), d.form.Email, d.now).
					Return(d.credentialsData, d.credentialsResetCode, d.credentialsErr)
			}

			if d.shouldCallIdentityService {
				identityService.
					On("Read", context.TODO(), d.userID).
					Return(d.identityData, d.identityErr)
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", d.shouldCallMailerWithEmail, d.shouldCallMailerWithTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			provider := NewProvider(Config{
				CredentialsService:    credentialsService,
				IdentityService:       identityService,
				ProfileService:        profileService,
				Mailer:                mailerService,
				Time:                  test_utils.GetTimeNow(d.now),
				PasswordResetLink:     d.passwordResetLink,
				PasswordResetTemplate: d.passwordResetTemplate,
			})

			deferred, err := provider.ResetPassword(context.TODO(), d.form)
			test_utils.RequireError(t, d.expectErr, err)

			if d.shouldReturnDeferred {
				require.NotNil(t, deferred)
				test_utils.RequireError(t, d.expectDeferErr, deferred())
			}

			credentialsService.AssertExpectations(t)
			identityService.AssertExpectations(t)
			profileService.AssertExpectations(t)
			mailerService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_ValidateEmail(t *testing.T) {
	data := []struct {
		name string

		now  time.Time
		form models.UserValidateEmailForm

		credentialsErr error
		expectErr      error
	}{
		{
			name: "Success",
			now:  baseTime,
			form: models.UserValidateEmailForm{
				ID:   test_utils.NumberUUID(1),
				Code: "super_validation_code_9000",
			},
		},
		{
			name: "Error/CredentialsServiceFailure",
			now:  baseTime,
			form: models.UserValidateEmailForm{
				ID:   test_utils.NumberUUID(1),
				Code: "super_validation_code_9000",
			},
			credentialsErr: fooErr,
			expectErr:      fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)

			credentialsService.
				On("ValidateEmail", context.TODO(), d.form.ID, d.form.Code, d.now).
				Return(nil, d.credentialsErr)

			provider := NewProvider(Config{
				CredentialsService: credentialsService,
				Time:               test_utils.GetTimeNow(d.now),
			})

			err := provider.ValidateEmail(context.TODO(), d.form)
			test_utils.RequireError(t, d.expectErr, err)

			credentialsService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_ValidateNewEmail(t *testing.T) {
	data := []struct {
		name string

		now  time.Time
		form models.UserValidateEmailForm

		credentialsErr error
		expectErr      error
	}{
		{
			name: "Success",
			now:  baseTime,
			form: models.UserValidateEmailForm{
				ID:   test_utils.NumberUUID(1),
				Code: "super_validation_code_9000",
			},
		},
		{
			name: "Error/CredentialsServiceFailure",
			now:  baseTime,
			form: models.UserValidateEmailForm{
				ID:   test_utils.NumberUUID(1),
				Code: "super_validation_code_9000",
			},
			credentialsErr: fooErr,
			expectErr:      fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)

			credentialsService.
				On("ValidateNewEmail", context.TODO(), d.form.ID, d.form.Code, d.now).
				Return(nil, d.credentialsErr)

			provider := NewProvider(Config{
				CredentialsService: credentialsService,
				Time:               test_utils.GetTimeNow(d.now),
			})

			err := provider.ValidateNewEmail(context.TODO(), d.form)
			test_utils.RequireError(t, d.expectErr, err)

			credentialsService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_ResendEmailValidation(t *testing.T) {
	data := []struct {
		name string

		now                     time.Time
		keys                    []ed25519.PrivateKey
		userID                  uuid.UUID
		emailValidationLink     generics.URL
		emailValidationTemplate string

		token string

		shouldCallTokenServiceDecode bool
		shouldCallCredentialsService bool
		shouldCallIdentityService    bool
		shouldCallMailer             bool
		shouldReturnDeferred         bool

		shouldCallMailerWithEmail    *mail.Email
		shouldCallMailerWithData     map[string]interface{}
		shouldCallMailerWithTemplate string

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		credentialsData        *models.UserCredentials
		credentialsResetCode   string
		credentialsErr         error
		identityData           *models.UserIdentity
		identityErr            error
		mailerErr              error

		expectErr      error
		expectDeferErr error
	}{
		{
			name:   "Success",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailValidationTemplate:      "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "sylvester@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		{
			name:   "Error/MailerFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailValidationTemplate:      "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "sylvester@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			mailerErr:      fooErr,
			expectDeferErr: fooErr,
		},
		{
			name:   "Error/IdentityServiceFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailValidationTemplate:      "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityErr:          fooErr,
			expectErr:            fooErr,
		},
		{
			name:   "Error/CredentialsServiceFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailValidationTemplate:      "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			shouldCallCredentialsService: true,
			credentialsErr:               fooErr,
			expectErr:                    fooErr,
		},
		{
			name:   "Error/TokenServiceFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			emailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			emailValidationTemplate:      "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)
			identityService := identity_service.NewMockService(t)
			profileService := profile_service.NewMockService(t)
			mailerService := mailer.NewMockMailer(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallCredentialsService {
				credentialsService.
					On("UpdateEmailValidation", context.TODO(), d.userID, d.now).
					Return(d.credentialsData, d.credentialsResetCode, d.credentialsErr)
			}

			if d.shouldCallIdentityService {
				identityService.
					On("Read", context.TODO(), d.userID).
					Return(d.identityData, d.identityErr)
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", d.shouldCallMailerWithEmail, d.shouldCallMailerWithTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			provider := NewProvider(Config{
				CredentialsService:      credentialsService,
				IdentityService:         identityService,
				ProfileService:          profileService,
				TokenService:            tokenService,
				KeysService:             keysService,
				Mailer:                  mailerService,
				Time:                    test_utils.GetTimeNow(d.now),
				EmailValidationLink:     d.emailValidationLink,
				EmailValidationTemplate: d.emailValidationTemplate,
			})

			deferred, err := provider.ResendEmailValidation(context.TODO(), d.token)
			test_utils.RequireError(t, d.expectErr, err)

			if d.shouldReturnDeferred {
				require.NotNil(t, deferred)
				test_utils.RequireError(t, d.expectDeferErr, deferred())
			}

			credentialsService.AssertExpectations(t)
			identityService.AssertExpectations(t)
			profileService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
			mailerService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_ResendNewEmailValidation(t *testing.T) {
	data := []struct {
		name string

		now                        time.Time
		keys                       []ed25519.PrivateKey
		userID                     uuid.UUID
		newEmailValidationLink     generics.URL
		newEmailValidationTemplate string

		token string

		shouldCallTokenServiceDecode bool
		shouldCallCredentialsService bool
		shouldCallIdentityService    bool
		shouldCallMailer             bool
		shouldReturnDeferred         bool

		shouldCallMailerWithEmail    *mail.Email
		shouldCallMailerWithData     map[string]interface{}
		shouldCallMailerWithTemplate string

		tokenServiceDecodeData *models.UserToken
		tokenServiceDecodeErr  error
		credentialsData        *models.UserCredentials
		credentialsResetCode   string
		credentialsErr         error
		identityData           *models.UserIdentity
		identityErr            error
		mailerErr              error

		expectErr      error
		expectDeferErr error
	}{
		{
			name:   "Success",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailValidationTemplate:   "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "boss@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		{
			name:   "Error/MailerFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailValidationTemplate:   "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallMailer:             true,
			shouldReturnDeferred:         true,
			shouldCallMailerWithEmail:    mail.NewEmail("Sylvestre", "boss@worldcompany.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "Sylvestre",
				"validation_link": "https://foo.com/bar?code=super_validation_code_9000&id=01010101-0101-0101-0101-010101010101",
			},
			shouldCallMailerWithTemplate: "foo_template",
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityData: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				FirstName: "Sylvestre",
				LastName:  "Stalle",
				Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			mailerErr:      fooErr,
			expectDeferErr: fooErr,
		},
		{
			name:   "Error/IdentityServiceFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailValidationTemplate:   "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			credentialsData: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Email:     "sylvester@worldcompany.com",
				NewEmail:  "boss@worldcompany.com",
				Validated: true,
			},
			credentialsResetCode: "super_validation_code_9000",
			identityErr:          fooErr,
			expectErr:            fooErr,
		},
		{
			name:   "Error/CredentialsServiceFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailValidationTemplate:   "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeData: &models.UserToken{
				Header: models.UserTokenHeader{
					IAT: baseTime.Add(-time.Hour),
					EXP: baseTime.Add(time.Hour),
					ID:  test_utils.NumberUUID(1000),
				},
				Payload: models.UserTokenPayload{
					ID: test_utils.NumberUUID(1),
				},
			},
			shouldCallCredentialsService: true,
			credentialsErr:               fooErr,
			expectErr:                    fooErr,
		},
		{
			name:   "Error/TokenServiceFailure",
			now:    baseTime,
			userID: test_utils.NumberUUID(1),
			newEmailValidationLink: generics.URL{
				Host: "https://foo.com",
				Path: "/bar",
			},
			newEmailValidationTemplate:   "foo_template",
			token:                        "foo.bar.qux",
			shouldCallTokenServiceDecode: true,
			tokenServiceDecodeErr:        fooErr,
			expectErr:                    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)
			identityService := identity_service.NewMockService(t)
			profileService := profile_service.NewMockService(t)
			mailerService := mailer.NewMockMailer(t)
			tokenService := token_service.NewMockService(t)
			keysService := jwk_service.NewMockServiceCached(t)

			publicKeys := make([]ed25519.PublicKey, len(d.keys))
			for i, key := range d.keys {
				publicKeys[i] = key.Public().(ed25519.PublicKey)
			}

			if d.shouldCallTokenServiceDecode {
				keysService.
					On("ListPublic").
					Return(publicKeys)

				tokenService.
					On("Decode", d.token, publicKeys, d.now).
					Return(d.tokenServiceDecodeData, d.tokenServiceDecodeErr)
			}

			if d.shouldCallCredentialsService {
				credentialsService.
					On("UpdateNewEmailValidation", context.TODO(), d.userID, d.now).
					Return(d.credentialsData, d.credentialsResetCode, d.credentialsErr)
			}

			if d.shouldCallIdentityService {
				identityService.
					On("Read", context.TODO(), d.userID).
					Return(d.identityData, d.identityErr)
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", d.shouldCallMailerWithEmail, d.shouldCallMailerWithTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			provider := NewProvider(Config{
				CredentialsService:         credentialsService,
				IdentityService:            identityService,
				ProfileService:             profileService,
				TokenService:               tokenService,
				KeysService:                keysService,
				Mailer:                     mailerService,
				Time:                       test_utils.GetTimeNow(d.now),
				NewEmailValidationLink:     d.newEmailValidationLink,
				NewEMailValidationTemplate: d.newEmailValidationTemplate,
			})

			deferred, err := provider.ResendNewEmailValidation(context.TODO(), d.token)
			test_utils.RequireError(t, d.expectErr, err)

			if d.shouldReturnDeferred {
				require.NotNil(t, deferred)
				test_utils.RequireError(t, d.expectDeferErr, deferred())
			}

			credentialsService.AssertExpectations(t)
			identityService.AssertExpectations(t)
			profileService.AssertExpectations(t)
			tokenService.AssertExpectations(t)
			keysService.AssertExpectations(t)
			mailerService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_DoesSlugExist(t *testing.T) {
	data := []struct {
		name string

		slug string

		profileData bool
		profileErr  error

		expect    bool
		expectErr error
	}{
		{
			name:        "Success",
			slug:        "foo",
			profileData: true,
			expect:      true,
		},
		{
			name:       "Error/ProfileServiceFailure",
			slug:       "foo",
			profileErr: fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			profileService := profile_service.NewMockService(t)

			profileService.
				On("SlugExists", context.TODO(), d.slug).
				Return(d.profileData, d.profileErr)

			provider := NewProvider(Config{
				ProfileService: profileService,
			})

			res, err := provider.DoesSlugExist(context.TODO(), d.slug)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			profileService.AssertExpectations(t)
		})
	}
}

func TestAccountProvider_DoesEmailExist(t *testing.T) {
	data := []struct {
		name string

		email string

		credentialsData bool
		credentialsErr  error

		expect    bool
		expectErr error
	}{
		{
			name:            "Success",
			email:           "user@company.com",
			credentialsData: true,
			expect:          true,
		},
		{
			name:           "Error/CredentialsServiceFailure",
			email:          "user@company.com",
			credentialsErr: fooErr,
			expectErr:      fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsService := credentials_service.NewMockService(t)

			credentialsService.
				On("EmailExists", context.TODO(), d.email).
				Return(d.credentialsData, d.credentialsErr)

			provider := NewProvider(Config{
				CredentialsService: credentialsService,
			})

			res, err := provider.DoesEmailExist(context.TODO(), d.email)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			credentialsService.AssertExpectations(t)
		})
	}
}
