package Routers

import (
	"fmt"
	"net/http"
	"rinkudesu-tags/Controllers"
)

const path = "tags"

func handleTags(w http.ResponseWriter, r *http.Request) {
	controller := Controllers.TagsController{}
	controller.HandleProducts(w, r)
}

func SetupTagsRoutes(basePath string) {
	tagHandler := http.HandlerFunc(handleTags)
	http.Handle(fmt.Sprintf("%s/v1/%s", basePath, path), tagHandler)
}
