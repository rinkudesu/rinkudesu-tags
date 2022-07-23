package Repositories

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Mocks"
	"rinkudesu-tags/Models"
	"testing"
)

type linkTagsRepositoryTests struct {
	connection Data.DbConnector
	repo       *LinkTagsRepository
	linkRepo   *LinksRepository
	tagRepo    *TagsRepository
	dbName     string
}

func newLinkTagsRepositoryTests() *linkTagsRepositoryTests {
	database, dbName := Mocks.GetDatabase()
	repo := NewLinkTagsRepository(database)
	return &linkTagsRepositoryTests{
		connection: database,
		repo:       repo,
		linkRepo:   NewLinksRepository(&database),
		tagRepo:    NewTagsRepository(NewTagQueryExecutor(database)),
		dbName:     dbName,
	}
}

func (test *linkTagsRepositoryTests) close() {
	test.connection.Close()
}

func TestLinkTagsRepository_Create_Created(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	defer test.close()
	id, _ := uuid.NewV4()
	link := Models.Link{Id: id}
	tag := Models.Tag{
		Name:   "test",
		UserId: id,
	}
	_ = test.linkRepo.Create(&link)
	_, _ = test.tagRepo.Create(&tag)
	linkTag := Models.LinkTag{LinkId: link.Id, TagId: tag.Id}

	err := test.repo.Create(&linkTag)

	assert.Nil(t, err)
	linkTagRows, _ := test.connection.QueryRows("select * from link_tags")
	defer linkTagRows.Close()
	count := 0
	for linkTagRows.Next() {
		linkTag := Models.LinkTag{}
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
	defer test.close()
	id, _ := uuid.NewV4()
	link := Models.Link{Id: id}
	tag := Models.Tag{
		Name:   "test",
		UserId: id,
	}
	_ = test.linkRepo.Create(&link)
	_, _ = test.tagRepo.Create(&tag)
	linkTag := Models.LinkTag{LinkId: link.Id, TagId: tag.Id}
	_ = test.repo.Create(&linkTag)

	err := test.repo.Create(&linkTag)

	assert.NotNil(t, err)
	assert.Equal(t, AlreadyExistsErr, err)
	linkTagRows, _ := test.connection.QueryRows("select * from link_tags")
	defer linkTagRows.Close()
	count := 0
	for linkTagRows.Next() {
		linkTag := Models.LinkTag{}
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
	defer test.close()
	id, _ := uuid.NewV4()
	link := Models.Link{Id: id}
	tag := Models.Tag{
		Name:   "test",
		UserId: id,
	}
	_, _ = test.tagRepo.Create(&tag)
	linkTag := Models.LinkTag{LinkId: link.Id, TagId: tag.Id}

	err := test.repo.Create(&linkTag)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	linkTagRows, _ := test.connection.QueryRows("select * from link_tags")
	defer linkTagRows.Close()
	assert.False(t, linkTagRows.Next())
}

func TestLinkTagsRepository_Create_TagMissing_Fails(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	defer test.close()
	id, _ := uuid.NewV4()
	link := Models.Link{Id: id}
	tag := Models.Tag{
		Name:   "test",
		UserId: id,
	}
	_ = test.linkRepo.Create(&link)
	linkTag := Models.LinkTag{LinkId: link.Id, TagId: tag.Id}

	err := test.repo.Create(&linkTag)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	linkTagRows, _ := test.connection.QueryRows("select * from link_tags")
	defer linkTagRows.Close()
	assert.False(t, linkTagRows.Next())
}

func TestLinkTagsRepository_Remove_NotFound(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	defer test.close()
	id, _ := uuid.NewV4()

	result := test.repo.Remove(id, id)

	assert.NotNil(t, result)
	assert.Equal(t, NotFoundErr, result)
}

func TestLinkTagsRepository_Remove_FoundAndRemoved(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	defer test.close()
	id, _ := uuid.NewV4()
	link := Models.Link{Id: id}
	tag := Models.Tag{
		Name:   "test",
		UserId: id,
	}
	_ = test.linkRepo.Create(&link)
	_, _ = test.tagRepo.Create(&tag)
	linkTag := Models.LinkTag{LinkId: link.Id, TagId: tag.Id}
	_ = test.repo.Create(&linkTag)

	result := test.repo.Remove(link.Id, tag.Id)

	assert.Nil(t, result)
	linkTagRows, _ := test.connection.QueryRows("select * from link_tags")
	defer linkTagRows.Close()
	assert.False(t, linkTagRows.Next())
}

func TestLinkTagsRepository_GetLinksForTag_TagIdNotFound(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	defer test.close()
	id, _ := uuid.NewV4()

	result, err := test.repo.GetLinksForTag(id)

	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestLinkTagsRepository_GetLinksForTag_LinksArrayReturned(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	defer test.close()
	const linksCount = 5
	createdLinks := make([]Models.Link, linksCount)
	userId, _ := uuid.NewV4()
	tag := Models.Tag{Name: "test tag", UserId: userId}
	_, _ = test.tagRepo.Create(&tag)
	for i := 0; i < linksCount; i++ {
		id, _ := uuid.NewV4()
		link := Models.Link{Id: id}
		createdLinks[i] = link
		assert.Nil(t, test.linkRepo.Create(&link))
		assert.Nil(t, test.repo.Create(&Models.LinkTag{LinkId: link.Id, TagId: tag.Id}))
	}

	links, err := test.repo.GetLinksForTag(tag.Id)

	assert.Nil(t, err)
	assert.NotNil(t, links)
	assert.Equal(t, linksCount, len(*links))
	for _, link := range createdLinks {
		assert.Contains(t, *links, link)
	}
}

func TestLinkTagsRepository_GetTagsForLink_LinkIdNotFound(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	defer test.close()
	id, _ := uuid.NewV4()

	result, err := test.repo.GetTagsForLink(id)

	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestLinkTagsRepository_GetTagsForLink_TagsArrayReturned(t *testing.T) {
	test := newLinkTagsRepositoryTests()
	defer test.close()
	const tagsCount = 5
	createdTags := make([]Models.Tag, tagsCount)
	linkId, _ := uuid.NewV4()
	link := Models.Link{Id: linkId}
	assert.Nil(t, test.linkRepo.Create(&link))
	userId, _ := uuid.NewV4()
	for i := 0; i < tagsCount; i++ {
		tag := Models.Tag{Name: fmt.Sprintf("test tag %d", i), UserId: userId}
		_, _ = test.tagRepo.Create(&tag)
		createdTags[i] = tag
		assert.Nil(t, test.repo.Create(&Models.LinkTag{LinkId: link.Id, TagId: tag.Id}))
	}

	tags, err := test.repo.GetTagsForLink(linkId)

	assert.Nil(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, tagsCount, len(*tags))
	for _, tag := range createdTags {
		assert.Contains(t, *tags, tag)
	}
}
