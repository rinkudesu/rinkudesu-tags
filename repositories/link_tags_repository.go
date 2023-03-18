package repositories

import (
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/data"
	"rinkudesu-tags/models"
	"rinkudesu-tags/services"
)

type LinkTagsRepository struct {
	connection data.DbConnector
}

func NewLinkTagsRepository(state *services.GlobalState) *LinkTagsRepository {
	return &LinkTagsRepository{connection: state.DbConnection}
}

// Create a new link-tag association. The calling method is responsible for making sure that both entries exist and that the user is authorised to create it.
func (repo *LinkTagsRepository) Create(linkTag *models.LinkTag, userInfo *models.UserInfo) error {
	createdIdRow, err := repo.connection.QueryRow("insert into link_tags (link_id, tag_id, user_id) values ($1, $2, $3) returning id", linkTag.LinkId, linkTag.TagId, userInfo.UserId)
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

func (repo *LinkTagsRepository) Remove(linkId uuid.UUID, tagId uuid.UUID, userInfo *models.UserInfo) error {
	result, err := repo.connection.Exec("delete from link_tags where link_id = $1 and tag_id = $2 and user_id = $3", linkId, tagId, userInfo.UserId)
	if err != nil {
		log.Warningf("Failed to delete link tag: %s", err.Error())
		return err
	}
	if result.RowsAffected() <= 0 {
		return NotFoundErr
	}
	return nil
}

func (repo *LinkTagsRepository) RemoveAllOfUser(userId uuid.UUID) error {
	_, err := repo.connection.Exec("delete from link_tags where user_id = $1", userId)
	if err != nil {
		log.Warningf("Failed to delete linktags for user: %s", err.Error())
		return err
	}
	return nil
}

func (repo *LinkTagsRepository) GetLinksForTag(tagId uuid.UUID, userInfo *models.UserInfo) (*[]models.Link, error) {
	linkRows, err := repo.connection.QueryRows("select l.id from link_tags lt join links l on lt.link_id = l.id where lt.tag_id = $1 and l.user_id = $2 and lt.user_id = $2", tagId, userInfo.UserId)
	if err != nil {
		log.Warningf("Failed to query for link_tags: %s", err.Error())
		return nil, err
	}
	defer linkRows.Close()
	links := make([]models.Link, 0)
	for linkRows.Next() {
		newLink := models.Link{}
		err = linkRows.Scan(&newLink.Id)
		if err != nil {
			log.Warningf("Failed to scan link: %s", err.Error())
			return nil, err
		}
		links = append(links, newLink)
	}
	return &links, nil
}

func (repo *LinkTagsRepository) GetTagsForLink(linkId uuid.UUID, userInfo *models.UserInfo) (*[]models.Tag, error) {
	tagRows, err := repo.connection.QueryRows("select t.id, t.name, t.colour from link_tags lt join tags t on lt.tag_id = t.id where lt.link_id = $1 and t.user_id = $2 and lt.user_id = $2", linkId, userInfo.UserId)
	if err != nil {
		log.Warningf("Failed to query for link_tags: %s", err.Error())
		return nil, err
	}
	defer tagRows.Close()
	tags := make([]models.Tag, 0)
	for tagRows.Next() {
		newTag := models.Tag{}
		err = tagRows.Scan(&newTag.Id, &newTag.Name, &newTag.Colour)
		if err != nil {
			log.Warningf("Failed to scan link: %s", err.Error())
			return nil, err
		}
		tags = append(tags, newTag)
	}
	return &tags, nil
}
