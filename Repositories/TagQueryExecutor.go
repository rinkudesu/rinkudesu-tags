package Repositories

import (
	"github.com/gofrs/uuid"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Models"
)

type TagQueryExecutor struct {
	connection Data.DbConnector
}

func NewTagQueryExecutor(connection Data.DbConnector) *TagQueryExecutor {
	return &TagQueryExecutor{connection: connection}
}

func (executor TagQueryExecutor) GetAll() (Data.Rows, error) {
	return executor.connection.Query("select * from tags")
}

func (executor TagQueryExecutor) GetSingleById(id uuid.UUID) (Data.Row, error) {
	return executor.connection.QueryRow("select name, user_id from tags where id = $1", id)
}

func (executor TagQueryExecutor) Insert(tag *Models.Tag) (Data.Row, error) {
	return executor.connection.QueryRow("insert into tags (name, user_id) values ($1, $2) returning id", tag.Name, tag.UserId)
}

func (executor TagQueryExecutor) Update(tag *Models.Tag) (Data.ExecResult, error) {
	return executor.connection.Exec("update tags set name = $1, user_id = $2 where id = $3", tag.Name, tag.UserId, tag.Id)
}

func (executor TagQueryExecutor) Delete(id uuid.UUID) error {
	_, err := executor.connection.Exec("delete from tags where id = $1", id)
	return err
}

func (executor TagQueryExecutor) ScanIntoTag(row Data.Row, id uuid.UUID) (*Models.Tag, error) {
	var name string
	var userId uuid.UUID
	err := row.Scan(&name, &userId)
	if err != nil {
		return nil, err
	}
	return &Models.Tag{Id: id, Name: name, UserId: userId}, nil
}
