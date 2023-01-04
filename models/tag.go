package models

import "github.com/gofrs/uuid"

type Tag struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (tag *Tag) IsValid() bool {
	if tag.Name == "" {
		return false
	}
	return true
}
