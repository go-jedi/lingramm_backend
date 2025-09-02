package notification

import (
	"time"
)

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

//
// SWAGGER
//

type CreateSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID      int64 `json:"id" example:"1"`
		Message struct {
			Title string `json:"title" example:"some title"`
			Text  string `json:"text" example:"some text"`
		} `json:"message"`
		Type       string     `json:"type" example:"some type"`
		TelegramID string     `json:"telegram_id" example:"1"`
		Status     string     `json:"status" example:"PENDING"`
		CreatedAt  time.Time  `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		SentAt     *time.Time `json:"sent_at,omitempty" example:"2025-09-02T15:30:20.095307198+03:00"`
	} `json:"data"`
}

type AllSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		ID      int64 `json:"id" example:"1"`
		Message struct {
			Title string `json:"title" example:"some title"`
			Text  string `json:"text" example:"some text"`
		} `json:"message"`
		Type       string     `json:"type" example:"some type"`
		TelegramID string     `json:"telegram_id" example:"1"`
		Status     string     `json:"status" example:"PENDING"`
		CreatedAt  time.Time  `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		SentAt     *time.Time `json:"sent_at,omitempty" example:"2025-09-02T15:30:20.095307198+03:00"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
