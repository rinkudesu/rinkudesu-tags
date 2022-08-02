package Authorisation

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"rinkudesu-tags/Controllers"
	"rinkudesu-tags/Mocks"
	"testing"
)

type ginAuthorisationFilterTests struct {
	jwtValidator JWTValidator
	context      *gin.Context
}

func newGinAuthorisationFilterTests(returnedToken *oidc.IDToken) *ginAuthorisationFilterTests {
	return &ginAuthorisationFilterTests{jwtValidator: &MockTokenVerifier{ReturnedToken: returnedToken}, context: &gin.Context{Writer: &Mocks.NilResponseWriter{}}}
}

func newGinAuthorisationFilterTestsWithClaims(returnedToken *oidc.IDToken, returnedClaims *Claims) *ginAuthorisationFilterTests {
	return &ginAuthorisationFilterTests{jwtValidator: &MockTokenVerifier{ReturnedToken: returnedToken, ReturnedClaims: returnedClaims}, context: &gin.Context{Writer: &Mocks.NilResponseWriter{}}}
}

func TestGetGinAuthorisationFilter_FailedToValidateToken_Aborted(t *testing.T) {
	test := newGinAuthorisationFilterTests(nil)

	GetGinAuthorisationFilter(test.jwtValidator)(test.context)

	assert.True(t, test.context.IsAborted())
}

func TestGetGinAuthorisationFilter_UserIdIncorrect_Aborted(t *testing.T) {
	test := newGinAuthorisationFilterTestsWithClaims(&oidc.IDToken{}, &Claims{Id: "test"})

	GetGinAuthorisationFilter(test.jwtValidator)(test.context)

	assert.True(t, test.context.IsAborted())
	token, ok := test.context.Get("token")
	assert.NotNil(t, token)
	assert.True(t, ok)
	claims, ok := test.context.Get("claims")
	assert.NotNil(t, claims)
	assert.True(t, ok)
}

func TestGetGinAuthorisationFilter_UserIdValid_NotAborted(t *testing.T) {
	userId, _ := uuid.NewV4()
	test := newGinAuthorisationFilterTestsWithClaims(&oidc.IDToken{}, &Claims{Id: userId.String()})

	GetGinAuthorisationFilter(test.jwtValidator)(test.context)

	assert.False(t, test.context.IsAborted())
	token, ok := test.context.Get("token")
	assert.NotNil(t, token)
	assert.True(t, ok)
	claims, ok := test.context.Get("claims")
	assert.NotNil(t, claims)
	assert.True(t, ok)
	storedUserId, err := Controllers.GetUserId(test.context)
	assert.Nil(t, err)
	assert.Equal(t, userId, storedUserId)
}
