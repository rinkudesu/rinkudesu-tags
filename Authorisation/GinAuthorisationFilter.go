package Authorisation

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetGinAuthorisationFilter(jwtValidator JWTValidator) gin.HandlerFunc {
	return func(context *gin.Context) {
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
