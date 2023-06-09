package credentials_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-novel/agora-backend/framework"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
	"time"
)

var (
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 9, 0, 0, 0, time.UTC)
	fooErr     = errors.New("it broken")
)

var (
	elonBezosStorage = &credentials_storage.Model{
		BaseModel: bun.BaseModel{},
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		Core: credentials_storage.Core{
			Email: models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			Password: models.Password{Hashed: "foobarqux"},
		},
	}
	elonBezosModel = &models.UserCredentials{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		Email:     "elon.bezos@gmail.com",
		Validated: true,
	}
)

func TestCredentialsService_PrepareRegistration(t *testing.T) {
	repository := credentials_storage.NewMockRepository(t)

	data := []struct {
		name string

		data                      *models.UserCredentialsLoginForm
		generateCodeError         error
		generateFromPasswordError error

		expect    *models.UserCredentialsRegistrationForm
		expectErr error
	}{
		{
			name: "Success",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: "123456",
			},
			expect: &models.UserCredentialsRegistrationForm{
				Email: models.Email{
					Validation: "code_hashed",
					User:       "elon.bezos",
					Domain:     "gmail.com",
				},
				Password:                  models.Password{Hashed: "password_hashed"},
				EmailPublicValidationCode: "code",
			},
		},
		{
			name: "Error/MissingEmail",
			data: &models.UserCredentialsLoginForm{
				Password: "123456",
			},
			expectErr: validation.ErrNil,
		},
		{
			name: "Error/MissingPassword",
			data: &models.UserCredentialsLoginForm{
				Email: "elon.bezos@gmail.com",
			},
			expectErr: validation.ErrNil,
		},
		{
			name: "Error/PasswordTooShort",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: "1",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/PasswordTooLong",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: strings.Repeat("1", MaxPasswordLength+1),
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedEmail",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos",
				Password: "123456",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedEmail#2",
			data: &models.UserCredentialsLoginForm{
				Email:    "@gmail.com",
				Password: "123456",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/GenerateEmailValidationFailure",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: "123456",
			},
			generateCodeError: fooErr,
			expectErr:         fooErr,
		},
		{
			name: "Error/EncryptPasswordFailure",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: "123456",
			},
			generateFromPasswordError: fooErr,
			expectErr:                 fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			service := NewService(
				repository,
				test_utils.GetSecurityGenerateCode("code", "code_hashed", d.generateCodeError),
				nil,
				test_utils.GetBcryptGenerateFromPassword("password_hashed", d.generateFromPasswordError),
				nil,
			)

			res, err := service.PrepareRegistration(context.TODO(), d.data)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_Authenticate(t *testing.T) {
	data := []struct {
		name string

		data                 *models.UserCredentialsLoginForm
		getUserData          *credentials_storage.Model
		comparePasswordError error
		getUserError         error

		shouldCallRepositoryWith *models.Email

		expect    *models.UserCredentials
		expectErr error
	}{
		{
			name: "Success",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: "foobarqux",
			},
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserData: elonBezosStorage,
			expect:      elonBezosModel,
		},
		{
			name: "Error/MissingEmail",
			data: &models.UserCredentialsLoginForm{
				Password: "foobarqux",
			},
			getUserData: elonBezosStorage,
			expectErr:   validation.ErrInvalidEntity,
		},
		{
			name: "Error/MissingPassword",
			data: &models.UserCredentialsLoginForm{
				Email: "elon.bezos@gmail.com",
			},
			getUserData: elonBezosStorage,
			expectErr:   validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedEmail",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@",
				Password: "foobarqux",
			},
			getUserData: elonBezosStorage,
			expectErr:   validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedEmail#2",
			data: &models.UserCredentialsLoginForm{
				Email:    "gmail.com",
				Password: "foobarqux",
			},
			getUserData: elonBezosStorage,
			expectErr:   validation.ErrInvalidEntity,
		},
		{
			name: "Error/GetUserFailure",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: "foobarqux",
			},
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserError: validation.ErrNotFound,
			expectErr:    validation.ErrNotFound,
		},
		{
			name: "Error/BadPassword",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: "foobarqux",
			},
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserData:          elonBezosStorage,
			comparePasswordError: bcrypt.ErrMismatchedHashAndPassword,
			expectErr:            validation.ErrInvalidCredentials,
		},
		{
			name: "Error/CheckPasswordFailure",
			data: &models.UserCredentialsLoginForm{
				Email:    "elon.bezos@gmail.com",
				Password: "foobarqux",
			},
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserData:          elonBezosStorage,
			comparePasswordError: fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallRepositoryWith != nil {
				repository.
					On("ReadEmail", context.TODO(), *d.shouldCallRepositoryWith).
					Return(d.getUserData, d.getUserError)
			}

			service := NewService(
				repository, nil, nil, nil,
				test_utils.GetBcryptCompareHashAndPassword(d.comparePasswordError),
			)

			res, err := service.Authenticate(context.TODO(), d.data)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_Read(t *testing.T) {
	data := []struct {
		name string

		id           uuid.UUID
		getUserData  *credentials_storage.Model
		getUserError error

		expect    *models.UserCredentials
		expectErr error
	}{
		{
			name:        "Success",
			id:          test_utils.NumberUUID(1000),
			getUserData: elonBezosStorage,
			expect:      elonBezosModel,
		},
		{
			name:         "Error/RepositoryFailure",
			id:           test_utils.NumberUUID(1000),
			getUserError: validation.ErrNotFound,
			expectErr:    validation.ErrNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)
			repository.
				On("Read", context.TODO(), d.id).
				Return(d.getUserData, d.getUserError)

			service := NewService(repository, nil, nil, nil, nil)

			res, err := service.Read(context.TODO(), d.id)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_ReadEmail(t *testing.T) {
	data := []struct {
		name string

		email        string
		getUserData  *credentials_storage.Model
		getUserError error

		shouldCallRepositoryWith *models.Email

		expect    *models.UserCredentials
		expectErr error
	}{
		{
			name:  "Success",
			email: "elon.bezos@gmail.com",
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserData: elonBezosStorage,
			expect:      elonBezosModel,
		},
		{
			name:  "Error/RepositoryFailure",
			email: "elon.bezos@gmail.com",
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserError: validation.ErrNotFound,
			expectErr:    validation.ErrNotFound,
		},
		{
			name:      "Error/NoEmail",
			expectErr: validation.ErrNil,
		},
		{
			name:      "Error/MalformedEmail#1",
			email:     "elon.bezos@",
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:      "Error/MalformedEmail#2",
			email:     "gmail.com",
			expectErr: validation.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallRepositoryWith != nil {
				repository.
					On("ReadEmail", context.TODO(), *d.shouldCallRepositoryWith).
					Return(d.getUserData, d.getUserError)
			}

			service := NewService(repository, nil, nil, nil, nil)

			res, err := service.ReadEmail(context.TODO(), d.email)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_EmailExists(t *testing.T) {
	data := []struct {
		name string

		email         string
		getUserExists bool
		getUserError  error

		shouldCallRepositoryWith *models.Email

		expect    bool
		expectErr error
	}{
		{
			name:  "Success/Exists",
			email: "elon.bezos@gmail.com",
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserExists: true,
			expect:        true,
		},
		{
			name:  "Success/DoesNotExists",
			email: "elon.bezos@gmail.com",
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserExists: false,
			expect:        false,
		},
		{
			name:  "Error/RepositoryFailure",
			email: "elon.bezos@gmail.com",
			shouldCallRepositoryWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			getUserError: validation.ErrNotFound,
			expectErr:    validation.ErrNotFound,
		},
		{
			name:      "Error/NoEmail",
			expectErr: validation.ErrNil,
		},
		{
			name:      "Error/MalformedEmail#1",
			email:     "elon.bezos@",
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:      "Error/MalformedEmail#2",
			email:     "gmail.com",
			expectErr: validation.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallRepositoryWith != nil {
				repository.
					On("EmailExists", context.TODO(), *d.shouldCallRepositoryWith).
					Return(d.getUserExists, d.getUserError)
			}

			service := NewService(repository, nil, nil, nil, nil)

			res, err := service.EmailExists(context.TODO(), d.email)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_UpdateEmail(t *testing.T) {
	data := []struct {
		name string

		email string
		id    uuid.UUID
		now   time.Time

		getEmailExists      bool
		getUserData         *credentials_storage.Model
		getEmailExistsError error
		getUserDataError    error
		generateCodeError   error

		shouldCallRepositoryEmailExistsWith      *models.Email
		shouldCallRepositoryUpdateEmailWithEmail *models.Email

		expect     *models.UserCredentials
		expectCode string
		expectErr  error
	}{
		{
			name:           "Success",
			email:          "anna.banana@coco.nut",
			id:             test_utils.NumberUUID(1000),
			now:            updateTime,
			getEmailExists: false,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRepositoryEmailExistsWith: &models.Email{
				User:   "anna.banana",
				Domain: "coco.nut",
			},
			shouldCallRepositoryUpdateEmailWithEmail: &models.Email{
				User:   "anna.banana",
				Domain: "coco.nut",
			},
			expect: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Email:     "elon.bezos@gmail.com",
				NewEmail:  "anna.banana@coco.nut",
				Validated: true,
			},
			expectCode: "code",
		},
		{
			name:           "Error/NoEmail",
			id:             test_utils.NumberUUID(1000),
			now:            updateTime,
			getEmailExists: false,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: validation.ErrNil,
		},
		{
			name:           "Error/MalformedEmail#1",
			email:          "anna.banana@",
			id:             test_utils.NumberUUID(1000),
			now:            updateTime,
			getEmailExists: false,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:           "Error/MalformedEmail#2",
			email:          "gmail.com",
			id:             test_utils.NumberUUID(1000),
			now:            updateTime,
			getEmailExists: false,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:           "Error/EmailTaken",
			email:          "anna.banana@coco.nut",
			id:             test_utils.NumberUUID(1000),
			now:            updateTime,
			getEmailExists: true,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRepositoryEmailExistsWith: &models.Email{
				User:   "anna.banana",
				Domain: "coco.nut",
			},
			expectErr: validation.ErrUniqConstraintViolation,
		},
		{
			name:                "Error/RepositoryFailureOnEmailExists",
			email:               "anna.banana@coco.nut",
			id:                  test_utils.NumberUUID(1000),
			now:                 updateTime,
			getEmailExists:      true,
			getEmailExistsError: fooErr,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRepositoryEmailExistsWith: &models.Email{
				User:   "anna.banana",
				Domain: "coco.nut",
			},
			expectErr: fooErr,
		},
		{
			name:              "Error/GenerateCodeFailure",
			email:             "anna.banana@coco.nut",
			id:                test_utils.NumberUUID(1000),
			now:               updateTime,
			getEmailExists:    false,
			generateCodeError: fooErr,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRepositoryEmailExistsWith: &models.Email{
				User:   "anna.banana",
				Domain: "coco.nut",
			},
			expectErr: fooErr,
		},
		{
			name:           "Error/UpdateEmailFailure",
			email:          "anna.banana@coco.nut",
			id:             test_utils.NumberUUID(1000),
			now:            updateTime,
			getEmailExists: false,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			getUserDataError: fooErr,
			shouldCallRepositoryEmailExistsWith: &models.Email{
				User:   "anna.banana",
				Domain: "coco.nut",
			},
			shouldCallRepositoryUpdateEmailWithEmail: &models.Email{
				User:   "anna.banana",
				Domain: "coco.nut",
			},
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallRepositoryEmailExistsWith != nil {
				repository.
					On("EmailExists", context.TODO(), *d.shouldCallRepositoryEmailExistsWith).
					Return(d.getEmailExists, d.getEmailExistsError)
			}

			if d.shouldCallRepositoryUpdateEmailWithEmail != nil {
				repository.
					On(
						"UpdateEmail",
						context.TODO(), *d.shouldCallRepositoryUpdateEmailWithEmail, "code_hashed", d.id, d.now,
					).
					Return(d.getUserData, d.getUserDataError)
			}

			service := NewService(
				repository,
				test_utils.GetSecurityGenerateCode("code", "code_hashed", d.generateCodeError),
				nil, nil, nil,
			)

			res, validationCode, err := service.UpdateEmail(context.TODO(), d.email, d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)
			require.Equal(t, d.expectCode, validationCode)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_ValidateEmail(t *testing.T) {
	data := []struct {
		name string

		id   uuid.UUID
		code string
		now  time.Time

		verifyCodeStatus bool
		verifyCodeError  error

		getUserData        *credentials_storage.Model
		getUserDataError   error
		validateEmailData  *credentials_storage.Model
		validateEmailError error

		shouldCallReadWith          *uuid.UUID
		shouldCallValidateEmailWith *uuid.UUID

		expect    *models.UserCredentials
		expectErr error
	}{
		{
			name:                        "Success",
			id:                          test_utils.NumberUUID(100),
			code:                        "code",
			now:                         updateTime,
			verifyCodeStatus:            true,
			shouldCallReadWith:          framework.ToPTR(test_utils.NumberUUID(100)),
			shouldCallValidateEmailWith: &elonBezosStorage.ID,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			validateEmailData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expect: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Email:     "elon.bezos@gmail.com",
				Validated: true,
			},
		},
		{
			name:      "Error/NoCode",
			id:        test_utils.NumberUUID(100),
			now:       updateTime,
			expectErr: validation.ErrNil,
		},
		{
			name:               "Error/ReadEmailFailure",
			id:                 test_utils.NumberUUID(100),
			code:               "code",
			now:                updateTime,
			shouldCallReadWith: framework.ToPTR(test_utils.NumberUUID(100)),
			getUserDataError:   fooErr,
			expectErr:          fooErr,
		},
		{
			name:               "Error/NoPendingValidation",
			id:                 test_utils.NumberUUID(100),
			code:               "code",
			now:                updateTime,
			shouldCallReadWith: framework.ToPTR(test_utils.NumberUUID(100)),
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: validation.ErrValidated,
		},
		{
			name:               "Error/VerifyCodeFailure",
			id:                 test_utils.NumberUUID(100),
			code:               "code",
			now:                updateTime,
			verifyCodeError:    fooErr,
			shouldCallReadWith: framework.ToPTR(test_utils.NumberUUID(100)),
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: fooErr,
		},
		{
			name:               "Error/WrongCode",
			id:                 test_utils.NumberUUID(100),
			code:               "code",
			now:                updateTime,
			verifyCodeStatus:   false,
			shouldCallReadWith: framework.ToPTR(test_utils.NumberUUID(100)),
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: validation.ErrInvalidCredentials,
		},
		{
			name:                        "Error/ValidateEmailFailure",
			id:                          test_utils.NumberUUID(100),
			code:                        "code",
			now:                         updateTime,
			verifyCodeStatus:            true,
			shouldCallReadWith:          framework.ToPTR(test_utils.NumberUUID(100)),
			shouldCallValidateEmailWith: &elonBezosStorage.ID,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			validateEmailError: fooErr,
			expectErr:          fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallReadWith != nil {
				repository.
					On("Read", context.TODO(), *d.shouldCallReadWith).
					Return(d.getUserData, d.getUserDataError)
			}

			if d.shouldCallValidateEmailWith != nil {
				repository.
					On("ValidateEmail", context.TODO(), *d.shouldCallValidateEmailWith, d.now).
					Return(d.validateEmailData, d.validateEmailError)
			}

			service := NewService(
				repository, nil,
				test_utils.GetSecurityVerifyCode(d.verifyCodeStatus, d.verifyCodeError),
				nil, nil,
			)

			res, err := service.ValidateEmail(context.TODO(), d.id, d.code, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_ValidateNewEmail(t *testing.T) {
	data := []struct {
		name string

		id   uuid.UUID
		code string
		now  time.Time

		verifyCodeStatus bool
		verifyCodeError  error

		getUserData        *credentials_storage.Model
		getUserDataError   error
		validateEmailData  *credentials_storage.Model
		validateEmailError error

		shouldCallReadWith          *uuid.UUID
		shouldCallValidateEmailWith *uuid.UUID

		expect    *models.UserCredentials
		expectErr error
	}{
		{
			name:                        "Success",
			id:                          test_utils.NumberUUID(100),
			code:                        "code",
			now:                         updateTime,
			verifyCodeStatus:            true,
			shouldCallReadWith:          framework.ToPTR(test_utils.NumberUUID(100)),
			shouldCallValidateEmailWith: &elonBezosStorage.ID,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			validateEmailData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "anna.banana",
						Domain: "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expect: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Email:     "anna.banana@coco.nut",
				Validated: true,
			},
		},
		{
			name:      "Error/NoCode",
			id:        test_utils.NumberUUID(100),
			now:       updateTime,
			expectErr: validation.ErrNil,
		},
		{
			name:               "Error/ReadEmailFailure",
			id:                 test_utils.NumberUUID(100),
			code:               "code",
			now:                updateTime,
			shouldCallReadWith: framework.ToPTR(test_utils.NumberUUID(100)),
			getUserDataError:   fooErr,
			expectErr:          fooErr,
		},
		{
			name:               "Error/NoPendingValidation",
			id:                 test_utils.NumberUUID(100),
			code:               "code",
			now:                updateTime,
			shouldCallReadWith: framework.ToPTR(test_utils.NumberUUID(100)),
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: validation.ErrValidated,
		},
		{
			name:               "Error/VerifyCodeFailure",
			id:                 test_utils.NumberUUID(100),
			code:               "code",
			now:                updateTime,
			verifyCodeError:    fooErr,
			shouldCallReadWith: framework.ToPTR(test_utils.NumberUUID(100)),
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: fooErr,
		},
		{
			name:               "Error/WrongCode",
			id:                 test_utils.NumberUUID(100),
			code:               "code",
			now:                updateTime,
			verifyCodeStatus:   false,
			shouldCallReadWith: framework.ToPTR(test_utils.NumberUUID(100)),
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expectErr: validation.ErrInvalidCredentials,
		},
		{
			name:                        "Error/ValidateEmailFailure",
			id:                          test_utils.NumberUUID(100),
			code:                        "code",
			now:                         updateTime,
			verifyCodeStatus:            true,
			shouldCallReadWith:          framework.ToPTR(test_utils.NumberUUID(100)),
			shouldCallValidateEmailWith: &elonBezosStorage.ID,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			validateEmailError: fooErr,
			expectErr:          fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallReadWith != nil {
				repository.
					On("Read", context.TODO(), *d.shouldCallReadWith).
					Return(d.getUserData, d.getUserDataError)
			}

			if d.shouldCallValidateEmailWith != nil {
				repository.
					On("ValidateNewEmail", context.TODO(), *d.shouldCallValidateEmailWith, d.now).
					Return(d.validateEmailData, d.validateEmailError)
			}

			service := NewService(
				repository, nil,
				test_utils.GetSecurityVerifyCode(d.verifyCodeStatus, d.verifyCodeError),
				nil, nil,
			)

			res, err := service.ValidateNewEmail(context.TODO(), d.id, d.code, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_UpdateEmailValidation(t *testing.T) {
	data := []struct {
		name string

		id  uuid.UUID
		now time.Time

		getUserData          *credentials_storage.Model
		getUserDataError     error
		updateEmailData      *credentials_storage.Model
		updateEmailDataError error
		generateCodeError    error

		shouldCallRead   bool
		shouldCallUpdate bool

		expect     *models.UserCredentials
		expectCode string
		expectErr  error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			updateEmailData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "code_hashed",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRead:   true,
			shouldCallUpdate: true,
			expect: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Email:     "elon.bezos@gmail.com",
				Validated: false,
			},
			expectCode: "code",
		},
		{
			name:             "Error/RepositoryReadFailure",
			id:               test_utils.NumberUUID(1000),
			now:              updateTime,
			getUserDataError: fooErr,
			shouldCallRead:   true,
			expectErr:        fooErr,
		},
		{
			name: "Error/NoPendingValidation",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRead: true,
			expectErr:      validation.ErrValidated,
		},
		{
			name: "Error/GenerateCodeFailure",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRead:    true,
			generateCodeError: fooErr,
			expectErr:         fooErr,
		},
		{
			name: "Error/UpdateEmailValidationFailure",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			updateEmailData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "code_hashed",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRead:       true,
			shouldCallUpdate:     true,
			updateEmailDataError: fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallRead {
				repository.
					On("Read", context.TODO(), d.id).
					Return(d.getUserData, d.getUserDataError)
			}

			if d.shouldCallUpdate {
				repository.
					On("UpdateEmailValidation", context.TODO(), "code_hashed", d.id, d.now).
					Return(d.updateEmailData, d.updateEmailDataError)
			}

			service := NewService(
				repository,
				test_utils.GetSecurityGenerateCode("code", "code_hashed", d.generateCodeError),
				nil, nil, nil,
			)

			res, validationCode, err := service.UpdateEmailValidation(context.TODO(), d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)
			require.Equal(t, d.expectCode, validationCode)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_UpdateNewEmailValidation(t *testing.T) {
	data := []struct {
		name string

		id  uuid.UUID
		now time.Time

		getUserData          *credentials_storage.Model
		getUserDataError     error
		updateEmailData      *credentials_storage.Model
		updateEmailDataError error
		generateCodeError    error

		shouldCallRead   bool
		shouldCallUpdate bool

		expect     *models.UserCredentials
		expectCode string
		expectErr  error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			updateEmailData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRead:   true,
			shouldCallUpdate: true,
			expect: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Email:     "elon.bezos@gmail.com",
				NewEmail:  "anna.banana@coco.nut",
				Validated: true,
			},
			expectCode: "code",
		},
		{
			name:             "Error/RepositoryReadFailure",
			id:               test_utils.NumberUUID(1000),
			now:              updateTime,
			getUserDataError: fooErr,
			shouldCallRead:   true,
			expectErr:        fooErr,
		},
		{
			name: "Error/NoPendingValidation",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRead: true,
			expectErr:      validation.ErrValidated,
		},
		{
			name: "Error/GenerateCodeFailure",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRead:    true,
			generateCodeError: fooErr,
			expectErr:         fooErr,
		},
		{
			name: "Error/UpdateEmailValidationFailure",
			id:   test_utils.NumberUUID(1000),
			now:  updateTime,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "youshallpass",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			updateEmailData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					NewEmail: models.Email{
						Validation: "code_hashed",
						User:       "anna.banana",
						Domain:     "coco.nut",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			shouldCallRead:       true,
			shouldCallUpdate:     true,
			updateEmailDataError: fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallRead {
				repository.
					On("Read", context.TODO(), d.id).
					Return(d.getUserData, d.getUserDataError)
			}

			if d.shouldCallUpdate {
				repository.
					On("UpdateNewEmailValidation", context.TODO(), "code_hashed", d.id, d.now).
					Return(d.updateEmailData, d.updateEmailDataError)
			}

			service := NewService(
				repository,
				test_utils.GetSecurityGenerateCode("code", "code_hashed", d.generateCodeError),
				nil, nil, nil,
			)

			res, validationCode, err := service.UpdateNewEmailValidation(context.TODO(), d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)
			require.Equal(t, d.expectCode, validationCode)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_CancelNewEmail(t *testing.T) {
	data := []struct {
		name string

		id  uuid.UUID
		now time.Time

		cancelEmailData      *credentials_storage.Model
		cancelEmailDataError error

		expect    *models.UserCredentials
		expectErr error
	}{
		{
			name:            "Success",
			id:              test_utils.NumberUUID(1000),
			now:             updateTime,
			cancelEmailData: elonBezosStorage,
			expect:          elonBezosModel,
		},
		{
			name:                 "Error/RepositoryFailure",
			id:                   test_utils.NumberUUID(1000),
			now:                  updateTime,
			cancelEmailDataError: fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			repository.
				On("CancelNewEmail", context.TODO(), d.id, d.now).
				Return(d.cancelEmailData, d.cancelEmailDataError)

			service := NewService(repository, nil, nil, nil, nil)

			res, err := service.CancelNewEmail(context.TODO(), d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_UpdatePassword(t *testing.T) {
	data := []struct {
		name string

		oldPassword string
		newPassword string
		id          uuid.UUID
		now         time.Time

		getUserData             *credentials_storage.Model
		getUserDataError        error
		updatePasswordData      *credentials_storage.Model
		updatePasswordDataError error

		shouldCallRead   bool
		shouldCallUpdate bool

		verifyCodeStatus bool
		verifyCodeError  error

		compareHashAndPasswordStatuses map[string]error

		expect    *models.UserCredentials
		expectErr error
	}{
		{
			name:                           "Success",
			id:                             test_utils.NumberUUID(1000),
			now:                            updateTime,
			oldPassword:                    "foobarqux",
			newPassword:                    "quxbarfoo",
			getUserData:                    elonBezosStorage,
			updatePasswordData:             elonBezosStorage,
			shouldCallRead:                 true,
			shouldCallUpdate:               true,
			compareHashAndPasswordStatuses: map[string]error{"foobarqux": nil},
			expect:                         elonBezosModel,
		},
		{
			name:             "Success/WithResetCode",
			id:               test_utils.NumberUUID(1000),
			now:              updateTime,
			oldPassword:      "code",
			newPassword:      "quxbarfoo",
			verifyCodeStatus: true,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux", Validation: "code_hashed"},
				},
			},
			updatePasswordData: elonBezosStorage,
			shouldCallRead:     true,
			shouldCallUpdate:   true,
			expect:             elonBezosModel,
		},
		{
			name:        "Success/WithCurrentPasswordWithActiveReset",
			id:          test_utils.NumberUUID(1000),
			now:         updateTime,
			oldPassword: "foobarqux",
			newPassword: "quxbarfoo",
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux", Validation: "code_hashed"},
				},
			},
			updatePasswordData: elonBezosStorage,
			shouldCallRead:     true,
			shouldCallUpdate:   true,
			compareHashAndPasswordStatuses: map[string]error{
				"code_hashed": bcrypt.ErrMismatchedHashAndPassword,
				"foobarqux":   nil,
			},
			expect: elonBezosModel,
		},
		{
			name:            "Error/ResetCodeValidationFailure",
			id:              test_utils.NumberUUID(1000),
			now:             updateTime,
			oldPassword:     "foobarqux",
			newPassword:     "quxbarfoo",
			verifyCodeError: fooErr,
			getUserData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux", Validation: "code_hashed"},
				},
			},
			updatePasswordData: elonBezosStorage,
			shouldCallRead:     true,
			expectErr:          fooErr,
		},
		{
			name:        "Error/NoNewPassword",
			id:          test_utils.NumberUUID(1000),
			now:         updateTime,
			oldPassword: "foobarqux",
			expectErr:   validation.ErrNil,
		},
		{
			name:        "Error/NoOldPassword",
			id:          test_utils.NumberUUID(1000),
			now:         updateTime,
			newPassword: "quxbarfoo",
			expectErr:   validation.ErrNil,
		},
		{
			name:        "Error/NewPasswordTooShort",
			id:          test_utils.NumberUUID(1000),
			now:         updateTime,
			oldPassword: "foobarqux",
			newPassword: "f",
			expectErr:   validation.ErrInvalidEntity,
		},
		{
			name:        "Error/NewPasswordTooLong",
			id:          test_utils.NumberUUID(1000),
			now:         updateTime,
			oldPassword: "foobarqux",
			newPassword: strings.Repeat("f", MaxPasswordLength+1),
			expectErr:   validation.ErrInvalidEntity,
		},
		{
			name:             "Error/RepositoryReadError",
			id:               test_utils.NumberUUID(1000),
			now:              updateTime,
			oldPassword:      "foobarqux",
			newPassword:      "quxbarfoo",
			getUserDataError: fooErr,
			shouldCallRead:   true,
			expectErr:        fooErr,
		},
		{
			name:                           "Error/RepositoryUpdateError",
			id:                             test_utils.NumberUUID(1000),
			now:                            updateTime,
			oldPassword:                    "foobarqux",
			newPassword:                    "quxbarfoo",
			getUserData:                    elonBezosStorage,
			updatePasswordDataError:        fooErr,
			shouldCallRead:                 true,
			shouldCallUpdate:               true,
			compareHashAndPasswordStatuses: map[string]error{"foobarqux": nil},
			expectErr:                      fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallRead {
				repository.
					On("Read", context.TODO(), d.id).
					Return(d.getUserData, d.getUserDataError)
			}

			if d.shouldCallUpdate {
				repository.
					On("UpdatePassword", context.TODO(), "password_hashed", d.id, d.now).
					Return(d.updatePasswordData, d.updatePasswordDataError)
			}

			service := NewService(
				repository, nil,
				test_utils.GetSecurityVerifyCode(d.verifyCodeStatus, d.verifyCodeError),
				test_utils.GetBcryptGenerateFromPassword("password_hashed", nil),
				func(hashedPassword []byte, password []byte) error {
					err, ok := d.compareHashAndPasswordStatuses[string(hashedPassword)]
					if ok {
						return err
					}

					return fmt.Errorf(
						"unexpected call to compareHashAndPassword with hashedPassword %s", string(hashedPassword),
					)
				},
			)

			res, err := service.UpdatePassword(context.TODO(), d.oldPassword, d.newPassword, d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_ResetPassword(t *testing.T) {
	data := []struct {
		name string

		email string
		now   time.Time

		getUserData            *credentials_storage.Model
		getUserDataError       error
		resetPasswordData      *credentials_storage.Model
		resetPasswordDataError error
		generateCodeError      error

		shouldCallReadWith   *models.Email
		shouldCallUpdateWith *models.Email

		expect     *models.UserCredentials
		expectCode string
		expectErr  error
	}{
		{
			name:        "Success",
			email:       "elon.bezos@gmail.com",
			now:         updateTime,
			getUserData: elonBezosStorage,
			resetPasswordData: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						User:   "elon.bezos",
						Domain: "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux", Validation: "youshallpass"},
				},
			},
			shouldCallReadWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			shouldCallUpdateWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			expect:     elonBezosModel,
			expectCode: "code",
		},
		{
			name:      "Error/NoEmail",
			now:       updateTime,
			expectErr: validation.ErrNil,
		},
		{
			name:      "Error/MalformedEmail#1",
			email:     "elon.bezos@",
			now:       updateTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:      "Error/MalformedEmail#2",
			email:     "gmail.com",
			now:       updateTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:             "Error/RepositoryReadFailure",
			email:            "elon.bezos@gmail.com",
			now:              updateTime,
			getUserDataError: fooErr,
			shouldCallReadWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			expectErr: fooErr,
		},
		{
			name:        "Error/GenerateCodeFailure",
			email:       "elon.bezos@gmail.com",
			now:         updateTime,
			getUserData: elonBezosStorage,
			shouldCallReadWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			generateCodeError: fooErr,
			expectErr:         fooErr,
		},
		{
			name:                   "Error/RepositoryResetPasswordFailure",
			email:                  "elon.bezos@gmail.com",
			now:                    updateTime,
			getUserData:            elonBezosStorage,
			resetPasswordDataError: fooErr,
			shouldCallReadWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			shouldCallUpdateWith: &models.Email{
				User:   "elon.bezos",
				Domain: "gmail.com",
			},
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := credentials_storage.NewMockRepository(t)

			if d.shouldCallReadWith != nil {
				repository.
					On("ReadEmail", context.TODO(), *d.shouldCallReadWith).
					Return(d.getUserData, d.getUserDataError)
			}

			if d.shouldCallUpdateWith != nil {
				repository.
					On("ResetPassword", context.TODO(), "code_hashed", *d.shouldCallUpdateWith, d.now).
					Return(d.resetPasswordData, d.resetPasswordDataError)
			}

			service := NewService(
				repository,
				test_utils.GetSecurityGenerateCode("code", "code_hashed", d.generateCodeError),
				nil, nil, nil,
			)

			res, validationCode, err := service.ResetPassword(context.TODO(), d.email, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)
			require.Equal(t, d.expectCode, validationCode)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestCredentialsService_StorageToModel(t *testing.T) {
	repository := credentials_storage.NewMockRepository(t)

	data := []struct {
		name   string
		data   *credentials_storage.Model
		expect *models.UserCredentials
	}{
		{
			name:   "Success",
			data:   elonBezosStorage,
			expect: elonBezosModel,
		},
		{
			name: "Success/Nil",
		},
		{
			name: "Success/EmailNotValidated",
			data: &credentials_storage.Model{
				BaseModel: bun.BaseModel{},
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Core: credentials_storage.Core{
					Email: models.Email{
						Validation: "youshallpass",
						User:       "elon.bezos",
						Domain:     "gmail.com",
					},
					Password: models.Password{Hashed: "foobarqux"},
				},
			},
			expect: &models.UserCredentials{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				Email:     "elon.bezos@gmail.com",
				Validated: false,
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			service := NewService(repository, nil, nil, nil, nil)

			res := service.StorageToModel(d.data)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}
