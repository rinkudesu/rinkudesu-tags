package Controllers

import "github.com/gin-gonic/gin"

type Routable interface {
	SetupRoutes(engine *gin.Engine, basePath string)
}
