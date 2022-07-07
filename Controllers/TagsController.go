package Controllers

import (
	json2 "encoding/json"
	"github.com/gofrs/uuid"
	"io"
	"net/http"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Repositories"
)

//todo: look into some sort of DI

type TagsController struct {
	repository Repositories.TagsRepository
}

func (controller *TagsController) GetTags(w http.ResponseWriter) {
	tags, err := controller.repository.GetTags()
	if err != nil {
		InternalServerError(w)
		return
	}

	writeJsonResponse(w, 200, tags)
	return
}

func (controller *TagsController) GetTag(w http.ResponseWriter, id string) {
	tagUuid, err := uuid.FromString(id)
	if err != nil {
		BadRequest(w)
		return
	}
	tag, err := controller.repository.GetTag(tagUuid)
	if err != nil {
		InternalServerError(w)
		return
	}
	if tag == nil {
		NotFound(w)
		return
	}

	writeJsonResponse(w, 200, *tag)
	return
}

func (controller *TagsController) CreateTag(w http.ResponseWriter, tagBody io.ReadCloser) {
	defer closeBody(tagBody)
	body, err := io.ReadAll(tagBody)
	if err != nil {
		BadRequest(w)
		return
	}
	var tag Models.Tag
	err = json2.Unmarshal(body, &tag)
	if err != nil {
		BadRequest(w)
		return
	}
	if !tag.IsValid() {
		BadRequest(w)
		return
	}
	returnedTag, err := controller.repository.Create(&tag)
	if err != nil {
		BadRequest(w)
		return
	}
	writeJsonResponse(w, 201, returnedTag)
}

func (controller *TagsController) UpdateTag(w http.ResponseWriter, tagBody io.ReadCloser) {
	defer closeBody(tagBody)

	body, err := io.ReadAll(tagBody)
	if err != nil {
		BadRequest(w)
		return
	}
	var tag Models.Tag
	err = json2.Unmarshal(body, &tag)
	if err != nil {
		BadRequest(w)
		return
	}

	returnedTag, err := controller.repository.Update(&tag)
	if err != nil {
		BadRequest(w)
		return
	}
	if returnedTag == nil {
		NotFound(w)
		return
	}

	writeJsonResponse(w, 200, returnedTag)
}

func (controller *TagsController) DeleteTag(w http.ResponseWriter, id string) {
	uuidValue, err := uuid.FromString(id)
	if err != nil {
		BadRequest(w)
		return
	}

	err = controller.repository.Delete(uuidValue)
	if err != nil {
		BadRequest(w)
		return
	}
	Ok(w)
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		//todo: log
	}
}

func writeJsonResponse(w http.ResponseWriter, code int, tags interface{}) {
	json, jsonErr := json2.Marshal(tags)
	if jsonErr != nil {
		InternalServerError(w)
		return
	}
	w.WriteHeader(code)
	_, err := w.Write(json)
	if err != nil {
		InternalServerError(w)
		return
	}
}

func (controller *TagsController) Init(initRepository Repositories.TagsRepository) {
	controller.repository = initRepository
}
