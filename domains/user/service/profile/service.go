package profile_service

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"regexp"
	"time"
)

var (
	usernameRegexp = regexp.MustCompile(`^[\p{L}\p{N}\p{S}\p{P}\p{M}]+( ([\p{L}\p{N}\p{S}\p{P}\p{M}]+))*$`)
	slugRegexp     = regexp.MustCompile(`^[a-z\d]+(-[a-z\d]+)*$`)
)

const (
	MaxUsernameLength = 64
	MaxSlugLength     = 64
)

type Service interface {
	PrepareRegistration(ctx context.Context, data *models.UserProfileUpdateForm) (*models.UserProfileRegistrationForm, error)

	Read(ctx context.Context, id uuid.UUID) (*models.UserProfile, error)
	ReadSlug(ctx context.Context, slug string) (*models.UserProfile, error)
	SlugExists(ctx context.Context, slug string) (bool, error)

	Update(ctx context.Context, data *models.UserProfileUpdateForm, id uuid.UUID, now time.Time) (*models.UserProfile, error)

	StorageToModel(source *profile_storage.Model) *models.UserProfile
}

type serviceImpl struct {
	repository profile_storage.Repository
}

func NewService(repository profile_storage.Repository) Service {
	return &serviceImpl{repository: repository}
}

func (service *serviceImpl) PrepareRegistration(_ context.Context, data *models.UserProfileUpdateForm) (*models.UserProfileRegistrationForm, error) {
	if err := validation.CheckRequire("slug", data.Slug); err != nil {
		return nil, err
	}

	if err := validation.CheckMinMax("username", data.Username, -1, MaxUsernameLength); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("slug", data.Slug, 1, MaxSlugLength); err != nil {
		return nil, err
	}

	if data.Username != "" {
		if err := validation.CheckRegexp("username", data.Username, usernameRegexp); err != nil {
			return nil, err
		}
	}
	if err := validation.CheckRegexp("slug", data.Slug, slugRegexp); err != nil {
		return nil, err
	}

	return &models.UserProfileRegistrationForm{
		Username: data.Username,
		Slug:     data.Slug,
	}, nil
}

func (service *serviceImpl) Read(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) ReadSlug(ctx context.Context, slug string) (*models.UserProfile, error) {
	storageModel, err := service.repository.ReadSlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) SlugExists(ctx context.Context, slug string) (bool, error) {
	ok, err := service.repository.SlugExists(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("failed to get user profile: %w", err)
	}

	return ok, nil
}

func (service *serviceImpl) Update(ctx context.Context, data *models.UserProfileUpdateForm, id uuid.UUID, now time.Time) (*models.UserProfile, error) {
	registerModel, err := service.PrepareRegistration(ctx, data)
	if err != nil {
		return nil, err
	}

	storageCore := &profile_storage.Core{
		Slug:     registerModel.Slug,
		Username: registerModel.Username,
	}
	storageModel, err := service.repository.Update(ctx, storageCore, id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) StorageToModel(source *profile_storage.Model) *models.UserProfile {
	if source == nil {
		return nil
	}

	return &models.UserProfile{
		ID:        source.ID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		Username:  source.Username,
		Slug:      source.Slug,
	}
}
