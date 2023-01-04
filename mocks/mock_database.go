package mocks

import (
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"rinkudesu-tags/data"
	"rinkudesu-tags/data/migrations"
)

func GetDatabase() (data.DbConnector, string) {
	database := &data.DbConnection{}
	dbName, _ := uuid.NewV4()
	baseConnectionString := getBaseConnectionString()
	err := database.Initialise(baseConnectionString + dbName.String())
	if err != nil {
		_ = database.Initialise(baseConnectionString + "postgres")
		_, _ = database.Exec(fmt.Sprintf("create database \"%s\"", dbName.String()))
		database.Close()
		database = &data.DbConnection{}
		_ = database.Initialise(baseConnectionString + dbName.String())
	}
	executor := migrations.NewExecutor(database)
	executor.Migrate()
	return database, dbName.String()
}

func DropDatabase(existingConnection data.DbConnector, dbName string) {
	existingConnection.Close()
	database := &data.DbConnection{}
	defer database.Close()
	_ = database.Initialise(getBaseConnectionString() + "postgres")
	_, _ = database.Exec(fmt.Sprintf("drop database \"%s\"", dbName))
}

func getBaseConnectionString() string {
	baseConnectionString := os.Getenv("TEST_POSTGRES")
	if baseConnectionString == "" {
		baseConnectionString = "postgres://postgres:postgres@localhost:5432/"
	}
	return baseConnectionString
}
