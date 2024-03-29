package repositories

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"rinkudesu-tags/data"
	"rinkudesu-tags/mocks"
	"rinkudesu-tags/models"
	"rinkudesu-tags/services"
	"testing"
)

type linkTagsRepositoryTests struct {
	connection data.DbConnector
	repo       *LinkTagsRepository
	linkRepo   *LinksRepository
	tagRepo    *TagsRepository
	dbName     string
	userInfo   *models.UserInfo
}

func newLinkTagsRepositoryTests() *linkTagsRepositoryTests {
	database, dbName := mocks.GetDatabase()
	globalState := services.NewGlobalState(database)
	repo := NewLinkTagsRepository(globalState)
	userId, _ := uuid.NewV4()
	return &linkTagsRepositoryTests{
		connection: database,
		repo:       repo,
		linkRepo:   CreateLinksRepository(globalState),
		tagRepo:    NewTagsRepository(globalState),
		dbName:     dbName,
		userInfo:   &models.UserInfo{UserId: userId},
	}
}

func (test *linkTagsRepositoryTests) close() {
	mocks.DropDatabase(test.connection, test.dbName)
}

func TestLinkTagsRepository_Create_Created(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	link := models.Link{Id: id}
	tag := models.Tag{
		Name: "test",
	}
	_ = test.linkRepo.Create(&link, test.userInfo)
	_, _ = test.tagRepo.Create(&tag, test.userInfo)
	linkTag := models.LinkTag{LinkId: link.Id, TagId: tag.Id}

	err := test.repo.Create(&linkTag, test.userInfo)

	assert.Nil(t, err)
	linkTagRows, _ := test.connection.QueryRows("select id, link_id, tag_id from link_tags")
	defer linkTagRows.Close()
	count := 0
	for linkTagRows.Next() {
		linkTag := models.LinkTag{}
		_ = linkTagRows.Scan(&linkTag.Id, &linkTag.LinkId, &linkTag.TagId)
		assert.NotEqual(t, uuid.Nil, linkTag.Id)
		assert.Equal(t, link.Id, linkTag.LinkId)
		assert.Equal(t, tag.Id, linkTag.TagId)
		count++
	}
	assert.Equal(t, 1, count)
}

func TestLinkTagsRepository_Create_PairAlreadyExists_Fails(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	link := models.Link{Id: id}
	tag := models.Tag{
		Name: "test",
	}
	_ = test.linkRepo.Create(&link, test.userInfo)
	_, _ = test.tagRepo.Create(&tag, test.userInfo)
	linkTag := models.LinkTag{LinkId: link.Id, TagId: tag.Id}
	_ = test.repo.Create(&linkTag, test.userInfo)

	err := test.repo.Create(&linkTag, test.userInfo)

	assert.NotNil(t, err)
	assert.Equal(t, AlreadyExistsErr, err)
	linkTagRows, _ := test.connection.QueryRows("select id, link_id, tag_id from link_tags")
	defer linkTagRows.Close()
	count := 0
	for linkTagRows.Next() {
		linkTag := models.LinkTag{}
		_ = linkTagRows.Scan(&linkTag.Id, &linkTag.LinkId, &linkTag.TagId)
		assert.NotEqual(t, uuid.Nil, linkTag.Id)
		assert.Equal(t, link.Id, linkTag.LinkId)
		assert.Equal(t, tag.Id, linkTag.TagId)
		count++
	}
	assert.Equal(t, 1, count)
}

func TestLinkTagsRepository_Create_LinkMissing_Fails(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	link := models.Link{Id: id}
	tag := models.Tag{
		Name: "test",
	}
	_, _ = test.tagRepo.Create(&tag, test.userInfo)
	linkTag := models.LinkTag{LinkId: link.Id, TagId: tag.Id}

	err := test.repo.Create(&linkTag, test.userInfo)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	linkTagRows, _ := test.connection.QueryRows("select * from link_tags")
	defer linkTagRows.Close()
	assert.False(t, linkTagRows.Next())
}

func TestLinkTagsRepository_Create_TagMissing_Fails(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	link := models.Link{Id: id}
	tag := models.Tag{
		Name: "test",
	}
	_ = test.linkRepo.Create(&link, test.userInfo)
	linkTag := models.LinkTag{LinkId: link.Id, TagId: tag.Id}

	err := test.repo.Create(&linkTag, test.userInfo)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	linkTagRows, _ := test.connection.QueryRows("select * from link_tags")
	defer linkTagRows.Close()
	assert.False(t, linkTagRows.Next())
}

func TestLinkTagsRepository_Remove_NotFound(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()

	result := test.repo.Remove(id, id, test.userInfo)

	assert.NotNil(t, result)
	assert.Equal(t, NotFoundErr, result)
}

func TestLinkTagsRepository_Remove_CreatedByAnotherUser_NotFound(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := models.UserInfo{UserId: anotherUserId}
	link := models.Link{Id: id}
	tag := models.Tag{
		Name: "test",
	}
	_ = test.linkRepo.Create(&link, test.userInfo)
	_, _ = test.tagRepo.Create(&tag, test.userInfo)
	linkTag := models.LinkTag{LinkId: link.Id, TagId: tag.Id}
	_ = test.repo.Create(&linkTag, test.userInfo)

	result := test.repo.Remove(id, id, &anotherUserInfo)

	assert.NotNil(t, result)
	assert.Equal(t, NotFoundErr, result)
}

func TestLinkTagsRepository_Remove_FoundAndRemoved(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	link := models.Link{Id: id}
	tag := models.Tag{
		Name: "test",
	}
	_ = test.linkRepo.Create(&link, test.userInfo)
	_, _ = test.tagRepo.Create(&tag, test.userInfo)
	linkTag := models.LinkTag{LinkId: link.Id, TagId: tag.Id}
	_ = test.repo.Create(&linkTag, test.userInfo)

	result := test.repo.Remove(link.Id, tag.Id, test.userInfo)

	assert.Nil(t, result)
	linkTagRows, _ := test.connection.QueryRows("select * from link_tags")
	defer linkTagRows.Close()
	assert.False(t, linkTagRows.Next())
}

func TestLinkTagsRepository_GetLinksForTag_TagIdNotFound(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()

	result, err := test.repo.GetLinksForTag(id, test.userInfo)

	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestLinkTagsRepository_GetLinksForTag_LinksArrayReturned(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	const linksCount = 5
	createdLinks := make([]models.Link, linksCount)
	tag := models.Tag{Name: "test tag"}
	_, _ = test.tagRepo.Create(&tag, test.userInfo)
	for i := 0; i < linksCount; i++ {
		id, _ := uuid.NewV4()
		link := models.Link{Id: id}
		createdLinks[i] = link
		assert.Nil(t, test.linkRepo.Create(&link, test.userInfo))
		assert.Nil(t, test.repo.Create(&models.LinkTag{LinkId: link.Id, TagId: tag.Id}, test.userInfo))
	}

	links, err := test.repo.GetLinksForTag(tag.Id, test.userInfo)

	assert.Nil(t, err)
	assert.NotNil(t, links)
	assert.Equal(t, linksCount, len(*links))
	for _, link := range createdLinks {
		assert.Contains(t, *links, link)
	}
}

func TestLinkTagsRepository_GetTagsForLink_LinkIdNotFound(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()

	result, err := test.repo.GetTagsForLink(id, test.userInfo)

	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestLinkTagsRepository_GetTagsForLink_TagsArrayReturned(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	const tagsCount = 5
	createdTags := make([]models.Tag, tagsCount)
	linkId, _ := uuid.NewV4()
	link := models.Link{Id: linkId}
	assert.Nil(t, test.linkRepo.Create(&link, test.userInfo))
	for i := 0; i < tagsCount; i++ {
		tag := models.Tag{Name: fmt.Sprintf("test tag %d", i)}
		_, _ = test.tagRepo.Create(&tag, test.userInfo)
		createdTags[i] = tag
		assert.Nil(t, test.repo.Create(&models.LinkTag{LinkId: link.Id, TagId: tag.Id}, test.userInfo))
	}

	tags, err := test.repo.GetTagsForLink(linkId, test.userInfo)

	assert.Nil(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, tagsCount, len(*tags))
	for _, tag := range createdTags {
		assert.Contains(t, *tags, tag)
	}
}

func TestLinkTagsRepository_RemoveAllOfUser_NoneExistForUser(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := models.UserInfo{UserId: anotherUserId}
	link := models.Link{Id: id}
	tag := models.Tag{
		Name: "test",
	}
	_ = test.linkRepo.Create(&link, test.userInfo)
	_, _ = test.tagRepo.Create(&tag, test.userInfo)
	linkTag := models.LinkTag{LinkId: link.Id, TagId: tag.Id}
	_ = test.repo.Create(&linkTag, test.userInfo)

	err := test.repo.RemoveAllOfUser(anotherUserInfo.UserId)

	assert.Nil(t, err)
	existing, _ := test.repo.GetTagsForLink(link.Id, test.userInfo)
	assert.NotEmpty(t, existing)
}

func TestLinkTagsRepository_RemoveAllOfUser_ExistForUser(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()
	link := models.Link{Id: id}
	tag := models.Tag{
		Name: "test",
	}
	_ = test.linkRepo.Create(&link, test.userInfo)
	_, _ = test.tagRepo.Create(&tag, test.userInfo)
	linkTag := models.LinkTag{LinkId: link.Id, TagId: tag.Id}
	_ = test.repo.Create(&linkTag, test.userInfo)

	err := test.repo.RemoveAllOfUser(test.userInfo.UserId)

	assert.Nil(t, err)
	existing, _ := test.repo.GetTagsForLink(link.Id, test.userInfo)
	assert.Empty(t, existing)
}
