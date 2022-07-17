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

type TagsRouter struct {
	connection Data.DbConnector
}

func NewTagsRouter(connection Data.DbConnector, basePath string) *TagsRouter {
	router := TagsRouter{connection: connection}
	tagHandler := http.HandlerFunc(router.handleTags)
	http.Handle(fmt.Sprintf("%s/v1/%s", basePath, path), tagHandler)
	return &router
}

func (router *TagsRouter) handleTags(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got %s request to %s", r.Method, r.URL)
	controller := router.getController()
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

//todo: this is so bad...
func (router *TagsRouter) getController() Controllers.TagsController {
	var repository = Repositories.NewTagsRepository(Repositories.NewTagQueryExecutor(router.connection))
	var controller = Controllers.NewTagsController(*repository)
	return *controller
}
