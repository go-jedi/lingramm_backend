package admin

import "time"

// Admin represents admin in the system.
type Admin struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TelegramID string    `json:"telegram_id"`
}

//
// SWAGGER
//

type AddAdminUserSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID         int64     `json:"id" example:"1"`
		CreatedAt  time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
		UpdatedAt  time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
		TelegramID string    `json:"telegram_id" example:"1"`
	} `json:"data"`
}

type ExistsSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    bool   `json:"data" example:"true"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
