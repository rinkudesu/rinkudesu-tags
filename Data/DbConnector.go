package Data

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type DbConnector interface {
	InitialiseEnv() error
	Initialise(connectionString string) error
	QueryRow(sql string, args ...interface{}) (pgx.Row, error)
	QueryRows(sql string, args ...interface{}) (pgx.Rows, error)
	QueryFunc(sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	Query(sql string) (pgx.Rows, error)
	Exec(sql string, args ...interface{}) (pgconn.CommandTag, error)
	Close()
}
