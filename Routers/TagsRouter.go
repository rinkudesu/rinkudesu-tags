package Routers

import (
	"fmt"
	"net/http"
	"rinkudesu-tags/Controllers"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Repositories"
)

const path = "tags"

var (
	database Data.DbConnection
)

func handleTags(w http.ResponseWriter, r *http.Request) {
	controller := getController()
	controller.HandleProducts(w, r)
}

func SetupTagsRoutes(basePath string) {
	tagHandler := http.HandlerFunc(handleTags)
	http.Handle(fmt.Sprintf("%s/v1/%s", basePath, path), tagHandler)
}

func SetupTagsDatabase(initDatabase Data.DbConnection) {
	database = initDatabase
}

//todo: this is so bad...
func getController() Controllers.TagsController {
	var repository = Repositories.TagsRepository{}
	repository.Init(database)
	var controller = Controllers.TagsController{}
	controller.Init(repository)
	return controller
}
