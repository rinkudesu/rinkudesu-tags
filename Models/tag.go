package Models

import "github.com/gofrs/uuid"

type Tag struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	UserId uuid.UUID `json:"userId"`
}

func (tag *Tag) IsValid() bool {
	if tag.Name == "" {
		return false
	}
	if tag.UserId == uuid.Nil {
		return false
	}
	return true
}
