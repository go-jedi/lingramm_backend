package apperrors

import "errors"

var (
	ErrEventTypeDoesNotExist  = errors.New("event type does not exist")
	ErrEventTypeAlreadyExists = errors.New("event type already exists")
)
