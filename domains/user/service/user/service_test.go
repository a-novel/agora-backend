package user_service

import (
	"context"
	"errors"
	"github.com/a-novel/agora-backend/domains/user/service/credentials"
	"github.com/a-novel/agora-backend/domains/user/service/identity"
	"github.com/a-novel/agora-backend/domains/user/service/profile"
	"github.com/a-novel/agora-backend/domains/user/storage/credentials"
	"github.com/a-novel/agora-backend/domains/user/storage/identity"
	"github.com/a-novel/agora-backend/domains/user/storage/profile"
	"github.com/a-novel/agora-backend/domains/user/storage/user"
	"github.com/a-novel/agora-backend/framework/test"
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

func TestUserService_Create(t *testing.T) {
	data := []struct {
		name string

		data *models.UserCreateForm
		id   uuid.UUID
		now  time.Time

		shouldCallCredentialsService bool
		shouldCallIdentityService    bool
		shouldCallProfileService     bool
		shouldCallRepository         bool

		expectCredentialsModel *models.UserCredentialsRegistrationForm
		expectIdentityModel    *models.UserIdentityRegistrationForm
		expectProfileModel     *models.UserProfileRegistrationForm

		expectCredentialsError error
		expectIdentityError    error
		expectProfileError     error

		expectCreateUserData  *user_storage.Model
		expectCreateUserError error

		expectModel         *models.User
		expectRegisterModel *models.UserPostRegistration
		expectError         error
	}{
		{
			name: "Success",
			data: &models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "anna.banana@fruit.basket",
					Password: "123456",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Anna",
					LastName:  "Banana",
					Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "BananaSplit",
					Slug:     "fruit-basket-anna",
				},
			},
			id:                           test_utils.NumberUUID(1),
			now:                          baseTime,
			shouldCallRepository:         true,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallProfileService:     true,
			expectCredentialsModel: &models.UserCredentialsRegistrationForm{
				Email: models.Email{
					User:   "anna.banana",
					Domain: "fruit.basket",
				},
				Password: models.Password{
					Hashed: "bulibublablub",
				},
				EmailPublicValidationCode: "foobarqux",
			},
			expectIdentityModel: &models.UserIdentityRegistrationForm{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			expectProfileModel: &models.UserProfileRegistrationForm{
				Username: "BananaSplit",
				Slug:     "fruit-basket-anna",
			},
			expectCreateUserData: &user_storage.Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Core: user_storage.Core{
					Credentials: credentials_storage.Core{
						Email: models.Email{
							User:   "anna.banana",
							Domain: "fruit.basket",
						},
						Password: models.Password{
							Hashed: "bulibublablub",
						},
					},
					Identity: identity_storage.Core{
						FirstName: "Anna",
						LastName:  "Banana",
						Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
						Sex:       models.SexFemale,
					},
					Profile: profile_storage.Core{
						Username: "BananaSplit",
						Slug:     "fruit-basket-anna",
					},
				},
			},
			expectModel: &models.User{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				Credentials: models.UserCredentials{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					Email:     "anna.banana@fruit.basket",
				},
				Identity: models.UserIdentity{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					FirstName: "Anna",
					LastName:  "Banana",
					Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
				Profile: models.UserProfile{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					Username:  "BananaSplit",
					Slug:      "fruit-basket-anna",
				},
			},
			expectRegisterModel: &models.UserPostRegistration{EmailValidationCode: "foobarqux"},
		},
		{
			name: "Error/CredentialsServiceFailure",
			data: &models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "anna.banana@fruit.basket",
					Password: "123456",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Anna",
					LastName:  "Banana",
					Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "BananaSplit",
					Slug:     "fruit-basket-anna",
				},
			},
			id:                           test_utils.NumberUUID(1),
			now:                          baseTime,
			shouldCallCredentialsService: true,
			expectCredentialsError:       fooErr,
			expectError:                  fooErr,
		},
		{
			name: "Error/IdentityServiceFailure",
			data: &models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "anna.banana@fruit.basket",
					Password: "123456",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Anna",
					LastName:  "Banana",
					Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "BananaSplit",
					Slug:     "fruit-basket-anna",
				},
			},
			id:                           test_utils.NumberUUID(1),
			now:                          baseTime,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			expectCredentialsModel: &models.UserCredentialsRegistrationForm{
				Email: models.Email{
					User:   "anna.banana",
					Domain: "fruit.basket",
				},
				Password: models.Password{
					Hashed: "bulibublablub",
				},
				EmailPublicValidationCode: "foobarqux",
			},
			expectIdentityError: fooErr,
			expectError:         fooErr,
		},
		{
			name: "Error/ProfileServiceFailure",
			data: &models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "anna.banana@fruit.basket",
					Password: "123456",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Anna",
					LastName:  "Banana",
					Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "BananaSplit",
					Slug:     "fruit-basket-anna",
				},
			},
			id:                           test_utils.NumberUUID(1),
			now:                          baseTime,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallProfileService:     true,
			expectCredentialsModel: &models.UserCredentialsRegistrationForm{
				Email: models.Email{
					User:   "anna.banana",
					Domain: "fruit.basket",
				},
				Password: models.Password{
					Hashed: "bulibublablub",
				},
				EmailPublicValidationCode: "foobarqux",
			},
			expectIdentityModel: &models.UserIdentityRegistrationForm{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			expectProfileError: fooErr,
			expectError:        fooErr,
		},
		{
			name: "Error/RepositoryFailure",
			data: &models.UserCreateForm{
				Credentials: models.UserCredentialsLoginForm{
					Email:    "anna.banana@fruit.basket",
					Password: "123456",
				},
				Identity: models.UserIdentityUpdateForm{
					FirstName: "Anna",
					LastName:  "Banana",
					Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
				Profile: models.UserProfileUpdateForm{
					Username: "BananaSplit",
					Slug:     "fruit-basket-anna",
				},
			},
			id:                           test_utils.NumberUUID(1),
			now:                          baseTime,
			shouldCallRepository:         true,
			shouldCallCredentialsService: true,
			shouldCallIdentityService:    true,
			shouldCallProfileService:     true,
			expectCredentialsModel: &models.UserCredentialsRegistrationForm{
				Email: models.Email{
					User:   "anna.banana",
					Domain: "fruit.basket",
				},
				Password: models.Password{
					Hashed: "bulibublablub",
				},
				EmailPublicValidationCode: "foobarqux",
			},
			expectIdentityModel: &models.UserIdentityRegistrationForm{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			expectProfileModel: &models.UserProfileRegistrationForm{
				Username: "BananaSplit",
				Slug:     "fruit-basket-anna",
			},
			expectCreateUserError: fooErr,
			expectError:           fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := user_storage.NewMockRepository(t)
			credentialsService := credentials_service.NewMockService(t)
			identityService := identity_service.NewMockService(t)
			profileService := profile_service.NewMockService(t)

			service := NewService(repository, credentialsService, identityService, profileService)

			if d.shouldCallCredentialsService {
				credentialsService.
					On("PrepareRegistration", context.TODO(), &d.data.Credentials).
					Return(d.expectCredentialsModel, d.expectCredentialsError)
			}

			if d.shouldCallIdentityService {
				identityService.
					On("PrepareRegistration", context.TODO(), &d.data.Identity, d.now).
					Return(d.expectIdentityModel, d.expectIdentityError)
			}

			if d.shouldCallProfileService {
				profileService.
					On("PrepareRegistration", context.TODO(), &d.data.Profile).
					Return(d.expectProfileModel, d.expectProfileError)
			}

			if d.shouldCallRepository {
				repository.
					On("Create", context.TODO(), mock.Anything, d.id, d.now).
					Return(d.expectCreateUserData, d.expectCreateUserError)
			}

			if d.expectError == nil {
				credentialsService.On("StorageToModel", mock.Anything).Return(&d.expectModel.Credentials)
				identityService.On("StorageToModel", mock.Anything).Return(&d.expectModel.Identity)
				profileService.On("StorageToModel", mock.Anything).Return(&d.expectModel.Profile)
			}

			model, registerModel, err := service.Create(context.TODO(), d.data, d.id, d.now)
			test_utils.RequireError(st, d.expectError, err)
			require.Equal(st, d.expectModel, model)
			require.Equal(st, d.expectRegisterModel, registerModel)

			require.True(st, repository.AssertExpectations(st))
			require.True(st, credentialsService.AssertExpectations(st))
			require.True(st, identityService.AssertExpectations(st))
			require.True(st, profileService.AssertExpectations(st))
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	data := []struct {
		name string
		id   uuid.UUID
		now  time.Time

		deleteUserData  *user_storage.Model
		deleteUserError error

		expect    *models.User
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1),
			now:  baseTime,
			deleteUserData: &user_storage.Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				DeletedAt: &baseTime,
				Core: user_storage.Core{
					Credentials: credentials_storage.Core{},
					Identity: identity_storage.Core{
						Birthday: time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
						Sex:      models.SexFemale,
					},
				},
			},
			expect: &models.User{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				DeletedAt: &baseTime,
				Credentials: models.UserCredentials{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
				},
				Identity: models.UserIdentity{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					Birthday:  time.Date(2000, time.June, 15, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
				Profile: models.UserProfile{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
				},
			},
		},
		{
			name:            "Error/RepositoryFailure",
			id:              test_utils.NumberUUID(1),
			now:             baseTime,
			deleteUserError: fooErr,
			expectErr:       fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := user_storage.NewMockRepository(t)
			credentialsService := credentials_service.NewMockService(t)
			identityService := identity_service.NewMockService(t)
			profileService := profile_service.NewMockService(t)

			service := NewService(repository, credentialsService, identityService, profileService)

			repository.
				On("Delete", context.TODO(), d.id, d.now).
				Return(d.deleteUserData, d.deleteUserError)

			if d.expectErr == nil {
				credentialsService.On("StorageToModel", mock.Anything).Return(&d.expect.Credentials)
				identityService.On("StorageToModel", mock.Anything).Return(&d.expect.Identity)
				profileService.On("StorageToModel", mock.Anything).Return(&d.expect.Profile)
			}

			model, err := service.Delete(context.TODO(), d.id, d.now)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, model)

			require.True(st, repository.AssertExpectations(st))
			require.True(st, credentialsService.AssertExpectations(st))
			require.True(st, identityService.AssertExpectations(st))
			require.True(st, profileService.AssertExpectations(st))
		})
	}
}

func TestUserService_Search(t *testing.T) {
	data := []struct {
		name string

		query  string
		limit  int
		offset int

		repositoryErr   error
		repositoryData  []*user_storage.PublicPreview
		repositoryCount int64

		expect      []*models.UserPublicPreview
		expectCount int64
		expectError error
	}{
		{
			name:   "Success",
			query:  "foo",
			limit:  10,
			offset: 5,
			repositoryData: []*user_storage.PublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Username:  "BigBrother",
					Slug:      "blue-x",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime,
				},
				{
					Slug:      "notyetwritten",
					FirstName: "Eleonore",
					LastName:  "Payet",
					CreatedAt: baseTime,
				},
			},
			repositoryCount: 200,
			expect: []*models.UserPublicPreview{
				{
					Slug:      "spaceorigin",
					FirstName: "Elon",
					LastName:  "Bezos",
					CreatedAt: baseTime.Add(time.Hour),
				},
				{
					Slug:      "blue-x",
					Username:  "BigBrother",
					CreatedAt: baseTime,
				},
				{
					Slug:      "notyetwritten",
					FirstName: "Eleonore",
					LastName:  "Payet",
					CreatedAt: baseTime,
				},
			},
			expectCount: 200,
		},
		{
			name:          "Error/RepositoryFailure",
			query:         "foo",
			limit:         10,
			offset:        5,
			repositoryErr: fooErr,
			expectError:   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := user_storage.NewMockRepository(t)

			service := NewService(repository, nil, nil, nil)

			repository.
				On("Search", mock.Anything, d.query, d.limit, d.offset).
				Return(d.repositoryData, d.repositoryCount, d.repositoryErr)

			results, count, err := service.Search(context.TODO(), d.query, d.limit, d.offset)
			test_utils.RequireError(st, d.expectError, err)
			require.Equal(st, d.expect, results)
			require.Equal(st, d.expectCount, count)

			require.True(st, repository.AssertExpectations(st))
		})
	}
}

func TestUserService_GetPreview(t *testing.T) {
	data := []struct {
		name string

		id uuid.UUID

		getData  *user_storage.Preview
		getError error

		expect    *models.UserPreview
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1),
			getData: &user_storage.Preview{
				ID:        test_utils.NumberUUID(1),
				FirstName: "Elon",
				LastName:  "Bezos",
				Username:  "BigBrother",
				Email: models.Email{
					User:   "elon.bezos",
					Domain: "spaceorigin.x",
				},
				Sex: models.SexMale,
			},
			expect: &models.UserPreview{
				ID:        test_utils.NumberUUID(1),
				FirstName: "Elon",
				LastName:  "Bezos",
				Username:  "BigBrother",
				Email:     "elon.bezos@spaceorigin.x",
				Sex:       models.SexMale,
			},
		},
		{
			name: "Success/NoUsername",
			id:   test_utils.NumberUUID(1),
			getData: &user_storage.Preview{
				ID:        test_utils.NumberUUID(1),
				FirstName: "Elon",
				LastName:  "Bezos",
				Email: models.Email{
					User:   "elon.bezos",
					Domain: "spaceorigin.x",
				},
				Sex: models.SexMale,
			},
			expect: &models.UserPreview{
				ID:        test_utils.NumberUUID(1),
				FirstName: "Elon",
				LastName:  "Bezos",
				Email:     "elon.bezos@spaceorigin.x",
				Sex:       models.SexMale,
			},
		},
		{
			name:      "Error/RepositoryFailure",
			id:        test_utils.NumberUUID(1),
			getError:  fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := user_storage.NewMockRepository(st)
			service := NewService(repository, nil, nil, nil)

			repository.
				On("GetPreview", context.TODO(), d.id).
				Return(d.getData, d.getError)

			result, err := service.GetPreview(context.TODO(), d.id)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, result)

			require.True(st, repository.AssertExpectations(st))
		})
	}
}

func TestUserService_GetPublic(t *testing.T) {
	data := []struct {
		name string

		slug string

		getData  *user_storage.Public
		getError error

		expect    *models.UserPublic
		expectErr error
	}{
		{
			name: "Success",
			slug: "spaceorigin",
			getData: &user_storage.Public{
				ID:        test_utils.NumberUUID(1),
				FirstName: "Elon",
				LastName:  "Bezos",
				Username:  "BigBrother",
				CreatedAt: baseTime,
				Sex:       models.SexMale,
			},
			expect: &models.UserPublic{
				ID:        test_utils.NumberUUID(1),
				Username:  "BigBrother",
				CreatedAt: baseTime,
				Sex:       models.SexMale,
			},
		},
		{
			name: "Success/NoUsername",
			slug: "spaceorigin",
			getData: &user_storage.Public{
				ID:        test_utils.NumberUUID(1),
				FirstName: "Elon",
				LastName:  "Bezos",
				CreatedAt: baseTime,
				Sex:       models.SexMale,
			},
			expect: &models.UserPublic{
				ID:        test_utils.NumberUUID(1),
				FirstName: "Elon",
				LastName:  "Bezos",
				CreatedAt: baseTime,
				Sex:       models.SexMale,
			},
		},
		{
			name:      "Error/RepositoryFailure",
			slug:      "spaceorigin",
			getError:  fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := user_storage.NewMockRepository(st)
			service := NewService(repository, nil, nil, nil)

			repository.
				On("GetPublic", context.TODO(), d.slug).
				Return(d.getData, d.getError)

			result, err := service.GetPublic(context.TODO(), d.slug)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, result)

			require.True(st, repository.AssertExpectations(st))
		})
	}
}

func TestUserService_GetAuthorizations(t *testing.T) {
	data := []struct {
		name string

		id uuid.UUID

		getData  *models.UserCredentials
		getError error

		expect    []string
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1),
			getData: &models.UserCredentials{
				Validated: true,
			},
			expect: []string{"account-validated"},
		},
		{
			name:      "Error/CredentialsServiceFailure",
			id:        test_utils.NumberUUID(1),
			getError:  fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			credentialsService := credentials_service.NewMockService(st)
			service := NewService(nil, credentialsService, nil, nil)

			credentialsService.
				On("Read", context.TODO(), d.id).
				Return(d.getData, d.getError)

			result, err := service.GetAuthorizations(context.TODO(), d.id)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, result)

			require.True(st, credentialsService.AssertExpectations(st))
		})
	}
}

func TestUserService_GetPublicPreviews(t *testing.T) {
	data := []struct {
		name string

		ids []uuid.UUID

		getData  []*user_storage.PublicPreview
		getError error

		expect    []*models.UserPublicPreview
		expectErr error
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			getData: []*user_storage.PublicPreview{
				{
					ID:        test_utils.NumberUUID(1),
					Slug:      "blue-x",
					FirstName: "Elon",
					LastName:  "Bezos",
					Username:  "BigBrother",
					CreatedAt: baseTime,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Slug:      "notyetwritten",
					FirstName: "Eleonore",
					LastName:  "Payet",
					CreatedAt: baseTime,
				},
			},
			expect: []*models.UserPublicPreview{
				{
					ID:        test_utils.NumberUUID(1),
					Slug:      "blue-x",
					Username:  "BigBrother",
					CreatedAt: baseTime,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Slug:      "notyetwritten",
					FirstName: "Eleonore",
					LastName:  "Payet",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name: "Error/RepositoryFailure",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			getError:  fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := user_storage.NewMockRepository(st)
			service := NewService(repository, nil, nil, nil)

			repository.
				On("GetPublicPreviews", context.TODO(), d.ids).
				Return(d.getData, d.getError)

			result, err := service.GetPublicPreviews(context.TODO(), d.ids)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, result)

			require.True(st, repository.AssertExpectations(st))
		})
	}
}
