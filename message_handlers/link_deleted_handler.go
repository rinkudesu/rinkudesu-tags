package message_handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/repositories"
)

type LinkDeletedHandler struct {
	linksRepository *repositories.LinksRepository
}

func NewLinkDeletedHandler(linksRepository *repositories.LinksRepository) *LinkDeletedHandler {
	return &LinkDeletedHandler{linksRepository: linksRepository}
}

func (handler *LinkDeletedHandler) HandleMessage(message []byte) bool {
	parsedMessage := handler.parseMessage(message)
	if parsedMessage == nil {
		return false
	}
	err := handler.linksRepository.DeleteForce(parsedMessage.LinkId)
	if err != nil && err != repositories.NotFoundErr {
		return false
	}
	return true
}

func (handler *LinkDeletedHandler) parseMessage(message []byte) *LinkDeletedMessage {
	var messageObject LinkDeletedMessage
	err := json.Unmarshal(message, &messageObject)
	if err != nil {
		log.Warningf("Failed to parse link deleted message: %s", err.Error())
		return nil
	}
	return &messageObject
}

func (handler *LinkDeletedHandler) GetTopic() string {
	return LinkDeletedTopic
}
