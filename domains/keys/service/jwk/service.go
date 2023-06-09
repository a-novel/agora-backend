package jwk_service

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/google/uuid"
	"sync"
)

const (
	MinBackups = 1
	MaxBackups = 32
)

// Service of the current layer. You can instantiate a new one with NewService.
type Service interface {
	Refresh(ctx context.Context, key ed25519.PrivateKey, id uuid.UUID, maxBackups int) error
	ReadOnly() ServiceCached

	ServiceCached
}

type serviceImpl struct {
	repository jwk_storage.Repository

	serviceCachedImpl
}

// NewService returns a new implementation of Service.
//
//	jwk_service.NewService(repository)
func NewService(repository jwk_storage.Repository) Service {
	return &serviceImpl{
		repository:        repository,
		serviceCachedImpl: serviceCachedImpl{repository: repository},
	}
}

func (service *serviceImpl) Refresh(ctx context.Context, key ed25519.PrivateKey, id uuid.UUID, maxBackups int) error {
	if err := validation.CheckRequire("key", key); err != nil {
		return err
	}
	if err := validation.CheckMinMax("maxBackups", maxBackups, MinBackups, MaxBackups); err != nil {
		return err
	}

	if _, err := service.repository.Write(ctx, key, id.String()); err != nil {
		return fmt.Errorf("failed to save new key: %w", err)
	}

	keys, err := service.repository.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list existing keys: %w", err)
	}

	if len(keys) > maxBackups {
		for _, extraKey := range keys[maxBackups:] {
			if err = service.repository.Delete(ctx, extraKey.Name); err != nil {
				return fmt.Errorf("failed to delete extra keys: %w", err)
			}
		}
	}

	return nil
}

func (service *serviceImpl) ReadOnly() ServiceCached {
	return &service.serviceCachedImpl
}

// ServiceCached is a read-only version of the Service interface.
type ServiceCached interface {
	ListPublic() []ed25519.PublicKey
	GetPrivate() ed25519.PrivateKey
	RefreshCache(ctx context.Context) error
}

type serviceCachedImpl struct {
	repository jwk_storage.Repository
	cached     []ed25519.PrivateKey
	mu         sync.RWMutex
}

// NewServiceCached returns a new implementation of Service.
//
//	jwk_service.NewServiceCached(repository)
func NewServiceCached(repository jwk_storage.Repository) ServiceCached {
	return &serviceCachedImpl{repository: repository}
}

func (service *serviceCachedImpl) ListPublic() []ed25519.PublicKey {
	output := make([]ed25519.PublicKey, len(service.cached))
	for i, v := range service.cached {
		output[i] = v.Public().(ed25519.PublicKey)
	}
	return output
}

func (service *serviceCachedImpl) GetPrivate() ed25519.PrivateKey {
	if service.cached == nil || len(service.cached) == 0 {
		return nil
	}
	return service.cached[0]
}

func (service *serviceCachedImpl) RefreshCache(ctx context.Context) error {
	keys, err := service.repository.List(ctx)
	if err != nil {
		return err
	}

	service.mu.Lock()
	defer service.mu.Unlock()

	service.cached = make([]ed25519.PrivateKey, len(keys))
	for i, key := range keys {
		service.cached[i] = key.Key
	}

	return nil
}
