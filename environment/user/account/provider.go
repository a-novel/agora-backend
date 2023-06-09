package account

import (
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/domains/generics"
	"github.com/a-novel/agora-backend/domains/keys/service/jwk"
	"github.com/a-novel/agora-backend/domains/user/service/credentials"
	"github.com/a-novel/agora-backend/domains/user/service/identity"
	"github.com/a-novel/agora-backend/domains/user/service/profile"
	"github.com/a-novel/agora-backend/domains/user/service/token"
	"github.com/a-novel/agora-backend/domains/user/service/user"
	"github.com/a-novel/agora-backend/environment"
	"github.com/a-novel/agora-backend/environment/user/authentication"
	"github.com/a-novel/agora-backend/framework"
	"github.com/a-novel/agora-backend/framework/mailer"
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"time"
)

type Provider interface {
	Register(ctx context.Context, form models.UserCreateForm) (*models.UserFlat, string, environment.Deferred, error)

	GetAccountInfo(ctx context.Context, token string) (*models.UserInfo, error)
	GetAccountPreview(ctx context.Context, token string) (*models.UserPreview, error)
	GetEmailValidationStatus(ctx context.Context, token string) (*models.UserEmailValidationStatus, error)
	GetAuthorizations(ctx context.Context, token string) ([]string, error)

	UpdateIdentity(ctx context.Context, token string, form models.UserIdentityUpdateForm) (*models.UserInfoIdentity, error)
	UpdateProfile(ctx context.Context, token string, form models.UserProfileUpdateForm) (*models.UserInfoProfile, error)
	UpdatePassword(ctx context.Context, form models.UserPasswordUpdateForm) error
	UpdateEmail(ctx context.Context, token string, form models.UserEmailUpdateForm) (*models.UserEmailValidationStatus, environment.Deferred, error)
	CancelNewEmail(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, form models.UserPasswordResetForm) (environment.Deferred, error)

	ValidateEmail(ctx context.Context, form models.UserValidateEmailForm) error
	ValidateNewEmail(ctx context.Context, form models.UserValidateEmailForm) error
	ResendEmailValidation(ctx context.Context, token string) (environment.Deferred, error)
	ResendNewEmailValidation(ctx context.Context, token string) (environment.Deferred, error)

	DoesSlugExist(ctx context.Context, slug string) (bool, error)
	DoesEmailExist(ctx context.Context, email string) (bool, error)
}

type Config struct {
	CredentialsService credentials_service.Service
	IdentityService    identity_service.Service
	ProfileService     profile_service.Service
	UserService        user_service.Service
	TokenService       token_service.Service
	KeysService        jwk_service.ServiceCached
	Mailer             mailer.Mailer

	Time func() time.Time
	ID   func() uuid.UUID

	TokenTTL                   time.Duration
	TokenRenewDelta            time.Duration
	EmailValidationLink        generics.URL
	NewEmailValidationLink     generics.URL
	PasswordResetLink          generics.URL
	EmailValidationTemplate    string
	NewEMailValidationTemplate string
	PasswordResetTemplate      string
}

type providerImpl struct {
	credentialsService credentials_service.Service
	identityService    identity_service.Service
	profileService     profile_service.Service
	userService        user_service.Service
	tokenService       token_service.Service
	keysService        jwk_service.ServiceCached
	mailer             mailer.Mailer

	time func() time.Time
	id   func() uuid.UUID

	tokenTTL                   time.Duration
	tokenRenewDelta            time.Duration
	emailValidationLink        generics.URL
	newEmailValidationLink     generics.URL
	passwordResetLink          generics.URL
	emailValidationTemplate    string
	newEmailValidationTemplate string
	passwordResetTemplate      string
}

func NewProvider(cfg Config) Provider {
	return &providerImpl{
		credentialsService: cfg.CredentialsService,
		identityService:    cfg.IdentityService,
		profileService:     cfg.ProfileService,
		userService:        cfg.UserService,
		tokenService:       cfg.TokenService,
		keysService:        cfg.KeysService,
		mailer:             cfg.Mailer,

		time: cfg.Time,
		id:   cfg.ID,

		tokenTTL:                   cfg.TokenTTL,
		tokenRenewDelta:            cfg.TokenRenewDelta,
		emailValidationLink:        cfg.EmailValidationLink,
		newEmailValidationLink:     cfg.NewEmailValidationLink,
		passwordResetLink:          cfg.PasswordResetLink,
		emailValidationTemplate:    cfg.EmailValidationTemplate,
		newEmailValidationTemplate: cfg.NewEMailValidationTemplate,
		passwordResetTemplate:      cfg.PasswordResetTemplate,
	}
}

func (provider *providerImpl) Register(ctx context.Context, form models.UserCreateForm) (*models.UserFlat, string, environment.Deferred, error) {
	now := provider.time()
	id := provider.id()

	// Generate token first, so if it fails, we don't insert useless data.
	token, err := provider.tokenService.Encode(
		models.UserTokenPayload{ID: id},
		provider.tokenTTL,
		provider.keysService.GetPrivate(),
		provider.id(),
		provider.time(),
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to generate token for new user: %w", err)
	}

	user, postRegistration, err := provider.userService.Create(ctx, &form, id, now)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to create new user with email %q: %w", form.Credentials.Email, err)
	}

	output := &models.UserFlat{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Credentials.Email,
		NewEmail:  user.Credentials.NewEmail,
		Validated: user.Credentials.Validated,
		FirstName: user.Identity.FirstName,
		LastName:  user.Identity.LastName,
		Birthday:  user.Identity.Birthday,
		Sex:       user.Identity.Sex,
		Username:  user.Profile.Username,
		Slug:      user.Profile.Slug,
	}

	return output, token, func() error {
		name := user.Identity.FirstName
		toEmail := mail.NewEmail(name, user.Credentials.Email)
		templateData := map[string]interface{}{
			"name": name,
			"validation_link": provider.emailValidationLink.WithQuery(map[string]interface{}{
				"id":   user.ID,
				"code": postRegistration.EmailValidationCode,
			}).String(),
		}

		if err := provider.mailer.Send(toEmail, provider.emailValidationTemplate, templateData); err != nil {
			return fmt.Errorf("failed to send email validation link to user %q: %w", user.Credentials.Email, err)
		}

		return nil
	}, nil
}

func (provider *providerImpl) GetAccountInfo(ctx context.Context, token string) (*models.UserInfo, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	credentials, err := provider.credentialsService.Read(ctx, claims.Payload.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch credentials for user %q: %w", claims.Payload.ID, err)
	}
	identity, err := provider.identityService.Read(ctx, claims.Payload.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch identity for user %q: %w", claims.Payload.ID, err)
	}
	profile, err := provider.profileService.Read(ctx, claims.Payload.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile for user %q: %w", claims.Payload.ID, err)
	}

	return &models.UserInfo{
		ID:        claims.Payload.ID,
		CreatedAt: credentials.CreatedAt,
		UpdatedAt: framework.MostRecent(credentials.UpdatedAt, identity.UpdatedAt, profile.UpdatedAt),
		Email:     credentials.Email,
		NewEmail:  credentials.NewEmail,
		Identity: models.UserInfoIdentity{
			FirstName: identity.FirstName,
			LastName:  identity.LastName,
			Birthday:  identity.Birthday,
			Sex:       identity.Sex,
		},
		Profile: models.UserInfoProfile{
			Username: profile.Username,
			Slug:     profile.Slug,
		},
	}, nil
}

func (provider *providerImpl) GetAccountPreview(ctx context.Context, token string) (*models.UserPreview, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	preview, err := provider.userService.GetPreview(ctx, claims.Payload.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch preview for user %q: %w", claims.Payload.ID, err)
	}

	return preview, nil
}

func (provider *providerImpl) GetEmailValidationStatus(ctx context.Context, token string) (*models.UserEmailValidationStatus, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	credentials, err := provider.credentialsService.Read(ctx, claims.Payload.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch credentials for user %q: %w", claims.Payload.ID, err)
	}

	return &models.UserEmailValidationStatus{
		Email:     credentials.Email,
		NewEmail:  credentials.NewEmail,
		Validated: credentials.Validated,
	}, nil
}

func (provider *providerImpl) GetAuthorizations(ctx context.Context, token string) ([]string, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	authorizations, err := provider.userService.GetAuthorizations(ctx, claims.Payload.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch authorizations for user %q: %w", claims.Payload.ID, err)
	}

	return authorizations, nil
}

func (provider *providerImpl) UpdateIdentity(ctx context.Context, token string, form models.UserIdentityUpdateForm) (*models.UserInfoIdentity, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	identity, err := provider.identityService.Update(ctx, &form, claims.Payload.ID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update identity for user %q: %w", claims.Payload.ID, err)
	}

	return &models.UserInfoIdentity{
		FirstName: identity.FirstName,
		LastName:  identity.LastName,
		Birthday:  identity.Birthday,
		Sex:       identity.Sex,
	}, nil
}

func (provider *providerImpl) UpdateProfile(ctx context.Context, token string, form models.UserProfileUpdateForm) (*models.UserInfoProfile, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	profile, err := provider.profileService.Update(ctx, &form, claims.Payload.ID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile for user %q: %w", claims.Payload.ID, err)
	}

	return &models.UserInfoProfile{
		Username: profile.Username,
		Slug:     profile.Slug,
	}, nil
}

func (provider *providerImpl) UpdatePassword(ctx context.Context, form models.UserPasswordUpdateForm) error {
	_, err := provider.credentialsService.UpdatePassword(ctx, form.OldPassword, form.Password, form.ID, provider.time())
	return err
}

func (provider *providerImpl) UpdateEmail(ctx context.Context, token string, form models.UserEmailUpdateForm) (*models.UserEmailValidationStatus, environment.Deferred, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, nil, err
	}

	credentials, validationCode, err := provider.credentialsService.UpdateEmail(ctx, form.Email, claims.Payload.ID, now)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update email to %q for user %q: %w", form.Email, claims.Payload.ID, err)
	}
	identity, err := provider.identityService.Read(ctx, claims.Payload.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch identity for user %q: %w", claims.Payload.ID, err)
	}

	return &models.UserEmailValidationStatus{
			Email:     credentials.Email,
			NewEmail:  credentials.NewEmail,
			Validated: credentials.Validated,
		}, func() error {
			name := identity.FirstName
			toEmail := mail.NewEmail(name, credentials.NewEmail)
			templateData := map[string]interface{}{
				"name": name,
				"validation_link": provider.newEmailValidationLink.WithQuery(map[string]interface{}{
					"id":   claims.Payload.ID,
					"code": validationCode,
				}).String(),
			}

			if err := provider.mailer.Send(toEmail, provider.newEmailValidationTemplate, templateData); err != nil {
				return fmt.Errorf("failed to send new email validation link to user %q: %w", credentials.NewEmail, err)
			}

			return nil
		}, nil
}

func (provider *providerImpl) CancelNewEmail(ctx context.Context, token string) error {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return err
	}

	if _, err = provider.credentialsService.CancelNewEmail(ctx, claims.Payload.ID, now); err != nil {
		return fmt.Errorf("failed to cancel email update for user %q: %w", claims.Payload.ID, err)
	}

	return nil
}

func (provider *providerImpl) ResetPassword(ctx context.Context, form models.UserPasswordResetForm) (environment.Deferred, error) {
	credentials, resetLink, err := provider.credentialsService.ResetPassword(ctx, form.Email, provider.time())
	if err != nil {
		return nil, fmt.Errorf("failed to reset password for user %q: %w", form.Email, err)
	}
	identity, err := provider.identityService.Read(ctx, credentials.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch identity for user %q: %w", credentials.ID, err)
	}

	return func() error {
		name := identity.FirstName
		toEmail := mail.NewEmail(name, credentials.Email)
		templateData := map[string]interface{}{
			"name": name,
			"validation_link": provider.passwordResetLink.WithQuery(map[string]interface{}{
				"id":   credentials.ID,
				"code": resetLink,
			}).String(),
		}

		if err := provider.mailer.Send(toEmail, provider.passwordResetTemplate, templateData); err != nil {
			return fmt.Errorf("failed to send password reset link to user %q: %w", credentials.Email, err)
		}

		return nil
	}, nil
}

func (provider *providerImpl) ValidateEmail(ctx context.Context, form models.UserValidateEmailForm) error {
	if _, err := provider.credentialsService.ValidateEmail(ctx, form.ID, form.Code, provider.time()); err != nil {
		return fmt.Errorf("failed to validate email for user %q: %w", form.ID, err)
	}

	return nil
}

func (provider *providerImpl) ValidateNewEmail(ctx context.Context, form models.UserValidateEmailForm) error {
	if _, err := provider.credentialsService.ValidateNewEmail(ctx, form.ID, form.Code, provider.time()); err != nil {
		return fmt.Errorf("failed to validate new email for user %q: %w", form.ID, err)
	}

	return nil
}

func (provider *providerImpl) ResendEmailValidation(ctx context.Context, token string) (environment.Deferred, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	credentials, validationLink, err := provider.credentialsService.UpdateEmailValidation(ctx, claims.Payload.ID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update email validation link for user %q: %w", claims.Payload.ID, err)
	}
	identity, err := provider.identityService.Read(ctx, claims.Payload.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch identity for user %q: %w", claims.Payload.ID, err)
	}

	return func() error {
		name := identity.FirstName
		toEmail := mail.NewEmail(name, credentials.Email)
		templateData := map[string]interface{}{
			"name": name,
			"validation_link": provider.emailValidationLink.WithQuery(map[string]interface{}{
				"id":   claims.Payload.ID,
				"code": validationLink,
			}).String(),
		}

		if err := provider.mailer.Send(toEmail, provider.emailValidationTemplate, templateData); err != nil {
			return fmt.Errorf("failed to resend email validation link to user %q: %w", credentials.Email, err)
		}

		return nil
	}, nil
}

func (provider *providerImpl) ResendNewEmailValidation(ctx context.Context, token string) (environment.Deferred, error) {
	now := provider.time()
	claims, err := authentication.ForceAuthentication(token, provider.tokenService, provider.keysService, now)
	if err != nil {
		return nil, err
	}

	credentials, validationLink, err := provider.credentialsService.UpdateNewEmailValidation(ctx, claims.Payload.ID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update new email validation link for user %q: %w", claims.Payload.ID, err)
	}
	identity, err := provider.identityService.Read(ctx, claims.Payload.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch identity for user %q: %w", claims.Payload.ID, err)
	}

	return func() error {
		name := identity.FirstName
		toEmail := mail.NewEmail(name, credentials.NewEmail)
		templateData := map[string]interface{}{
			"name": name,
			"validation_link": provider.newEmailValidationLink.WithQuery(map[string]interface{}{
				"id":   claims.Payload.ID,
				"code": validationLink,
			}).String(),
		}

		if err := provider.mailer.Send(toEmail, provider.newEmailValidationTemplate, templateData); err != nil {
			return fmt.Errorf("failed to resend new email validation link to user %q: %w", credentials.NewEmail, err)
		}

		return nil
	}, nil
}

func (provider *providerImpl) DoesSlugExist(ctx context.Context, slug string) (bool, error) {
	ok, err := provider.profileService.SlugExists(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("failed to check if slug %q exists: %w", slug, err)
	}

	return ok, nil
}

func (provider *providerImpl) DoesEmailExist(ctx context.Context, email string) (bool, error) {
	ok, err := provider.credentialsService.EmailExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if email %q exists: %w", email, err)
	}

	return ok, nil
}
