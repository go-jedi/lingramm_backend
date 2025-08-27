package apperrors

import "errors"

var (
	ErrAchievementTypeDoesNotExist  = errors.New("achievement type does not exist")
	ErrAchievementTypeAlreadyExists = errors.New("achievement type already exists")
)
