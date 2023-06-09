package secretsapi

import (
	"github.com/a-novel/agora-backend/api"
	"github.com/a-novel/agora-backend/environment/secrets"
	"github.com/a-novel/agora-backend/environment/user/authentication"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWKsAPI(basePath string, r gin.IRouter, provider secrets.Provider, allowedUsers []string) {
	api.LoadAPI(r, basePath, api.Config{
		"/": {
			http.MethodPost: func(c *gin.Context) {
				var auth *authentication.BackendServiceAuth

				if len(allowedUsers) > 0 {
					auth = &authentication.BackendServiceAuth{
						UserAgent:     c.Request.UserAgent(),
						Authorization: c.GetHeader("Authorization"),
						AllowedUsers:  allowedUsers,
					}
				}

				if err := provider.RotateJWKs(c, auth); err != nil {
					_ = c.AbortWithError(http.StatusInternalServerError, err)
				}

				c.AbortWithStatus(http.StatusNoContent)
			},
		},
	})
}
