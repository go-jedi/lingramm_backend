package notification

import "time"

const (
	PendingStatus = "PENDING"
	SentStatus    = "SENT"
	FailedStatus  = "FAILED"
)

const (
	AchievementType      = "achievement"
	InternalCurrencyType = "internal_currency"
	LevelType            = "level"
	MiniGameType         = "mini_game"
)

// Notification represents notification in the system.
type Notification struct {
	ID         int64      `json:"id"`
	Message    Message    `json:"message"`
	Type       string     `json:"type"`
	TelegramID string     `json:"telegram_id"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	SentAt     *time.Time `json:"sent_at,omitempty"`
}

type Message struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

//
// CREATE
//

type CreateDTO struct {
	Message    Message `json:"message"`
	Type       string  `json:"type"`
	TelegramID string  `json:"telegram_id"`
}

//
// SEND NOTIFICATION
//

type SendNotificationDTO struct {
	ID         int64     `json:"id"`
	Message    Message   `json:"message"`
	Type       string    `json:"type"`
	TelegramID string    `json:"telegram_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type ACKMessage struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}
