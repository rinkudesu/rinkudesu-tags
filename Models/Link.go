package Models

import "github.com/gofrs/uuid"

type Link struct {
	Id uuid.UUID `json:"id"`
	//userID is not here as it's only ever used on the database-side for validation and there's no need for it to be ever sent by or to the user
}
