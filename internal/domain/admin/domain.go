package admin

import "time"

// Admin represents a admin in the system.
type Admin struct {
	ID         int64     `json:"id"`
	TelegramID string    `json:"telegram_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
