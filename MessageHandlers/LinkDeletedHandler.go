package MessageHandlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/Repositories"
)

type LinkDeletedHandler struct {
	linksRepository *Repositories.LinksRepository
}

func NewLinkDeletedHandler(linksRepository *Repositories.LinksRepository) *LinkDeletedHandler {
	return &LinkDeletedHandler{linksRepository: linksRepository}
}

func (handler *LinkDeletedHandler) HandleMessage(message []byte) bool {
	parsedMessage := handler.parseMessage(message)
	if parsedMessage == nil {
		return false
	}
	err := handler.linksRepository.DeleteForce(parsedMessage.LinkId)
	if err != nil && err != Repositories.NotFoundErr {
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
