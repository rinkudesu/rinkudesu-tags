//go:build wireinject
// +build wireinject

package repositories

import (
	"github.com/google/wire"
	"rinkudesu-tags/services"
)

var (
	LinkRepositorySet     = wire.NewSet(NewLinksRepository)
	TagsRepositorySet     = wire.NewSet(NewTagsRepository)
	LinkTagsRepositorySet = wire.NewSet(NewLinkTagsRepository)

	AllRepositoriesSet = wire.NewSet(LinkRepositorySet, TagsRepositorySet, LinkTagsRepositorySet)
)

func CreateLinksRepository(state *services.GlobalState) *LinksRepository {
	wire.Build(LinkRepositorySet)
	return nil
}

func CreateTagsRepository(state *services.GlobalState) *TagsRepository {
	wire.Build(TagsRepositorySet)
	return nil
}

func CreateLinkTagsRepository(state *services.GlobalState) *LinkTagsRepository {
	wire.Build(LinkTagsRepositorySet)
	return nil
}
