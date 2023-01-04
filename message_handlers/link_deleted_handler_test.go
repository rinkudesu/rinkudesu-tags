package message_handlers

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"rinkudesu-tags/data"
	"rinkudesu-tags/mocks"
	"rinkudesu-tags/models"
	"rinkudesu-tags/repositories"
	"rinkudesu-tags/services"
	"testing"
)

type linkDeletedHandlerTests struct {
	connection data.DbConnector
	repo       *repositories.LinksRepository
	dbName     string
	handler    *LinkDeletedHandler
}

func newLinkDeletedHandlerTests() *linkDeletedHandlerTests {
	database, name := mocks.GetDatabase()
	repo := repositories.CreateLinksRepository(services.NewGlobalState(database))
	return &linkDeletedHandlerTests{
		connection: database,
		dbName:     name,
		repo:       repo,
		handler:    NewLinkDeletedHandler(repo),
	}
}

func (test *linkDeletedHandlerTests) close() {
	mocks.DropDatabase(test.connection, test.dbName)
}

func TestLinkDeletedHandler_GetTopic(t *testing.T) {
	handler := NewLinkDeletedHandler(nil)

	assert.Equal(t, "links-delete", handler.GetTopic())
}

func TestLinkDeletedHandler_HandleMessage_LinkExists(t *testing.T) {
	test := newLinkDeletedHandlerTests()
	t.Cleanup(func() {
		mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	_ = test.repo.Create(&testLink, &models.UserInfo{UserId: id})
	message := LinkDeletedMessage{LinkId: testLink.Id}
	messageBytes, _ := json.Marshal(message)

	result := test.handler.HandleMessage(messageBytes)

	assert.True(t, result)
	linksRows, _ := test.connection.QueryRows("select * from links")
	assert.False(t, linksRows.Next())
}

func TestLinkDeletedHandler_HandleMessage_LinkDoesntExist(t *testing.T) {
	test := newLinkDeletedHandlerTests()
	t.Cleanup(func() {
		mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	wrongId, _ := uuid.NewV4()
	testLink := models.Link{Id: id}
	_ = test.repo.Create(&testLink, &models.UserInfo{UserId: id})
	message := LinkDeletedMessage{LinkId: wrongId}
	messageBytes, _ := json.Marshal(message)

	result := test.handler.HandleMessage(messageBytes)

	assert.True(t, result)
	linksRows, _ := test.connection.QueryRows("select * from links")
	defer linksRows.Close()
	assert.True(t, linksRows.Next())
}
