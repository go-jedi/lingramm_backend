package studiedlanguage

import "time"

type StudiedLanguage struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Lang        string    `json:"lang"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//
// CREATE
//

type CreateDTO struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=1"`
	Lang        string `json:"lang" validate:"required,len=2"`
}

//
// SWAGGER
//

type CreateSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID          int64     `json:"id" example:"1"`
		Name        string    `json:"name" example:"some name"`
		Description string    `json:"description" example:"some description"`
		Lang        string    `json:"lang" example:"eu"`
		CreatedAt   time.Time `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		UpdatedAt   time.Time `json:"updated_at" example:"2025-09-02T15:30:20.095307198+03:00"`
	} `json:"data"`
}

type AllSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		ID          int64     `json:"id" example:"1"`
		Name        string    `json:"name" example:"some name"`
		Description string    `json:"description" example:"some description"`
		Lang        string    `json:"lang" example:"eu"`
		CreatedAt   time.Time `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		UpdatedAt   time.Time `json:"updated_at" example:"2025-09-02T15:30:20.095307198+03:00"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
