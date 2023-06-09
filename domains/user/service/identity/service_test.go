package identity_service

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
	elonBezosStorage = &identity_storage.Model{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		Core: identity_storage.Core{
			FirstName: "Elon",
			LastName:  "Bezos",
			Birthday:  time.Date(1970, time.June, 28, 18, 0, 0, 0, time.UTC),
			Sex:       models.SexMale,
		},
	}
	elonBezosModel = &models.UserIdentity{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		FirstName: "Elon",
		LastName:  "Bezos",
		Birthday:  time.Date(1970, time.June, 28, 18, 0, 0, 0, time.UTC),
		Sex:       models.SexMale,
	}
)

func TestIdentityService_PrepareRegistration(t *testing.T) {
	repository := identity_storage.NewMockRepository(t)

	data := []struct {
		name string

		data *models.UserIdentityUpdateForm
		now  time.Time

		expect    *models.UserIdentityRegistrationForm
		expectErr error
	}{
		{
			name: "Success",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now: baseTime,
			expect: &models.UserIdentityRegistrationForm{
				FirstName: "Camille",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		{
			name: "Success/ComposedName",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Louise-Françoise Leblanc",
				LastName:  "De La-Vallière",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now: baseTime,
			expect: &models.UserIdentityRegistrationForm{
				FirstName: "Louise-Françoise Leblanc",
				LastName:  "De La-Vallière",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		{
			name: "Success/NonRomanizedName",
			data: &models.UserIdentityUpdateForm{
				FirstName: "ルイズ フランソワーズ ル ブラン",
				LastName:  "ド ラ ヴァリエール",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now: baseTime,
			expect: &models.UserIdentityRegistrationForm{
				FirstName: "ルイズ フランソワーズ ル ブラン",
				LastName:  "ド ラ ヴァリエール",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		{
			name: "Error/MissingFirstName",
			data: &models.UserIdentityUpdateForm{
				LastName: "Dupont",
				Birthday: time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:      models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name: "Error/MissingLastName",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name: "Error/MissingSex",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.Sex(""),
			},
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name: "Error/InvalidFirstName",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille!",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidLastName",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille",
				LastName:  "Dupont!",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidSex",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.Sex("cat"),
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/TooYoung",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille",
				LastName:  "Dupont",
				Birthday:  time.Date(2010, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/TooOld",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille",
				LastName:  "Dupont",
				Birthday:  time.Date(1800, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidFirstName#2",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille ",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidFirstName#3",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Camille-",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidFirstName#4",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Louise--Françoise",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidFirstName#5",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Louise  Françoise",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidFirstName#6",
			data: &models.UserIdentityUpdateForm{
				FirstName: " Camille",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidFirstName#7",
			data: &models.UserIdentityUpdateForm{
				FirstName: "-Camille",
				LastName:  "Dupont",
				Birthday:  time.Date(1990, time.June, 15, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			service := NewService(repository)

			res, err := service.PrepareRegistration(context.TODO(), d.data, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestIdentityService_Read(t *testing.T) {
	data := []struct {
		name string

		id           uuid.UUID
		getUserData  *identity_storage.Model
		getUserError error

		expect    *models.UserIdentity
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
			repository := identity_storage.NewMockRepository(t)
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

func TestIdentityService_Update(t *testing.T) {
	data := []struct {
		name string

		data                *models.UserIdentityUpdateForm
		id                  uuid.UUID
		now                 time.Time
		updateUserData      *identity_storage.Model
		updateUserDataError error

		shouldCallRepository bool

		expect    *models.UserIdentity
		expectErr error
	}{
		{
			name: "Success",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			id:  test_utils.NumberUUID(1000),
			now: updateTime,
			updateUserData: &identity_storage.Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: identity_storage.Core{
					FirstName: "Anna",
					LastName:  "Banana",
					Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
			},
			shouldCallRepository: true,
			expect: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		{
			name: "Error/InvalidData",
			data: &models.UserIdentityUpdateForm{
				FirstName: " Anna",
				LastName:  "Banana",
				Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			id:        test_utils.NumberUUID(1000),
			now:       updateTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/RepositoryFailure",
			data: &models.UserIdentityUpdateForm{
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
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
			repository := identity_storage.NewMockRepository(t)

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
	repository := identity_storage.NewMockRepository(t)

	data := []struct {
		name   string
		data   *identity_storage.Model
		expect *models.UserIdentity
	}{
		{
			name: "Standard",
			data: &identity_storage.Model{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				Core: identity_storage.Core{
					FirstName: "Anna",
					LastName:  "Banana",
					Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
			},
			expect: &models.UserIdentity{
				ID:        test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				FirstName: "Anna",
				LastName:  "Banana",
				Birthday:  time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
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

func TestIdentityService_Age(t *testing.T) {
	repository := identity_storage.NewMockRepository(t)

	data := []struct {
		name string

		data *models.UserIdentity
		now  time.Time

		expect int
	}{
		{
			name: "BeforeBirthday",
			now:  baseTime,
			data: &models.UserIdentity{
				Birthday: time.Date(2000, time.June, 28, 12, 0, 0, 0, time.UTC),
			},
			expect: 19,
		},
		{
			name: "AfterBirthday",
			now:  baseTime,
			data: &models.UserIdentity{
				Birthday: time.Date(2000, time.February, 28, 12, 0, 0, 0, time.UTC),
			},
			expect: 20,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			service := NewService(repository)

			res := service.Age(d.data, d.now)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}
