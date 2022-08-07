package Repositories

import "errors"

var (
	AlreadyExistsErr = errors.New("this data already exists")
	NotFoundErr      = errors.New("data was not found")
)
