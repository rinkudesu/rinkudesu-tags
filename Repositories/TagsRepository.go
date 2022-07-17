package Repositories

import (
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/Models"
)

type TagsRepository struct {
	executor TagQueryExecutable
}

func NewTagsRepository(executor TagQueryExecutable) *TagsRepository {
	return &TagsRepository{executor: executor}
}

func (repository *TagsRepository) GetTags() ([]*Models.Tag, error) {
	rows, err := repository.executor.GetAll()
	defer rows.Close()
	if err != nil {
		log.Warningf("Failed to query for all tags: %s", err.Error())
		return nil, err
	}
	tags := make([]*Models.Tag, 0)
	for rows.Next() {
		var id uuid.UUID
		var name string
		var userId uuid.UUID

		scanErr := rows.Scan(&id, &name, &userId)
		if scanErr != nil {
			log.Warningf("Failed to scan tag: %s", scanErr.Error())
			return nil, scanErr
		}
		tags = append(tags, &Models.Tag{Id: id, Name: name, UserId: userId})
	}
	return tags, nil
}

func (repository *TagsRepository) GetTag(id uuid.UUID) (*Models.Tag, error) {
	row, err := repository.executor.GetSingleById(id)
	if err != nil {
		log.Warningf("Failed to query for tag: %s", err.Error())
		return nil, err
	}
	tag, err := repository.executor.ScanIntoTag(row, id)
	if err != nil {
		if IsPostgresNotFoundError(err) {
			return nil, NotFoundErr
		}
		log.Warningf("Unexpected error when scanning tag: %s", err.Error())
		return nil, err
	}
	return tag, nil
}

func (repository *TagsRepository) Create(tag *Models.Tag) (*Models.Tag, error) {
	result, err := repository.executor.Insert(tag)
	if err != nil {
		log.Warningf("Error when inserting tag: %s", err.Error())
		return nil, err
	}
	var newId uuid.UUID
	err = result.Scan(&newId)
	if err != nil {
		if IsPostgresDuplicateValue(err) {
			return nil, AlreadyExistsErr
		}
		log.Warningf("Unexpected error when scanning inserted id: %s", err.Error())
		return nil, err
	}
	tag.Id = newId
	return tag, nil
}

func (repository *TagsRepository) Update(tag *Models.Tag) (*Models.Tag, error) {
	result, err := repository.executor.Update(tag)
	if err != nil {
		if IsPostgresDuplicateValue(err) {
			return nil, AlreadyExistsErr
		}
		log.Errorf("Unexpected error when updating tag: %s", err.Error())
		return nil, err
	}
	if result.RowsAffected() <= 0 {
		return nil, NotFoundErr
	}
	return tag, nil
}

func (repository *TagsRepository) Delete(id uuid.UUID) error {
	result, err := repository.executor.Delete(id)
	if err != nil {
		log.Warningf("Failed to delete tag: %s", err.Error())
		return err
	}
	if result.RowsAffected() <= 0 {
		return NotFoundErr
	}
	return err
}
