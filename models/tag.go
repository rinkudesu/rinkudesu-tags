package models

import "github.com/gofrs/uuid"

type Tag struct {
	Id   uuid.UUID `json:"id" binding:"required"`
	Name string    `json:"name" binding:"required,max=50"`
}

type TagCreateViewModel struct {
	Name string `json:"name" binding:"required,max=50"`
}

func (vm *TagCreateViewModel) GetTag() Tag {
	return Tag{
		Name: vm.Name,
	}
}
