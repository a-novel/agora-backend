package secrets

import (
	"context"
	"crypto/ed25519"
	"errors"
	"github.com/a-novel/agora-backend/framework"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	fooErr   = errors.New("it broken")
)

func TestSecretsProvider_RotateJWKs(t *testing.T) {
	data := []struct {
		name string

		now time.Time
		id  uuid.UUID

		maxBackups int

		shouldCallKeysService bool

		keyGenErr  error
		keyGenData ed25519.PrivateKey
		keysErr    error

		expectErr error
	}{
		{
			name:                  "Success",
			now:                   baseTime,
			id:                    test_utils.NumberUUID(1),
			maxBackups:            7,
			shouldCallKeysService: true,
			keyGenData:            jwk_storage.MockedKeys[0],
		},
		{
			name:                  "Error/KeysServiceFailure",
			now:                   baseTime,
			id:                    test_utils.NumberUUID(1),
			maxBackups:            7,
			shouldCallKeysService: true,
			keysErr:               fooErr,
			expectErr:             fooErr,
		},
		{
			name:       "Error/KeyGenFailure",
			now:        baseTime,
			id:         test_utils.NumberUUID(1),
			maxBackups: 7,
			keyGenErr:  fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			keysService := jwk_service.NewMockService(t)

			if d.shouldCallKeysService {
				keysService.
					On("Refresh", context.TODO(), d.keyGenData, d.id, d.maxBackups).
					Return(d.keysErr)
			}

			provider := NewProvider(Config{
				KeysService: keysService,
				KeyGen:      test_utils.GetJWKKeyGen(d.keyGenData, d.keyGenErr),
				Now:         test_utils.GetTimeNow(d.now),
				ID:          test_utils.GetUUID(d.id),
				MaxBackups:  d.maxBackups,
			})

			err := provider.RotateJWKs(context.TODO(), nil)
			test_utils.RequireError(t, d.expectErr, err)

			keysService.AssertExpectations(t)
		})
	}
}

func TestSecretsProvider_UpdateCache(t *testing.T) {
	data := []struct {
		name string

		now time.Time

		updateInterval time.Duration
		lastUpdated    *time.Time

		shouldCallKeysService bool

		keysServiceErr error

		expectLastUpdated *time.Time
		expectErr         error
	}{
		{
			name:                  "Success",
			now:                   baseTime,
			updateInterval:        10 * time.Minute,
			lastUpdated:           nil,
			shouldCallKeysService: true,
			expectLastUpdated:     &baseTime,
		},
		{
			name:              "Success/UpdatedRecently",
			now:               baseTime,
			updateInterval:    10 * time.Minute,
			lastUpdated:       framework.ToPTR(baseTime.Add(-5 * time.Minute)),
			expectLastUpdated: framework.ToPTR(baseTime.Add(-5 * time.Minute)),
		},
		{
			name:                  "Success/UpdatedLongAgo",
			now:                   baseTime,
			updateInterval:        10 * time.Minute,
			lastUpdated:           framework.ToPTR(baseTime.Add(-15 * time.Minute)),
			shouldCallKeysService: true,
			expectLastUpdated:     &baseTime,
		},
		{
			name:                  "Error/KeysServiceFailure",
			now:                   baseTime,
			updateInterval:        10 * time.Minute,
			lastUpdated:           framework.ToPTR(baseTime.Add(-15 * time.Minute)),
			shouldCallKeysService: true,
			keysServiceErr:        fooErr,
			expectLastUpdated:     framework.ToPTR(baseTime.Add(-15 * time.Minute)),
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			keysService := jwk_service.NewMockService(t)

			if d.shouldCallKeysService {
				keysService.
					On("RefreshCache", context.TODO()).
					Return(d.keysServiceErr)
			}

			provider := NewProvider(Config{
				KeysService:    keysService,
				Now:            test_utils.GetTimeNow(d.now),
				UpdateInterval: d.updateInterval,
			})

			tProvider := provider.(*providerImpl)
			tProvider.lastUpdated = d.lastUpdated

			err := provider.UpdateCache(context.TODO())
			test_utils.RequireError(t, d.expectErr, err)
			require.Equal(t, d.expectLastUpdated, tProvider.lastUpdated)

			keysService.AssertExpectations(t)
		})
	}
}

func TestSecretsProvider_UpdateCacheStress(t *testing.T) {
	keysRepository := jwk_storage.NewFileSystemRepository(t.TempDir(), "foo")
	keysService := jwk_service.NewService(keysRepository)

	// We use a mock wrapper to count calls. Only one call should be made, if it works correctly.
	keysServiceMocked := jwk_service.NewMockService(t)
	keysServiceMocked.On("RefreshCache", context.TODO()).Return(keysService.RefreshCache(context.TODO()))

	provider := NewProvider(Config{
		KeysService:    keysServiceMocked,
		Now:            time.Now,
		UpdateInterval: 10 * time.Minute,
	})

	wg := new(sync.WaitGroup)
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			err := provider.UpdateCache(context.TODO())
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}
