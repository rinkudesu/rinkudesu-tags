package Controllers

import (
	json2 "encoding/json"
	"github.com/gofrs/uuid"
	"net/http"
	"rinkudesu-tags/Repositories"
)

//todo: look into some sort of DI

type TagsController struct {
	repository Repositories.TagsRepository
}

func (controller *TagsController) GetTags(w http.ResponseWriter) {
	tags, err := controller.repository.GetTags()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	writeJsonResponse(w, tags)
	return
}

func (controller *TagsController) GetTag(w http.ResponseWriter, id string) {
	tagUuid, err := uuid.FromString(id)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	tag, err := controller.repository.GetTag(tagUuid)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if tag == nil {
		w.WriteHeader(404)
		return
	}

	writeJsonResponse(w, *tag)
	return
}

func writeJsonResponse(w http.ResponseWriter, tags interface{}) {
	json, jsonErr := json2.Marshal(tags)
	if jsonErr != nil {
		w.WriteHeader(500)
		return
	}
	_, err := w.Write(json)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func (controller *TagsController) Init(initRepository Repositories.TagsRepository) {
	controller.repository = initRepository
}
