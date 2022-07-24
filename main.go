package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"rinkudesu-tags/Controllers"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
	"rinkudesu-tags/Services"
)

var (
	routables []Controllers.Routable
	router    *gin.Engine
	config    *Configuration
)

func init() {
	config = NewConfiguration()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel)
}

func main() {
	var connection = Data.DbConnection{}
	err := connection.Initialise(config.DbConnection)
	if err != nil {
		log.Panicf("Failed to initialise database connection: %s", err.Error())
	}
	defer connection.Close()
	migrate(&connection)

	createControllers(&connection)
	setupRouter()
}

func migrate(connection Data.DbConnector) {
	migrator := Migrations.NewExecutor(connection)
	migrator.Migrate()
}

func createControllers(connection Data.DbConnector) {
	routables = make([]Controllers.Routable, 3)
	routables[0] = Controllers.CreateLinksController(connection)
	routables[1] = Controllers.CreateTagsController(connection)
	routables[2] = Controllers.CreateLinkTagsController(connection)
}

func setupRouter() {
	router = gin.New()
	router.Use(gin.Recovery())
	router.Use(Services.GetGinLogger())
	err := router.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		log.Panicf("Failed to set trusted proxies: %s", err.Error())
	}
	//todo: GIN_MODE=release

	for _, routable := range routables {
		routable.SetupRoutes(router, config.BasePath)
	}

	if err := router.Run(config.ListenAddress); err != nil {
		log.Panicf("Server failed while listening: %s", err.Error())
	}
}
