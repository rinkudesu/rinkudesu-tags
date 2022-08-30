package MessageHandlers

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Mocks"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Repositories"
	"rinkudesu-tags/Services"
	"testing"
)

type linkDeletedHandlerTests struct {
	connection Data.DbConnector
	repo       *Repositories.LinksRepository
	dbName     string
	handler    *LinkDeletedHandler
}

func newLinkDeletedHandlerTests() *linkDeletedHandlerTests {
	database, name := Mocks.GetDatabase()
	repo := Repositories.CreateLinksRepository(Services.NewGlobalState(database))
	return &linkDeletedHandlerTests{
		connection: database,
		dbName:     name,
		repo:       repo,
		handler:    NewLinkDeletedHandler(repo),
	}
}

func (test *linkDeletedHandlerTests) close() {
	Mocks.DropDatabase(test.connection, test.dbName)
}

func TestLinkDeletedHandler_GetTopic(t *testing.T) {
	handler := NewLinkDeletedHandler(nil)

	assert.Equal(t, "links-delete", handler.GetTopic())
}

func TestLinkDeletedHandler_HandleMessage_LinkExists(t *testing.T) {
	test := newLinkDeletedHandlerTests()
	t.Cleanup(func() {
		Mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	testLink := Models.Link{Id: id}
	_ = test.repo.Create(&testLink, &Models.UserInfo{UserId: id})
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
		Mocks.DropDatabase(test.connection, test.dbName)
	})
	id, _ := uuid.NewV4()
	wrongId, _ := uuid.NewV4()
	testLink := Models.Link{Id: id}
	_ = test.repo.Create(&testLink, &Models.UserInfo{UserId: id})
	message := LinkDeletedMessage{LinkId: wrongId}
	messageBytes, _ := json.Marshal(message)

	result := test.handler.HandleMessage(messageBytes)

	assert.True(t, result)
	linksRows, _ := test.connection.QueryRows("select * from links")
	defer linksRows.Close()
	assert.True(t, linksRows.Next())
}
