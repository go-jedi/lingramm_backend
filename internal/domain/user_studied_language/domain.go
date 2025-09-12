package userstudiedlanguage

import "time"

type UserStudiedLanguage struct {
	ID                int64     `json:"id"`
	StudiedLanguageID int64     `json:"studied_language_id"`
	TelegramID        string    `json:"telegram_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

//
// CREATE
//

type CreateDTO struct {
	StudiedLanguageID int64  `json:"studied_language_id" validate:"required,gt=0"`
	TelegramID        string `json:"telegram_id" validate:"required,min=1"`
}

//
// UPDATE
//

type UpdateDTO struct {
	StudiedLanguageID int64  `json:"studied_language_id" validate:"required,gt=0"`
	TelegramID        string `json:"telegram_id" validate:"required,min=1"`
}

//
// GET BY TELEGRAM ID
//

type GetByTelegramIDResponse struct {
	ID                int64     `json:"id"`
	StudiedLanguageID int64     `json:"studied_language_id"`
	TelegramID        string    `json:"telegram_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Lang              string    `json:"lang"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

//
// SWAGGER
//

type CreateSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID                int64     `json:"id" example:"1"`
		StudiedLanguageID int64     `json:"studied_language_id" example:"1"`
		TelegramID        string    `json:"telegram_id" example:"1"`
		CreatedAt         time.Time `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		UpdatedAt         time.Time `json:"updated_at" example:"2025-09-02T15:30:20.095307198+03:00"`
	} `json:"data"`
}

type ExistsByTelegramIDSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    bool   `json:"data" example:"true"`
}

type GetByTelegramIDSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID                int64     `json:"id" example:"1"`
		StudiedLanguageID int64     `json:"studied_language_id" example:"1"`
		TelegramID        string    `json:"telegram_id" example:"1"`
		Name              string    `json:"name" example:"some name"`
		Description       string    `json:"description" example:"some description"`
		Lang              string    `json:"lang" example:"en"`
		CreatedAt         time.Time `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		UpdatedAt         time.Time `json:"updated_at" example:"2025-09-02T15:30:20.095307198+03:00"`
	} `json:"data"`
}

type UpdateSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID                int64     `json:"id" example:"1"`
		StudiedLanguageID int64     `json:"studied_language_id" example:"1"`
		TelegramID        string    `json:"telegram_id" example:"1"`
		CreatedAt         time.Time `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		UpdatedAt         time.Time `json:"updated_at" example:"2025-09-02T15:30:20.095307198+03:00"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
