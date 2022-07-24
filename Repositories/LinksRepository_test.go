package Repositories

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Mocks"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Services"
	"testing"
)

type linksRepositoryTests struct {
	connection *Data.DbConnector
	repo       *LinksRepository
	dbName     string
}

func newLinksRepositoryTests() *linksRepositoryTests {
	database, name := Mocks.GetDatabase()
	return &linksRepositoryTests{
		connection: &database,
		dbName:     name,
		repo:       CreateLinksRepository(Services.NewGlobalState(database)),
	}
}

func TestLinksRepository_Create_DataCreated(t *testing.T) {
	test := newLinksRepositoryTests()
	defer Mocks.DropDatabase(test.connection, test.dbName)
	id, _ := uuid.NewV4()
	testLink := Models.Link{Id: id}

	result := test.repo.Create(&testLink)

	assert.Nil(t, result)
	linksRows, _ := (*test.connection).QueryRows("select * from links")
	defer linksRows.Close()
	loaded := false
	for linksRows.Next() {
		var loadedId uuid.UUID
		_ = linksRows.Scan(&loadedId)
		assert.Equal(t, id, loadedId)
		assert.False(t, loaded)
		loaded = true
	}
}

func TestLinksRepository_Create_DuplicateData(t *testing.T) {
	test := newLinksRepositoryTests()
	defer Mocks.DropDatabase(test.connection, test.dbName)
	id, _ := uuid.NewV4()
	testLink := Models.Link{Id: id}
	_ = test.repo.Create(&testLink)

	result := test.repo.Create(&testLink)

	assert.NotNil(t, result)
	assert.Equal(t, AlreadyExistsErr, result)
	linksRows, _ := (*test.connection).QueryRows("select * from links")
	defer linksRows.Close()
	loaded := false
	for linksRows.Next() {
		var loadedId uuid.UUID
		_ = linksRows.Scan(&loadedId)
		assert.Equal(t, id, loadedId)
		assert.False(t, loaded)
		loaded = true
	}
}

func TestLinksRepository_Delete_LinkExists(t *testing.T) {
	test := newLinksRepositoryTests()
	defer Mocks.DropDatabase(test.connection, test.dbName)
	id, _ := uuid.NewV4()
	testLink := Models.Link{Id: id}
	_ = test.repo.Create(&testLink)

	result := test.repo.Delete(id)

	assert.Nil(t, result)
	linksRows, _ := (*test.connection).QueryRows("select * from links")
	assert.False(t, linksRows.Next())
}

func TestLinksRepository_Delete_LinkDoesntExist(t *testing.T) {
	test := newLinksRepositoryTests()
	defer Mocks.DropDatabase(test.connection, test.dbName)
	id, _ := uuid.NewV4()

	result := test.repo.Delete(id)

	assert.NotNil(t, result)
	assert.Equal(t, NotFoundErr, result)
	linksRows, _ := (*test.connection).QueryRows("select * from links")
	assert.False(t, linksRows.Next())
}
