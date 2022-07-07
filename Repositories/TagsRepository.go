package Repositories

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Models"
)

type TagsRepository struct {
	connection Data.DbConnector
}

func NewTagsRepository(connection Data.DbConnector) *TagsRepository {
	return &TagsRepository{connection: connection}
}

func (repository *TagsRepository) GetTags() ([]Models.Tag, error) {
	rows, err := repository.connection.Query("select * from tags")
	if err != nil {
		return nil, err
	}
	tags := make([]Models.Tag, len(rows.RawValues()))
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
	row, err := repository.connection.QueryRow("select name, user_id from tags where id = $1", id)
	if err != nil {
		return nil, err
	}
	var name string
	var userId uuid.UUID
	err = row.Scan(&name, &userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &Models.Tag{Id: id, Name: name, UserId: userId}, nil
}

func (repository *TagsRepository) Create(tag *Models.Tag) (*Models.Tag, error) {
	result, err := repository.connection.QueryRow("insert into tags (name, user_id) values ($1, $2) returning id", tag.Name, tag.UserId)
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
	result, err := repository.connection.Exec("update tags set name = $1, user_id = $2 where id = $3", tag.Name, tag.UserId, tag.Id)
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
	_, err := repository.connection.Exec("call delete_tag($1::uuid);", id)
	return err
}
