package Data

import (
	"context"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
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

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Panicln("Unable to create database connection config")
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}

	localPool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Println(err)
		return err
	}

	connection.pool = localPool
	connection.closed = false
	return nil
}

func (connection *DbConnection) QueryRow(sql string, args ...interface{}) (Row, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	return connection.pool.QueryRow(context.Background(), sql, args...), nil
}

func (connection *DbConnection) QueryRows(sql string, args ...interface{}) (Rows, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	return connection.pool.Query(context.Background(), sql, args...)
}

func (connection *DbConnection) Query(sql string) (Rows, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	return connection.pool.Query(context.Background(), sql)
}

func (connection *DbConnection) Exec(sql string, args ...interface{}) (ExecResult, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	return connection.pool.Exec(context.Background(), sql, args...)
}

func (connection *DbConnection) Close() {
	if connection.closed {
		return
	}

	connection.pool.Close()
	connection.pool = nil
	connection.closed = true
}

func (connection *DbConnection) ensureOpen() error {
	if connection.closed || connection.pool == nil {
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
