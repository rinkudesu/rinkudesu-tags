package Authorisation

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type MockTokenVerifier struct {
	ReturnedToken  *oidc.IDToken
	ReturnedClaims *Claims
}

func (m *MockTokenVerifier) ValidateTokenFromHeader(_ *gin.Context) (*oidc.IDToken, *Claims, error) {
	if m.ReturnedToken == nil {
		return nil, nil, errors.New("failed to verify token")
	}
	return m.ReturnedToken, m.ReturnedClaims, nil
}

func (m *MockTokenVerifier) ValidateToken(_ string) (*oidc.IDToken, *Claims, error) {
	if m.ReturnedToken == nil {
		return nil, nil, errors.New("failed to verify token")
	}
	return m.ReturnedToken, m.ReturnedClaims, nil
}

func (m *MockTokenVerifier) Verify(_ context.Context, _ string) (*oidc.IDToken, error) {
	if m.ReturnedToken == nil {
		return nil, errors.New("failed to verify token")
	}
	return m.ReturnedToken, nil
}

type MockClaimsReader struct {
	Claims *Claims
}

func (m *MockClaimsReader) GetClaims(_ *oidc.IDToken) (*Claims, error) {
	if m.Claims == nil {
		return nil, errors.New("whatever")
	}
	return m.Claims, nil
}
