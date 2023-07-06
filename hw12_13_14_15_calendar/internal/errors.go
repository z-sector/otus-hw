package internal

import "errors"

var (
	ErrStorageNotFound = errors.New("event not found")
	ErrStorageConflict = errors.New("state conflict")
)

type ValidationError struct {
	Message string
}

func (v *ValidationError) Error() string {
	return v.Message
}
