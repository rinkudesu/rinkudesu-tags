package Repositories

import (
	"fmt"
	"github.com/gofrs/uuid"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
	"rinkudesu-tags/Models"
	"testing"
)

func TestTagQueryExecutor_GetAll(t *testing.T) {
	repo, executor, database, dbName := getRepository()
	defer dropDatabase(database, dbName)
	newId, _ := uuid.NewV4()
	insertResult, err := executor.Insert(&Models.Tag{
		Name:   "this is a test",
		UserId: newId,
	})
	if err != nil {
		t.Fatalf("failed to setup test data")
	}
	var insertedId uuid.UUID
	err = insertResult.Scan(&insertedId)
	if err != nil {
		t.Fatalf("failed to setup test data")
	}

	result, err := repo.GetTags()

	if err != nil {
		t.Fatalf("err was not nil")
	}
	if len(result) != 1 {
		t.Fatalf("length of result was not 1")
	}
	if result[0].Name != "this is a test" || result[0].UserId != newId {
		t.Fatalf("unexpected value loaded")
	}
}

func getRepository() (*TagsRepository, TagQueryExecutable, *Data.DbConnector, string) {
	database, name := getDatabase()
	executor := NewTagQueryExecutor(&database)
	return NewTagsRepository(executor), executor, &database, name
}

func getDatabase() (Data.DbConnector, string) {
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

func dropDatabase(existingConnection *Data.DbConnector, dbName string) {
	(*existingConnection).Close()
	database := &Data.DbConnection{}
	defer database.Close()
	_ = database.Initialise("postgres://postgres:postgres@localhost:5432/postgres")
	_, _ = database.Exec(fmt.Sprintf("drop database \"%s\"", dbName))
}
