package internal

import "errors"

var (
	ErrStorageNotFound      = errors.New("event not found")
	ErrStorageAlreadyExists = errors.New("event already exists")
)
