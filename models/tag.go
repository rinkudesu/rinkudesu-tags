package models

import "github.com/gofrs/uuid"

type Tag struct {
	Id   uuid.UUID `json:"id" binding:"required"`
	Name string    `json:"name" binding:"required,max=50"`
	// Colour represents a visual tag identifier. Should be provided in HTML format, as in "#aabbcc".
	Colour string `json:"colour" binding:"required,max=7"`
}

type TagCreateViewModel struct {
	Name   string `json:"name" binding:"required,max=50"`
	Colour string `json:"colour" binding:"required,max=7"`
}

func (vm *TagCreateViewModel) GetTag() Tag {
	return Tag{
		Name:   vm.Name,
		Colour: vm.Colour,
	}
}
