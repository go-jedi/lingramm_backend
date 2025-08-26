package apperrors

import "errors"

var (
	ErrUserStudiedLanguageDoesNotExist  = errors.New("user studied language does not exist")
	ErrUserStudiedLanguageAlreadyExists = errors.New("user studied language already exists")
)
