package admin

import "time"

// Admin represents admin in the system.
type Admin struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TelegramID string    `json:"telegram_id"`
}
