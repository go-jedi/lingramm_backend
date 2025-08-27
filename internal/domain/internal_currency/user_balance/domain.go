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
	Description string          `json:"description"`
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
