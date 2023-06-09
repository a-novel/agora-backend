package secrets

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"github.com/a-novel/agora-backend/environment/user/authentication"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Provider interface {
	RotateJWKs(ctx context.Context, auth *authentication.BackendServiceAuth) error
	UpdateCache(ctx context.Context) error
}

type providerImpl struct {
	keysService jwk_service.Service
	keyGen      func() (ed25519.PrivateKey, error)
	now         func() time.Time
	id          func() uuid.UUID

	maxBackups     int
	updateInterval time.Duration
	lastUpdated    *time.Time
	mu             sync.RWMutex
}

func NewProvider(cfg Config) Provider {
	return &providerImpl{
		keysService: cfg.KeysService,
		keyGen:      cfg.KeyGen,
		now:         cfg.Now,
		id:          cfg.ID,

		maxBackups:     cfg.MaxBackups,
		updateInterval: cfg.UpdateInterval,
	}
}

func (provider *providerImpl) RotateJWKs(ctx context.Context, auth *authentication.BackendServiceAuth) error {
	if err := authentication.ForceBackendService(ctx, auth); err != nil {
		return err
	}

	key, err := provider.keyGen()
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	return provider.keysService.Refresh(ctx, key, provider.id(), provider.maxBackups)
}

func (provider *providerImpl) UpdateCache(ctx context.Context) error {
	now := provider.now()

	provider.mu.RLock()
	last := provider.lastUpdated
	provider.mu.RUnlock()

	if last != nil && now.Sub(*last) < provider.updateInterval {
		return nil
	}

	provider.mu.Lock()
	provider.lastUpdated = &now
	provider.mu.Unlock()

	if err := provider.keysService.RefreshCache(ctx); err != nil {
		provider.mu.Lock()
		provider.lastUpdated = last
		provider.mu.Unlock()

		return err
	}

	return nil
}
