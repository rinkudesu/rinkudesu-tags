package Routers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/Controllers"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Repositories"
)

const path = "tags"

var (
	database Data.DbConnector
)

func handleTags(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got %s request to %s", r.Method, r.URL)
	controller := getController()
	switch r.Method {
	case http.MethodGet:
		{
			if tagId := r.URL.Query().Get("id"); r.URL.Query().Has("id") {
				controller.GetTag(w, tagId)
				break
			}
			controller.GetTags(w)
			break
		}
	case http.MethodPost:
		{
			controller.CreateTag(w, r.Body)
			break
		}
	case http.MethodPut:
		{
			controller.UpdateTag(w, r.Body)
			break
		}
	case http.MethodDelete:
		{
			if tagId := r.URL.Query().Get("id"); r.URL.Query().Has("id") {
				controller.DeleteTag(w, tagId)
				break
			}
			Controllers.BadRequest(w)
			break
		}
	}

}

func SetupTagsRoutes(basePath string) {
	tagHandler := http.HandlerFunc(handleTags)
	http.Handle(fmt.Sprintf("%s/v1/%s", basePath, path), tagHandler)
}

func SetupTagsDatabase(initDatabase Data.DbConnector) {
	database = initDatabase
}

//todo: this is so bad...
func getController() Controllers.TagsController {
	var repository = Repositories.NewTagsRepository(Repositories.NewTagQueryExecutor(&database))
	var controller = Controllers.NewTagsController(*repository)
	return *controller
}
