﻿package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
	"rinkudesu-tags/Routers"
)

//todo: base path and port should be configurable
const basePath = "/api"

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) //todo: this should be configurable
}

func main() {
	var connection = Data.DbConnection{}
	_ = connection.Initialise("postgres://postgres:postgres@localhost:5432/postgres")
	defer connection.Close()
	migrate(connection)

	setupRoutes(&connection)

	err := http.ListenAndServe(":5000", nil) //todo: this nil should probably be *something*
	if err != nil {
		log.Panic("Unable to listen on port 5000")
	}
}

func migrate(connection Data.DbConnection) {
	migrator := Migrations.NewExecutor(&connection)
	migrator.Migrate()
}

func setupRoutes(connection Data.DbConnector) {
	Routers.SetupTagsRoutes(basePath)
	Routers.SetupTagsDatabase(connection)
}
