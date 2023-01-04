package authorisation

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type JWTValidator interface {
	ValidateTokenFromHeader(c *gin.Context) (*oidc.IDToken, *Claims, error)
	ValidateToken(rawToken string) (*oidc.IDToken, *Claims, error)
}

type JWTClaimsReader interface {
	GetClaims(token *oidc.IDToken) (*Claims, error)
}

type TokenVerifier interface {
	Verify(ctx context.Context, rawToken string) (*oidc.IDToken, error)
}
