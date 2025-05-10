package apperrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDoesNotExist  = errors.New("user does not exist")
)
