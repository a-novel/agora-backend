package jwk_service

import (
	"context"
	"crypto/ed25519"
	"errors"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	fooErr   = errors.New("it broken")
)

func TestServiceCached_RefreshCacheAndListPublic(t *testing.T) {
	data := []struct {
		name string

		listData []*jwk_storage.Model
		listErr  error

		expect                []ed25519.PublicKey
		expectRefreshCacheErr error
	}{
		{
			name: "Success",
			listData: []*jwk_storage.Model{
				{
					Key:  jwk_storage.MockedKeys[0],
					Date: baseTime,
					Name: "test-1",
				},
				{
					Key:  jwk_storage.MockedKeys[2],
					Date: baseTime,
					Name: "test-3",
				},
				{
					Key:  jwk_storage.MockedKeys[1],
					Date: baseTime,
					Name: "test-2",
				},
			},
			expect: []ed25519.PublicKey{
				jwk_storage.MockedKeys[0].Public().(ed25519.PublicKey),
				jwk_storage.MockedKeys[2].Public().(ed25519.PublicKey),
				jwk_storage.MockedKeys[1].Public().(ed25519.PublicKey),
			},
		},
		{
			name:     "Success/NoKeys",
			listData: []*jwk_storage.Model(nil),
			expect:   []ed25519.PublicKey{},
		},
		{
			name:                  "Error/RepositoryFailure",
			listErr:               fooErr,
			expectRefreshCacheErr: fooErr,
			expect:                []ed25519.PublicKey{},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			repository := jwk_storage.NewMockRepository(t)

			repository.
				On("List", context.TODO()).
				Return(d.listData, d.listErr)

			service := NewServiceCached(repository)

			test_utils.RequireError(t, d.expectRefreshCacheErr, service.RefreshCache(context.TODO()))

			keys := service.ListPublic()
			require.Len(t, keys, len(d.expect))
			for i, key := range keys {
				require.True(t, key.Equal(d.expect[i]))
			}
		})
	}
}

func TestServiceCached_RefreshCacheAndGetPrivate(t *testing.T) {
	data := []struct {
		name string

		listData []*jwk_storage.Model
		listErr  error

		expect                ed25519.PrivateKey
		expectRefreshCacheErr error
	}{
		{
			name: "Success",
			listData: []*jwk_storage.Model{
				{
					Key:  jwk_storage.MockedKeys[0],
					Date: baseTime,
					Name: "test-1",
				},
				{
					Key:  jwk_storage.MockedKeys[2],
					Date: baseTime,
					Name: "test-3",
				},
				{
					Key:  jwk_storage.MockedKeys[1],
					Date: baseTime,
					Name: "test-2",
				},
			},
			expect: jwk_storage.MockedKeys[0],
		},
		{
			name:     "Success/NoKeys",
			listData: []*jwk_storage.Model(nil),
		},
		{
			name:                  "Error/RepositoryFailure",
			listErr:               fooErr,
			expectRefreshCacheErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			repository := jwk_storage.NewMockRepository(t)

			repository.
				On("List", context.TODO()).
				Return(d.listData, d.listErr)

			service := NewServiceCached(repository)

			test_utils.RequireError(t, d.expectRefreshCacheErr, service.RefreshCache(context.TODO()))

			key := service.GetPrivate()
			require.True(t, key.Equal(d.expect))
		})
	}
}

func TestService_Refresh(t *testing.T) {
	data := []struct {
		name string

		key        ed25519.PrivateKey
		id         uuid.UUID
		maxBackups int

		listData    []*jwk_storage.Model
		listErr     error
		writeErr    error
		deleteErr   error
		deleteErrOn int

		shouldCallListData  bool
		shouldCallWrite     bool
		shouldCallDeleteFor []string

		expectRefreshErr error
	}{
		{
			name:       "Success",
			maxBackups: 10,
			key:        jwk_storage.MockedKeys[4],
			id:         test_utils.NumberUUID(1),
			listData: []*jwk_storage.Model{
				{
					Key:  jwk_storage.MockedKeys[0],
					Date: baseTime,
					Name: "test-1",
				},
				{
					Key:  jwk_storage.MockedKeys[2],
					Date: baseTime,
					Name: "test-3",
				},
				{
					Key:  jwk_storage.MockedKeys[1],
					Date: baseTime,
					Name: "test-2",
				},
			},
			shouldCallListData: true,
			shouldCallWrite:    true,
		},
		{
			name:               "Success/NoKeyReturned",
			maxBackups:         10,
			key:                jwk_storage.MockedKeys[4],
			id:                 test_utils.NumberUUID(1),
			listData:           []*jwk_storage.Model(nil),
			shouldCallListData: true,
			shouldCallWrite:    true,
		},
		{
			name:       "Success/TooMuchKeys",
			maxBackups: 2,
			key:        jwk_storage.MockedKeys[4],
			id:         test_utils.NumberUUID(1),
			listData: []*jwk_storage.Model{
				{
					Key:  jwk_storage.MockedKeys[4],
					Date: baseTime,
					Name: "test-4",
				},
				{
					Key:  jwk_storage.MockedKeys[0],
					Date: baseTime,
					Name: "test-1",
				},
				{
					Key:  jwk_storage.MockedKeys[2],
					Date: baseTime,
					Name: "test-3",
				},
				{
					Key:  jwk_storage.MockedKeys[1],
					Date: baseTime,
					Name: "test-2",
				},
			},
			shouldCallListData:  true,
			shouldCallWrite:     true,
			shouldCallDeleteFor: []string{"test-2", "test-3"},
		},
		{
			name:             "Error/MissingKey",
			maxBackups:       10,
			id:               test_utils.NumberUUID(1),
			expectRefreshErr: validation.ErrInvalidEntity,
		},
		{
			name:             "Error/NotEnoughBackups",
			maxBackups:       0,
			key:              jwk_storage.MockedKeys[4],
			id:               test_utils.NumberUUID(1),
			expectRefreshErr: validation.ErrInvalidEntity,
		},
		{
			name:             "Error/TooMuchBackups",
			maxBackups:       1000,
			key:              jwk_storage.MockedKeys[4],
			id:               test_utils.NumberUUID(1),
			expectRefreshErr: validation.ErrInvalidEntity,
		},
		{
			name:               "Error/RepositoryListFailure",
			maxBackups:         10,
			key:                jwk_storage.MockedKeys[4],
			id:                 test_utils.NumberUUID(1),
			listErr:            fooErr,
			shouldCallListData: true,
			shouldCallWrite:    true,
			expectRefreshErr:   fooErr,
		},
		{
			name:             "Error/RepositoryWriteFailure",
			maxBackups:       10,
			key:              jwk_storage.MockedKeys[4],
			id:               test_utils.NumberUUID(1),
			shouldCallWrite:  true,
			writeErr:         fooErr,
			expectRefreshErr: fooErr,
		},
		{
			name:       "Error/RepositoryDeleteFailure",
			maxBackups: 1,
			key:        jwk_storage.MockedKeys[4],
			id:         test_utils.NumberUUID(1),
			listData: []*jwk_storage.Model{
				{
					Key:  jwk_storage.MockedKeys[4],
					Date: baseTime,
					Name: "test-4",
				},
				{
					Key:  jwk_storage.MockedKeys[0],
					Date: baseTime,
					Name: "test-1",
				},
				{
					Key:  jwk_storage.MockedKeys[2],
					Date: baseTime,
					Name: "test-3",
				},
				{
					Key:  jwk_storage.MockedKeys[1],
					Date: baseTime,
					Name: "test-2",
				},
			},
			shouldCallListData:  true,
			shouldCallWrite:     true,
			shouldCallDeleteFor: []string{"test-1", "test-3"},
			deleteErrOn:         1,
			deleteErr:           fooErr,
			expectRefreshErr:    fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			repository := jwk_storage.NewMockRepository(t)

			if d.shouldCallWrite {
				repository.
					On("Write", context.TODO(), d.key, d.id.String()).
					// Done this way cuz for now, we don't care about the returned model on this call.
					// To change if we do later.
					Return(new(jwk_storage.Model), d.writeErr)

			}

			if d.shouldCallListData {
				repository.
					On("List", context.TODO()).
					Return(d.listData, d.listErr)
			}

			for i, name := range d.shouldCallDeleteFor {
				var err error
				if i == d.deleteErrOn {
					err = d.deleteErr
				}

				repository.
					On("Delete", context.TODO(), name).
					Return(err)
			}

			service := NewService(repository)

			test_utils.RequireError(t, d.expectRefreshErr, service.Refresh(context.TODO(), d.key, d.id, d.maxBackups))
		})
	}
}
