package apperrors

import "errors"

var (
	ErrTextContentAlreadyExists     = errors.New("text content already exists")
	ErrTextContentDoesNotExist      = errors.New("text content does not exist")
	ErrTextTranslationAlreadyExists = errors.New("text translation already exists")
)
