package Repositories

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func IsPostgresDuplicateValue(err error) bool {
	var pgErr *pgconn.PgError
	if err == nil || !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "23505"
}

func IsPostgresNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if err == pgx.ErrNoRows {
		return true
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return true
	}
	return false
}
