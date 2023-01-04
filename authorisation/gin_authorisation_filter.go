package authorisation

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rinkudesu-tags/models"
)

func GetGinAuthorisationFilter(jwtValidator JWTValidator, config *models.Configuration) gin.HandlerFunc {
	return func(context *gin.Context) {
		if config.IgnoreAuthorisation {
			log.Warningf("User authorisation is disabled by a config value. Using %s as user id", uuid.Nil.String())
			context.Set("userId", uuid.Nil)
			context.Next()
			return
		}

		token, claims, err := jwtValidator.ValidateTokenFromHeader(context)
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		context.Set("token", token)
		context.Set("claims", claims)
		userId, err := uuid.FromString(claims.Id)
		if err != nil {
			log.Warnf("Failed to parse user id from JWT: %s", err.Error())
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		context.Set("userId", userId)

		context.Next()
	}
}
