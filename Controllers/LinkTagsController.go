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

type LinkTagsController struct {
	repository *Repositories.LinkTagsRepository
}

func NewLinkTagsController(repository *Repositories.LinkTagsRepository) *LinkTagsController {
	return &LinkTagsController{repository: repository}
}

func (controller *LinkTagsController) Create(w http.ResponseWriter, linkTagBody io.ReadCloser) {
	defer CloseBody(linkTagBody)
	linkTagBytes, err := ReadBody(linkTagBody)
	if err != nil {
		InternalServerError(w)
		return
	}

	var linkTag Models.LinkTag
	err = json2.Unmarshal(linkTagBytes, &linkTag)
	if err != nil {
		log.Warningf("Failed to parse linkTag from json: %s", err.Error())
		BadRequest(w)
		return
	}

	err = controller.repository.Create(&linkTag)
	if err != nil {
		if err == Repositories.AlreadyExistsErr {
			BadRequest(w)
		} else if err == Repositories.NotFoundErr {
			NotFound(w)
		} else {
			InternalServerError(w)
		}
		return
	}

	WriteJsonResponse(w, 201, &linkTag)
}

func (controller *LinkTagsController) Delete(w http.ResponseWriter, stringId string) {
	id, err := uuid.FromString(stringId)
	if err != nil {
		log.Infof("Failed to parse %s as uuid", stringId)
		BadRequest(w)
		return
	}

	err = controller.repository.Remove(id)
	if err != nil {
		if err == Repositories.NotFoundErr {
			NotFound(w)
		} else {
			InternalServerError(w)
		}
		return
	}

	Ok(w)
}

func (controller *LinkTagsController) GetLinksForTag(w http.ResponseWriter, stringId string) {
	id, err := uuid.FromString(stringId)
	if err != nil {
		log.Infof("Failed to parse %s as uuid", stringId)
		BadRequest(w)
		return
	}

	result, err := controller.repository.GetLinksForTag(id)
	if err != nil {
		if err == Repositories.NotFoundErr {
			NotFound(w)
		} else {
			InternalServerError(w)
		}
		return
	}

	WriteJsonResponse(w, 200, result)
}

func (controller *LinkTagsController) GetTagsForLink(w http.ResponseWriter, stringId string) {
	id, err := uuid.FromString(stringId)
	if err != nil {
		log.Infof("Failed to parse %s as uuid", stringId)
		BadRequest(w)
		return
	}

	result, err := controller.repository.GetTagsForLink(id)
	if err != nil {
		if err == Repositories.NotFoundErr {
			NotFound(w)
		} else {
			InternalServerError(w)
		}
		return
	}

	WriteJsonResponse(w, 200, result)
}
