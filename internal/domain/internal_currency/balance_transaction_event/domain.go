package balancetransactionevent

import (
	"time"

	"github.com/shopspring/decimal"
)

type BalanceTransactionEvent struct {
	ID          int64           `json:"id"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
