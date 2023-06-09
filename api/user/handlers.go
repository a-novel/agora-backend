package userapi

import (
	"github.com/a-novel/agora-backend/api"
	"github.com/a-novel/agora-backend/environment/user/account"
	"github.com/a-novel/agora-backend/environment/user/authentication"
	"github.com/a-novel/agora-backend/environment/user/profile"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/a-novel/agora-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func accountRegisterAPI(c *gin.Context, token string, body RegisterForm, provider account.Provider) (api.CallbackResponse, error) {
	user, token, deferred, err := provider.Register(c, models.UserCreateForm{
		Credentials: models.UserCredentialsLoginForm{
			Email:    body.Email,
			Password: body.Password,
		},
		Identity: models.UserIdentityUpdateForm{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Birthday:  body.Birthday,
			Sex:       body.Sex,
		},
		Profile: models.UserProfileUpdateForm{
			Username: body.Username,
			Slug:     body.Slug,
		},
	})

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		// This information is used for the post-registration modal.
		Body: map[string]interface{}{
			"email":     user.Email,
			"username":  user.Username,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"sex":       user.Sex,
			"token":     token,
		},
		CTXData: map[string]interface{}{
			"action":     "register user",
			"email":      user.Email,
			"username":   user.Username,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
		Deferred: deferred,
	}, nil
}

func accountInfoAPI(c *gin.Context, token string, _ interface{}, provider account.Provider) (api.CallbackResponse, error) {
	info, err := provider.GetAccountInfo(c, token)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"id":        info.ID,
			"createdAt": info.CreatedAt,
			"updatedAt": info.UpdatedAt,
			"email":     info.Email,
			"newEmail":  info.NewEmail,
			"username":  info.Profile.Username,
			"slug":      info.Profile.Slug,
			"firstName": info.Identity.FirstName,
			"lastName":  info.Identity.LastName,
			"birthday":  info.Identity.Birthday,
			"sex":       info.Identity.Sex,
		},
	}, nil
}

func accountPreviewAPI(c *gin.Context, token string, _ interface{}, provider account.Provider) (api.CallbackResponse, error) {
	info, err := provider.GetAccountPreview(c, token)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: info,
	}, nil
}

func accountAuthorizationsAPI(c *gin.Context, token string, _ interface{}, provider account.Provider) (api.CallbackResponse, error) {
	authorizations, err := provider.GetAuthorizations(c, token)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: authorizations,
	}, nil
}

func accountEmailValidationStatusAPI(c *gin.Context, token string, _ interface{}, provider account.Provider) (api.CallbackResponse, error) {
	info, err := provider.GetEmailValidationStatus(c, token)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: info,
	}, nil
}

func accountIdentityUpdateAPI(c *gin.Context, token string, body IdentityUpdateForm, provider account.Provider) (api.CallbackResponse, error) {
	user, err := provider.UpdateIdentity(c, token, models.UserIdentityUpdateForm{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Birthday:  body.Birthday,
		Sex:       body.Sex,
	})

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: user,
		CTXData: map[string]interface{}{
			"action":     "update identity",
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}, nil
}

func accountProfileUpdateAPI(c *gin.Context, token string, body ProfileUpdateForm, provider account.Provider) (api.CallbackResponse, error) {
	user, err := provider.UpdateProfile(c, token, models.UserProfileUpdateForm{
		Username: body.Username,
		Slug:     body.Slug,
	})

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: user,
		CTXData: map[string]interface{}{
			"action": "update profile",
			"slug":   user.Slug,
		},
	}, nil
}

func accountPasswordUpdateAPI(c *gin.Context, _ string, body PasswordUpdateForm, provider account.Provider) (api.CallbackResponse, error) {
	err := provider.UpdatePassword(c, models.UserPasswordUpdateForm{
		ID:          body.ID,
		Password:    body.Password,
		OldPassword: body.OldPassword,
	})

	return api.CallbackResponse{
		MaskErrorsWithStatus: map[error]int{
			// If a user update a password (from its settings), not found can never be returned under
			// normal circumstances.
			// The only way to have a not found is from the post reset-password page, where user
			// uses a bad ID. Since the ID is generated from our backend, this should never happen.
			// To prevent hackers from knowing if the id is valid, we return forbidden under all
			// circumstances.
			validation.ErrNotFound: http.StatusForbidden,
		},
	}, err
}

func accountEmailUpdateAPI(c *gin.Context, token string, body EmailUpdateForm, provider account.Provider) (api.CallbackResponse, error) {
	status, deferred, err := provider.UpdateEmail(c, token, models.UserEmailUpdateForm{Email: body.Email})

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: status,
		CTXData: map[string]interface{}{
			"action":    "update email",
			"email":     status.Email,
			"new_email": status.NewEmail,
		},
		Deferred: deferred,
	}, nil
}

func accountEmailCancelUpdateAPI(c *gin.Context, token string, _ interface{}, provider account.Provider) (api.CallbackResponse, error) {
	err := provider.CancelNewEmail(c, token)

	return api.CallbackResponse{}, err
}

func accountPasswordResetAPI(c *gin.Context, _ string, body PasswordResetForm, provider account.Provider) (api.CallbackResponse, error) {
	deferred, err := provider.ResetPassword(c, models.UserPasswordResetForm{Email: body.Email})

	return api.CallbackResponse{
		Deferred: deferred,
	}, err
}

func accountValidateEmailAPI(c *gin.Context, _ string, body ValidateEmailForm, provider account.Provider) (api.CallbackResponse, error) {
	err := provider.ValidateEmail(c, models.UserValidateEmailForm{
		ID:   body.ID,
		Code: body.Code,
	})

	return api.CallbackResponse{
		MaskErrorsWithStatus: map[error]int{
			validation.ErrInvalidEntity:      http.StatusBadRequest,
			validation.ErrInvalidCredentials: http.StatusBadRequest,
			validation.ErrNotFound:           http.StatusBadRequest,
			validation.ErrValidated:          http.StatusBadRequest,
		},
	}, err
}

func accountValidateNewEmailAPI(c *gin.Context, _ string, body ValidateEmailForm, provider account.Provider) (api.CallbackResponse, error) {
	err := provider.ValidateNewEmail(c, models.UserValidateEmailForm{
		ID:   body.ID,
		Code: body.Code,
	})

	return api.CallbackResponse{
		MaskErrorsWithStatus: map[error]int{
			validation.ErrInvalidEntity:      http.StatusBadRequest,
			validation.ErrInvalidCredentials: http.StatusBadRequest,
			validation.ErrNotFound:           http.StatusBadRequest,
			validation.ErrValidated:          http.StatusBadRequest,
		},
	}, err
}

func accountResendEmailValidationAPI(c *gin.Context, token string, _ interface{}, provider account.Provider) (api.CallbackResponse, error) {
	deferred, err := provider.ResendEmailValidation(c, token)

	return api.CallbackResponse{
		Deferred: deferred,
	}, err
}

func accountResendNewEmailValidationAPI(c *gin.Context, token string, _ interface{}, provider account.Provider) (api.CallbackResponse, error) {
	deferred, err := provider.ResendNewEmailValidation(c, token)

	return api.CallbackResponse{
		Deferred: deferred,
	}, err
}

func accountEmailExistsAPI(c *gin.Context, _ string, body EmailExistsForm, provider account.Provider) (api.CallbackResponse, error) {
	exists, err := provider.DoesEmailExist(c, body.Email)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"exists": exists,
		},
	}, nil
}

func accountSlugExistsAPI(c *gin.Context, _ string, body SlugExistsForm, provider account.Provider) (api.CallbackResponse, error) {
	exists, err := provider.DoesSlugExist(c, body.Slug)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"exists": exists,
		},
	}, nil
}

func authenticationStatusAPI(c *gin.Context, token string, _ interface{}, provider authentication.Provider) (api.CallbackResponse, error) {
	res, err := provider.Authenticate(c, token, true)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"token": res,
		},
	}, nil
}

func authenticationLoginAPI(c *gin.Context, _ string, body LoginForm, provider authentication.Provider) (api.CallbackResponse, error) {
	token, err := provider.Login(c, models.UserCredentialsLoginForm{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"token": token,
		},
	}, nil
}

func profileReadAPI(c *gin.Context, _ string, body ReadProfileForm, provider profile.Provider) (api.CallbackResponse, error) {
	res, err := provider.Read(c, body.Slug)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: res,
	}, nil
}

func profileSearchAPI(c *gin.Context, _ string, body SearchProfileForm, provider profile.Provider) (api.CallbackResponse, error) {
	res, total, err := provider.Search(c, body.Query, body.Limit, body.Offset)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data":  res,
			"total": total,
		},
	}, nil
}

func profilePreviewsAPI(c *gin.Context, _ string, body PreviewProfilesForm, provider profile.Provider) (api.CallbackResponse, error) {
	previews, err := provider.Previews(c, body.IDs)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": previews,
		},
	}, nil
}
