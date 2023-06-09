package improve_suggestion_service

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
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 9, 0, 0, 0, time.UTC)
	fooErr     = errors.New("it broken")
)

func TestImproveSuggestionService_Read(t *testing.T) {
	data := []struct {
		name string

		id       uuid.UUID
		getData  *improve_suggestion_storage.Model
		getError error

		expect    *models.ImproveSuggestion
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1000),
			getData: &improve_suggestion_storage.Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: false,
				UpVotes:   17,
				DownVotes: 3,
				Core: improve_suggestion_storage.Core{
					RequestID: test_utils.NumberUUID(11),
					Title:     "Dummy post",
					Content:   "Foo bar qux.",
				},
			},
			expect: &models.ImproveSuggestion{
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
			repository := improve_suggestion_storage.NewMockRepository(t)
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

func TestImproveSuggestionService_Create(t *testing.T) {
	data := []struct {
		name string

		userID   uuid.UUID
		sourceID uuid.UUID
		data     *models.ImproveSuggestionUpsert
		id       uuid.UUID
		now      time.Time

		shouldCallRepository     bool
		shouldCallRepositoryWith *improve_suggestion_storage.Core
		createData               *improve_suggestion_storage.Model
		createError              error

		expect    *models.ImproveSuggestion
		expectErr error
	}{
		{
			name:     "Success",
			userID:   test_utils.NumberUUID(100),
			sourceID: test_utils.NumberUUID(10),
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			id:                   test_utils.NumberUUID(1),
			now:                  baseTime,
			shouldCallRepository: true,
			shouldCallRepositoryWith: &improve_suggestion_storage.Core{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			createData: &improve_suggestion_storage.Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: false,
				UpVotes:   17,
				DownVotes: 3,
				Core: improve_suggestion_storage.Core{
					RequestID: test_utils.NumberUUID(11),
					Title:     "Dummy post",
					Content:   "Foo bar qux.",
				},
			},
			expect: &models.ImproveSuggestion{
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
		},
		{
			name:     "Error/NoTitle",
			userID:   test_utils.NumberUUID(100),
			sourceID: test_utils.NumberUUID(10),
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Content:   "Foo bar qux.",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name:     "Error/NoContent",
			userID:   test_utils.NumberUUID(100),
			sourceID: test_utils.NumberUUID(10),
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name:     "Error/TitleTooShort",
			userID:   test_utils.NumberUUID(100),
			sourceID: test_utils.NumberUUID(10),
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "D",
				Content:   "Foo bar qux.",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:     "Error/ContentTooShort",
			userID:   test_utils.NumberUUID(100),
			sourceID: test_utils.NumberUUID(10),
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "F",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:     "Error/TitleInvalid",
			userID:   test_utils.NumberUUID(100),
			sourceID: test_utils.NumberUUID(10),
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy\n post",
				Content:   "Foo bar qux.",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name:     "Error/RepositoryFailure",
			userID:   test_utils.NumberUUID(100),
			sourceID: test_utils.NumberUUID(10),
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			id:                   test_utils.NumberUUID(1),
			now:                  baseTime,
			shouldCallRepository: true,
			shouldCallRepositoryWith: &improve_suggestion_storage.Core{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			createError: fooErr,
			expectErr:   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_suggestion_storage.NewMockRepository(t)

			if d.shouldCallRepository {
				repository.
					On("Create", context.TODO(), d.shouldCallRepositoryWith, d.userID, d.sourceID, d.id, d.now).
					Return(d.createData, d.createError)
			}

			service := NewService(repository)

			res, err := service.Create(context.TODO(), d.data, d.userID, d.sourceID, d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveSuggestionService_Update(t *testing.T) {
	data := []struct {
		name string

		data *models.ImproveSuggestionUpsert
		id   uuid.UUID
		now  time.Time

		shouldCallRepository     bool
		shouldCallRepositoryWith *improve_suggestion_storage.Core
		createData               *improve_suggestion_storage.Model
		createError              error

		expect    *models.ImproveSuggestion
		expectErr error
	}{
		{
			name: "Success",
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			id:                   test_utils.NumberUUID(1),
			now:                  baseTime,
			shouldCallRepository: true,
			shouldCallRepositoryWith: &improve_suggestion_storage.Core{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			createData: &improve_suggestion_storage.Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: false,
				UpVotes:   17,
				DownVotes: 3,
				Core: improve_suggestion_storage.Core{
					RequestID: test_utils.NumberUUID(11),
					Title:     "Dummy post",
					Content:   "Foo bar qux.",
				},
			},
			expect: &models.ImproveSuggestion{
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
		},
		{
			name: "Error/NoTitle",
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Content:   "Foo bar qux.",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name: "Error/NoContent",
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrNil,
		},
		{
			name: "Error/TitleTooShort",
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "D",
				Content:   "Foo bar qux.",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/ContentTooShort",
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "F",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/TitleInvalid",
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy\n post",
				Content:   "Foo bar qux.",
			},
			id:        test_utils.NumberUUID(1),
			now:       baseTime,
			expectErr: validation.ErrInvalidEntity,
		},
		{
			name: "Error/RepositoryFailure",
			data: &models.ImproveSuggestionUpsert{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			id:                   test_utils.NumberUUID(1),
			now:                  baseTime,
			shouldCallRepository: true,
			shouldCallRepositoryWith: &improve_suggestion_storage.Core{
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
			createError: fooErr,
			expectErr:   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_suggestion_storage.NewMockRepository(t)

			if d.shouldCallRepository {
				repository.
					On("Update", context.TODO(), d.shouldCallRepositoryWith, d.id, d.now).
					Return(d.createData, d.createError)
			}

			service := NewService(repository)

			res, err := service.Update(context.TODO(), d.data, d.id, d.now)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveSuggestionService_Delete(t *testing.T) {
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
			repository := improve_suggestion_storage.NewMockRepository(t)

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

func TestImproveSuggestionService_Validate(t *testing.T) {
	data := []struct {
		name string

		validated bool
		id        uuid.UUID

		validateData  *improve_suggestion_storage.Model
		validateError error

		expect    *models.ImproveSuggestion
		expectErr error
	}{
		{
			name:      "Success",
			id:        test_utils.NumberUUID(1),
			validated: true,
			validateData: &improve_suggestion_storage.Model{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   17,
				DownVotes: 3,
				Core: improve_suggestion_storage.Core{
					RequestID: test_utils.NumberUUID(11),
					Title:     "Dummy post",
					Content:   "Foo bar qux.",
				},
			},
			expect: &models.ImproveSuggestion{
				ID:        test_utils.NumberUUID(1),
				CreatedAt: baseTime,
				UpdatedAt: &updateTime,
				SourceID:  test_utils.NumberUUID(10),
				UserID:    test_utils.NumberUUID(100),
				Validated: true,
				UpVotes:   17,
				DownVotes: 3,
				RequestID: test_utils.NumberUUID(11),
				Title:     "Dummy post",
				Content:   "Foo bar qux.",
			},
		},
		{
			name:          "Error/RepositoryFailure",
			id:            test_utils.NumberUUID(1),
			validateError: fooErr,
			expectErr:     fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_suggestion_storage.NewMockRepository(t)

			repository.
				On("Validate", context.TODO(), d.validated, d.id).
				Return(d.validateData, d.validateError)

			service := NewService(repository)

			res, err := service.Validate(context.TODO(), d.validated, d.id)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveSuggestionService_List(t *testing.T) {
	data := []struct {
		name string

		query  models.ImproveSuggestionsList
		limit  int
		offset int

		shouldCallRepositoryWith improve_suggestion_storage.ListQuery
		listData                 []*improve_suggestion_storage.Model
		listCount                int64
		listError                error

		expect      []*models.ImproveSuggestion
		expectCount int64
		expectErr   error
	}{
		{
			name: "Success",
			query: models.ImproveSuggestionsList{
				UserID:    framework.ToPTR(test_utils.NumberUUID(100)),
				SourceID:  framework.ToPTR(test_utils.NumberUUID(10)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(11)),
				Validated: framework.ToPTR(true),
			},
			limit:  10,
			offset: 20,
			shouldCallRepositoryWith: improve_suggestion_storage.ListQuery{
				UserID:    framework.ToPTR(test_utils.NumberUUID(100)),
				SourceID:  framework.ToPTR(test_utils.NumberUUID(10)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(11)),
				Validated: framework.ToPTR(true),
			},
			listData: []*improve_suggestion_storage.Model{
				{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(10),
					UserID:    test_utils.NumberUUID(100),
					Validated: false,
					UpVotes:   17,
					DownVotes: 3,
					Core: improve_suggestion_storage.Core{
						RequestID: test_utils.NumberUUID(11),
						Title:     "Dummy post",
						Content:   "Foo bar qux.",
					},
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
					Core: improve_suggestion_storage.Core{
						RequestID: test_utils.NumberUUID(14),
						Title:     "Smart post",
						Content:   "Cats on a nap.",
					},
				},
			},
			listCount: 123,
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
			expectCount: 123,
		},
		{
			name: "Error/RepositoryFailure",
			query: models.ImproveSuggestionsList{
				UserID:    framework.ToPTR(test_utils.NumberUUID(100)),
				SourceID:  framework.ToPTR(test_utils.NumberUUID(10)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(11)),
				Validated: framework.ToPTR(true),
			},
			limit:  10,
			offset: 20,
			shouldCallRepositoryWith: improve_suggestion_storage.ListQuery{
				UserID:    framework.ToPTR(test_utils.NumberUUID(100)),
				SourceID:  framework.ToPTR(test_utils.NumberUUID(10)),
				RequestID: framework.ToPTR(test_utils.NumberUUID(11)),
				Validated: framework.ToPTR(true),
			},
			listError: fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_suggestion_storage.NewMockRepository(t)

			repository.
				On("List", context.TODO(), d.shouldCallRepositoryWith, d.limit, d.offset).
				Return(d.listData, d.listCount, d.listError)

			service := NewService(repository)

			res, count, err := service.List(context.TODO(), d.query, d.limit, d.offset)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)
			require.Equal(t, d.expectCount, count)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveSuggestionService_IsCreator(t *testing.T) {
	data := []struct {
		name string

		id       uuid.UUID
		userID   uuid.UUID
		getData  bool
		getError error

		expect    bool
		expectErr error
	}{
		{
			name:    "Success",
			id:      test_utils.NumberUUID(1000),
			userID:  test_utils.NumberUUID(100),
			getData: true,
			expect:  true,
		},
		{
			name:      "Error/RepositoryFailure",
			id:        test_utils.NumberUUID(1000),
			userID:    test_utils.NumberUUID(100),
			getError:  fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_suggestion_storage.NewMockRepository(t)
			repository.
				On("IsCreator", context.TODO(), d.userID, d.id).
				Return(d.getData, d.getError)

			service := NewService(repository)

			res, err := service.IsCreator(context.TODO(), d.userID, d.id)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}

func TestImproveSuggestionService_GetPreviews(t *testing.T) {
	data := []struct {
		name string

		ids []uuid.UUID

		listData  []*improve_suggestion_storage.Model
		listError error

		expect    []*models.ImproveSuggestion
		expectErr error
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			listData: []*improve_suggestion_storage.Model{
				{
					ID:        test_utils.NumberUUID(1),
					CreatedAt: baseTime,
					UpdatedAt: &updateTime,
					SourceID:  test_utils.NumberUUID(10),
					UserID:    test_utils.NumberUUID(100),
					Validated: false,
					UpVotes:   17,
					DownVotes: 3,
					Core: improve_suggestion_storage.Core{
						RequestID: test_utils.NumberUUID(11),
						Title:     "Dummy post",
						Content:   "Foo bar qux.",
					},
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
					Core: improve_suggestion_storage.Core{
						RequestID: test_utils.NumberUUID(14),
						Title:     "Smart post",
						Content:   "Cats on a nap.",
					},
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
			name: "Error/RepositoryFailure",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1),
				test_utils.NumberUUID(2),
			},
			listError: fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(st *testing.T) {
			repository := improve_suggestion_storage.NewMockRepository(t)

			repository.
				On("GetPreviews", context.TODO(), d.ids).
				Return(d.listData, d.listError)

			service := NewService(repository)

			res, err := service.GetPreviews(context.TODO(), d.ids)
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expect, res)

			require.True(st, repository.AssertExpectations(t))
		})
	}
}
