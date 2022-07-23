package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	id := c.Param("id")
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
