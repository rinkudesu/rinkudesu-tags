package Routers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/Controllers"
)

type LinksRouter struct {
	controller Controllers.LinksController
}

func NewLinksRouter(controller *Controllers.LinksController, basePath string) *LinksRouter {
	const path = "links"
	router := LinksRouter{controller: *controller}
	tagHandler := http.HandlerFunc(router.handleTags)
	http.Handle(fmt.Sprintf("%s/v1/%s", basePath, path), tagHandler)
	return &router
}

func (router *LinksRouter) handleTags(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got %s request to %s", r.Method, r.URL)
	switch r.Method {
	case http.MethodPost:
		{
			router.controller.CreateLink(w, r.Body)
			break
		}
	case http.MethodDelete:
		{
			if id := r.URL.Query().Get("id"); r.URL.Query().Has("id") {
				router.controller.DeleteLink(w, id)
				break
			}
			Controllers.BadRequest(w)
			break
		}
	default:
		{
			Controllers.MethodNotAllowed(w)
			break
		}
	}
}
