package repositories

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/rinkudesu/curry"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/data"
	"rinkudesu-tags/models"
	"rinkudesu-tags/services"
	"strings"
)

type TagsRepository struct {
	connection data.DbConnector
}

func NewTagsRepository(state *services.GlobalState) *TagsRepository {
	return &TagsRepository{connection: state.DbConnection}
}

func (repository *TagsRepository) GetTags(userInfo *models.UserInfo, name string) ([]*models.Tag, error) {
	userIdFilter := curry.NewWhereItem("user_id", "=", curry.NewParameter(userInfo.UserId))
	nameLike := fmt.Sprintf("%%%s%%", strings.ToUpper(name))
	tagNameFilter := curry.NewWhereItem("name_normalised", "like", curry.NewOptionalParameter(nameLike, "'%%'"))
	query := curry.Select("id, name", "tags", "").Where(curry.WhereBegin(userIdFilter).And(tagNameFilter))
	sql, parameters, err := query.ToExecutable()
	if err != nil {
		log.Warnf("Failed to generate sql query: %s", err.Error())
	}

	rows, err := repository.connection.QueryRows(sql, parameters...)
	defer rows.Close()
	if err != nil {
		log.Warningf("Failed to query for all tags: %s", err.Error())
		return nil, err
	}
	tags := make([]*models.Tag, 0)
	for rows.Next() {
		var id uuid.UUID
		var name string

		scanErr := rows.Scan(&id, &name)
		if scanErr != nil {
			log.Warningf("Failed to scan tag: %s", scanErr.Error())
			return nil, scanErr
		}
		tags = append(tags, &models.Tag{Id: id, Name: name})
	}
	return tags, nil
}

func (repository *TagsRepository) GetTag(id uuid.UUID, userInfo *models.UserInfo) (*models.Tag, error) {
	row, err := repository.connection.QueryRow("select name from tags where id = $1 and user_id = $2", id, userInfo.UserId)
	if err != nil {
		log.Warningf("Failed to query for tag: %s", err.Error())
		return nil, err
	}
	tag, err := repository.scanIntoTag(row, id)
	if err != nil {
		if IsPostgresNotFoundError(err) {
			return nil, NotFoundErr
		}
		log.Warningf("Unexpected error when scanning tag: %s", err.Error())
		return nil, err
	}
	return tag, nil
}

func (repository *TagsRepository) Create(tag *models.Tag, userInfo *models.UserInfo) (*models.Tag, error) {
	result, err := repository.connection.QueryRow("insert into tags (name, user_id) values ($1, $2) returning id", tag.Name, userInfo.UserId)
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

func (repository *TagsRepository) Update(tag *models.Tag, userInfo *models.UserInfo) (*models.Tag, error) {
	result, err := repository.connection.Exec("update tags set name = $1 where id = $3 and user_id = $2", tag.Name, userInfo.UserId, tag.Id)
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

func (repository *TagsRepository) Delete(id uuid.UUID, userInfo *models.UserInfo) error {
	result, err := repository.connection.Exec("delete from tags where id = $1 and user_id = $2", id, userInfo.UserId)
	if err != nil {
		log.Warningf("Failed to delete tag: %s", err.Error())
		return err
	}
	if result.RowsAffected() <= 0 {
		return NotFoundErr
	}
	return err
}

func (repository *TagsRepository) DeleteAllOfUser(userId uuid.UUID) error {
	_, err := repository.connection.Exec("delete from tags where user_id = $1", userId)
	if err != nil {
		log.Warningf("Failed to delete all tags for user: %s", err.Error())
		return err
	}
	return nil
}

func (repository *TagsRepository) Exists(id uuid.UUID, userInfo *models.UserInfo) (bool, error) {
	result, err := repository.connection.QueryRow("select count(*) from tags where id = $1 and user_id = $2", id, userInfo.UserId)
	if err != nil {
		log.Warningf("Failed to count tags: %s", err.Error())
		return false, err
	}
	var tagCount int
	err = result.Scan(&tagCount)
	if err != nil {
		log.Warningf("Failed to count tags: %s", err.Error())
		return false, err
	}
	return tagCount > 0, nil
}

func (repository *TagsRepository) scanIntoTag(row data.Row, id uuid.UUID) (*models.Tag, error) {
	var name string
	err := row.Scan(&name)
	if err != nil {
		return nil, err
	}
	return &models.Tag{Id: id, Name: name}, nil
}
