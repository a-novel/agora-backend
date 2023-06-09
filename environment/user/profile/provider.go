package profile

import (
	"context"
	"github.com/a-novel/agora-backend/domains/user/service/user"
	"github.com/google/uuid"
)

type Config struct {
	UserService user_service.Service
}

type Provider interface {
	Read(ctx context.Context, slug string) (*Model, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*Preview, int64, error)
	Previews(ctx context.Context, ids []uuid.UUID) ([]*Preview, error)
}

type providerImpl struct {
	userService user_service.Service
}

func NewProvider(cfg Config) Provider {
	return &providerImpl{
		userService: cfg.UserService,
	}
}

func (provider *providerImpl) Read(ctx context.Context, slug string) (*Model, error) {
	user, err := provider.userService.GetPublic(ctx, slug)
	if err != nil {
		return nil, err
	}

	return &Model{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		Sex:       user.Sex,
	}, nil
}

func (provider *providerImpl) Search(ctx context.Context, query string, limit, offset int) ([]*Preview, int64, error) {
	users, count, err := provider.userService.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	profiles := make([]*Preview, 0, len(users))
	for _, user := range users {
		profiles = append(profiles, &Preview{
			ID:        user.ID,
			Slug:      user.Slug,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
		})
	}

	return profiles, count, nil
}

func (provider *providerImpl) Previews(ctx context.Context, ids []uuid.UUID) ([]*Preview, error) {
	users, err := provider.userService.GetPublicPreviews(ctx, ids)
	if err != nil {
		return nil, err
	}

	profiles := make([]*Preview, 0, len(users))
	for _, user := range users {
		profiles = append(profiles, &Preview{
			ID:        user.ID,
			Slug:      user.Slug,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
		})
	}

	return profiles, nil
}
