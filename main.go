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

//todo: base path and port should be configurable
const basePath = "/api"

var (
	routables []Controllers.Routable

	router *gin.Engine
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) //todo: this should be configurable
}

func main() {
	var connection = Data.DbConnection{}
	_ = connection.Initialise("postgres://postgres:postgres@localhost:5432/postgres") //todo: this should be configurable
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
	_ = router.SetTrustedProxies(nil) //todo: this should read from config (and probably handle error then as well)
	//todo: GIN_MODE=release

	for _, routable := range routables {
		routable.SetupRoutes(router, basePath)
	}

	if err := router.Run("localhost:5000"); err != nil { //todo: make url configurable
		log.Panicf("Server failed while listening: %s", err.Error())
	}
}
