package apperrors

import "errors"

var (
	ErrAchievementAlreadyExists          = errors.New("achievement already exists")
	ErrAchievementDoesNotExist           = errors.New("achievement does not exist")
	ErrAchievementConditionAlreadyExists = errors.New("achievement condition already exists")
)
