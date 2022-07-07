package main

import (
	"log"
	"net/http"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
	"rinkudesu-tags/Routers"
)

//todo: base path and port should be configurable
const basePath = "/api"

func main() {
	var connection = Data.DbConnection{}
	_ = connection.Initialise("postgres://postgres:postgres@localhost:5432/postgres")
	defer connection.Close()
	migrate(connection)

	setupRoutes(connection)

	err := http.ListenAndServe(":5000", nil) //todo: this nil should probably be *something*
	if err != nil {
		log.Panicln("Unable to listen on port 5000")
	}
}

func migrate(connection Data.DbConnection) {
	migrator := Migrations.NewExecutor(connection)
	migrator.Migrate()
}

func setupRoutes(connection Data.DbConnector) {
	Routers.SetupTagsRoutes(basePath)
	Routers.SetupTagsDatabase(connection)
}
