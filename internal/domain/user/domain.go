package user

import "time"

// User represents a user in the system.
type User struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TelegramID string    `json:"telegram_id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
}

//
// CREATE
//

// CreateDTO represents the data required to create a new user.
type CreateDTO struct {
	TelegramID string `json:"telegram_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}

//
// SWAGGER
//

type CreateDailyTaskSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID         int64     `json:"id" example:"1"`
		CreatedAt  time.Time `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		UpdatedAt  time.Time `json:"updated_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		TelegramID string    `json:"telegram_id" example:"1"`
		Username   string    `json:"username" example:"some username"`
		FirstName  string    `json:"first_name" example:"some first name"`
		LastName   string    `json:"last_name" example:"some last name"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
