package models

import (
	"github.com/google/uuid"
	"time"
)

type UserTokenHeader struct {
	IAT time.Time `json:"iat"`
	EXP time.Time `json:"exp"`
	ID  uuid.UUID `json:"id"`
}

type UserTokenPayload struct {
	ID uuid.UUID `json:"id"`
}

type UserToken struct {
	Header  UserTokenHeader  `json:"header"`
	Payload UserTokenPayload `json:"token"`
}
