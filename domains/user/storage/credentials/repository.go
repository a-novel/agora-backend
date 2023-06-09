package credentials_storage

import (
	"context"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Repository of the current layer. You can instantiate a new one with NewRepository.
type Repository interface {
	// Read reads a credentials object, based on a user id.
	Read(ctx context.Context, id uuid.UUID) (*Model, error)
	// ReadEmail reads a credentials object, based on a user email.
	// The match should be exact, thus requiring to pass a full Email object, rather than a string representation.
	// The Email.Validation field is ignored, and the Core.NewEmail is not used for matching.
	ReadEmail(ctx context.Context, email models.Email) (*Model, error)
	// EmailExists looks if a given email is already used by a user as their main email (Core.Email).
	// The match should be exact, thus requiring to pass a full Email object, rather than a string representation.
	// The Email.Validation field is ignored, and the Core.NewEmail is not used for matching.
	EmailExists(ctx context.Context, email models.Email) (bool, error)

	// UpdateEmail updates the email of a user. The new email value is set as Core.NewEmail.
	// The validation code should not be set on the Email.Validation field, as it will be filtered. The code value
	// MUST be hashed.
	// To make the new email the primary email of the user, you must call ValidateNewEmail.
	UpdateEmail(ctx context.Context, email models.Email, code string, id uuid.UUID, now time.Time) (*Model, error)
	// ValidateEmail nullifies the Email.Validation value of Core.Email, for the targeted user.
	ValidateEmail(ctx context.Context, id uuid.UUID, now time.Time) (*Model, error)
	// ValidateNewEmail sets the email in argument as the primary email (Core.Email) for the targeted user.
	// The Core.NewEmail value is nullified in the process, and Email.Validation is filtered.
	ValidateNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*Model, error)

	// UpdateEmailValidation sets a new Email.Validation code for the targeted user Core.Email.
	// The code value MUST be hashed.
	UpdateEmailValidation(ctx context.Context, code string, id uuid.UUID, now time.Time) (*Model, error)
	// UpdateNewEmailValidation sets a new Email.Validation code for the targeted user Core.NewEmail.
	// The code value MUST be hashed. This method fails with sql.ErrNoRows if Core.NewEmail contains an empty email
	// value.
	UpdateNewEmailValidation(ctx context.Context, code string, id uuid.UUID, now time.Time) (*Model, error)
	// CancelNewEmail nullifies the Core.NewEmail value for the targeted user. It does not fail if this field was
	// already empty.
	CancelNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*Model, error)

	// UpdatePassword updates the password of the targeted user. The password value MUST be hashed in order to be
	// saved properly.
	UpdatePassword(ctx context.Context, newPassword string, id uuid.UUID, now time.Time) (*Model, error)
	// ResetPassword sets Password.Validation field. The code value MUST be hashed. This does not nullify the
	// Password.Hashed field, so authentication can still work while password is being reset.
	ResetPassword(ctx context.Context, code string, email models.Email, now time.Time) (*Model, error)
}

func NewRepository(db bun.IDB) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db bun.IDB
}

func (repository *repositoryImpl) Read(ctx context.Context, id uuid.UUID) (*Model, error) {
	model := &Model{ID: id}

	if err := repository.db.NewSelect().Model(model).WherePK().Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) ReadEmail(ctx context.Context, email models.Email) (*Model, error) {
	model := new(Model)

	if err := repository.db.NewSelect().Model(model).Where(WhereEmail("email", email)).Scan(ctx); err != nil {
		return nil, validation.HandlePGError(err)
	}

	return model, nil
}

func (repository *repositoryImpl) EmailExists(ctx context.Context, email models.Email) (bool, error) {
	count, err := repository.db.NewSelect().Model(new(Model)).Where(WhereEmail("email", email)).Count(ctx)
	return count > 0, validation.HandlePGError(err)
}

func (repository *repositoryImpl) UpdateEmail(ctx context.Context, email models.Email, code string, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{
		ID:        id,
		UpdatedAt: &now,
		// Set new email with the given validation code. The main email remains unchanged until this email is
		// validated.
		Core: Core{
			NewEmail: models.Email{User: email.User, Domain: email.Domain, Validation: code},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		Column("new_email_user", "new_email_domain", "new_email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	if err = validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *repositoryImpl) ValidateEmail(ctx context.Context, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{ID: id, UpdatedAt: &now}
	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// User must have a pending email validation.
		Where("email_validation_code != ''").
		Column("email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	if err = validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *repositoryImpl) ValidateNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{ID: id}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// User must have a pending email update.
		Where("new_email_validation_code != ''").
		// Use the pending update ONLY to update the main email.
		SetColumn("email_user", "new_email_user").
		SetColumn("email_domain", "new_email_domain").
		SetColumn("email_validation_code", "''").
		// Empty the new_email columns, and update timestamps.
		SetColumn("new_email_user", "''").
		SetColumn("new_email_domain", "''").
		SetColumn("new_email_validation_code", "''").
		SetColumn("updated_at", "?", now).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	if err = validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *repositoryImpl) UpdatePassword(ctx context.Context, newPassword string, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{
		ID:        id,
		UpdatedAt: &now,
		Core: Core{
			Password: models.Password{Hashed: newPassword},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// "password_validation_code" is important to invalidate any pending reset, since a new known password is
		// now available.
		Column("password_hashed", "password_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	if err = validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *repositoryImpl) ResetPassword(ctx context.Context, code string, email models.Email, now time.Time) (*Model, error) {
	model := &Model{
		UpdatedAt: &now,
		Core: Core{
			Password: models.Password{Validation: code},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		Where(WhereEmail("email", email)).
		Column("password_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	if err = validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *repositoryImpl) UpdateEmailValidation(ctx context.Context, code string, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{
		ID:        id,
		UpdatedAt: &now,
		Core: Core{
			Email: models.Email{Validation: code},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// User must have a pending validation update.
		Where("email_validation_code != ''").
		Column("email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	if err = validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *repositoryImpl) UpdateNewEmailValidation(ctx context.Context, code string, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{
		ID:        id,
		UpdatedAt: &now,
		Core: Core{
			NewEmail: models.Email{Validation: code},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// User must have a pending email update.
		Where("new_email_validation_code != ''").
		Column("new_email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	if err = validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *repositoryImpl) CancelNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*Model, error) {
	model := &Model{ID: id, UpdatedAt: &now}
	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		Column("new_email_user", "new_email_domain", "new_email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, validation.HandlePGError(err)
	}

	if err = validation.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}
