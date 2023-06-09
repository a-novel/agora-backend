package bookmarkapi

import (
	"github.com/a-novel/agora-backend/api"
	"github.com/a-novel/agora-backend/environment/bookmark/improve_post"
	"github.com/gin-gonic/gin"
)

func improvePostCreateAPI(c *gin.Context, token string, form CreateImprovePostForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.Bookmark(c, token, form.RequestID, form.Target, form.Level)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func improvePostDeleteAPI(c *gin.Context, token string, form DeleteImprovePostForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	return api.CallbackResponse{}, provider.UnBookmark(c, token, form.RequestID, form.Target)
}

func improvePostReadAPI(c *gin.Context, _ string, form ReadImprovePostForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.IsBookmarked(c, form.UserID, form.RequestID, form.Target)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, err
}

func improvePostReadSearchAPI(c *gin.Context, _ string, form SearchImprovePostForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, total, err := provider.List(c, form.UserID, form.Level, form.Target, form.Limit, form.Offset)

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
