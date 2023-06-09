package improve_request_service

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

func TestImproveRequestService_Read(t *testing.T) {
	data := []struct {
		name string

		id       uuid.UUID
		getData  *improve_request_storage.Model
		getError error

		expect    *models.ImproveRequest
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1000),
			getData: &improve_request_storage.Model{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(10),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(100),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
				UpVotes:   10,
				DownVotes: 3,
			},
			expect: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(10),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(100),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
				UpVotes:   10,
				DownVotes: 3,
			},
		},
		{
			name:      "Error/RepositoryFailure",
			id:        test_utils.NumberUUID(1000),
			getError:  fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_request_storage.NewMockRepository(t)
			repository.
				On("Read", context.TODO(), d.id).
				Return(d.getData, d.getError)

			service := NewService(repository)

			res, err := service.Read(context.TODO(), d.id)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveRequestService_ReadRevisions(t *testing.T) {
	data := []struct {
		name string

		id       uuid.UUID
		getData  []*improve_request_storage.Model
		getError error

		expect    []*models.ImproveRequest
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1000),
			getData: []*improve_request_storage.Model{
				{
					ID:        test_utils.NumberUUID(1),
					Source:    test_utils.NumberUUID(10),
					CreatedAt: baseTime.Add(time.Minute),
					UserID:    test_utils.NumberUUID(100),
					Title:     "Dummy post updated",
					Content:   "Foo bar qux.",
					UpVotes:   2,
					DownVotes: 0,
				},
				{
					ID:        test_utils.NumberUUID(1),
					Source:    test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UserID:    test_utils.NumberUUID(100),
					Title:     "Dummy post",
					Content:   "Foo bar qux.",
					UpVotes:   10,
					DownVotes: 3,
				},
			},
			expect: []*models.ImproveRequest{
				{
					ID:        test_utils.NumberUUID(1),
					Source:    test_utils.NumberUUID(10),
					CreatedAt: baseTime.Add(time.Minute),
					UserID:    test_utils.NumberUUID(100),
					Title:     "Dummy post updated",
					Content:   "Foo bar qux.",
					UpVotes:   2,
					DownVotes: 0,
				},
				{
					ID:        test_utils.NumberUUID(1),
					Source:    test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UserID:    test_utils.NumberUUID(100),
					Title:     "Dummy post",
					Content:   "Foo bar qux.",
					UpVotes:   10,
					DownVotes: 3,
				},
			},
		},
		{
			name:      "Error/RepositoryFailure",
			id:        test_utils.NumberUUID(1000),
			getError:  fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_request_storage.NewMockRepository(t)
			repository.
				On("ReadRevisions", context.TODO(), d.id).
				Return(d.getData, d.getError)

			service := NewService(repository)

			res, err := service.ReadRevisions(context.TODO(), d.id)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveRequestService_Create(t *testing.T) {
	data := []struct {
		name string

		userID  uuid.UUID
		title   string
		content string
		id      uuid.UUID
		now     time.Time

		createData  *improve_request_storage.Model
		createError error

		shouldCallRepository bool

		expect    *models.ImproveRequest
		expectErr error
	}{
		{
			name:                 "Success",
			userID:               test_utils.NumberUUID(100),
			title:                "Dummy post",
			content:              "Foo bar qux.",
			id:                   test_utils.NumberUUID(1),
			now:                  baseTime,
			shouldCallRepository: true,
			createData: &improve_request_storage.Model{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(100),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			expect: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(100),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
		},
		{
			name:      "Error/NoTitle",
			userID:    test_utils.NumberUUID(100),
			content:   "Foo bar qux.",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name:      "Error/NoContent",
			userID:    test_utils.NumberUUID(100),
			title:     "Dummy post",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name:      "Error/TitleTooShort",
			userID:    test_utils.NumberUUID(100),
			title:     "D",
			content:   "Foo bar qux.",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:      "Error/ContentTooShort",
			userID:    test_utils.NumberUUID(100),
			title:     "Dummy post",
			content:   "F",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:      "Error/TitleInvalid",
			userID:    test_utils.NumberUUID(100),
			title:     "Dummy\n post.",
			content:   "Foo bar qux.",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:                 "Error/RepositoryFailure",
			userID:               test_utils.NumberUUID(100),
			title:                "Dummy post",
			content:              "Foo bar qux.",
			id:                   test_utils.NumberUUID(1),
			now:                  baseTime,
			shouldCallRepository: true,
			createError:          fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_request_storage.NewMockRepository(t)

			if d.shouldCallRepository {
				repository.
					On("Create", context.TODO(), d.userID, d.title, d.content, d.id, d.now).
					Return(d.createData, d.createError)
			}

			service := NewService(repository)

			res, err := service.Create(context.TODO(), d.userID, d.title, d.content, d.id, d.now)
			test_utils.RequireError(st, d.expectErr, err)
			require.Equal(st, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveRequestService_CreateRevision(t *testing.T) {
	data := []struct {
		name string

		userID   uuid.UUID
		sourceID uuid.UUID
		title    string
		content  string
		id       uuid.UUID
		now      time.Time

		createData  *improve_request_storage.Model
		createError error

		shouldCallRepository bool

		expect    *models.ImproveRequest
		expectErr error
	}{
		{
			name:                 "Success",
			userID:               test_utils.NumberUUID(100),
			sourceID:             test_utils.NumberUUID(2),
			title:                "Dummy post",
			content:              "Foo bar qux.",
			id:                   test_utils.NumberUUID(1),
			now:                  baseTime,
			shouldCallRepository: true,
			createData: &improve_request_storage.Model{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(100),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			expect: &models.ImproveRequest{
				ID:        test_utils.NumberUUID(1),
				Source:    test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(100),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
		},
		{
			name:      "Error/NoTitle",
			userID:    test_utils.NumberUUID(100),
			sourceID:  test_utils.NumberUUID(2),
			content:   "Foo bar qux.",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name:      "Error/NoContent",
			userID:    test_utils.NumberUUID(100),
			sourceID:  test_utils.NumberUUID(2),
			title:     "Dummy post",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name:      "Error/TitleTooShort",
			userID:    test_utils.NumberUUID(100),
			sourceID:  test_utils.NumberUUID(2),
			title:     "D",
			content:   "Foo bar qux.",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:      "Error/ContentTooShort",
			userID:    test_utils.NumberUUID(100),
			sourceID:  test_utils.NumberUUID(2),
			title:     "Dummy post",
			content:   "F",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:      "Error/TitleInvalid",
			userID:    test_utils.NumberUUID(100),
			sourceID:  test_utils.NumberUUID(2),
			title:     "Dummy\n post.",
			content:   "Foo bar qux.",
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:                 "Error/RepositoryFailure",
			userID:               test_utils.NumberUUID(100),
			sourceID:             test_utils.NumberUUID(2),
			title:                "Dummy post",
			content:              "Foo bar qux.",
			id:                   test_utils.NumberUUID(1),
			now:                  baseTime,
			shouldCallRepository: true,
			createError:          fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_request_storage.NewMockRepository(t)

			if d.shouldCallRepository {
				repository.
					On("CreateRevision", context.TODO(), d.userID, d.sourceID, d.title, d.content, d.id, d.now).
					Return(d.createData, d.createError)
			}

			service := NewService(repository)

			res, err := service.CreateRevision(context.TODO(), d.userID, d.sourceID, d.title, d.content, d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveRequestService_Delete(t *testing.T) {
	data := []struct {
		name string

		id uuid.UUID

		deleteError error
		expectErr   error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1),
		},
		{
			name:        "Error/RepositoryFailure",
			id:          test_utils.NumberUUID(1),
			deleteError: fooErr,
			expectErr:   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_request_storage.NewMockRepository(t)

			repository.
				On("Delete", context.TODO(), d.id).
				Return(d.deleteError)

			service := NewService(repository)

			err := service.Delete(context.TODO(), d.id)
			test_utils.RequireError(t, d.expectErr, err)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveRequestService_Search(t *testing.T) {
	data := []struct {
		name string

		query  models.ImproveRequestSearch
		limit  int
		offset int

		shouldCallRepositoryWithQuery improve_request_storage.SearchQuery
		searchData                    []*improve_request_storage.Preview
		searchCount                   int64
		searchError                   error

		expect      []*models.ImproveRequestPreview
		expectCount int64
		expectErr   error
	}{
		{
			name: "Success",
			query: models.ImproveRequestSearch{
				UserID: framework.ToPTR(test_utils.NumberUUID(1)),
				Query:  "foo bar",
			},
			limit:  10,
			offset: 20,
			shouldCallRepositoryWithQuery: improve_request_storage.SearchQuery{
				UserID: framework.ToPTR(test_utils.NumberUUID(1)),
				Query:  "foo bar",
			},
			searchData: []*improve_request_storage.Preview{
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
			searchCount: 123,
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
			expectCount: 123,
		},
		{
			name: "Error/RepositoryFailure",
			query: models.ImproveRequestSearch{
				UserID: framework.ToPTR(test_utils.NumberUUID(1)),
				Query:  "foo bar",
			},
			limit:  10,
			offset: 20,
			shouldCallRepositoryWithQuery: improve_request_storage.SearchQuery{
				UserID: framework.ToPTR(test_utils.NumberUUID(1)),
				Query:  "foo bar",
			},
			searchError: fooErr,
			expectErr:   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_request_storage.NewMockRepository(t)

			repository.
				On("Search", context.TODO(), d.shouldCallRepositoryWithQuery, d.limit, d.offset).
				Return(d.searchData, d.searchCount, d.searchError)

			service := NewService(repository)

			res, count, err := service.Search(context.TODO(), d.query, d.limit, d.offset)
			test_utils.RequireError(t, d.expectErr, err)
			require.EqualValues(t, d.expect, res)
			require.Equal(t, d.expectCount, count)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveRequestService_IsCreator(t *testing.T) {
	data := []struct {
		name string

		id     uuid.UUID
		userID uuid.UUID
		strict bool

		getData  bool
		getError error

		expect    bool
		expectErr error
	}{
		{
			name:    "Success",
			id:      test_utils.NumberUUID(1000),
			userID:  test_utils.NumberUUID(100),
			strict:  true,
			getData: true,
			expect:  true,
		},
		{
			name:      "Error/RepositoryFailure",
			id:        test_utils.NumberUUID(1000),
			userID:    test_utils.NumberUUID(100),
			strict:    true,
			getError:  fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_request_storage.NewMockRepository(t)
			repository.
				On("IsCreator", context.TODO(), d.userID, d.id, d.strict).
				Return(d.getData, d.getError)

			service := NewService(repository)

			res, err := service.IsCreator(context.TODO(), d.userID, d.id, d.strict)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveRequestService_GetPreviews(t *testing.T) {
	data := []struct {
		name string

		ids []uuid.UUID

		searchData  []*improve_request_storage.Preview
		searchError error

		expect    []*models.ImproveRequestPreview
		expectErr error
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			searchData: []*improve_request_storage.Preview{
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
			name: "Error/RepositoryFailure",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			searchError: fooErr,
			expectErr:   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_request_storage.NewMockRepository(t)

			repository.
				On("GetPreviews", context.TODO(), d.ids).
				Return(d.searchData, d.searchError)

			service := NewService(repository)

			res, err := service.GetPreviews(context.TODO(), d.ids)
			test_utils.RequireError(t, d.expectErr, err)
			require.EqualValues(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}
