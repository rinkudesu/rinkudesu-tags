package Models

import "github.com/gofrs/uuid"

type LinkTag struct {
	Id     uuid.UUID `json:"id"`
	LinkId uuid.UUID `json:"linkId"`
	TagId  uuid.UUID `json:"tagId"`
}
