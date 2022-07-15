package Mocks

import (
	"fmt"
	"github.com/gofrs/uuid"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
)

func GetDatabase() (Data.DbConnector, string) {
	database := &Data.DbConnection{}
	dbName, _ := uuid.NewV4()
	err := database.Initialise("postgres://postgres:postgres@localhost:5432/" + dbName.String())
	if err != nil {
		_ = database.Initialise("postgres://postgres:postgres@localhost:5432/postgres")
		_, _ = database.Exec(fmt.Sprintf("create database \"%s\"", dbName.String()))
		database.Close()
		database = &Data.DbConnection{}
		_ = database.Initialise("postgres://postgres:postgres@localhost:5432/" + dbName.String())
	}
	executor := Migrations.NewExecutor(database)
	executor.Migrate()
	return database, dbName.String()
}

func DropDatabase(existingConnection *Data.DbConnector, dbName string) {
	(*existingConnection).Close()
	database := &Data.DbConnection{}
	defer database.Close()
	_ = database.Initialise("postgres://postgres:postgres@localhost:5432/postgres")
	_, _ = database.Exec(fmt.Sprintf("drop database \"%s\"", dbName))
}
