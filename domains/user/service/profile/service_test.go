package profile_service

import (
	"context"
	"errors"
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

var (
	elonBezosStorage = &profile_storage.Model{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		Core: profile_storage.Core{
			Username: "SpaceOrigin",
			Slug:     "nyan-cat-girl",
		},
	}
	elonBezosModel = &models.UserProfile{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		Username:  "SpaceOrigin",
		Slug:      "nyan-cat-girl",
	}
)

func TestProfileService_PrepareRegistration(t *testing.T) {
	repository := profile_storage.NewMockRepository(t)

	data := []struct {
		name string

		data *models.UserProfileUpdateForm

		expect    *models.UserProfileRegistrationForm
		expectErr error
	}{
		{
			name: "Success",
			data: &models.UserProfileUpdateForm{
				Slug:     "anabanana",
				Username: "BananaSplit",
			},
			expect: &models.UserProfileRegistrationForm{
				Slug:     "anabanana",
				Username: "BananaSplit",
			},
		},
		{
			name: "Success/NoUsername",
			data: &models.UserProfileUpdateForm{
				Slug: "anabanana",
			},
			expect: &models.UserProfileRegistrationForm{
				Slug: "anabanana",
			},
		},
		{
			name: "Success/ComposedSlug",
			data: &models.UserProfileUpdateForm{
				Slug:     "ana-banana-tapioca",
				Username: "BananaSplit",
			},
			expect: &models.UserProfileRegistrationForm{
				Slug:     "ana-banana-tapioca",
				Username: "BananaSplit",
			},
		},
		{
			name: "Success/WTFUsername",
			data: &models.UserProfileUpdateForm{
				Slug:     "anabanana",
				Username: "A|na\\B@nan ana-- Split",
			},
			expect: &models.UserProfileRegistrationForm{
				Slug:     "anabanana",
				Username: "A|na\\B@nan ana-- Split",
			},
		},
		{
			name: "Error/NoSlug",
			data: &models.UserProfileUpdateForm{
				Username: "BananaSplit",
			},
			expectErr: validation.ErrNil,
		},
		{
			name: "Error/MalformedSlug#1",
			data: &models.UserProfileUpdateForm{
				Slug:     "ana banana",
				Username: "BananaSplit",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedSlug#2",
			data: &models.UserProfileUpdateForm{
				Slug:     "-anabanana",
				Username: "BananaSplit",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedSlug#3",
			data: &models.UserProfileUpdateForm{
				Slug:     "ana--banana",
				Username: "BananaSplit",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedSlug#4",
			data: &models.UserProfileUpdateForm{
				Slug:     "anabanana-",
				Username: "BananaSplit",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedSlug#5",
			data: &models.UserProfileUpdateForm{
				Slug:     "anaBanana",
				Username: "BananaSplit",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedSlug#6",
			data: &models.UserProfileUpdateForm{
				Slug:     "an@banana",
				Username: "BananaSplit",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedUsername#1",
			data: &models.UserProfileUpdateForm{
				Slug:     "anabanana",
				Username: "BananaS  plit",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedUsername#2",
			data: &models.UserProfileUpdateForm{
				Slug:     "anabanana",
				Username: "BananaSplit ",
			},
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/MalformedUsername#3",
			data: &models.UserProfileUpdateForm{
				Slug:     "anabanana",
				Username: " BananaSplit",
			},
			expectErr: validation.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			service := NewService(repository)

			res, err := service.PrepareRegistration(context.TODO(), d.data)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestProfileService_Read(t *testing.T) {
	data := []struct {
		name string

		id           uuid.UUID
		getUserData  *profile_storage.Model
		getUserError error

		expect    *models.UserProfile
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
			repository := profile_storage.NewMockRepository(t)
			repository.
				On("Read", context.TODO(), d.id).
				Return(d.getUserData, d.getUserError)

			service := NewService(repository)

			res, err := service.Read(context.TODO(), d.id)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestProfileService_ReadSlug(t *testing.T) {
	data := []struct {
		name string

		slug         string
		getUserData  *profile_storage.Model
		getUserError error

		expect    *models.UserProfile
		expectErr error
	}{
		{
			name:        "Success",
			slug:        "nyan-cat-girl",
			getUserData: elonBezosStorage,
			expect:      elonBezosModel,
		},
		{
			name:         "Error/RepositoryFailure",
			slug:         "zqkur",
			getUserError: validation.ErrNotFound,
			expectErr:    validation.ErrNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := profile_storage.NewMockRepository(t)
			repository.
				On("ReadSlug", context.TODO(), d.slug).
				Return(d.getUserData, d.getUserError)

			service := NewService(repository)

			res, err := service.ReadSlug(context.TODO(), d.slug)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestProfileService_SlugExists(t *testing.T) {
	data := []struct {
		name string

		slug         string
		getUserData  bool
		getUserError error

		expect    bool
		expectErr error
	}{
		{
			name:        "Success",
			slug:        "nyan-cat-girl",
			getUserData: true,
			expect:      true,
		},
		{
			name:        "Success/DoesntExist",
			slug:        "zqkur",
			getUserData: false,
			expect:      false,
		},
		{
			name:         "Error/RepositoryFailure",
			slug:         "zqkur",
			getUserError: fooErr,
			expectErr:    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := profile_storage.NewMockRepository(t)
			repository.
				On("SlugExists", context.TODO(), d.slug).
				Return(d.getUserData, d.getUserError)

			service := NewService(repository)

			res, err := service.SlugExists(context.TODO(), d.slug)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestProfileService_Update(t *testing.T) {
	data := []struct {
		name string

		data                *models.UserProfileUpdateForm
		id                  uuid.UUID
		now                 time.Time
		updateUserData      *profile_storage.Model
		updateUserDataError error

		shouldCallRepository bool

		expect    *models.UserProfile
		expectErr error
	}{
		{
			name: "Success",
			data: &models.UserProfileUpdateForm{
				Username: "AnnaBanana",
				Slug:     "anna-banana",
			},
			id:  test_utils.NumberUUID(1000),
			now: updateTime,
			updateUserData: &profile_storage.Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: profile_storage.Core{
					Username: "AnnaBanana",
					Slug:     "anna-banana",
				},
			},
			shouldCallRepository: true,
			expect: &models.UserProfile{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Username:  "AnnaBanana",
				Slug:      "anna-banana",
			},
		},
		{
			name: "Error/InvalidData",
			data: &models.UserProfileUpdateForm{
				Username: "AnnaBanana",
				Slug:     "anna-banana-",
			},
			id:        test_utils.NumberUUID(1000),
			now:       updateTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/RepositoryFailure",
			data: &models.UserProfileUpdateForm{
				Username: "AnnaBanana",
				Slug:     "anna-banana",
			},
			id:                   test_utils.NumberUUID(1000),
			now:                  updateTime,
			updateUserDataError:  fooErr,
			shouldCallRepository: true,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := profile_storage.NewMockRepository(t)

			if d.shouldCallRepository {
				repository.
					On("Update", context.TODO(), mock.Anything, d.id, d.now).
					Return(d.updateUserData, d.updateUserDataError)
			}

			service := NewService(repository)

			res, err := service.Update(context.TODO(), d.data, d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestIdentityService_StorageToModel(t *testing.T) {
	repository := profile_storage.NewMockRepository(t)

	data := []struct {
		name   string
		data   *profile_storage.Model
		expect *models.UserProfile
	}{
		{
			name: "Standard",
			data: &profile_storage.Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: profile_storage.Core{
					Username: "AnnaBanana",
					Slug:     "anna-banana",
				},
			},
			expect: &models.UserProfile{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Username:  "AnnaBanana",
				Slug:      "anna-banana",
			},
		},
		{
			name: "Nil",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			service := NewService(repository)

			res := service.StorageToModel(d.data)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}
