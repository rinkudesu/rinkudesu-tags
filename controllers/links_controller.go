package controllers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/models"
	"rinkudesu-tags/repositories"

	"github.com/gin-gonic/gin"
)

type LinksController struct {
	repository *repositories.LinksRepository
}

func NewLinksController(repository *repositories.LinksRepository) *LinksController {
	return &LinksController{repository: repository}
}

func (controller *LinksController) CreateLink(c *gin.Context) {
	var link models.Link
	err := BindJson(c, &link)
	if err != nil {
		log.Infof("Link model is invalid: %s", err.Error())
		return
	}

	err = controller.repository.Create(&link, GetUserInfo(c))
	if err != nil {
		if err == repositories.AlreadyExistsErr {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, &link)
}

func (controller *LinksController) DeleteLink(c *gin.Context) {
	linkUuid, err := ParseUuidFromParam(c)
	if err != nil {
		return
	}

	err = controller.repository.Delete(linkUuid, GetUserInfo(c))
	if err != nil {
		if err == repositories.NotFoundErr {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
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
