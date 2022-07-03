package Models

import "github.com/google/uuid"

type LinkTag struct {
	Id     uuid.UUID `json:"id"`
	LinkId uuid.UUID `json:"linkId"`
	TagId  uuid.UUID `json:"tagId"`
}
