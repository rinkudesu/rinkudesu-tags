package Data

import (
	"context"
	"errors"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	ConnectionClosedError   = errors.New("this connection to the database has already been closed")
	AlreadyInitialisedError = errors.New("this connection to the database has already been initialised")
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
		return AlreadyInitialisedError
	}

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Panicf("Unable to create database connection config: %s", err.Error())
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
		log.Errorf("Failed to connect to database: %s", err.Error())
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

	log.Info("Closing database connection")
	connection.pool.Close()
	connection.pool = nil
	connection.closed = true
}

func (connection *DbConnection) ensureOpen() error {
	if connection.closed || connection.pool == nil {
		return ConnectionClosedError
	}
	return nil
}
