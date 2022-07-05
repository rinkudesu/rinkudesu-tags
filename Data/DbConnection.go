package Data

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

var (
	pool   *pgxpool.Pool
	closed bool
)

type DbConnection struct {
}

func (connection *DbConnection) InitialiseEnv() error {
	connectionString := os.Getenv("RINKU_TAGS_CONNECTIONSTRING")
	return connection.Initialise(connectionString)
}

func (connection *DbConnection) Initialise(connectionString string) error {
	if pool != nil {
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

	pool = localPool
	return nil
}

func (connection *DbConnection) QueryRow(sql string, args ...interface{}) (pgx.Row, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	if len(args) == 0 {
		return pool.QueryRow(context.Background(), sql), nil
	}
	return pool.QueryRow(context.Background(), sql, args...), nil
}

func (connection *DbConnection) QueryRows(sql string, args ...interface{}) (pgx.Rows, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	if len(args) == 0 {
		return pool.Query(context.Background(), sql)
	}
	return pool.Query(context.Background(), sql, args...)
}

func (connection *DbConnection) QueryFunc(sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	return pool.QueryFunc(context.Background(), sql, args, scans, f)
}

func (connection *DbConnection) Query(sql string) (pgx.Rows, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	return pool.Query(context.Background(), sql)
}

func (connection *DbConnection) Exec(sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if openErr := connection.ensureOpen(); openErr != nil {
		return nil, openErr
	}

	if len(args) == 0 {
		return pool.Exec(context.Background(), sql)
	}
	return pool.Exec(context.Background(), sql, args...)
}

func (connection *DbConnection) Close() {
	if closed {
		return
	}

	pool.Close()
	closed = true
}

func (connection *DbConnection) ensureOpen() error {
	if closed || pool == nil {
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
