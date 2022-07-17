package Repositories

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Mocks"
	"rinkudesu-tags/Models"
	"testing"
)

type tagsRepositoryTests struct {
	repo     *TagsRepository
	executor TagQueryExecutable
	database *Data.DbConnector
	dbName   string
}

func newTagsRepositoryTests() *tagsRepositoryTests {
	database, name := Mocks.GetDatabase()
	executor := NewTagQueryExecutor(database)
	return &tagsRepositoryTests{
		repo:     NewTagsRepository(executor),
		executor: executor,
		database: &database,
		dbName:   name,
	}
}

func TestTagQueryExecutor_GetAll_TagsPresent(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()
	tags := []*Models.Tag{
		{Name: "tag 1", UserId: userId},
		{Name: "tag 2", UserId: userId},
		{Name: "tag 3", UserId: userId},
	}
	tagIds := addTags(test.executor, t, tags)

	result, err := test.repo.GetTags()

	assert.Nil(t, err)
	assert.Equal(t, 3, len(result))
	for i := 0; i < 3; i++ {
		assert.Contains(t, tagIds, result[i].Id)
		assert.True(t, containsTag(tags[i], result))
	}
}

func TestTagQueryExecutor_GetAll_NoTagsReturnsEmpty(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)

	result, err := test.repo.GetTags()

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestTagsRepository_GetTag_Found(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()
	tag := Models.Tag{Name: "test", UserId: userId}
	tagId := addTag(test.executor, t, &tag)

	result, err := test.repo.GetTag(tagId)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tagId, result.Id)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, userId, result.UserId)
}

func TestTagsRepository_GetTag_NotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	id, _ := uuid.NewV4()

	result, err := test.repo.GetTag(id)
	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	assert.Nil(t, result)
}

func TestTagsRepository_Create_Creates(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()
	tag := Models.Tag{Name: "test", UserId: userId}

	result, err := test.repo.Create(&tag)

	assert.Nil(t, err)
	assert.Equal(t, &tag, result)
	assert.NotEqual(t, uuid.Nil, result.Id)
}

func TestTagsRepository_Create_DuplicateName(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()
	tag := Models.Tag{Name: "test", UserId: userId}
	_, _ = test.repo.Create(&tag)

	result, err := test.repo.Create(&tag)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, AlreadyExistsErr, err)
	tags, _ := test.repo.GetTags()
	assert.Equal(t, 1, len(tags))
}

func TestTagsRepository_Update_TagExists(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()
	tag := Models.Tag{Name: "test", UserId: userId}
	_, _ = test.repo.Create(&tag)
	tag.Name = "new name"

	result, err := test.repo.Update(&tag)

	assert.Nil(t, err)
	assert.Equal(t, &tag, result)
	loadedTag, _ := test.repo.GetTag(tag.Id)
	assert.Equal(t, "new name", loadedTag.Name)
}

func TestTagsRepository_Update_NameAlreadyExists(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()
	tag := Models.Tag{Name: "test", UserId: userId}
	tag2 := Models.Tag{Name: "new name", UserId: userId}
	_, _ = test.repo.Create(&tag)
	_, _ = test.repo.Create(&tag2)
	tag.Name = "new name"

	result, err := test.repo.Update(&tag)

	assert.NotNil(t, err)
	assert.Equal(t, AlreadyExistsErr, err)
	assert.Nil(t, result)
	loadedTag, _ := test.repo.GetTag(tag.Id)
	assert.Equal(t, "test", loadedTag.Name)
}

func TestTagsRepository_Update_TagNotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()
	tag := Models.Tag{Name: "test", UserId: userId}

	result, err := test.repo.Update(&tag)
	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	assert.Nil(t, result)
	loadedTags, _ := test.repo.GetTags()
	assert.Empty(t, loadedTags)
}

func TestTagsRepository_Delete_Deletes(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()
	tag := Models.Tag{Name: "test", UserId: userId}
	_, _ = test.repo.Create(&tag)

	err := test.repo.Delete(tag.Id)

	assert.Nil(t, err)
	tags, _ := test.repo.GetTags()
	assert.Empty(t, tags)
}

func TestTagsRepository_Delete_NotFound(t *testing.T) {
	test := newTagsRepositoryTests()
	defer Mocks.DropDatabase(test.database, test.dbName)
	userId, _ := uuid.NewV4()

	err := test.repo.Delete(userId)

	assert.NotNil(t, err)
	assert.Equal(t, NotFoundErr, err)
	tags, _ := test.repo.GetTags()
	assert.Empty(t, tags)
}

func containsTag(expected *Models.Tag, collection []*Models.Tag) bool {
	for _, tag := range collection {
		if tag.Name == expected.Name && tag.UserId == expected.UserId {
			return true
		}
	}
	return false
}

func addTag(executor TagQueryExecutable, t *testing.T, tag *Models.Tag) uuid.UUID {
	tagArray := []*Models.Tag{tag}
	return addTags(executor, t, tagArray)[0]
}

func addTags(executor TagQueryExecutable, t *testing.T, tags []*Models.Tag) []uuid.UUID {
	ids := make([]uuid.UUID, len(tags))
	for i, tag := range tags {
		insertResult, err := executor.Insert(tag)
		if err != nil {
			t.Fatalf("failed to setup test data")
		}
		var insertedId uuid.UUID
		err = insertResult.Scan(&insertedId)
		if err != nil {
			t.Fatalf("failed to setup test data")
		}
		ids[i] = insertedId
	}
	return ids
}
