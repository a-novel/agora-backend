package credentials_service

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"
)

var (
	emailUserRegexp   = regexp.MustCompile("^[a-zA-Z\\d!#$%&'*+/=?^_`{|}.~-]+$")
	emailDomainRegexp = regexp.MustCompile("^[a-z\\d]{2,}(.[a-z\\d]{2,})+$")
)

const (
	MinPasswordLength    = 2
	MaxPasswordLength    = 256
	MaxEmailUserLength   = 128
	MaxEmailDomainLength = 128
)

// Service of the current layer. You can instantiate a new one with NewService.
type Service interface {
	// PrepareRegistration computes the UserCredentialsLoginForm before sending it to user_service.Service.
	PrepareRegistration(ctx context.Context, data *models.UserCredentialsLoginForm) (*models.UserCredentialsRegistrationForm, error)
	// Authenticate verifies that the claims contained in UserCredentialsLoginForm match an existing user, and returns this user on
	// success.
	Authenticate(ctx context.Context, data *models.UserCredentialsLoginForm) (*models.UserCredentials, error)

	// Read reads a user, based on its ID.
	Read(ctx context.Context, id uuid.UUID) (*models.UserCredentials, error)
	// ReadEmail reads a user, based on its email. The match must be exact.
	// This method only searches for the main user email, it ignores pending state emails.
	ReadEmail(ctx context.Context, email string) (*models.UserCredentials, error)
	// EmailExists looks if a given email is already used by a user as their main email. The match must be exact.
	// This method only searches for the main user email, it ignores pending state emails.
	EmailExists(ctx context.Context, email string) (bool, error)

	// UpdateEmail updates the email of a user. The new email value is set in a special pending state.
	// The old email remains active as the main address, until the new value is validated.
	// To make the new email the primary email of the user, you must call ValidateNewEmail with the correct code.
	UpdateEmail(ctx context.Context, email string, id uuid.UUID, now time.Time) (*models.UserCredentials, string, error)
	// ValidateEmail validates the main email of the targeted user. The user is searched based on its main email.
	ValidateEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) (*models.UserCredentials, error)
	// ValidateNewEmail validates the pending update email of the targeted user.
	// The user is searched based on its MAIN email. Once validated, the new email becomes the main.
	ValidateNewEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) (*models.UserCredentials, error)
	// UpdateEmailValidation generates a new validation code for the main email of the targeted user. The main email
	// must be pending validation already.
	UpdateEmailValidation(ctx context.Context, id uuid.UUID, now time.Time) (*models.UserCredentials, string, error)
	// UpdateNewEmailValidation generates a new validation code for the pending update email of the targeted user.
	// There must be an email pending validation.
	UpdateNewEmailValidation(ctx context.Context, id uuid.UUID, now time.Time) (*models.UserCredentials, string, error)
	// CancelNewEmail cancels the pending validation email. It does not fail in case no email update is pending.
	CancelNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*models.UserCredentials, error)

	// UpdatePassword updates the targeted user password. The current password is required as an extra security.
	// If ResetPassword has been called, the code returned may be used in place of the old password, once only.
	UpdatePassword(ctx context.Context, oldPassword, newPassword string, id uuid.UUID, now time.Time) (*models.UserCredentials, error)
	// ResetPassword creates a code to securely update the password when the current one has been forgotten.
	ResetPassword(ctx context.Context, email string, now time.Time) (*models.UserCredentials, string, error)

	StorageToModel(source *credentials_storage.Model) *models.UserCredentials
}

type serviceImpl struct {
	repository credentials_storage.Repository

	generateCode           func() (string, string, error)
	verifyCode             func(code string, encrypted string) (bool, error)
	generateFromPassword   func(password []byte, cost int) ([]byte, error)
	compareHashAndPassword func(hashedPassword []byte, password []byte) error
}

// NewService returns a new implementation of Service.
//
//	credentials_service.NewService(
//	 	repository,
//	  	security.GenerateCode,
//	  	security.VerifyCode,
//	  	bcrypt.GenerateFromPassword,
//	  	bcrypt.CompareHashAndPassword,
//	)
func NewService(
	repository credentials_storage.Repository,
	generateCode func() (string, string, error),
	verifyCode func(code string, encrypted string) (bool, error),
	generateFromPassword func(password []byte, cost int) ([]byte, error),
	compareHashAndPassword func(hashedPassword []byte, password []byte) error,
) Service {
	return &serviceImpl{
		repository:             repository,
		generateCode:           generateCode,
		verifyCode:             verifyCode,
		generateFromPassword:   generateFromPassword,
		compareHashAndPassword: compareHashAndPassword,
	}
}

func (service *serviceImpl) PrepareRegistration(_ context.Context, data *models.UserCredentialsLoginForm) (*models.UserCredentialsRegistrationForm, error) {
	if err := validation.CheckRequire("password", data.Password); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("email", data.Email); err != nil {
		return nil, err
	}

	email, err := service.parseEmail(data.Email)
	if err != nil {
		return nil, err
	}
	if err := service.validateEmail(email); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("password", data.Password, MinPasswordLength, MaxPasswordLength); err != nil {
		return nil, err
	}

	// Generate the code to validate user email. The private (hashed) key goes in the database. The public key will
	// be sent to the user address, to ensure it is valid.
	publicEmailValidationCode, privateEmailValidationCode, err := service.generateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate email validation code: %w", err)
	}
	email.Validation = privateEmailValidationCode

	// Hash password before saving it to database.
	passwordHashed, err := service.generateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash user password: %w", err)
	}

	return &models.UserCredentialsRegistrationForm{
		Email:                     email,
		Password:                  models.Password{Hashed: string(passwordHashed)},
		EmailPublicValidationCode: publicEmailValidationCode,
	}, nil
}

func (service *serviceImpl) Authenticate(ctx context.Context, data *models.UserCredentialsLoginForm) (*models.UserCredentials, error) {
	if err := validation.CheckRequire("password", data.Password); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("email", data.Email); err != nil {
		return nil, err
	}

	parsedEmail, err := service.parseEmail(data.Email)
	if err != nil {
		return nil, err
	}
	if err := service.validateEmail(parsedEmail); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.ReadEmail(ctx, parsedEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credentials: %w", err)
	}

	// Validate the user provided the correct password.
	// If set, password reset key MUST NOT be used here.
	err = service.compareHashAndPassword([]byte(storageModel.Password.Hashed), []byte(data.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, validation.NewErrInvalidCredentials("password does not match the one in database")
		}
		return nil, fmt.Errorf("failed to verify user password: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) Read(ctx context.Context, id uuid.UUID) (*models.UserCredentials, error) {
	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credentials: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) ReadEmail(ctx context.Context, email string) (*models.UserCredentials, error) {
	if err := validation.CheckRequire("email", email); err != nil {
		return nil, err
	}

	parsedEmail, err := service.parseEmail(email)
	if err != nil {
		return nil, err
	}
	if err := service.validateEmail(parsedEmail); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.ReadEmail(ctx, parsedEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to get user credentials: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	if err := validation.CheckRequire("email", email); err != nil {
		return false, err
	}

	parsedEmail, err := service.parseEmail(email)
	if err != nil {
		return false, err
	}
	if err = service.validateEmail(parsedEmail); err != nil {
		return false, err
	}

	ok, err := service.repository.EmailExists(ctx, parsedEmail)
	if err != nil {
		return false, fmt.Errorf("failed to verify if email exists: %w", err)
	}

	return ok, nil
}

func (service *serviceImpl) UpdateEmail(ctx context.Context, email string, id uuid.UUID, now time.Time) (*models.UserCredentials, string, error) {
	if err := validation.CheckRequire("email", email); err != nil {
		return nil, "", err
	}

	// Ensure the new email is valid.
	parsedEmail, err := service.parseEmail(email)
	if err != nil {
		return nil, "", err
	}
	if err := service.validateEmail(parsedEmail); err != nil {
		return nil, "", err
	}

	// The new email must not already be taken by a user.
	exists, err := service.repository.EmailExists(ctx, parsedEmail)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check if email is taken: %w", err)
	}
	if exists {
		return nil, "", validation.ErrUniqConstraintViolation
	}

	// Generate the code to validate the new user email. The private (hashed) key goes in the database. The public key
	// will be sent to the user address, to ensure it is valid.
	publicEmailValidationCode, privateEmailValidationCode, err := service.generateCode()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate email validation code: %w", err)
	}

	storageModel, err := service.repository.UpdateEmail(ctx, parsedEmail, privateEmailValidationCode, id, now)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update user credentials: %w", err)
	}

	return service.StorageToModel(storageModel), publicEmailValidationCode, nil
}

func (service *serviceImpl) ValidateEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) (*models.UserCredentials, error) {
	if err := validation.CheckRequire("validation_code", code); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read user credentials: %w", err)
	}

	// The email must be pending validation, to avoid errors.
	if storageModel.Email.Validation == "" {
		return nil, validation.ErrValidated
	}

	// Ensure the code is correct.
	ok, err := service.verifyCode(code, storageModel.Email.Validation)
	if err != nil {
		return nil, fmt.Errorf("failed to verify email code: %w", err)
	}
	if !ok {
		return nil, validation.NewErrInvalidCredentials("validation code does not match the one in database")
	}

	storageModel, err = service.repository.ValidateEmail(ctx, storageModel.ID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update credentials: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) ValidateNewEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) (*models.UserCredentials, error) {
	if err := validation.CheckRequire("validation_code", code); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read user credentials: %w", err)
	}

	// The new email must be pending validation, to avoid errors.
	if storageModel.NewEmail.Validation == "" {
		return nil, validation.ErrValidated
	}

	// Ensure the code is correct.
	ok, err := service.verifyCode(code, storageModel.NewEmail.Validation)
	if err != nil {
		return nil, fmt.Errorf("failed to verify email code: %w", err)
	}
	if !ok {
		return nil, validation.NewErrInvalidCredentials("validation code does not match the one in database")
	}

	storageModel, err = service.repository.ValidateNewEmail(ctx, storageModel.ID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update credentials: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) UpdateEmailValidation(ctx context.Context, id uuid.UUID, now time.Time) (*models.UserCredentials, string, error) {
	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user credentials: %w", err)
	}

	if storageModel.Email.Validation == "" {
		return nil, "", validation.ErrValidated
	}

	publicEmailValidationCode, privateEmailValidationCode, err := service.generateCode()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate email validation code: %w", err)
	}

	storageModel, err = service.repository.UpdateEmailValidation(ctx, privateEmailValidationCode, id, now)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update user credentials: %w", err)
	}

	return service.StorageToModel(storageModel), publicEmailValidationCode, nil
}

func (service *serviceImpl) UpdateNewEmailValidation(ctx context.Context, id uuid.UUID, now time.Time) (*models.UserCredentials, string, error) {
	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user credentials: %w", err)
	}

	if storageModel.NewEmail.Validation == "" {
		return nil, "", validation.ErrValidated
	}

	publicEmailValidationCode, privateEmailValidationCode, err := service.generateCode()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate email validation code: %w", err)
	}

	storageModel, err = service.repository.UpdateNewEmailValidation(ctx, privateEmailValidationCode, id, now)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update user credentials: %w", err)
	}

	return service.StorageToModel(storageModel), publicEmailValidationCode, nil
}

func (service *serviceImpl) CancelNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*models.UserCredentials, error) {
	storageModel, err := service.repository.CancelNewEmail(ctx, id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update user credentials: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) UpdatePassword(ctx context.Context, oldPassword, newPassword string, id uuid.UUID, now time.Time) (*models.UserCredentials, error) {
	var resetCodeValidated bool

	if err := validation.CheckRequire("new_password", newPassword); err != nil {
		return nil, err
	}
	if err := validation.CheckRequire("old_password", oldPassword); err != nil {
		return nil, err
	}
	if err := validation.CheckMinMax("password", newPassword, MinPasswordLength, MaxPasswordLength); err != nil {
		return nil, err
	}

	storageModel, err := service.repository.Read(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	// If a password reset was asked, the reset code can be used in place of the old password. Try this first, and
	// fallback on standard password validation otherwise.
	if storageModel.Password.Validation != "" {
		ok, err := service.verifyCode(oldPassword, storageModel.Password.Validation)
		if err != nil {
			return nil, fmt.Errorf("failed to verify user password: %w", err)
		}

		resetCodeValidated = ok
	}

	if !resetCodeValidated {
		// Verify if the hashed password in database matches the user value.
		err = service.compareHashAndPassword([]byte(storageModel.Password.Hashed), []byte(oldPassword))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return nil, validation.NewErrInvalidCredentials("the password entered does not match the database value")
			}
			return nil, fmt.Errorf("failed to verify user password: %w", err)
		}
	}

	passwordHashed, err := service.generateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash user password: %w", err)
	}

	storageModel, err = service.repository.UpdatePassword(ctx, string(passwordHashed), id, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update user credentials: %w", err)
	}

	return service.StorageToModel(storageModel), nil
}

func (service *serviceImpl) ResetPassword(ctx context.Context, email string, now time.Time) (*models.UserCredentials, string, error) {
	if err := validation.CheckRequire("email", email); err != nil {
		return nil, "", err
	}

	parsedEmail, err := service.parseEmail(email)
	if err != nil {
		return nil, "", err
	}
	if err := service.validateEmail(parsedEmail); err != nil {
		return nil, "", err
	}

	storageModel, err := service.repository.ReadEmail(ctx, parsedEmail)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user credentials: %w", err)
	}

	publicPasswordValidationCode, privatePasswordValidationCode, err := service.generateCode()
	if err != nil {
		return nil, "", fmt.Errorf("failed to geenerate password validation code: %w", err)
	}

	storageModel, err = service.repository.ResetPassword(ctx, privatePasswordValidationCode, parsedEmail, now)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update user credentials: %w", err)
	}

	return service.StorageToModel(storageModel), publicPasswordValidationCode, nil
}

func (service *serviceImpl) StorageToModel(source *credentials_storage.Model) *models.UserCredentials {
	if source == nil {
		return nil
	}

	return &models.UserCredentials{
		ID:        source.ID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		Email:     source.Email.String(),
		NewEmail:  source.NewEmail.String(),
		Validated: source.Email.Validation == "",
	}
}

func (service *serviceImpl) parseEmail(source string) (models.Email, error) {
	var model models.Email

	parts := strings.Split(source, "@")
	if len(parts) != 2 {
		return model, validation.NewErrInvalidEntity("email", "must be of format [user]@[domain]")
	}

	model.User = parts[0]
	model.Domain = parts[1]

	return model, nil
}

func (service *serviceImpl) validateEmail(source models.Email) error {
	if err := validation.CheckMinMax("email.user", source.User, 1, MaxEmailUserLength); err != nil {
		return err
	}
	if err := validation.CheckMinMax("email.domain", source.Domain, 1, MaxEmailDomainLength); err != nil {
		return err
	}
	if err := validation.CheckRegexp("email.user", source.User, emailUserRegexp); err != nil {
		return err
	}
	if err := validation.CheckRegexp("email.domain", source.Domain, emailDomainRegexp); err != nil {
		return err
	}

	return nil
}
