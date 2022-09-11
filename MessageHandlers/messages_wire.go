//go:build wireinject
// +build wireinject

package MessageHandlers

import (
	"github.com/google/wire"
	"rinkudesu-tags/Repositories"
	"rinkudesu-tags/Services"
)

var (
	UserDeletedHandlerSet = wire.NewSet(NewUserDeletedHandler, Repositories.AllRepositoriesSet)
)

func CreateUserDeletedHandler(state *Services.GlobalState) *UserDeletedHandler {
	wire.Build(UserDeletedHandlerSet)
	return nil
}
