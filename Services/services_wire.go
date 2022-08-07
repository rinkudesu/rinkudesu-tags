//go:build wireinject
// +build wireinject

package Services

import "github.com/google/wire"

var (
	HealthcheckSet = wire.NewSet(NewHealthcheck)
)

func CreateHealthcheck(state *GlobalState) *Healthcheck {
	wire.Build(HealthcheckSet)
	return nil
}
