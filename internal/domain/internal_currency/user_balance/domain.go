package userbalance

import (
	"time"

	"github.com/shopspring/decimal"
)

// UserBalance represents a user balance in the system.
type UserBalance struct {
	ID         int64           `json:"id"`
	TelegramID string          `json:"telegram_id"`
	Balance    decimal.Decimal `json:"balance"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

//
// ADD USER BALANCE
//

type AddUserBalanceDTO struct {
	EventID     int64           `json:"event_id"`
	TelegramID  string          `json:"telegram_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}

//
// REDUCE USER BALANCE
//

type ReduceUserBalanceDTO struct {
	EventID     int64           `json:"event_id"`
	TelegramID  string          `json:"telegram_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}

//
// UPDATE USER BALANCE
//

type UpdateUserBalanceDTO struct {
	EventID     int64           `json:"event_id"`
	TelegramID  string          `json:"telegram_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}
