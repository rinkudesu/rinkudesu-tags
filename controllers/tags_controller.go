package controllers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/models"
	"rinkudesu-tags/repositories"
)

type TagsController struct {
	repository *repositories.TagsRepository
}

func NewTagsController(repository *repositories.TagsRepository) *TagsController {
	return &TagsController{repository: repository}
}

func (controller *TagsController) GetTags(c *gin.Context) {
	var query struct {
		Name   string `form:"name"`
		Offset int    `form:"offset"`
		Limit  int    `form:"limit"`
	}
	if err := c.BindQuery(&query); err != nil {
		log.Warnf("Failed to bind tags query string: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tags, err := controller.repository.GetTags(GetUserInfo(c), query.Name, query.Offset, query.Limit)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (controller *TagsController) GetTag(c *gin.Context) {
	tagUuid, err := ParseUuidFromParam(c)
	if err != nil {
		return
	}

	tag, err := controller.repository.GetTag(tagUuid, GetUserInfo(c))
	if err != nil {
		if err == repositories.NotFoundErr {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (controller *TagsController) CreateTag(c *gin.Context) {
	var tagVm models.TagCreateViewModel
	err := BindJson(c, &tagVm)
	if err != nil {
		log.Infof("Log object is not valid: %s", err.Error())
		return
	}

	tag := tagVm.GetTag()
	returnedTag, err := controller.repository.Create(&tag, GetUserInfo(c))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, returnedTag)
}

func (controller *TagsController) UpdateTag(c *gin.Context) {
	var tag models.Tag
	err := BindJson(c, &tag)
	if err != nil {
		log.Infof("Log object is not valid: %s", err.Error())
		return
	}

	returnedTag, err := controller.repository.Update(&tag, GetUserInfo(c))
	if err != nil {
		if err == repositories.NotFoundErr {
			c.AbortWithStatus(http.StatusNotFound)
		} else if err == repositories.AlreadyExistsErr {
			c.String(http.StatusBadRequest, "Tag already exists")
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, returnedTag)
}

func (controller *TagsController) DeleteTag(c *gin.Context) {
	uuidValue, err := ParseUuidFromParam(c)
	if err != nil {
		return
	}

	err = controller.repository.Delete(uuidValue, GetUserInfo(c))
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

func (controller *TagsController) SetupRoutes(engine *gin.Engine, basePath string) {
	const apiVersion = "v0"
	url := GetUrl(basePath, apiVersion, "tags")

	engine.GET(url, controller.GetTags)
	engine.GET(url+"/:id", controller.GetTag)
	engine.POST(url, controller.CreateTag)
	engine.PUT(url, controller.UpdateTag)
	engine.DELETE(url+"/:id", controller.DeleteTag)
}
