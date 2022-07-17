package Controllers

import (
	json2 "encoding/json"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Repositories"
)

//todo: look into some sort of DI

type TagsController struct {
	repository Repositories.TagsRepository
}

func NewTagsController(repository Repositories.TagsRepository) *TagsController {
	return &TagsController{repository: repository}
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
		log.Infof("Unable to parse '%s' as uuid", id)
		BadRequest(w)
		return
	}
	tag, err := controller.repository.GetTag(tagUuid)
	if err != nil {
		if err == Repositories.NotFoundErr {
			NotFound(w)
			return
		}
		InternalServerError(w)
		return
	}

	writeJsonResponse(w, 200, *tag)
	return
}

func (controller *TagsController) CreateTag(w http.ResponseWriter, tagBody io.ReadCloser) {
	defer closeBody(tagBody)
	body, err := readBody(tagBody)
	if err != nil {
		BadRequest(w)
		return
	}
	var tag Models.Tag
	err = parseJson(body, &tag)
	if err != nil {
		BadRequest(w)
		return
	}
	if !tag.IsValid() {
		log.Info("Log object is not valid")
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

	body, err := readBody(tagBody)
	if err != nil {
		BadRequest(w)
		return
	}
	var tag Models.Tag
	err = parseJson(body, &tag)
	if err != nil {
		BadRequest(w)
		return
	}

	returnedTag, err := controller.repository.Update(&tag)
	if err != nil {
		if err == Repositories.NotFoundErr {
			NotFound(w)
			return
		}
		BadRequest(w)
		return
	}

	writeJsonResponse(w, 200, returnedTag)
}

func (controller *TagsController) DeleteTag(w http.ResponseWriter, id string) {
	uuidValue, err := uuid.FromString(id)
	if err != nil {
		log.Infof("Unable to parse '%s' as uuid", id)
		BadRequest(w)
		return
	}

	err = controller.repository.Delete(uuidValue)
	if err != nil {
		if err == Repositories.NotFoundErr {
			NotFound(w)
			return
		}
		BadRequest(w)
		return
	}
	Ok(w)
}

func parseJson(json []byte, tag *Models.Tag) error {
	err := json2.Unmarshal(json, tag)
	if err != nil {
		log.Warningf("Failed to parse tag json: %s", err.Error())
	}
	return err
}

func readBody(body io.ReadCloser) ([]byte, error) {
	array, err := io.ReadAll(body)
	if err != nil {
		log.Warningf("Failed to read from body %s", err.Error())
	}
	return array, err
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Warningf("Failed to close request body: %s", err.Error())
	}
}

func writeJsonResponse(w http.ResponseWriter, code int, tags interface{}) {
	json, jsonErr := json2.Marshal(tags)
	if jsonErr != nil {
		log.Warningf("Failed to serialise to json: %s", jsonErr.Error())
		InternalServerError(w)
		return
	}
	w.WriteHeader(code)
	_, err := w.Write(json)
	if err != nil {
		log.Warningf("Failed to write response: %s", err.Error())
		InternalServerError(w)
		return
	}
}
