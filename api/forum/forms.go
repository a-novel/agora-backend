package forumapi

import (
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
)

type ReadImproveRequestForm struct {
	PostID uuid.UUID `json:"postID"`
}

type CreateImproveRequestForm struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateImproveRequestForm struct {
	SourceID uuid.UUID `json:"sourceID"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
}

type DeleteImproveRequestForm struct {
	PostID uuid.UUID `json:"postID"`
}

type SearchImproveRequestForm struct {
	UserID *uuid.UUID                        `json:"userID"`
	Query  string                            `json:"query"`
	Limit  int                               `json:"limit"`
	Offset int                               `json:"offset"`
	Order  *models.ImproveRequestSearchOrder `json:"order"`
}

type PreviewImproveRequestsForm struct {
	IDs []uuid.UUID `json:"ids"`
}

type ReadImproveSuggestionForm struct {
	PostID uuid.UUID `json:"postID"`
}

type CreateImproveSuggestionForm struct {
	RequestID uuid.UUID `json:"requestID"`
	SourceID  uuid.UUID `json:"sourceID"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

type UpdateImproveSuggestionForm struct {
	PostID    uuid.UUID `json:"postID"`
	RequestID uuid.UUID `json:"requestID"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

type DeleteImproveSuggestionForm struct {
	PostID uuid.UUID `json:"postID"`
}

type SearchImproveSuggestionForm struct {
	UserID    *uuid.UUID                           `json:"userID"`
	SourceID  *uuid.UUID                           `json:"sourceID"`
	RequestID *uuid.UUID                           `json:"requestID"`
	Validated *bool                                `json:"validated"`
	Limit     int                                  `json:"limit"`
	Offset    int                                  `json:"offset"`
	Order     *models.ImproveSuggestionSearchOrder `json:"order"`
}

type PreviewImproveSuggestionsForm struct {
	IDs []uuid.UUID `json:"ids"`
}

type VoteForm struct {
	PostID uuid.UUID         `json:"postID"`
	Target models.VoteTarget `json:"target"`
	Vote   models.VoteValue  `json:"vote"`
}

type ReadVoteForm struct {
	PostID uuid.UUID         `json:"postID"`
	Target models.VoteTarget `json:"target"`
}

type SearchVotesForm struct {
	UserID uuid.UUID         `json:"userID"`
	Target models.VoteTarget `json:"target"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}
