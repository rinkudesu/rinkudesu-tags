//go:build wireinject
// +build wireinject

package Controllers

import (
	"github.com/google/wire"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Repositories"
)

var (
	LinksControllerSet    = wire.NewSet(NewLinksController, Repositories.LinkRepositorySet)
	TagsControllerSet     = wire.NewSet(NewTagsController, Repositories.TagsRepositorySet)
	LinkTagsControllerSet = wire.NewSet(NewLinkTagsController, Repositories.LinkTagsRepositorySet)
)

func CreateLinksController(dbConnection Data.DbConnector) *LinksController {
	wire.Build(LinksControllerSet)
	return nil
}

func CreateTagsController(dbConnection Data.DbConnector) *TagsController {
	wire.Build(TagsControllerSet)
	return nil
}

func CreateLinkTagsController(dbConnection Data.DbConnector) *LinkTagsController {
	wire.Build(LinkTagsControllerSet)
	return nil
}
