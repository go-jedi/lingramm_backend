package apperrors

import "errors"

var (
	ErrNoRowsWereAffected      = errors.New("no rows were affected")
	ErrMoreThanOneRowsAffected = errors.New("more than one rows were affected")
	ErrRequestTimedOut         = errors.New("request timed out")
)
