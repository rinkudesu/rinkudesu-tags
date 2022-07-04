package Repositories

import (
	"errors"
	"github.com/google/uuid"
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

func (repository *TagsRepository) Init(initConnection Data.DbConnection) {
	repository.connection = initConnection
}
