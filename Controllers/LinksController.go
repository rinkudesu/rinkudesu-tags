package Controllers

import (
	"net/http"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Repositories"

	"github.com/gin-gonic/gin"
)

type LinksController struct {
	repository *Repositories.LinksRepository
}

func NewLinksController(repository *Repositories.LinksRepository) *LinksController {
	return &LinksController{repository: repository}
}

func (controller *LinksController) CreateLink(c *gin.Context) {
	var link Models.Link
	err := BindJson(c, &link)
	if err != nil {
		return
	}

	err = controller.repository.Create(&link)
	if err != nil {
		if err == Repositories.AlreadyExistsErr {
			c.Status(http.StatusBadRequest)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, &link)
}

func (controller *LinksController) DeleteLink(c *gin.Context) {
	linkUuid, err := ParseUuidFromParam("id", c)
	if err != nil {
		return
	}

	err = controller.repository.Delete(linkUuid)
	if err != nil {
		if err == Repositories.NotFoundErr {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}
	c.Status(http.StatusOK)
}

func (controller *LinksController) SetupRoutes(router *gin.Engine, basePath string) {
	const apiVersion = "v0"
	url := GetUrl(basePath, apiVersion, "links")

	router.POST(url, controller.CreateLink)
	router.DELETE(url+"/:id", controller.DeleteLink)
}
