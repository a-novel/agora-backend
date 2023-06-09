package improve_request_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
)

var Fixtures = []*Model{
	{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		Source:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(2000),
		UpVotes:   10,
		Title:     "Test",
		Content:   "Dummy content.",
	},
	{
		ID:        test_utils.NumberUUID(1001),
		CreatedAt: baseTime.Add(10 * time.Minute),
		Source:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(2001),
		UpVotes:   7,
		DownVotes: 2,
		Title:     "Test",
		Content:   "Dummy content updated.",
	},
	{
		ID:        test_utils.NumberUUID(1002),
		CreatedAt: baseTime.Add(time.Minute),
		Source:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(2000),
		UpVotes:   21,
		DownVotes: 8,
		Title:     "New Test",
		Content:   "Dummy content updated.",
	},
	{
		ID:        test_utils.NumberUUID(5000),
		CreatedAt: baseTime,
		Source:    test_utils.NumberUUID(5000),
		UserID:    test_utils.NumberUUID(3000),
		UpVotes:   34,
		DownVotes: 52,
		Title:     "Lorem Ipsum",
		Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a.",
	},
	{
		ID:        test_utils.NumberUUID(6000),
		CreatedAt: baseTime.Add(30 * time.Minute),
		Source:    test_utils.NumberUUID(6000),
		UserID:    test_utils.NumberUUID(2000),
		UpVotes:   11,
		DownVotes: 3,
		Title:     "New title Updated.",
		Content:   "qwertyuiopasdfghjklzxcvbnm",
	},
}

func TestImproveRequestRepository_Read(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id uuid.UUID

		expect    *Model
		expectErr error
	}{
		{
			name:   "Success",
			id:     test_utils.NumberUUID(1000),
			expect: Fixtures[0],
		},
		{
			name:   "Success/Revision",
			id:     test_utils.NumberUUID(1001),
			expect: Fixtures[1],
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1010),
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.Read(ctx, d.id)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveRequestRepository_ReadRevisions(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id uuid.UUID

		expect    []*Model
		expectErr error
	}{
		{
			name: "Success",
			id:   test_utils.NumberUUID(1000),
			expect: []*Model{
				Fixtures[1],
				Fixtures[2],
				Fixtures[0],
			},
		},
		{
			name:      "Error/ProvidingRevisionID",
			id:        test_utils.NumberUUID(1001),
			expectErr: validation.ErrNotFound,
		},
		{
			name:      "Error/NotFound",
			id:        test_utils.NumberUUID(1010),
			expectErr: validation.ErrNotFound,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.ReadRevisions(ctx, d.id)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveRequestRepository_Create(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		userID  uuid.UUID
		title   string
		content string
		id      uuid.UUID
		now     time.Time

		expect    *Model
		expectErr error
	}{
		{
			name:   "Success",
			userID: test_utils.NumberUUID(1),
			title:  "FooBar Symphony",
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:  test_utils.NumberUUID(2),
			now: baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(2),
				Source:    test_utils.NumberUUID(2),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(1),
				Title:     "FooBar Symphony",
				Content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			},
		},
		{
			name:   "Error/Exists",
			userID: test_utils.NumberUUID(1),
			title:  "FooBar Symphony",
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:        test_utils.NumberUUID(1001),
			now:       baseTime,
			expectErr: validation.ErrUniqConstraintViolation,
		},
		{
			name:   "Error/NoTitle",
			userID: test_utils.NumberUUID(1),
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:        test_utils.NumberUUID(2000),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name:      "Error/NoContent",
			userID:    test_utils.NumberUUID(1),
			title:     "FooBar Symphony",
			id:        test_utils.NumberUUID(2000),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx, 10)

				res, err := repository.Create(ctx, d.userID, d.title, d.content, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveRequestRepository_CreateRevision(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		userID   uuid.UUID
		sourceID uuid.UUID
		title    string
		content  string
		id       uuid.UUID
		now      time.Time

		expect    *Model
		expectErr error
	}{
		{
			name:     "Success",
			userID:   test_utils.NumberUUID(2000),
			sourceID: test_utils.NumberUUID(1000),
			title:    "FooBar Symphony",
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:  test_utils.NumberUUID(2),
			now: baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(2),
				Source:    test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(2000),
				Title:     "FooBar Symphony",
				Content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			},
		},
		{
			name:     "Success/DifferentUser",
			userID:   test_utils.NumberUUID(3),
			sourceID: test_utils.NumberUUID(1000),
			title:    "FooBar Symphony",
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:  test_utils.NumberUUID(2),
			now: baseTime,
			expect: &Model{
				ID:        test_utils.NumberUUID(2),
				Source:    test_utils.NumberUUID(1000),
				CreatedAt: baseTime,
				UserID:    test_utils.NumberUUID(3),
				Title:     "FooBar Symphony",
				Content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			},
		},
		{
			name:     "Error/Exists",
			userID:   test_utils.NumberUUID(1),
			sourceID: test_utils.NumberUUID(1000),
			title:    "FooBar Symphony",
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:        test_utils.NumberUUID(1001),
			now:       baseTime,
			expectErr: validation.ErrUniqConstraintViolation,
		},
		{
			name:     "Error/NoTitle",
			userID:   test_utils.NumberUUID(1),
			sourceID: test_utils.NumberUUID(1000),
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:        test_utils.NumberUUID(2000),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name:      "Error/NoContent",
			userID:    test_utils.NumberUUID(1),
			sourceID:  test_utils.NumberUUID(1000),
			title:     "FooBar Symphony",
			id:        test_utils.NumberUUID(2000),
			now:       baseTime,
			expectErr: validation.ErrConstraintViolation,
		},
		{
			name:     "Error/MissingSource",
			userID:   test_utils.NumberUUID(2000),
			sourceID: test_utils.NumberUUID(100),
			title:    "FooBar Symphony",
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:        test_utils.NumberUUID(2),
			now:       baseTime,
			expectErr: validation.ErrMissingRelation,
		},
		{
			name:     "Error/SourceCannotBeARevision",
			userID:   test_utils.NumberUUID(2000),
			sourceID: test_utils.NumberUUID(1001),
			title:    "FooBar Symphony",
			content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a. Quisque venenatis hendrerit laoreet. Praesent egestas turpis imperdiet felis vulputate, a eleifend turpis luctus. Aliquam at varius metus, eu placerat orci. Pellentesque vel convallis nisl. Pellentesque porta tellus nec vulputate efficitur. Ut eleifend, quam ut ultricies vulputate, nibh felis sollicitudin ante, convallis tincidunt urna nisl sed erat.

Proin rutrum commodo tincidunt. Sed convallis risus ut justo egestas vestibulum. Nullam tincidunt sed quam a viverra. Cras eu nulla at dui varius cursus ut at turpis. Phasellus nec pellentesque nisi. Aenean est dolor, facilisis a eros eu, elementum sagittis nulla. Duis pulvinar sed augue nec fermentum. Duis eu malesuada justo, a porttitor urna. Aliquam justo mi, aliquam in sem sit amet, vulputate tristique felis. Donec molestie accumsan nunc a facilisis.`,
			id:        test_utils.NumberUUID(2),
			now:       baseTime,
			expectErr: validation.ErrMissingRelation,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx, 10)

				res, err := repository.CreateRevision(ctx, d.userID, d.sourceID, d.title, d.content, d.id, d.now)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveRequestRepository_Delete(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id uuid.UUID

		expectErr  error
		expectRows int
	}{
		{
			name:       "Success",
			id:         test_utils.NumberUUID(1000),
			expectRows: 2,
		},
		{
			name:       "Success/Revision",
			id:         test_utils.NumberUUID(1001),
			expectRows: 4,
		},
		{
			name:       "Success/NotFound",
			id:         test_utils.NumberUUID(100),
			expectRows: 5,
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.Begin()
				require.NoError(st, err)
				defer stx.Rollback()
				repository := NewRepository(stx, 10)

				err = repository.Delete(ctx, d.id)
				test_utils.RequireError(t, d.expectErr, err)

				count, err := stx.NewSelect().Model(new(Model)).Count(ctx)
				require.NoError(st, err)
				require.Equal(st, d.expectRows, count)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveRequestRepository_Search(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	user3000ID := test_utils.NumberUUID(3000)

	data := []struct {
		name string

		query  SearchQuery
		limit  int
		offset int

		expect      []*Preview
		expectCount int64
		expectErr   error
	}{
		{
			name: "Success",
			query: SearchQuery{
				Query: "beau nitescence",
			},
			expectCount: 2,
			limit:       10,
			offset:      0,
			expect: []*Preview{
				// Priority to title
				{
					ID:            test_utils.NumberUUID(2000),
					CreatedAt:     baseTime.Add(2 * time.Minute),
					UserID:        test_utils.NumberUUID(2000),
					Source:        test_utils.NumberUUID(2000),
					Title:         "Beauté nitescente dans la nuit",
					Content:       "Une mère d",
					UpVotes:       4,
					DownVotes:     1,
					RevisionCount: 1,
				},
				// Only last revision
				{
					ID:            test_utils.NumberUUID(1002),
					CreatedAt:     baseTime.Add(10 * time.Minute),
					UserID:        test_utils.NumberUUID(2000),
					Source:        test_utils.NumberUUID(1000),
					Title:         "Coup de foudre au premier regard",
					Content:       "Aussi, qua",
					UpVotes:       38,
					DownVotes:     10,
					RevisionCount: 3,
				},
			},
		},
		{
			name:        "Success/NoQuery",
			query:       SearchQuery{},
			expectCount: 4,
			limit:       2,
			offset:      1,
			// Ordered by date.
			expect: []*Preview{
				{
					ID:            test_utils.NumberUUID(4000),
					CreatedAt:     baseTime.Add(7 * time.Minute),
					UserID:        test_utils.NumberUUID(5000),
					Source:        test_utils.NumberUUID(4000),
					UpVotes:       17,
					DownVotes:     6,
					Title:         "Lois robotiques",
					Content:       "Les trois ",
					RevisionCount: 1,
				},
				{
					ID:            test_utils.NumberUUID(3001),
					CreatedAt:     baseTime.Add(4 * time.Minute),
					UserID:        test_utils.NumberUUID(2000),
					Source:        test_utils.NumberUUID(3000),
					UpVotes:       45,
					DownVotes:     55,
					Title:         "Fascination étrange",
					Content:       "Alors que ",
					RevisionCount: 2,
				},
			},
		},
		{
			name: "Success/User",
			query: SearchQuery{
				Query:  "beau nitescence",
				UserID: &user3000ID,
			},
			expectCount: 2,
			limit:       10,
			offset:      0,
			expect: []*Preview{
				// Priority to most recent
				{
					ID:                  test_utils.NumberUUID(1001),
					CreatedAt:           baseTime.Add(5 * time.Minute),
					UserID:              test_utils.NumberUUID(3000),
					Source:              test_utils.NumberUUID(1000),
					Title:               "Calixte dans la lumière",
					Content:             "Aussi, qua",
					UpVotes:             38,
					DownVotes:           10,
					RevisionCount:       3,
					MoreRecentRevisions: 1,
				},
				// Only last revision
				{
					ID:                  test_utils.NumberUUID(3000),
					CreatedAt:           baseTime.Add(3 * time.Minute),
					UserID:              test_utils.NumberUUID(3000),
					Source:              test_utils.NumberUUID(3000),
					Title:               "Fascination étrange",
					Content:             "Alors que ",
					UpVotes:             45,
					DownVotes:           55,
					RevisionCount:       2,
					MoreRecentRevisions: 1,
				},
			},
		},
	}

	fixtures := []interface{}{
		// Corpus 1.
		&Model{
			ID:        test_utils.NumberUUID(1000),
			CreatedAt: baseTime,
			UserID:    test_utils.NumberUUID(2000),
			Source:    test_utils.NumberUUID(1000),
			UpVotes:   10,
			Title:     "Calixte dans la lumière",
			Content:   "Aussi, quand il rencontra Sombreval et sa fille, fut-il frappé d'un éblouissement qui ne venait pas seulement de l'incroyablement beauté nitescente Calixte, marchant dans l'éclat solaire.",
		},
		&Model{
			ID:        test_utils.NumberUUID(1001),
			CreatedAt: baseTime.Add(5 * time.Minute),
			UserID:    test_utils.NumberUUID(3000),
			Source:    test_utils.NumberUUID(1000),
			UpVotes:   7,
			DownVotes: 2,
			Title:     "Calixte dans la lumière",
			Content:   "Aussi, quand il rencontra Sombreval et sa fille traversant le cimetière, fut-il frappé d'un éblouissement qui ne venait pas seulement de la beauté nitescente de Calixte, marchant dans l'éclat solaire d'un jour d'été.",
		},
		&Model{
			ID:        test_utils.NumberUUID(1002),
			CreatedAt: baseTime.Add(10 * time.Minute),
			UserID:    test_utils.NumberUUID(2000),
			Source:    test_utils.NumberUUID(1000),
			UpVotes:   21,
			DownVotes: 8,
			Title:     "Coup de foudre au premier regard",
			Content:   "Aussi, quand il rencontra Sombreval et sa fille traversant le cimetière, fut-il frappé d'un éblouissement qui ne venait pas seulement de la beauté nitescente de Calixte, marchant dans l'éclat solaire d'un jour d'été.",
		},
		// Corpus 2.
		&Model{
			ID:        test_utils.NumberUUID(2000),
			CreatedAt: baseTime.Add(2 * time.Minute),
			UserID:    test_utils.NumberUUID(2000),
			Source:    test_utils.NumberUUID(2000),
			UpVotes:   4,
			DownVotes: 1,
			Title:     "Beauté nitescente dans la nuit",
			Content:   "Une mère dont la vie n'a de sens que l'existence de son fils voit son monde basculer dans les méandres incertains d'une autre réalité. Là, règne une guerre, cachée aux yeux des mortels, entre des créatures mythiques dont seuls ses rêves pouvaient lui souffler l'existence.",
		},
		// Corpus 3.
		&Model{
			ID:        test_utils.NumberUUID(3000),
			CreatedAt: baseTime.Add(3 * time.Minute),
			UserID:    test_utils.NumberUUID(3000),
			Source:    test_utils.NumberUUID(3000),
			UpVotes:   34,
			DownVotes: 52,
			Title:     "Fascination étrange",
			Content:   "Alors que ses paupières s'alourdissaient, se mirent à danser devant elle les arabesques à la beauté nitescente des esprits de la nuit, envoutantes et menaçantes à la fois. Comme si la réalité perdait de son sens.",
		},
		&Model{
			ID:        test_utils.NumberUUID(3001),
			CreatedAt: baseTime.Add(4 * time.Minute),
			UserID:    test_utils.NumberUUID(2000),
			Source:    test_utils.NumberUUID(3000),
			UpVotes:   11,
			DownVotes: 3,
			Title:     "Fascination étrange",
			Content:   "Alors que ses paupières s'alourdissaient, se mirent à danser devant elle les arabesques nitescentes des esprits de la nuit, envoutantes et menaçantes à la fois. Comme si la réalité perdait de son sens.",
		},
		// Corpus 4.
		&Model{
			ID:        test_utils.NumberUUID(4000),
			CreatedAt: baseTime.Add(7 * time.Minute),
			UserID:    test_utils.NumberUUID(5000),
			Source:    test_utils.NumberUUID(4000),
			UpVotes:   17,
			DownVotes: 6,
			Title:     "Lois robotiques",
			Content:   "Les trois Lois constituent les principes directeurs essentiels d'une grande partie des systèmes moraux du monde. Evidemment, chaque être humain possède, en principe, l'instinct de conservation. C'est la Troisième Loi de la robotique. De même, chacun des bons êtres humains, possédant une conscience sociale et le sens de la responsabilité, doit obéir aux autorités établies, écouter son docteur, son patron, son gouvernement, son psychiatre, son semblable... même lorsque ceux-ci troublent son confort ou sa sécurité. C'est ce qui correspond à la Seconde Loi de la robotique. Chaque bon humain doit également aimer son prochain comme lui-même, risquer sa vie pour sauver celle d'un autre. Telle est la Première Loi de la robotique.",
		},
	}

	err := test_utils.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, count, err := repository.Search(ctx, d.query, d.limit, d.offset)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
				require.Equal(t, d.expectCount, count)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveRequestRepository_IsCreator(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		id     uuid.UUID
		userID uuid.UUID
		strict bool

		expect    bool
		expectErr error
	}{
		{
			name:   "Success",
			id:     test_utils.NumberUUID(1001),
			userID: test_utils.NumberUUID(2001),
			expect: true,
		},
		{
			name:   "Success/AnotherRevision",
			id:     test_utils.NumberUUID(1002),
			userID: test_utils.NumberUUID(2001),
			expect: true,
		},
		{
			name:   "Success/StrictModeWrongRevision",
			id:     test_utils.NumberUUID(1002),
			userID: test_utils.NumberUUID(2001),
			strict: true,
		},
		{
			name:   "Success/StrictModeCorrectRevision",
			id:     test_utils.NumberUUID(1001),
			userID: test_utils.NumberUUID(2001),
			strict: true,
			expect: true,
		},
		{
			name:   "Success/NotFound",
			id:     test_utils.NumberUUID(2001),
			userID: test_utils.NumberUUID(2001),
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.IsCreator(ctx, d.userID, d.id, d.strict)
				test_utils.RequireError(t, d.expectErr, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestImproveRequestRepository_GetPreviews(t *testing.T) {
	db, sqlDB := test_utils.GetPostgres(t)
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name   string
		ids    []uuid.UUID
		expect []*Preview
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				test_utils.NumberUUID(1002),
				test_utils.NumberUUID(6000),
				test_utils.NumberUUID(8000),
			},
			expect: []*Preview{
				{
					ID:                  test_utils.NumberUUID(1002),
					CreatedAt:           baseTime.Add(time.Minute),
					Source:              test_utils.NumberUUID(1000),
					UserID:              test_utils.NumberUUID(2000),
					UpVotes:             38,
					DownVotes:           10,
					Title:               "New Test",
					Content:             "Dummy cont",
					RevisionCount:       3,
					MoreRecentRevisions: 1,
				},
				{
					ID:            test_utils.NumberUUID(6000),
					CreatedAt:     baseTime.Add(30 * time.Minute),
					Source:        test_utils.NumberUUID(6000),
					UserID:        test_utils.NumberUUID(2000),
					UpVotes:       11,
					DownVotes:     3,
					Title:         "New title Updated.",
					Content:       "qwertyuiop",
					RevisionCount: 1,
				},
			},
		},
	}

	err := test_utils.RunTransactionalTest(db, Fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := NewRepository(tx, 10)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.GetPreviews(ctx, d.ids)
				require.NoError(st, err)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}
