//go:build wireinject
// +build wireinject

package Repositories

import (
	"github.com/google/wire"
	"rinkudesu-tags/Services"
)

var (
	LinkRepositorySet     = wire.NewSet(NewLinksRepository)
	TagsRepositorySet     = wire.NewSet(NewTagsRepository)
	LinkTagsRepositorySet = wire.NewSet(NewLinkTagsRepository)
)

func CreateLinksRepository(state *Services.GlobalState) *LinksRepository {
	wire.Build(LinkRepositorySet)
	return nil
}

func CreateTagsRepository(state *Services.GlobalState) *TagsRepository {
	wire.Build(TagsRepositorySet)
	return nil
}

func CreateLinkTagsRepository(state *Services.GlobalState) *LinkTagsRepository {
	wire.Build(LinkTagsRepositorySet)
	return nil
}
