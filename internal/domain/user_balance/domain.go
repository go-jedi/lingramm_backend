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
