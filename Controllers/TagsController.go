package Controllers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

func (controller *TagsController) GetTags(c *gin.Context) {
	tags, err := controller.repository.GetTags()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (controller *TagsController) GetTag(c *gin.Context) {
	tagUuid, err := ParseUuidFromParam("id", c)
	if err != nil {
		return
	}

	tag, err := controller.repository.GetTag(tagUuid)
	if err != nil {
		if err == Repositories.NotFoundErr {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (controller *TagsController) CreateTag(c *gin.Context) {
	var tag Models.Tag
	err := BindJson(c, &tag)
	if err != nil {
		return
	}

	if !tag.IsValid() {
		log.Info("Log object is not valid")
		c.Status(http.StatusBadRequest)
		return
	}

	returnedTag, err := controller.repository.Create(&tag)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, returnedTag)
}

func (controller *TagsController) UpdateTag(c *gin.Context) {
	var tag Models.Tag
	err := BindJson(c, &tag)
	if err != nil {
		return
	}

	if !tag.IsValid() {
		log.Info("Log object is not valid")
		c.Status(http.StatusBadRequest)
		return
	}

	returnedTag, err := controller.repository.Update(&tag)
	if err != nil {
		if err == Repositories.NotFoundErr {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, returnedTag)
}

func (controller *TagsController) DeleteTag(c *gin.Context) {
	uuidValue, err := ParseUuidFromParam("id", c)
	if err != nil {
		return
	}

	err = controller.repository.Delete(uuidValue)
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

func (controller *TagsController) SetupRoutes(engine *gin.Engine, basePath string) {
	const apiVersion = "v0"
	url := GetUrl(basePath, apiVersion, "tags")

	engine.GET(url, controller.GetTags)
	engine.GET(url+"/:id", controller.GetTag)
	engine.POST(url, controller.CreateTag)
	engine.PUT(url, controller.UpdateTag)
	engine.DELETE(url+"/:id", controller.DeleteTag)
}
