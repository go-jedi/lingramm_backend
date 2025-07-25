package balancetransactionevent

import (
	"time"

	"github.com/shopspring/decimal"
)

// BalanceTransactionEvent represents balance transaction event in the system.
type BalanceTransactionEvent struct {
	ID          int64           `json:"id"`
	Amount      decimal.Decimal `json:"amount"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	IsActive    bool            `json:"is_active"`
}
