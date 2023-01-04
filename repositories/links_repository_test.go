package repositories

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"rinkudesu-tags/data"
	"rinkudesu-tags/mocks"
	"rinkudesu-tags/models"
	"rinkudesu-tags/services"
	"testing"
)

type linksRepositoryTests struct {
	connection data.DbConnector
	repo       *LinksRepository
	dbName     string
	userInfo   *models.UserInfo
}

func newLinksRepositoryTests() *linksRepositoryTests {
	database, name := mocks.GetDatabase()
	userId, _ := uuid.NewV4()
	return &linksRepositoryTests{
		connection: database,
		dbName:     name,
		repo:       CreateLinksRepository(services.NewGlobalState(database)),
		userInfo:   &models.UserInfo{UserId: userId},
	}
}

func (test *linksRepositoryTests) close() {
	mocks.DropDatabase(test.connection, test.dbName)
}

func TestLinksRepository_Create_DataCreated(t *testing.T) {
	test := newLinksRepositoryTests()
	defer mocks.DropDatabase(test.connection, test.dbName)
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}

	result := test.repo.Create(&testLink, test.userInfo)

	assert.Nil(t, result)
	linksRows, _ := test.connection.QueryRows("select id from links")
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
	defer mocks.DropDatabase(test.connection, test.dbName)
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	_ = test.repo.Create(&testLink, test.userInfo)

	result := test.repo.Create(&testLink, test.userInfo)

	assert.NotNil(t, result)
	assert.Equal(t, AlreadyExistsErr, result)
	linksRows, _ := test.connection.QueryRows("select id from links")
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
	t.Cleanup(func() {
		mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	_ = test.repo.Create(&testLink, test.userInfo)

	result := test.repo.Delete(id, test.userInfo)

	assert.Nil(t, result)
	linksRows, _ := test.connection.QueryRows("select * from links")
	assert.False(t, linksRows.Next())
}

func TestLinksRepository_Delete_LinkCreatedByAnotherUser_FailsToDelete(t *testing.T) {
	test := newLinksRepositoryTests()
	t.Cleanup(func() {
		mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := models.UserInfo{UserId: anotherUserId}
	_ = test.repo.Create(&testLink, &anotherUserInfo)

	result := test.repo.Delete(id, test.userInfo)

	assert.NotNil(t, result)
	assert.Equal(t, NotFoundErr, result)
}

func TestLinksRepository_Delete_LinkDoesntExist(t *testing.T) {
	test := newLinksRepositoryTests()
	defer mocks.DropDatabase(test.connection, test.dbName)
	id, _ := uuid.NewV4()

	result := test.repo.Delete(id, test.userInfo)

	assert.NotNil(t, result)
	assert.Equal(t, NotFoundErr, result)
	linksRows, _ := test.connection.QueryRows("select * from links")
	assert.False(t, linksRows.Next())
}

func TestLinksRepository_Exists_ExistsForDifferentUser_ReturnsFalse(t *testing.T) {
	test := newLinksRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := models.UserInfo{UserId: anotherUserId}
	_ = test.repo.Create(&testLink, &anotherUserInfo)

	result, err := test.repo.Exists(id, test.userInfo)

	assert.False(t, result)
	assert.Nil(t, err)
}

func TestLinksRepository_Exists_ExistsForCurrentUser_ReturnsTrue(t *testing.T) {
	test := newLinksRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	_ = test.repo.Create(&testLink, test.userInfo)

	result, err := test.repo.Exists(id, test.userInfo)

	assert.True(t, result)
	assert.Nil(t, err)
}

func TestLinksRepository_DeleteForce_LinkExists(t *testing.T) {
	test := newLinksRepositoryTests()
	t.Cleanup(func() {
		mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	_ = test.repo.Create(&testLink, test.userInfo)

	result := test.repo.DeleteForce(id)

	assert.Nil(t, result)
	linksRows, _ := test.connection.QueryRows("select * from links")
	assert.False(t, linksRows.Next())
}

func TestLinksRepository_DeleteForce_LinkDoesntExist(t *testing.T) {
	test := newLinksRepositoryTests()
	defer mocks.DropDatabase(test.connection, test.dbName)
	id, _ := uuid.NewV4()

	result := test.repo.DeleteForce(id)

	assert.NotNil(t, result)
	assert.Equal(t, NotFoundErr, result)
	linksRows, _ := test.connection.QueryRows("select * from links")
	assert.False(t, linksRows.Next())
}

func TestLinksRepository_DeleteForUser_NoneExistForUser(t *testing.T) {
	test := newLinksRepositoryTests()
	t.Cleanup(func() {
		mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := models.UserInfo{UserId: anotherUserId}
	_ = test.repo.Create(&testLink, &anotherUserInfo)

	err := test.repo.DeleteForUser(test.userInfo.UserId)

	assert.Nil(t, err)
	exists, _ := test.repo.Exists(id, &anotherUserInfo)
	assert.True(t, exists)
}

func TestLinksRepository_DeleteForUser_ExistsForUser(t *testing.T) {
	test := newLinksRepositoryTests()
	t.Cleanup(func() {
		mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	_ = test.repo.Create(&testLink, test.userInfo)

	err := test.repo.DeleteForUser(test.userInfo.UserId)

	assert.Nil(t, err)
	exists, _ := test.repo.Exists(id, test.userInfo)
	assert.False(t, exists)
}
