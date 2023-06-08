package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/models"
)

var (
	ErrUserIdNotAvailable = errors.New("user id is not available")
)

func GetUrl(basePath string, apiVersion string, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", basePath, apiVersion, endpoint)
}

func BindJson(c *gin.Context, obj any) error {
	err := c.BindJSON(obj)
	if err != nil {
		log.Infof("Failed to parse json: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return err
}

type id struct {
	Id string `form:"id" uri:"id" binding:"required,uuid"`
}

func ParseUuidFromParam(c *gin.Context) (uuid.UUID, error) {
	var id id
	err := c.BindUri(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.FromString(id.Id)
}

func GetUserInfo(c *gin.Context) *models.UserInfo {
	idValue, isPresent := c.Get("userId")
	if !isPresent {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Panic("User id missing from context")
	}
	id, ok := idValue.(uuid.UUID)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Panic("Unexpected type in gin context as user id")
	}
	return &models.UserInfo{UserId: id}
}

func ParseUuids(unparsed []string) ([]uuid.UUID, error) {
	parsed := make([]uuid.UUID, len(unparsed))
	for i, s := range unparsed {
		temp, err := uuid.FromString(s)
		if err != nil {
			return nil, err
		}
		parsed[i] = temp
	}
	return parsed, nil
}
