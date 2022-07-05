package Models

import "github.com/gofrs/uuid"

type Tag struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	UserId uuid.UUID `json:"userId"`
}
