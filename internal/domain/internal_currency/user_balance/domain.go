package userbalance

import (
	"time"

	"github.com/shopspring/decimal"
)

// UserBalance represents a user balance in the system.
type UserBalance struct {
	ID         int64           `json:"id"`
	Balance    decimal.Decimal `json:"balance"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	TelegramID string          `json:"telegram_id"`
}

//
// ADD USER BALANCE
//

type AddUserBalanceDTO struct {
	EventTypeID int64           `json:"event_type_id"`
	Amount      decimal.Decimal `json:"amount"`
	TelegramID  string          `json:"telegram_id"`
	Description *string         `json:"description,omitempty"`
}

//
// REDUCE USER BALANCE
//

type ReduceUserBalanceDTO struct {
	EventTypeID int64           `json:"event_type_id"`
	Amount      decimal.Decimal `json:"amount"`
	TelegramID  string          `json:"telegram_id"`
	Description string          `json:"description"`
}

//
// SWAGGER
//

type GetUserBalanceSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID         int64           `json:"id" example:"1"`
		Balance    decimal.Decimal `json:"balance" example:"100.00"`
		CreatedAt  time.Time       `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
		UpdatedAt  time.Time       `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
		TelegramID string          `json:"telegram_id" example:"1"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
