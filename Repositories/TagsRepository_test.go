package Repositories

import (
	"github.com/gofrs/uuid"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Mocks"
	"rinkudesu-tags/Models"
	"testing"
)

func TestTagQueryExecutor_GetAll(t *testing.T) {
	repo, executor, database, dbName := getRepository()
	defer Mocks.DropDatabase(database, dbName)
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
	database, name := Mocks.GetDatabase()
	executor := NewTagQueryExecutor(&database)
	return NewTagsRepository(executor), executor, &database, name
}
