package Data

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

type DbConnection struct {
	pool   *pgxpool.Pool
	closed bool
}

func (connection *DbConnection) InitialiseEnv() error {
	connectionString := os.Getenv("RINKU_TAGS_CONNECTIONSTRING")
	return connection.Initialise(connectionString)
}

func (connection *DbConnection) Initialise(connectionString string) error {
	if connection.pool != nil {
		return alreadyInitialisedError{}
	}

	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		log.Println(err)
		return err
	}

	connection.pool = pool
	return nil
}

func (connection *DbConnection) QueryRow(sql string, args ...interface{}) (pgx.Row, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	if len(args) == 0 {
		return connection.pool.QueryRow(context.Background(), sql), nil
	}
	return connection.pool.QueryRow(context.Background(), sql, args), nil
}

func (connection *DbConnection) QueryRows(sql string, args ...interface{}) (pgx.Rows, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	if len(args) == 0 {
		return connection.pool.Query(context.Background(), sql)
	}
	return connection.pool.Query(context.Background(), sql, args)
}

func (connection *DbConnection) QueryFunc(sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	return connection.pool.QueryFunc(context.Background(), sql, args, scans, f)
}

func (connection *DbConnection) Query(sql string) (pgx.Rows, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	return connection.pool.Query(context.Background(), sql)
}

func (connection *DbConnection) Exec(sql string) error {
	if openErr := connection.ensureOpen(); openErr != nil {
		return openErr
	}

	_, err := connection.pool.Exec(context.Background(), sql)
	return err
}

func (connection *DbConnection) Close() {
	connection.pool.Close()
	connection.closed = true
}

func (connection *DbConnection) ensureOpen() error {
	if connection.closed {
		return connectionClosedError{}
	}
	return nil
}

type connectionClosedError struct {
}

func (err connectionClosedError) Error() string {
	return "This connection to the database has been already closed"
}

type alreadyInitialisedError struct {
}

func (err alreadyInitialisedError) Error() string {
	return "This connection to the database has been already initialised"
}
