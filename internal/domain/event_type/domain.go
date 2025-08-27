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
