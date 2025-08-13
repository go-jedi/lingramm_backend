package subscription

import "time"

// Subscription represents a subscription in the system.
type Subscription struct {
	ID           int64      `json:"id"`
	TelegramID   string     `json:"telegram_id"`
	SubscribedAt *time.Time `json:"subscribed_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	IsActive     bool       `json:"is_active"`
}

// History represents a subscription history in the system.
type History struct {
	ID         int64     `json:"id"`
	TelegramID string    `json:"telegram_id"`
	ActionTime time.Time `json:"action_time"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

//
// CREATE SUBSCRIPTION HISTORY
//

type CreateSubscriptionHistoryDTO struct {
	TelegramID string     `json:"telegram_id"`
	ActionTime time.Time  `json:"action_time"`
	ExpiresAt  *time.Time `json:"expires_at"`
}
