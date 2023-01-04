package authorisation

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestJWTHandler_ValidateTokenFromHeader_HeaderNotSet(t *testing.T) {
	testData := []string{"", "invalid", "invalid "}
	handler := JWTHandler{}
	ginContext := gin.Context{Request: &http.Request{Header: http.Header{}}}
	for _, test := range testData {
		ginContext.Request.Header.Set("Authorization", test)
		_, _, err := handler.ValidateTokenFromHeader(&ginContext)
		assert.Equal(t, AuthTokenInvalid, err)
	}
}

func TestJWTHandler_ValidateToken_ValidationFailed(t *testing.T) {
	handler := JWTHandler{tokenVerifier: &MockTokenVerifier{}}

	token, claims, err := handler.ValidateToken("this is a test, please ignore")

	assert.Nil(t, token)
	assert.Nil(t, claims)
	assert.NotNil(t, err)
}

func TestJWTHandler_ValidateToken_FailedToReadClaims(t *testing.T) {
	handler := &JWTHandler{tokenVerifier: &MockTokenVerifier{ReturnedToken: &oidc.IDToken{}}, claimsReader: &MockClaimsReader{}}

	token, claims, err := handler.ValidateToken("this is a test, please ignore")

	assert.Nil(t, token)
	assert.Nil(t, claims)
	assert.NotNil(t, err)
}

func TestJWTHandler_ValidateToken_NoRinkudesuInAudience(t *testing.T) {
	handler := &JWTHandler{tokenVerifier: &MockTokenVerifier{ReturnedToken: &oidc.IDToken{}}, claimsReader: &MockClaimsReader{Claims: &Claims{Aud: []string{"test"}}}}

	token, claims, err := handler.ValidateToken("this is a test, please ignore")

	assert.Nil(t, token)
	assert.Nil(t, claims)
	assert.NotNil(t, err)
}

func TestJWTHandler_ValidateToken_ValidTokenProvided(t *testing.T) {
	handler := &JWTHandler{tokenVerifier: &MockTokenVerifier{ReturnedToken: &oidc.IDToken{}}, claimsReader: &MockClaimsReader{Claims: &Claims{Aud: []string{"rinkudesu"}}}}

	token, claims, err := handler.ValidateToken("this is a test, please ignore")

	assert.NotNil(t, token)
	assert.NotNil(t, claims)
	assert.Nil(t, err)
	assert.Contains(t, claims.Aud, "rinkudesu")
}
