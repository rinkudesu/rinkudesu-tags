package Mocks

import (
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
)

func GetDatabase() (Data.DbConnector, string) {
	database := &Data.DbConnection{}
	dbName, _ := uuid.NewV4()
	baseConnectionString := getBaseConnectionString()
	err := database.Initialise(baseConnectionString + dbName.String())
	if err != nil {
		_ = database.Initialise(baseConnectionString + "postgres")
		_, _ = database.Exec(fmt.Sprintf("create database \"%s\"", dbName.String()))
		database.Close()
		database = &Data.DbConnection{}
		_ = database.Initialise(baseConnectionString + dbName.String())
	}
	executor := Migrations.NewExecutor(database)
	executor.Migrate()
	return database, dbName.String()
}

func DropDatabase(existingConnection Data.DbConnector, dbName string) {
	existingConnection.Close()
	database := &Data.DbConnection{}
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
