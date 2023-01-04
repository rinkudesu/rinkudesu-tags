//go:build wireinject
// +build wireinject

package controllers

import (
	"github.com/google/wire"
	"rinkudesu-tags/repositories"
	"rinkudesu-tags/services"
)

var (
	LinksControllerSet    = wire.NewSet(NewLinksController, repositories.LinkRepositorySet)
	TagsControllerSet     = wire.NewSet(NewTagsController, repositories.TagsRepositorySet)
	LinkTagsControllerSet = wire.NewSet(NewLinkTagsController, repositories.LinkTagsRepositorySet, repositories.LinkRepositorySet, repositories.TagsRepositorySet)
)

func CreateLinksController(state *services.GlobalState) *LinksController {
	wire.Build(LinksControllerSet)
	return nil
}

func CreateTagsController(state *services.GlobalState) *TagsController {
	wire.Build(TagsControllerSet)
	return nil
}

func CreateLinkTagsController(state *services.GlobalState) *LinkTagsController {
	wire.Build(LinkTagsControllerSet)
	return nil
}
