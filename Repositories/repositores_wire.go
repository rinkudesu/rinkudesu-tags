//go:build wireinject
// +build wireinject

package Repositories

import (
	"github.com/google/wire"
	"rinkudesu-tags/Data"
)

var (
	LinkRepositorySet     = wire.NewSet(NewLinksRepository)
	TagsRepositorySet     = wire.NewSet(NewTagsRepository, NewTagQueryExecutor)
	LinkTagsRepositorySet = wire.NewSet(NewLinkTagsRepository)
)

func CreateLinksRepository(dbConnection Data.DbConnector) *LinksRepository {
	wire.Build(LinkRepositorySet)
	return nil
}

func CreateTagsRepository(dbConnection Data.DbConnector) *TagsRepository {
	wire.Build(TagsRepositorySet)
	return nil
}

func CreateLinkTagsRepository(dbConnection Data.DbConnector) *LinkTagsRepository {
	wire.Build(LinkTagsRepositorySet)
	return nil
}
