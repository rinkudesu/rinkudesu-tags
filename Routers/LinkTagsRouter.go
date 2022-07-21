package Routers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/Controllers"
)

type LinkTagsRouter struct {
	controller *Controllers.LinkTagsController
}

func NewLinkTagsRouter(controller *Controllers.LinkTagsController, basePath string) *LinkTagsRouter {
	const path = "linkTags"
	router := LinkTagsRouter{controller: controller}
	tagHandler := http.HandlerFunc(router.handleLinkTags)
	http.Handle(fmt.Sprintf("%s/v1/%s", basePath, path), tagHandler)
	tagsForLinkHandler := http.HandlerFunc(router.getTagsForLink)
	http.Handle(fmt.Sprintf("%s/v1/%s/getTagsForLink", basePath, path), tagsForLinkHandler)
	linksForTagHandler := http.HandlerFunc(router.getLinksForTag)
	http.Handle(fmt.Sprintf("%s/v1/%s/getLinksForTag", basePath, path), linksForTagHandler)
	return &router
}

func (router *LinkTagsRouter) getLinksForTag(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got %s request to %s", r.Method, r.URL)
	if r.Method != http.MethodGet {
		Controllers.MethodNotAllowed(w)
		return
	}

	if id := r.URL.Query().Get("id"); r.URL.Query().Has("id") {
		router.controller.GetLinksForTag(w, id)
	} else {
		Controllers.BadRequest(w)
	}
}

func (router *LinkTagsRouter) getTagsForLink(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got %s request to %s", r.Method, r.URL)
	if r.Method != http.MethodGet {
		Controllers.MethodNotAllowed(w)
		return
	}

	if id := r.URL.Query().Get("id"); r.URL.Query().Has("id") {
		router.controller.GetTagsForLink(w, id)
	} else {
		Controllers.BadRequest(w)
	}
}

func (router *LinkTagsRouter) handleLinkTags(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got %s request to %s", r.Method, r.URL)
	switch r.Method {
	case http.MethodPost:
		{
			router.controller.Create(w, r.Body)
			break
		}
	case http.MethodDelete:
		{
			if id := r.URL.Query().Get("id"); r.URL.Query().Has("id") {
				router.controller.Delete(w, id)
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
