package Repositories

import (
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Services"
)

type LinksRepository struct {
	connection Data.DbConnector
}

func NewLinksRepository(state *Services.GlobalState) *LinksRepository {
	return &LinksRepository{connection: state.DbConnection}
}

func (repo *LinksRepository) Create(link *Models.Link, userInfo *Models.UserInfo) error {
	_, err := repo.connection.Exec("insert into links values ($1, $2)", link.Id, userInfo.UserId)
	if err != nil {
		if IsPostgresDuplicateValue(err) {
			return AlreadyExistsErr
		}
		log.Warningf("Failed to create link: %s", err.Error())
		return err
	}
	return nil
}

func (repo *LinksRepository) Exists(id uuid.UUID, userInfo *Models.UserInfo) (bool, error) {
	result, err := repo.connection.QueryRow("select count(*) from links where id = $1 and user_id = $2", id, userInfo.UserId)
	if err != nil {
		log.Warningf("Failed to count links: %s", err.Error())
		return false, err
	}
	var linkCount int
	err = result.Scan(&linkCount)
	if err != nil {
		log.Warningf("Failed to count links: %s", err.Error())
		return false, err
	}
	return linkCount > 0, nil
}

func (repo *LinksRepository) Delete(id uuid.UUID, userInfo *Models.UserInfo) error {
	result, err := repo.connection.Exec("delete from links where id = $1 and user_id = $2", id, userInfo.UserId)
	if err != nil {
		log.Warningf("Failed to delete link: %s", err.Error())
		return err
	}
	if result.RowsAffected() <= 0 {
		log.Info("Link to delete was not found")
		return NotFoundErr
	}
	return nil
}
