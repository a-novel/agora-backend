package profile

import (
	"context"
	"errors"
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

func TestProfileProvider_Read(t *testing.T) {
	data := []struct {
		name string

		slug string

		userData *models.UserPublic
		userErr  error

		expect    *Model
		expectErr error
	}{
		{
			name: "Success",
			slug: "foobar",
			userData: &models.UserPublic{
				ID:        test_utils.NumberUUID(1),
				Username:  "qwerty",
				FirstName: "Foo",
				LastName:  "Bar",
				CreatedAt: baseTime,
				Sex:       models.SexFemale,
			},
			expect: &Model{
				ID:        test_utils.NumberUUID(1),
				Username:  "qwerty",
				FirstName: "Foo",
				LastName:  "Bar",
				CreatedAt: baseTime,
				Sex:       models.SexFemale,
			},
		},
		{
			name:      "Error/ServiceFailure",
			slug:      "foobar",
			userErr:   fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			userService := user_service.NewMockService(t)

			userService.On("GetPublic", context.TODO(), d.slug).Return(d.userData, d.userErr)

			provider := NewProvider(Config{UserService: userService})

			profile, err := provider.Read(context.TODO(), d.slug)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, profile)

			userService.AssertExpectations(t)
		})
	}
}

func TestProfileProvider_Search(t *testing.T) {
	data := []struct {
		name string

		query  string
		limit  int
		offset int

		userData  []*models.UserPublicPreview
		userCount int64
		userErr   error

		expect      []*Preview
		expectCount int64
		expectErr   error
	}{
		{
			name:   "Success",
			query:  "foo",
			limit:  10,
			offset: 20,
			userData: []*models.UserPublicPreview{
				{
					ID:        test_utils.NumberUUID(1),
					Slug:      "foobar",
					Username:  "qwerty",
					FirstName: "Foo",
					LastName:  "Bar",
					CreatedAt: baseTime,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Slug:      "foobaz",
					Username:  "azerty",
					FirstName: "Foo",
					LastName:  "Baz",
					CreatedAt: baseTime,
				},
			},
			userCount: 200,
			expect: []*Preview{
				{
					ID:        test_utils.NumberUUID(1),
					Slug:      "foobar",
					Username:  "qwerty",
					FirstName: "Foo",
					LastName:  "Bar",
					CreatedAt: baseTime,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Slug:      "foobaz",
					Username:  "azerty",
					FirstName: "Foo",
					LastName:  "Baz",
					CreatedAt: baseTime,
				},
			},
			expectCount: 200,
		},
		{
			name:      "Error/RepositoryFailure",
			query:     "foo",
			limit:     10,
			offset:    20,
			userErr:   fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			userService := user_service.NewMockService(t)

			userService.
				On("Search", context.TODO(), d.query, d.limit, d.offset).
				Return(d.userData, d.userCount, d.userErr)

			provider := NewProvider(Config{UserService: userService})

			profiles, count, err := provider.Search(context.TODO(), d.query, d.limit, d.offset)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, profiles)
			require.Equal(t, d.expectCount, count)

			userService.AssertExpectations(t)
		})
	}
}

func TestProfileProvider_Previews(t *testing.T) {
	data := []struct {
		name string

		ids []uuid.UUID

		userData []*models.UserPublicPreview
		userErr  error

		expect    []*Preview
		expectErr error
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			userData: []*models.UserPublicPreview{
				{
					ID:        test_utils.NumberUUID(1),
					Slug:      "foobar",
					Username:  "qwerty",
					FirstName: "Foo",
					LastName:  "Bar",
					CreatedAt: baseTime,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Slug:      "foobaz",
					Username:  "azerty",
					FirstName: "Foo",
					LastName:  "Baz",
					CreatedAt: baseTime,
				},
			},
			expect: []*Preview{
				{
					ID:        test_utils.NumberUUID(1),
					Slug:      "foobar",
					Username:  "qwerty",
					FirstName: "Foo",
					LastName:  "Bar",
					CreatedAt: baseTime,
				},
				{
					ID:        test_utils.NumberUUID(2),
					Slug:      "foobaz",
					Username:  "azerty",
					FirstName: "Foo",
					LastName:  "Baz",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name:      "Error/ServiceFailure",
			ids:       []uuid.UUID{test_utils.NumberUUID(1)},
			userErr:   fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			userService := user_service.NewMockService(t)

			userService.On("GetPublicPreviews", context.TODO(), d.ids).Return(d.userData, d.userErr)

			provider := NewProvider(Config{UserService: userService})

			profiles, err := provider.Previews(context.TODO(), d.ids)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, profiles)

			userService.AssertExpectations(t)
		})
	}
}
