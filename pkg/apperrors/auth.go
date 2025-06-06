package apperrors

import "errors"

var (
	ErrNoActiveSessionFound   = errors.New("no active session found")
	ErrTokenMismatchOrExpired = errors.New("token mismatch or expired")
)
