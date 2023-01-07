package models

import "github.com/gofrs/uuid"

type LinkTag struct {
	Id     uuid.UUID `json:"id"`
	LinkId uuid.UUID `json:"linkId" binding:"required"`
	TagId  uuid.UUID `json:"tagId" binding:"required"`
}
