package MessageHandlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/Repositories"
)

type UserDeletedHandler struct {
	linksRepository    *Repositories.LinksRepository
	tagsRepository     *Repositories.TagsRepository
	linkTagsRepository *Repositories.LinkTagsRepository
}

func NewUserDeletedHandler(linksRepository *Repositories.LinksRepository, tagsRepository *Repositories.TagsRepository, linkTagsRepository *Repositories.LinkTagsRepository) *UserDeletedHandler {
	return &UserDeletedHandler{linksRepository: linksRepository, tagsRepository: tagsRepository, linkTagsRepository: linkTagsRepository}
}

func (handler *UserDeletedHandler) HandleMessage(message []byte) bool {
	parsedMessage := handler.parseMessage(message)
	if parsedMessage == nil {
		return false
	}
	failed := false
	if err := handler.linksRepository.DeleteForUser(parsedMessage.UserId); err != nil {
		failed = true
	}
	if err := handler.tagsRepository.DeleteAllOfUser(parsedMessage.UserId); err != nil {
		failed = true
	}
	if err := handler.linkTagsRepository.RemoveAllOfUser(parsedMessage.UserId); err != nil {
		failed = true
	}
	return !failed
}

func (handler *UserDeletedHandler) parseMessage(message []byte) *UserDeletedMessage {
	var messageObject UserDeletedMessage
	err := json.Unmarshal(message, &messageObject)
	if err != nil {
		log.Warningf("Failed to parse user deleted message: %s", err.Error())
		return nil
	}
	return &messageObject
}

func (handler *UserDeletedHandler) GetTopic() string {
	return LinkDeletedTopic
}
