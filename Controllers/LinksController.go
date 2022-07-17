package Controllers

import (
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

func (controller *LinksController) CreateLink(w http.ResponseWriter, id string) {
	linkUuid, err := ParseUuid(id)
	if err != nil {
		BadRequest(w)
		return
	}
	link := Models.Link{Id: linkUuid}
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
