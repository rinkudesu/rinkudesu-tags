//go:build wireinject
// +build wireinject

package Controllers

import (
	"github.com/google/wire"
	"rinkudesu-tags/Repositories"
	"rinkudesu-tags/Services"
)

var (
	LinksControllerSet    = wire.NewSet(NewLinksController, Repositories.LinkRepositorySet)
	TagsControllerSet     = wire.NewSet(NewTagsController, Repositories.TagsRepositorySet)
	LinkTagsControllerSet = wire.NewSet(NewLinkTagsController, Repositories.LinkTagsRepositorySet)
)

func CreateLinksController(state *Services.GlobalState) *LinksController {
	wire.Build(LinksControllerSet)
	return nil
}

func CreateTagsController(state *Services.GlobalState) *TagsController {
	wire.Build(TagsControllerSet)
	return nil
}

func CreateLinkTagsController(state *Services.GlobalState) *LinkTagsController {
	wire.Build(LinkTagsControllerSet)
	return nil
}
