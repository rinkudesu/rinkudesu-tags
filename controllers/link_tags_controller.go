package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/models"
	"rinkudesu-tags/repositories"
)

type LinkTagsController struct {
	repository      *repositories.LinkTagsRepository
	linksRepository *repositories.LinksRepository
	tagsRepository  *repositories.TagsRepository
}

func NewLinkTagsController(repository *repositories.LinkTagsRepository, linksRepository *repositories.LinksRepository, tagsRepository *repositories.TagsRepository) *LinkTagsController {
	return &LinkTagsController{repository: repository, linksRepository: linksRepository, tagsRepository: tagsRepository}
}

func (controller *LinkTagsController) Create(c *gin.Context) {
	var linkTag models.LinkTag
	err := BindJson(c, &linkTag)
	if err != nil {
		log.Infof("LinkTag model is invalid: %s", err.Error())
		return
	}
	userInfo := GetUserInfo(c)

	// ignore duplicate values errors, since they just mean required data is already available
	if err = controller.linksRepository.Create(&models.Link{Id: linkTag.LinkId}, userInfo); err != repositories.AlreadyExistsErr {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// to create a link-tag the tag must already exist, since we utilise the id here
	// if the tag no longer exists, then assume it was deleted and return error
	if tagExists, _ := controller.tagsRepository.Exists(linkTag.TagId, userInfo); !tagExists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = controller.repository.Create(&linkTag, userInfo)
	if err != nil {
		if err == repositories.AlreadyExistsErr {
			c.AbortWithStatus(http.StatusBadRequest)
		} else if err == repositories.NotFoundErr {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, &linkTag)
}

func (controller *LinkTagsController) Delete(c *gin.Context) {
	var ids struct {
		LinkId string `form:"linkId" binding:"required,uuid"`
		TagId  string `form:"tagId" binding:"required,uuid"`
	}
	err := c.BindQuery(&ids)
	if err != nil {
		return
	}

	parsedIds, err := ParseUuids([]string{ids.LinkId, ids.TagId})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = controller.repository.Remove(parsedIds[0], parsedIds[1], GetUserInfo(c))
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

func (controller *LinkTagsController) GetLinksForTag(c *gin.Context) {
	id, err := ParseUuidFromParam(c)
	if err != nil {
		return
	}

	result, err := controller.repository.GetLinksForTag(id, GetUserInfo(c))
	if err != nil {
		if err == repositories.NotFoundErr {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *LinkTagsController) GetTagsForLink(c *gin.Context) {
	id, err := ParseUuidFromParam(c)
	if err != nil {
		return
	}

	result, err := controller.repository.GetTagsForLink(id, GetUserInfo(c))
	if err != nil {
		if err == repositories.NotFoundErr {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
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
