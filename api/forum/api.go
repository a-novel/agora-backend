package forumapi

import (
	"github.com/a-novel/agora-backend/api"
	"github.com/a-novel/agora-backend/environment/forum/improve_post"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ImproveRequestAPI(basePath string, r gin.IRouter, provider improve_post.Provider) {
	api.LoadAPI(r, basePath, api.Config{
		"/": {
			http.MethodPost: api.WithContext[ReadImproveRequestForm, improve_post.Provider](improveRequestReadAPI, provider),
		},
		"/edit": {
			http.MethodPost:   api.WithContext[CreateImproveRequestForm, improve_post.Provider](improveRequestCreateAPI, provider),
			http.MethodPut:    api.WithContext[UpdateImproveRequestForm, improve_post.Provider](improveRequestUpdateAPI, provider),
			http.MethodDelete: api.WithContext[DeleteImproveRequestForm, improve_post.Provider](improveRequestDeleteAPI, provider),
		},
		"/search": {
			http.MethodPost: api.WithContext[SearchImproveRequestForm, improve_post.Provider](improveRequestSearchAPI, provider),
		},
		"/previews": {
			http.MethodPost: api.WithContext[PreviewImproveRequestsForm, improve_post.Provider](improveRequestPreviewsAPI, provider),
		},
	})
}

func ImproveSuggestionAPI(basePath string, r gin.IRouter, provider improve_post.Provider) {
	api.LoadAPI(r, basePath, api.Config{
		"/": {
			http.MethodPost: api.WithContext[ReadImproveSuggestionForm, improve_post.Provider](improveSuggestionReadAPI, provider),
		},
		"/edit": {
			http.MethodPost:   api.WithContext[CreateImproveSuggestionForm, improve_post.Provider](improveSuggestionCreateAPI, provider),
			http.MethodPut:    api.WithContext[UpdateImproveSuggestionForm, improve_post.Provider](improveSuggestionUpdateAPI, provider),
			http.MethodDelete: api.WithContext[DeleteImproveSuggestionForm, improve_post.Provider](improveSuggestionDeleteAPI, provider),
		},
		"/search": {
			http.MethodPost: api.WithContext[SearchImproveSuggestionForm, improve_post.Provider](improveSuggestionSearchAPI, provider),
		},
		"/previews": {
			http.MethodPost: api.WithContext[PreviewImproveSuggestionsForm, improve_post.Provider](improveSuggestionPreviewsAPI, provider),
		},
	})
}

func VotesAPI(basePath string, r gin.IRouter, provider improve_post.Provider) {
	api.LoadAPI(r, basePath, api.Config{
		"/": {
			http.MethodPost: api.WithContext[VoteForm, improve_post.Provider](voteUpdateAPI, provider),
		},
		"/status": {
			http.MethodPost: api.WithContext[ReadVoteForm, improve_post.Provider](voteReadAPI, provider),
		},
		"/search": {
			http.MethodPost: api.WithContext[SearchVotesForm, improve_post.Provider](voteSearchAPI, provider),
		},
	})
}
