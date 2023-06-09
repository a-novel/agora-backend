package bookmarkapi

import (
	"github.com/a-novel/agora-backend/api"
	"github.com/a-novel/agora-backend/environment/bookmark/improve_post"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ImprovePostAPI(basePath string, r gin.IRouter, provider improve_post.Provider) {
	api.LoadAPI(r, basePath, api.Config{
		"/": {
			http.MethodPost:   api.WithContext[CreateImprovePostForm, improve_post.Provider](improvePostCreateAPI, provider),
			http.MethodDelete: api.WithContext[DeleteImprovePostForm, improve_post.Provider](improvePostDeleteAPI, provider),
		},
		"/status": {
			http.MethodPost: api.WithContext[ReadImprovePostForm, improve_post.Provider](improvePostReadAPI, provider),
		},
		"/search": {
			http.MethodPost: api.WithContext[SearchImprovePostForm, improve_post.Provider](improvePostReadSearchAPI, provider),
		},
	})
}
