package user_service

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/domains/user/service/credentials"
	"github.com/a-novel/agora-backend/domains/user/service/identity"
	"github.com/a-novel/agora-backend/domains/user/service/profile"
	"github.com/a-novel/agora-backend/domains/user/storage/credentials"
	"github.com/a-novel/agora-backend/domains/user/storage/identity"
	"github.com/a-novel/agora-backend/domains/user/storage/profile"
	"github.com/a-novel/agora-backend/domains/user/storage/user"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	Create(ctx context.Context, data *models.UserCreateForm, id uuid.UUID, now time.Time) (*models.User, *models.UserPostRegistration, error)
	Delete(ctx context.Context, id uuid.UUID, now time.Time) (*models.User, error)

	Search(ctx context.Context, query string, limit, offset int) ([]*models.UserPublicPreview, int64, error)
	GetPreview(ctx context.Context, id uuid.UUID) (*models.UserPreview, error)
	GetPublic(ctx context.Context, slug string) (*models.UserPublic, error)
	GetPublicPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.UserPublicPreview, error)
	GetAuthorizations(ctx context.Context, id uuid.UUID) ([]string, error)
	HasAuthorizations(ctx context.Context, id uuid.UUID, authorizations models.UserAuthorizations) (bool, error)

	StorageToModel(source *user_storage.Model) *models.User
}

type serviceImpl struct {
	repository         user_storage.Repository
	credentialsService credentials_service.Service
	identityService    identity_service.Service
	profileService     profile_service.Service
}

func NewService(
	repository user_storage.Repository,
	credentialsService credentials_service.Service,
	identityService identity_service.Service,
	profileService profile_service.Service,
) Service {
	return &serviceImpl{
		repository:         repository,
		credentialsService: credentialsService,
		identityService:    identityService,
		profileService:     profileService,
	}
}

func (service *serviceImpl) publicPreviewStorageToModel(storageModel *user_storage.PublicPreview) *models.UserPublicPreview {
	output := &models.UserPublicPreview{
		ID:        storageModel.ID,
		CreatedAt: storageModel.CreatedAt,
		Username:  storageModel.Username,
		Slug:      storageModel.Slug,
	}

	// Don't return real name when username is set.
	if storageModel.Username == "" {
		output.FirstName = storageModel.FirstName
		output.LastName = storageModel.LastName
	}

	return output
}

func (service *serviceImpl) publicStorageToModel(storageModel *user_storage.Public) *models.UserPublic {
	output := &models.UserPublic{
		ID:        storageModel.ID,
		CreatedAt: storageModel.CreatedAt,
		Username:  storageModel.Username,
		Sex:       storageModel.Sex,
	}

	// Don't return real name when username is set.
	if storageModel.Username == "" {
		output.FirstName = storageModel.FirstName
		output.LastName = storageModel.LastName
	}

	return output
}

func (service *serviceImpl) Create(ctx context.Context, data *models.UserCreateForm, id uuid.UUID, now time.Time) (*models.User, *models.UserPostRegistration, error) {
	credentialsRegisterModel, err := service.credentialsService.PrepareRegistration(ctx, &data.Credentials)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid user credentials: %w", err)
	}

	identityRegisterModel, err := service.identityService.PrepareRegistration(ctx, &data.Identity, now)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid user identity: %w", err)
	}

	profileRegisterModel, err := service.profileService.PrepareRegistration(ctx, &data.Profile)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid user profile: %w", err)
	}

	storageCore := &user_storage.Core{
		Credentials: credentials_storage.Core{
			Email:    credentialsRegisterModel.Email,
			Password: credentialsRegisterModel.Password,
		},
		Identity: identity_storage.Core{
			FirstName: identityRegisterModel.FirstName,
			LastName:  identityRegisterModel.LastName,
			Birthday:  identityRegisterModel.Birthday,
			Sex:       identityRegisterModel.Sex,
		},
		Profile: profile_storage.Core{
			Username: profileRegisterModel.Username,
			Slug:     profileRegisterModel.Slug,
		},
	}

	storageModel, err := service.repository.Create(ctx, storageCore, id, now)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	registerModel := &models.UserPostRegistration{
		EmailValidationCode: credentialsRegisterModel.EmailPublicValidationCode,
	}

	return service.StorageToModel(storageModel), registerModel, nil
}

func (service *serviceImpl) Delete(ctx context.Context, id uuid.UUID, now time.Time) (*models.User, error) {
	storageModel, err := service.repository.Delete(ctx, id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) Search(ctx context.Context, query string, limit, offset int) ([]*models.UserPublicPreview, int64, error) {
	storageModels, count, err := service.repository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	results := make([]*models.UserPublicPreview, len(storageModels))
	for i, storageModel := range storageModels {
		results[i] = service.publicPreviewStorageToModel(storageModel)
	}

	return results, count, nil
}

func (service *serviceImpl) GetPreview(ctx context.Context, id uuid.UUID) (*models.UserPreview, error) {
	storageModel, err := service.repository.GetPreview(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user preview: %w", err)
	}

	return &models.UserPreview{
		ID:        storageModel.ID,
		Username:  storageModel.Username,
		FirstName: storageModel.FirstName,
		LastName:  storageModel.LastName,
		Email:     storageModel.Email.String(),
		Slug:      storageModel.Slug,
		Sex:       storageModel.Sex,
	}, nil
}

func (service *serviceImpl) GetPublic(ctx context.Context, slug string) (*models.UserPublic, error) {
	storageModel, err := service.repository.GetPublic(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get user preview: %w", err)
	}

	return service.publicStorageToModel(storageModel), nil
}

func (service *serviceImpl) GetPublicPreviews(ctx context.Context, ids []uuid.UUID) ([]*models.UserPublicPreview, error) {
	storageModels, err := service.repository.GetPublicPreviews(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get user previews: %w", err)
	}

	results := make([]*models.UserPublicPreview, len(storageModels))
	for i, storageModel := range storageModels {
		results[i] = service.publicPreviewStorageToModel(storageModel)
	}

	return results, nil
}

func (service *serviceImpl) GetAuthorizations(ctx context.Context, id uuid.UUID) ([]string, error) {
	storageModel, err := service.credentialsService.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	authorizations := map[string]bool{
		models.UserAuthorizationsAccountValidated: storageModel.Validated,
	}

	var output []string

	for authorization, ok := range authorizations {
		if ok {
			output = append(output, authorization)
		}
	}

	return output, nil
}

func (service *serviceImpl) hasAllAuthorizations(authorizations []string, userAuthorizations map[string]bool) bool {
	for _, authorization := range authorizations {
		if _, ok := userAuthorizations[authorization]; !ok {
			return false
		}
	}

	return true
}

func (service *serviceImpl) HasAuthorizations(ctx context.Context, id uuid.UUID, authorizations models.UserAuthorizations) (bool, error) {
	userAuthorizations, err := service.GetAuthorizations(ctx, id)
	if err != nil {
		return false, err
	}

	parsedUserAuthorizations := make(map[string]bool, len(userAuthorizations))
	for _, authorization := range userAuthorizations {
		parsedUserAuthorizations[authorization] = true
	}

	for _, authorizationsGroup := range authorizations {
		if service.hasAllAuthorizations(authorizationsGroup, parsedUserAuthorizations) {
			return true, nil
		}
	}

	return false, nil
}

func (service *serviceImpl) StorageToModel(source *user_storage.Model) *models.User {
	if source == nil {
		return nil
	}

	return &models.User{
		ID:        source.ID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		DeletedAt: source.DeletedAt,
		Credentials: *service.credentialsService.StorageToModel(&credentials_storage.Model{
			ID:        source.ID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
			Core:      source.Credentials,
		}),
		Identity: *service.identityService.StorageToModel(&identity_storage.Model{
			ID:        source.ID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
			DeletedAt: source.DeletedAt,
			Core:      source.Identity,
		}),
		Profile: *service.profileService.StorageToModel(&profile_storage.Model{
			ID:        source.ID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
			Core:      source.Profile,
		}),
	}
}
