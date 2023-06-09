package bookmarkapi

import (
	"github.com/a-novel/agora-backend/models"
	"github.com/google/uuid"
)

type CreateImprovePostForm struct {
	RequestID uuid.UUID             `json:"requestID"`
	Target    models.BookmarkTarget `json:"target"`
	Level     models.BookmarkLevel  `json:"level"`
}

type DeleteImprovePostForm struct {
	RequestID uuid.UUID             `json:"requestID"`
	Target    models.BookmarkTarget `json:"target"`
}

type ReadImprovePostForm struct {
	UserID    uuid.UUID             `json:"userID"`
	RequestID uuid.UUID             `json:"requestID"`
	Target    models.BookmarkTarget `json:"target"`
}

type SearchImprovePostForm struct {
	UserID uuid.UUID             `json:"userID"`
	Target models.BookmarkTarget `json:"target"`
	Level  models.BookmarkLevel  `json:"level"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}
