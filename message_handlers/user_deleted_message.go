package message_handlers

import "github.com/gofrs/uuid"

type UserDeletedMessage struct {
	UserId uuid.UUID `json:"user_id"`
}
