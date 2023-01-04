//go:build wireinject
// +build wireinject

package message_handlers

import (
	"github.com/google/wire"
	"rinkudesu-tags/repositories"
	"rinkudesu-tags/services"
)

var (
	UserDeletedHandlerSet = wire.NewSet(NewUserDeletedHandler, repositories.AllRepositoriesSet)
)

func CreateUserDeletedHandler(state *services.GlobalState) *UserDeletedHandler {
	wire.Build(UserDeletedHandlerSet)
	return nil
}
