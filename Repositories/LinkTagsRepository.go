package Repositories

import (
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Services"
)

type LinkTagsRepository struct {
	connection Data.DbConnector
}

func NewLinkTagsRepository(state *Services.GlobalState) *LinkTagsRepository {
	return &LinkTagsRepository{connection: state.DbConnection}
}

func (repo *LinkTagsRepository) Create(linkTag *Models.LinkTag) error {
	createdIdRow, err := repo.connection.QueryRow("insert into link_tags (link_id, tag_id) values ($1, $2) returning id", linkTag.LinkId, linkTag.TagId)
	if err != nil {
		log.Warningf("Failed to create linkTag: %s", err.Error())
		return err
	}
	var createdId uuid.UUID
	err = createdIdRow.Scan(&createdId)
	if err != nil {
		if IsPostgresDuplicateValue(err) {
			return AlreadyExistsErr
		}
		if IsPostgresNotFoundError(err) {
			return NotFoundErr
		}
		log.Warningf("Failed to scan created linkTag id: %s", err.Error())
		return err
	}
	linkTag.Id = createdId
	return nil
}

func (repo *LinkTagsRepository) Remove(linkId uuid.UUID, tagId uuid.UUID) error {
	result, err := repo.connection.Exec("delete from link_tags where link_id = $1 and tag_id = $2", linkId, tagId)
	if err != nil {
		log.Warningf("Failed to delete link tag: %s", err.Error())
		return err
	}
	if result.RowsAffected() <= 0 {
		return NotFoundErr
	}
	return nil
}

func (repo *LinkTagsRepository) GetLinksForTag(tagId uuid.UUID) (*[]Models.Link, error) {
	linkRows, err := repo.connection.QueryRows("select l.id from link_tags lt join links l on lt.link_id = l.id where lt.tag_id = $1", tagId)
	if err != nil {
		log.Warningf("Failed to query for link_tags: %s", err.Error())
		return nil, err
	}
	defer linkRows.Close()
	links := make([]Models.Link, 0)
	for linkRows.Next() {
		newLink := Models.Link{}
		err = linkRows.Scan(&newLink.Id)
		if err != nil {
			log.Warningf("Failed to scan link: %s", err.Error())
			return nil, err
		}
		links = append(links, newLink)
	}
	return &links, nil
}

func (repo *LinkTagsRepository) GetTagsForLink(linkId uuid.UUID) (*[]Models.Tag, error) {
	tagRows, err := repo.connection.QueryRows("select t.id, t.name, t.user_id from link_tags lt join tags t on lt.tag_id = t.id where lt.link_id = $1", linkId)
	if err != nil {
		log.Warningf("Failed to query for link_tags: %s", err.Error())
		return nil, err
	}
	defer tagRows.Close()
	tags := make([]Models.Tag, 0)
	for tagRows.Next() {
		newTag := Models.Tag{}
		err = tagRows.Scan(&newTag.Id, &newTag.Name, &newTag.UserId)
		if err != nil {
			log.Warningf("Failed to scan link: %s", err.Error())
			return nil, err
		}
		tags = append(tags, newTag)
	}
	return &tags, nil
}
