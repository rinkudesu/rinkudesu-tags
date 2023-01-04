package services

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/data"
	"rinkudesu-tags/data/migrations"
	"strings"
)

type Healthcheck struct {
	db       data.DbConnector
	executor *migrations.Executor
}

type HealthcheckResult int

const (
	Unhealthy HealthcheckResult = iota
	Degraded
	Healthy
)

func NewHealthcheck(state *GlobalState) *Healthcheck {
	return &Healthcheck{db: state.DbConnection, executor: migrations.NewExecutor(state.DbConnection)}
}

func (healthcheck *Healthcheck) GetStatus() HealthcheckResult {
	result, err := healthcheck.executor.IsMigrated()

	if err != nil {
		log.Warningf("Failed database healthcheck: %s", err.Error())
		return Unhealthy
	}

	if !result {
		log.Warning("Database is not in the latest format, migration is necessary")
		return Degraded
	}

	return Healthy
}

func GetHealthcheck(healthcheck *Healthcheck) gin.HandlerFunc {
	return func(context *gin.Context) {
		if strings.EqualFold(context.Request.URL.Path, "/health") || strings.EqualFold(context.Request.URL.Path, "/health/") {
			result := healthcheck.GetStatus()
			if result == Unhealthy {
				context.AbortWithStatus(http.StatusServiceUnavailable)
			} else {
				context.AbortWithStatus(http.StatusOK)
			}
		}
	}
}
