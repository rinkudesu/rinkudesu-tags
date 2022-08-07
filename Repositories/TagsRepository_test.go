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

type tagsRepositoryTests struct {
	repo     *TagsRepository
	database Data.DbConnector
	dbName   string
	userInfo *Models.UserInfo
}

func newTagsRepositoryTests() *tagsRepositoryTests {
	database, name := Mocks.GetDatabase()
	globalState := Services.NewGlobalState(database)
	userId, _ := uuid.NewV4()
	return &tagsRepositoryTests{
		repo:     CreateTagsRepository(globalState),
		database: database,
		dbName:   name,
		userInfo: &Models.UserInfo{UserId: userId},
	}
}

func (test *tagsRepositoryTests) close() {
	Mocks.DropDatabase(test.database, test.dbName)
}

func TestTagQueryExecutor_GetAll_TagsPresent(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tags := []*Models.Tag{
		{Name: "tag 1"},
		{Name: "tag 2"},
		{Name: "tag 3"},
	}
	tagIds := test.addTags(t, tags)

	result, err := test.repo.GetTags(test.userInfo)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(result))
	for i := 0; i < 3; i++ {
		assert.Contains(t, tagIds, result[i].Id)
		assert.True(t, test.containsTag(tags[i], result))
	}
}

func TestTagQueryExecutor_GetAll_TagsCreatedByAnotherUser_ReturnsEmptySlice(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tags := []*Models.Tag{
		{Name: "tag 1"},
		{Name: "tag 2"},
		{Name: "tag 3"},
	}
	_ = test.addTags(t, tags)
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := Models.UserInfo{UserId: anotherUserId}

	result, err := test.repo.GetTags(&anotherUserInfo)

	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestTagQueryExecutor_GetAll_NoTagsReturnsEmpty(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)

	result, err := test.repo.GetTags(test.userInfo)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestTagsRepository_GetTag_Found(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	tagId := test.addTag(t, &tag)

	result, err := test.repo.GetTag(tagId, test.userInfo)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tagId, result.Id)
	assert.Equal(t, "test", result.Name)
}

func TestTagsRepository_GetTag_CreatedByAnotherUser_NotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	tagId := test.addTag(t, &tag)
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := Models.UserInfo{UserId: anotherUserId}

	result, err := test.repo.GetTag(tagId, &anotherUserInfo)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	assert.Nil(t, result)
}

func TestTagsRepository_GetTag_NotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	id, _ := uuid.NewV4()

	result, err := test.repo.GetTag(id, test.userInfo)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	assert.Nil(t, result)
}

func TestTagsRepository_Create_Creates(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}

	result, err := test.repo.Create(&tag, test.userInfo)

	assert.Nil(t, err)
	assert.Equal(t, &tag, result)
	assert.NotEqual(t, uuid.Nil, result.Id)
}

func TestTagsRepository_Create_DuplicateName(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	_, _ = test.repo.Create(&tag, test.userInfo)

	result, err := test.repo.Create(&tag, test.userInfo)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, AlreadyExistsErr, err)
	tags, _ := test.repo.GetTags(test.userInfo)
	assert.Equal(t, 1, len(tags))
}

func TestTagsRepository_Create_NameUserByAnotherUser_CreatedAnyway(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	_, _ = test.repo.Create(&tag, test.userInfo)
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := Models.UserInfo{UserId: anotherUserId}

	result, err := test.repo.Create(&tag, &anotherUserInfo)

	assert.Nil(t, err)
	assert.Equal(t, &tag, result)
	assert.NotEqual(t, uuid.Nil, result.Id)
}

func TestTagsRepository_Update_TagExists(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	_, _ = test.repo.Create(&tag, test.userInfo)
	tag.Name = "new name"

	result, err := test.repo.Update(&tag, test.userInfo)

	assert.Nil(t, err)
	assert.Equal(t, &tag, result)
	loadedTag, _ := test.repo.GetTag(tag.Id, test.userInfo)
	assert.Equal(t, "new name", loadedTag.Name)
}

func TestTagsRepository_Update_NameAlreadyExists(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	tag2 := Models.Tag{Name: "new name"}
	_, _ = test.repo.Create(&tag, test.userInfo)
	_, _ = test.repo.Create(&tag2, test.userInfo)
	tag.Name = "new name"

	result, err := test.repo.Update(&tag, test.userInfo)

	assert.NotNil(t, err)
	assert.Equal(t, AlreadyExistsErr, err)
	assert.Nil(t, result)
	loadedTag, _ := test.repo.GetTag(tag.Id, test.userInfo)
	assert.Equal(t, "test", loadedTag.Name)
}

func TestTagsRepository_Update_TagNotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}

	result, err := test.repo.Update(&tag, test.userInfo)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	assert.Nil(t, result)
	loadedTags, _ := test.repo.GetTags(test.userInfo)
	assert.Empty(t, loadedTags)
}

func TestTagsRepository_Update_TagCreatedByAnotherUser_NotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	_, _ = test.repo.Create(&tag, test.userInfo)
	tag.Name = "new name"
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := Models.UserInfo{UserId: anotherUserId}

	result, err := test.repo.Update(&tag, &anotherUserInfo)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	assert.Nil(t, result)
	loadedTags, _ := test.repo.GetTags(test.userInfo)
	assert.Equal(t, "test", loadedTags[0].Name)
}

func TestTagsRepository_Delete_Deletes(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	_, _ = test.repo.Create(&tag, test.userInfo)

	err := test.repo.Delete(tag.Id, test.userInfo)

	assert.Nil(t, err)
	tags, _ := test.repo.GetTags(test.userInfo)
	assert.Empty(t, tags)
}

func TestTagsRepository_Delete_NotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	userId, _ := uuid.NewV4()

	err := test.repo.Delete(userId, test.userInfo)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	tags, _ := test.repo.GetTags(test.userInfo)
	assert.Empty(t, tags)
}

func TestTagsRepository_Delete_CreatedByAnotherUser_NotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	t.Cleanup(test.close)
	tag := Models.Tag{Name: "test"}
	_, _ = test.repo.Create(&tag, test.userInfo)
	anotherUserId, _ := uuid.NewV4()
	anotherUserInfo := Models.UserInfo{UserId: anotherUserId}

	err := test.repo.Delete(tag.Id, &anotherUserInfo)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
}

func (test *tagsRepositoryTests) containsTag(expected *Models.Tag, collection []*Models.Tag) bool {
	for _, tag := range collection {
		if tag.Name == expected.Name {
			return true
		}
	}
	return false
}

func (test *tagsRepositoryTests) addTag(t *testing.T, tag *Models.Tag) uuid.UUID {
	tagArray := []*Models.Tag{tag}
	return test.addTags(t, tagArray)[0]
}

func (test *tagsRepositoryTests) addTags(t *testing.T, tags []*Models.Tag) []uuid.UUID {
	ids := make([]uuid.UUID, len(tags))
	for i, tag := range tags {
		insertResult, err := test.repo.Create(tag, test.userInfo)
		if err != nil {
			t.Fatalf("failed to setup test data")
		}
		ids[i] = insertResult.Id
	}
	return ids
}
