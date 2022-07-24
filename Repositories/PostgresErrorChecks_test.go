package Repositories

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsPostgresDuplicateValue_NilError_False(t *testing.T) {
	assert.False(t, IsPostgresDuplicateValue(nil))
}

func TestIsPostgresDuplicateValue_InvalidTypeError_False(t *testing.T) {
	assert.False(t, IsPostgresDuplicateValue(errors.New("test")))
}

func TestIsPostgresDuplicateValue_InvalidErrorCode_False(t *testing.T) {
	assert.False(t, IsPostgresDuplicateValue(&pgconn.PgError{Code: ":)"}))
}

func TestIsPostgresDuplicateValue_CorrectErrorCode_True(t *testing.T) {
	assert.True(t, IsPostgresDuplicateValue(&pgconn.PgError{Code: "23505"}))
}

func TestIsPostgresNotFoundError_ErrorNil_False(t *testing.T) {
	assert.False(t, IsPostgresNotFoundError(nil))
}

func TestIsPostgresNotFoundError_ErrorInvalidType_False(t *testing.T) {
	assert.False(t, IsPostgresNotFoundError(errors.New("test")))
}

func TestIsPostgresNotFoundError_ErrorNoRows_True(t *testing.T) {
	assert.True(t, IsPostgresNotFoundError(pgx.ErrNoRows))
}

func TestIsPostgresNotFoundError_ErrorPgErrInvalidCode_False(t *testing.T) {
	assert.False(t, IsPostgresNotFoundError(&pgconn.PgError{Code: ":)"}))
}

func TestIsPostgresNotFoundError_ErrorPgErrCorrectCode_True(t *testing.T) {
	assert.True(t, IsPostgresNotFoundError(&pgconn.PgError{Code: "23503"}))
}
