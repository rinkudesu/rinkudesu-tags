package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

func (controller *LinkTagsController) Create(c *gin.Context) {
	var linkTag Models.LinkTag
	err := BindJson(c, &linkTag)
	if err != nil {
		return
	}

	err = controller.repository.Create(&linkTag, GetUserInfo(c))
	if err != nil {
		if err == Repositories.AlreadyExistsErr {
			c.Status(http.StatusBadRequest)
		} else if err == Repositories.NotFoundErr {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, &linkTag)
}

func (controller *LinkTagsController) Delete(c *gin.Context) {
	linkId, err := ParseUuidFromQuery("linkId", c)
	if err != nil {
		return
	}
	tagId, err := ParseUuidFromQuery("tagId", c)
	if err != nil {
		return
	}

	err = controller.repository.Remove(linkId, tagId, GetUserInfo(c))
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

func (controller *LinkTagsController) GetLinksForTag(c *gin.Context) {
	id, err := ParseUuidFromParam("id", c)
	if err != nil {
		return
	}

	result, err := controller.repository.GetLinksForTag(id, GetUserInfo(c))
	if err != nil {
		if err == Repositories.NotFoundErr {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *LinkTagsController) GetTagsForLink(c *gin.Context) {
	id, err := ParseUuidFromParam("id", c)
	if err != nil {
		return
	}

	result, err := controller.repository.GetTagsForLink(id, GetUserInfo(c))
	if err != nil {
		if err == Repositories.NotFoundErr {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *LinkTagsController) SetupRoutes(engine *gin.Engine, basePath string) {
	const apiVersion = "v0"
	url := GetUrl(basePath, apiVersion, "linkTags")

	engine.POST(url, controller.Create)
	engine.DELETE(url, controller.Delete)
	engine.GET(fmt.Sprintf("%s/getLinksForTag/:id", url), controller.GetLinksForTag)
	engine.GET(fmt.Sprintf("%s/getTagsForLink/:id", url), controller.GetTagsForLink)
}
