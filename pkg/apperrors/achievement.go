package apperrors

import "errors"

var (
	ErrAchievementAlreadyExists          = errors.New("achievement already exists")
	ErrAchievementConditionAlreadyExists = errors.New("achievement condition already exists")
)
