package Controllers

import (
	json2 "encoding/json"
	"net/http"
	"rinkudesu-tags/Repositories"
)

//todo: look into some sort of DI

type TagsController struct {
	repository Repositories.TagsRepository
}

func (controller *TagsController) HandleProducts(w http.ResponseWriter, r *http.Request) {
	tags, err := controller.repository.GetTags()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	json, jsonErr := json2.Marshal(tags)
	if jsonErr != nil {
		w.WriteHeader(500)
		return
	}
	_, err = w.Write(json)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	return
}

func (controller *TagsController) Init(initRepository Repositories.TagsRepository) {
	controller.repository = initRepository
}
