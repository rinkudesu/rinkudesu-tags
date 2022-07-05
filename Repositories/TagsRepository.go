package Repositories

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Models"
)

type TagsRepository struct {
	connection Data.DbConnection
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

func (repository *TagsRepository) Init(initConnection Data.DbConnection) {
	repository.connection = initConnection
}
