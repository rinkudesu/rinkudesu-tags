package authorisation

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/models"
	"rinkudesu-tags/utils"
	"strings"
)

var (
	AuthTokenInvalid           = errors.New("authorisation token is invalid")
	MissingConfigurationValues = errors.New("SSO config values are missing")
)

type JWTHandler struct {
	oauthProvider *oidc.Provider
	tokenVerifier TokenVerifier
	claimsReader  JWTClaimsReader
}

func NewJWTHandler(config *models.Configuration) (*JWTHandler, error) {
	if config.SsoClientId == "" || config.SsoAuthority == "" {
		return nil, MissingConfigurationValues
	}
	oauthProvider, err := oidc.NewProvider(context.Background(), config.SsoAuthority)
	if err != nil {
		return nil, err
	}
	verifier := oauthProvider.Verifier(&oidc.Config{ClientID: config.SsoClientId})
	handler := &JWTHandler{
		oauthProvider: oauthProvider,
		tokenVerifier: verifier,
	}
	handler.claimsReader = handler
	return handler, nil
}

func (handler *JWTHandler) ValidateTokenFromHeader(c *gin.Context) (*oidc.IDToken, *Claims, error) {
	tokenHeader := c.GetHeader("Authorization")
	splitHeader := strings.Split(tokenHeader, " ")
	if len(splitHeader) < 2 || splitHeader[1] == "" {
		log.Warning("Failed to parse authorisation header")
		return nil, nil, AuthTokenInvalid
	}
	return handler.ValidateToken(splitHeader[1])
}

func (handler *JWTHandler) ValidateToken(rawToken string) (*oidc.IDToken, *Claims, error) {
	//verify with provider
	token, err := handler.tokenVerifier.Verify(context.Background(), rawToken)
	if err != nil {
		log.Warningf("Failed to verify JWT: %s", err.Error())
		return nil, nil, AuthTokenInvalid
	}

	//try reading claims
	claims, err := handler.claimsReader.GetClaims(token)
	if err != nil {
		return nil, nil, AuthTokenInvalid
	}

	//verify audience
	if !utils.Contains(claims.Aud, "rinkudesu") {
		log.Warning("JWT does not contain required audience")
		return nil, nil, AuthTokenInvalid
	}
	return token, claims, nil
}

func (handler *JWTHandler) GetClaims(token *oidc.IDToken) (*Claims, error) {
	readClaims := Claims{}
	err := token.Claims(&readClaims)
	if err != nil {
		log.Warningf("Failed to read claims from jwt: %s", err.Error())
		return nil, err
	}

	//attempt to parse claims with a single aud
	singleAud := claimsSingleAud{}
	err = token.Claims(&singleAud)
	if err == nil {
		readClaims.Aud = []string{singleAud.Aud}
		return &readClaims, nil
	}

	//if that failed, attempt to parse aud array
	multiAud := claimsArrayAud{}
	err = token.Claims(&multiAud)
	if err == nil {
		readClaims.Aud = multiAud.Aud
		return &readClaims, nil
	}

	return &readClaims, nil
}

type claimsSingleAud struct {
	Aud string `json:"aud"`
}

type claimsArrayAud struct {
	Aud []string `json:"aud"`
}
