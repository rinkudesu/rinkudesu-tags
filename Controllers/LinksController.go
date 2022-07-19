package Controllers

import (
	json2 "encoding/json"
	"io"
	"net/http"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Repositories"
)

type LinksController struct {
	repository *Repositories.LinksRepository
}

func NewLinksController(repository *Repositories.LinksRepository) *LinksController {
	return &LinksController{repository: repository}
}

func (controller *LinksController) CreateLink(w http.ResponseWriter, linkBody io.ReadCloser) {
	defer CloseBody(linkBody)
	linkBytes, err := ReadBody(linkBody)
	if err != nil {
		InternalServerError(w)
		return
	}
	var link Models.Link
	err = json2.Unmarshal(linkBytes, &link)
	if err != nil {
		BadRequest(w)
		return
	}
	err = controller.repository.Create(&link)
	if err != nil {
		if err == Repositories.AlreadyExistsErr {
			BadRequest(w)
		} else {
			InternalServerError(w)
		}
		return
	}
	WriteJsonResponse(w, 201, &link)
}

func (controller *LinksController) DeleteLink(w http.ResponseWriter, id string) {
	linkUuid, err := ParseUuid(id)
	if err != nil {
		BadRequest(w)
		return
	}
	err = controller.repository.Delete(linkUuid)
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
