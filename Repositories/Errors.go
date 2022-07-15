package Repositories

import "errors"

var (
	AlreadyExistsErr = errors.New("tag with this name already exists")
	NotFoundErr      = errors.New("data was not found")
)
