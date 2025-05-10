package apperrors

import "errors"

var (
	ErrParamIsRequired = errors.New("parameter is required")
	ErrQueryIsRequired = errors.New("query is required")
)
