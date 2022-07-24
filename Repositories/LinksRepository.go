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

func (repo *LinksRepository) Create(link *Models.Link) error {
	_, err := repo.connection.Exec("insert into links values ($1)", link.Id)
	if err != nil {
		if IsPostgresDuplicateValue(err) {
			return AlreadyExistsErr
		}
		log.Warningf("Failed to create link: %s", err.Error())
		return err
	}
	return nil
}

func (repo *LinksRepository) Delete(id uuid.UUID) error {
	result, err := repo.connection.Exec("delete from links where id = $1", id)
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
