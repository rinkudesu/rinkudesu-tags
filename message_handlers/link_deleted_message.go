package message_handlers

import "github.com/gofrs/uuid"

type LinkDeletedMessage struct {
	LinkId uuid.UUID `json:"link_id"`
}
