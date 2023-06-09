package forumapi

import (
	"github.com/a-novel/agora-backend/api"
	"github.com/a-novel/agora-backend/environment/forum/improve_post"
	"github.com/a-novel/agora-backend/models"
	"github.com/gin-gonic/gin"
)

func improveRequestReadAPI(c *gin.Context, _ string, form ReadImproveRequestForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	revisions, err := provider.ReadImproveRequest(c, form.PostID)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": revisions,
		},
	}, nil
}

func improveRequestCreateAPI(c *gin.Context, token string, form CreateImproveRequestForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.CreateImproveRequest(c, token, form.Title, form.Content)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func improveRequestUpdateAPI(c *gin.Context, token string, form UpdateImproveRequestForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.CreateImproveRequestRevision(c, token, form.SourceID, form.Title, form.Content)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func improveRequestDeleteAPI(c *gin.Context, token string, form DeleteImproveRequestForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	return api.CallbackResponse{}, provider.DeleteImproveRequest(c, token, form.PostID)
}

func improveRequestSearchAPI(c *gin.Context, _ string, form SearchImproveRequestForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, total, err := provider.SearchImproveRequests(c, models.ImproveRequestSearch{
		UserID: form.UserID,
		Query:  form.Query,
		Order:  form.Order,
	}, form.Limit, form.Offset)

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

func improveRequestPreviewsAPI(c *gin.Context, _ string, form PreviewImproveRequestsForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.GetImproveRequestPreviews(c, form.IDs)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func improveSuggestionReadAPI(c *gin.Context, _ string, form ReadImproveSuggestionForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.ReadImproveSuggestion(c, form.PostID)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func improveSuggestionCreateAPI(c *gin.Context, token string, form CreateImproveSuggestionForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.CreateImproveSuggestion(c, token, form.RequestID, form.SourceID, form.Title, form.Content)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func improveSuggestionUpdateAPI(c *gin.Context, token string, form UpdateImproveSuggestionForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.UpdateImproveSuggestion(c, token, form.PostID, form.RequestID, form.Title, form.Content)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func improveSuggestionDeleteAPI(c *gin.Context, token string, form DeleteImproveSuggestionForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	return api.CallbackResponse{}, provider.DeleteImproveSuggestion(c, token, form.PostID)
}

func improveSuggestionSearchAPI(c *gin.Context, _ string, form SearchImproveSuggestionForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, total, err := provider.ListImproveSuggestions(c, models.ImproveSuggestionsList{
		UserID:    form.UserID,
		SourceID:  form.SourceID,
		RequestID: form.RequestID,
		Validated: form.Validated,
		Order:     form.Order,
	}, form.Limit, form.Offset)

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

func improveSuggestionPreviewsAPI(c *gin.Context, _ string, form PreviewImproveSuggestionsForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.GetImproveSuggestionPreviews(c, form.IDs)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func voteUpdateAPI(c *gin.Context, token string, form VoteForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.Vote(c, token, form.PostID, form.Target, form.Vote)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func voteReadAPI(c *gin.Context, token string, form ReadVoteForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, err := provider.HasVoted(c, token, form.PostID, form.Target)

	if err != nil {
		return api.CallbackResponse{}, err
	}

	return api.CallbackResponse{
		Body: map[string]interface{}{
			"data": res,
		},
	}, nil
}

func voteSearchAPI(c *gin.Context, _ string, form SearchVotesForm, provider improve_post.Provider) (api.CallbackResponse, error) {
	res, total, err := provider.GetVotedPosts(c, form.UserID, form.Target, form.Limit, form.Offset)

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
