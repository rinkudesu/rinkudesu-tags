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
	return err != nil && err == pgx.ErrNoRows
}
