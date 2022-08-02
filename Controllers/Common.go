package Controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/Models"
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
		c.Status(http.StatusBadRequest)
	}
	return err
}

func ParseUuidFromParam(paramName string, c *gin.Context) (uuid.UUID, error) {
	id := c.Param(paramName)
	parsed, err := ParseUuid(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	return parsed, err
}

func ParseUuidFromQuery(paramName string, c *gin.Context) (uuid.UUID, error) {
	id := c.Query(paramName)
	parsed, err := ParseUuid(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	return parsed, err
}

func ParseUuid(id string) (uuid.UUID, error) {
	result, err := uuid.FromString(id)
	if err != nil {
		log.Infof("Unable to parse '%s' as uuid", id)
		return uuid.Nil, err
	}
	return result, nil
}

func GetUserInfo(c *gin.Context) *Models.UserInfo {
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
	return &Models.UserInfo{UserId: id}
}
