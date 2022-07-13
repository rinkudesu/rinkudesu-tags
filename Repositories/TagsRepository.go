package Repositories

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"rinkudesu-tags/Models"
)

type TagsRepository struct {
	executor TagQueryExecutable
}

func NewTagsRepository(executor TagQueryExecutable) *TagsRepository {
	return &TagsRepository{executor: executor}
}

func (repository *TagsRepository) GetTags() ([]Models.Tag, error) {
	rows, err := repository.executor.GetAll()
	if err != nil {
		return nil, err
	}
	tags := make([]Models.Tag, 0)
	for rows.Next() {
		var id uuid.UUID
		var name string
		var userId uuid.UUID

		scanErr := rows.Scan(&id, &name, &userId)
		if scanErr != nil {
			return nil, errors.New("failed to scan tag from database")
		}
		tags = append(tags, Models.Tag{Id: id, Name: name, UserId: userId})
	}
	return tags, nil
}

func (repository *TagsRepository) GetTag(id uuid.UUID) (*Models.Tag, error) {
	row, err := repository.executor.GetSingleById(id)
	if err != nil {
		return nil, err
	}
	tag, err := repository.executor.ScanIntoTag(row, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
	}
	return tag, nil
}

func (repository *TagsRepository) Create(tag *Models.Tag) (*Models.Tag, error) {
	result, err := repository.executor.Insert(tag)
	if err != nil {
		return nil, err
	}
	var newId uuid.UUID
	err = result.Scan(&newId)
	if err != nil {
		return nil, err //todo: figure out what to do when name/user is duplicated
	}
	tag.Id = newId
	return tag, nil
}

func (repository *TagsRepository) Update(tag *Models.Tag) (*Models.Tag, error) {
	result, err := repository.executor.Update(tag)
	//todo: figure out what to do when name duplicated
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() <= 0 {
		return nil, nil
	}
	return tag, nil
}

func (repository *TagsRepository) Delete(id uuid.UUID) error {
	return repository.executor.Delete(id)
}
