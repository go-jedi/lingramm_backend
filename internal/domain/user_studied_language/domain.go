package userstudiedlanguage

import "time"

type UserStudiedLanguage struct {
	ID                 int64     `json:"id"`
	StudiedLanguagesID int64     `json:"studied_languages_id"`
	TelegramID         string    `json:"telegram_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

//
// CREATE
//

type CreateDTO struct {
	StudiedLanguagesID int64  `json:"studied_languages_id" validate:"required,gt=0"`
	TelegramID         string `json:"telegram_id" validate:"required,min=1"`
}

//
// UPDATE
//

type UpdateDTO struct {
	StudiedLanguagesID int64  `json:"studied_languages_id" validate:"required,gt=0"`
	TelegramID         string `json:"telegram_id" validate:"required,min=1"`
}

//
// GET BY TELEGRAM ID
//

type GetByTelegramIDResponse struct {
	ID                 int64     `json:"id"`
	StudiedLanguagesID int64     `json:"studied_languages_id"`
	TelegramID         string    `json:"telegram_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Lang               string    `json:"lang"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
