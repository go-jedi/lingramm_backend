package eventtype

import (
	"time"

	"github.com/shopspring/decimal"
)

type EventType struct {
	ID                  int64            `json:"id"`
	Name                string           `json:"name"`
	Description         *string          `json:"description,omitempty"`
	XP                  int64            `json:"xp"`
	Amount              *decimal.Decimal `json:"amount,omitempty"`
	NotificationMessage *string          `json:"notification_message,omitempty"`
	IsSendNotification  bool             `json:"is_send_notification"`
	IsActive            bool             `json:"is_active"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

//
// CREATE
//

type CreateDTO struct {
	XP                  int64            `json:"xp" validate:"required,gt=0"`
	Name                string           `json:"name" validate:"required,min=1"`
	NotificationMessage *string          `json:"notification_message,omitempty" validate:"omitempty,min=1"`
	Description         *string          `json:"description,omitempty" validate:"omitempty,min=1"`
	Amount              *decimal.Decimal `json:"amount,omitempty" validate:"omitempty"`
	IsSendNotification  bool             `json:"is_send_notification"`
	IsActive            bool             `json:"is_active"`
}

//
// SWAGGER
//

type CreateSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID                  int64     `json:"id" example:"1"`
		Name                string    `json:"name" example:"some name"`
		Description         *string   `json:"description,omitempty" example:"some description"`
		XP                  int64     `json:"xp" example:"20"`
		Amount              *int64    `json:"amount,omitempty" example:"10"`
		NotificationMessage *string   `json:"notification_message,omitempty" example:"some message"`
		IsSendNotification  bool      `json:"is_send_notification" example:"true"`
		IsActive            bool      `json:"is_active" example:"true"`
		CreatedAt           time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
		UpdatedAt           time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
	} `json:"data"`
}

type AllSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		ID                  int64     `json:"id" example:"1"`
		Name                string    `json:"name" example:"some name"`
		Description         *string   `json:"description,omitempty" example:"some description"`
		XP                  int64     `json:"xp" example:"20"`
		Amount              *int64    `json:"amount,omitempty" example:"10"`
		NotificationMessage *string   `json:"notification_message,omitempty" example:"some message"`
		IsSendNotification  bool      `json:"is_send_notification" example:"true"`
		IsActive            bool      `json:"is_active" example:"true"`
		CreatedAt           time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
		UpdatedAt           time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
	} `json:"data"`
}

type GetByNameSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		ID                  int64     `json:"id" example:"1"`
		Name                string    `json:"name" example:"some name"`
		Description         *string   `json:"description,omitempty" example:"some description"`
		XP                  int64     `json:"xp" example:"20"`
		Amount              *int64    `json:"amount,omitempty" example:"10"`
		NotificationMessage *string   `json:"notification_message,omitempty" example:"some message"`
		IsSendNotification  bool      `json:"is_send_notification" example:"true"`
		IsActive            bool      `json:"is_active" example:"true"`
		CreatedAt           time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
		UpdatedAt           time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
