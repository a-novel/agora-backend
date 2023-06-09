package identity_service

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
	// A valid name:
	//   - must start with a letter character.
	//   - must end with a letter character.
	//   - may contain letters and hyphens.
	nameRegexp = regexp.MustCompile(`^\p{L}+([- ']\p{L}+)*$`)
)

const (
	MaxNameLength = 32
	MinUserAge    = 16
	MaxUserAge    = 150
)

type Service interface {
	PrepareRegistration(ctx context.Context, data *models.UserIdentityUpdateForm, now time.Time) (*models.UserIdentityRegistrationForm, error)

	Read(ctx context.Context, id uuid.UUID) (*models.UserIdentity, error)
	Update(ctx context.Context, data *models.UserIdentityUpdateForm, id uuid.UUID, now time.Time) (*models.UserIdentity, error)

	StorageToModel(source *identity_storage.Model) *models.UserIdentity
	Age(source *models.UserIdentity, now time.Time) int
}

type serviceImpl struct {
	repository identity_storage.Repository
}

func NewService(repository identity_storage.Repository) Service {
	return &serviceImpl{repository: repository}
}

func (service *serviceImpl) PrepareRegistration(_ context.Context, data *models.UserIdentityUpdateForm, now time.Time) (*models.UserIdentityRegistrationForm, error) {
	if err := validation.CheckRequire("firstName", data.FirstName); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("lastName", data.LastName); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("birthday", data.Birthday); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("sex", data.Sex); err != nil {
		return nil, err
	}

	if err := validation.CheckMinMax("firstName", data.FirstName, 1, MaxNameLength); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("lastName", data.LastName, 1, MaxNameLength); err != nil {
		return nil, err
	}

	if err := validation.CheckRegexp("firstName", data.FirstName, nameRegexp); err != nil {
		return nil, err
	}
	if err := validation.CheckRegexp("lastName", data.LastName, nameRegexp); err != nil {
		return nil, err
	}

	if data.Sex != models.SexMale && data.Sex != models.SexFemale {
		return nil, validation.NewErrInvalidEntity("sex", "can only be male or female")
	}

	age := service.Age(&models.UserIdentity{Birthday: data.Birthday}, now)
	if err := validation.CheckMinMax("age", age, MinUserAge, MaxUserAge); err != nil {
		return nil, err
	}

	return &models.UserIdentityRegistrationForm{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Birthday:  data.Birthday,
		Sex:       data.Sex,
	}, nil
}

func (service *serviceImpl) Read(ctx context.Context, id uuid.UUID) (*models.UserIdentity, error) {
	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user identity: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) Update(ctx context.Context, data *models.UserIdentityUpdateForm, id uuid.UUID, now time.Time) (*models.UserIdentity, error) {
	registerModel, err := service.PrepareRegistration(ctx, data, now)
	if err != nil {
		return nil, err
	}

	storageCore := &identity_storage.Core{
		FirstName: registerModel.FirstName,
		LastName:  registerModel.LastName,
		Birthday:  registerModel.Birthday,
		Sex:       registerModel.Sex,
	}
	storageModel, err := service.repository.Update(ctx, storageCore, id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update user identity: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) StorageToModel(source *identity_storage.Model) *models.UserIdentity {
	if source == nil {
		return nil
	}

	return &models.UserIdentity{
		ID:        source.ID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		FirstName: source.FirstName,
		LastName:  source.LastName,
		Birthday:  source.Birthday,
		Sex:       source.Sex,
	}
}

func (service *serviceImpl) Age(source *models.UserIdentity, now time.Time) int {
	return now.In(source.Birthday.Location()).AddDate(
		-source.Birthday.Year(),
		// Because month and day start at 1 rather than 0, we have to account for this difference.
		// -x+1 = -(x-1), to translate them back to 0 based values.
		-int(source.Birthday.Month())+1,
		-source.Birthday.Day()+1,
	).Year()
}
