// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package Controllers

import (
	"github.com/google/wire"
	"rinkudesu-tags/Repositories"
	"rinkudesu-tags/Services"
)

// Injectors from controllers_wire.go:

func CreateLinksController(state *Services.GlobalState) *LinksController {
	linksRepository := Repositories.NewLinksRepository(state)
	linksController := NewLinksController(linksRepository)
	return linksController
}

func CreateTagsController(state *Services.GlobalState) *TagsController {
	tagQueryExecutable := Repositories.NewTagQueryExecutor(state)
	tagsRepository := Repositories.NewTagsRepository(tagQueryExecutable)
	tagsController := NewTagsController(tagsRepository)
	return tagsController
}

func CreateLinkTagsController(state *Services.GlobalState) *LinkTagsController {
	linkTagsRepository := Repositories.NewLinkTagsRepository(state)
	linkTagsController := NewLinkTagsController(linkTagsRepository)
	return linkTagsController
}

// controllers_wire.go:

var (
	LinksControllerSet    = wire.NewSet(NewLinksController, Repositories.LinkRepositorySet)
	TagsControllerSet     = wire.NewSet(NewTagsController, Repositories.TagsRepositorySet)
	LinkTagsControllerSet = wire.NewSet(NewLinkTagsController, Repositories.LinkTagsRepositorySet)
)
