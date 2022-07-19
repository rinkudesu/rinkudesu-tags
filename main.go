package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"rinkudesu-tags/Controllers"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
	"rinkudesu-tags/Repositories"
	"rinkudesu-tags/Routers"
)

//todo: base path and port should be configurable
const basePath = "/api"

var (
	tagsRouter  *Routers.TagsRouter
	linksRouter *Routers.LinksRouter
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) //todo: this should be configurable
}

func main() {
	var connection = Data.DbConnection{}
	_ = connection.Initialise("postgres://postgres:postgres@localhost:5432/postgres")
	defer connection.Close()
	migrate(&connection)

	setupRoutes(&connection)

	err := http.ListenAndServe(":5000", nil) //todo: this nil should probably be *something*
	if err != nil {
		log.Panic("Unable to listen on port 5000")
	}
}

func migrate(connection Data.DbConnector) {
	migrator := Migrations.NewExecutor(connection)
	migrator.Migrate()
}

func setupRoutes(connection Data.DbConnector) {
	tagsRouter = Routers.NewTagsRouter(connection, basePath)
	linksRouter = Routers.NewLinksRouter(Controllers.NewLinksController(Repositories.NewLinksRepository(&connection)), basePath)
}
