package userapi

import (
	"github.com/a-novel/agora-backend/api"
	"github.com/a-novel/agora-backend/environment/user/account"
	"github.com/a-novel/agora-backend/environment/user/authentication"
	"github.com/a-novel/agora-backend/environment/user/profile"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AccountAPI(basePath string, r gin.IRouter, provider account.Provider) {
	api.LoadAPI(r, basePath, api.Config{
		"/": {
			http.MethodPost: api.WithContext[RegisterForm, account.Provider](accountRegisterAPI, provider),
		},
		"/info": {
			http.MethodGet: api.WithContext[any, account.Provider](accountInfoAPI, provider),
		},
		"/preview": {
			http.MethodGet: api.WithContext[any, account.Provider](accountPreviewAPI, provider),
		},
		"/authorizations": {
			http.MethodGet: api.WithContext[any, account.Provider](accountAuthorizationsAPI, provider),
		},
		"/identity": {
			http.MethodPut: api.WithContext[IdentityUpdateForm, account.Provider](accountIdentityUpdateAPI, provider),
		},
		"/profile": {
			http.MethodPut: api.WithContext[ProfileUpdateForm, account.Provider](accountProfileUpdateAPI, provider),
		},
		"/credentials/password": {
			http.MethodPatch:  api.WithContext[PasswordUpdateForm, account.Provider](accountPasswordUpdateAPI, provider),
			http.MethodDelete: api.WithContext[PasswordResetForm, account.Provider](accountPasswordResetAPI, provider),
		},
		"/credentials/email": {
			http.MethodGet:    api.WithContext[any, account.Provider](accountEmailValidationStatusAPI, provider),
			http.MethodPatch:  api.WithContext[EmailUpdateForm, account.Provider](accountEmailUpdateAPI, provider),
			http.MethodDelete: api.WithContext[any, account.Provider](accountEmailCancelUpdateAPI, provider),
		},
		"/credentials/email/validation": {
			http.MethodPost: api.WithContext[ValidateEmailForm, account.Provider](accountValidateEmailAPI, provider),
			http.MethodGet:  api.WithContext[any, account.Provider](accountResendEmailValidationAPI, provider),
		},
		"/credentials/new-email/validation": {
			http.MethodPost: api.WithContext[ValidateEmailForm, account.Provider](accountValidateNewEmailAPI, provider),
			http.MethodGet:  api.WithContext[any, account.Provider](accountResendNewEmailValidationAPI, provider),
		},
		"/credentials/email/exists": {
			http.MethodPost: api.WithContext[EmailExistsForm, account.Provider](accountEmailExistsAPI, provider),
		},
		"/profile/slug/exists": {
			http.MethodPost: api.WithContext[SlugExistsForm, account.Provider](accountSlugExistsAPI, provider),
		},
	})
}

func AuthenticationAPI(basePath string, r gin.IRouter, provider authentication.Provider) {
	api.LoadAPI(r, basePath, api.Config{
		"/": {
			http.MethodGet:  api.WithContext[any, authentication.Provider](authenticationStatusAPI, provider),
			http.MethodPost: api.WithContext[LoginForm, authentication.Provider](authenticationLoginAPI, provider),
		},
	})
}

func ProfileAPI(basePath string, r gin.IRouter, provider profile.Provider) {
	api.LoadAPI(r, basePath, api.Config{
		"/read/:slug": {
			http.MethodGet: api.WithContext[ReadProfileForm, profile.Provider](profileReadAPI, provider),
		},
		"/search": {
			http.MethodPost: api.WithContext[SearchProfileForm, profile.Provider](profileSearchAPI, provider),
		},
		"/previews": {
			http.MethodPost: api.WithContext[PreviewProfilesForm, profile.Provider](profilePreviewsAPI, provider),
		},
	})
}
