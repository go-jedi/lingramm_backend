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

//
// SWAGGER
//

type GetByTelegramIDSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID           int64      `json:"id" example:"1"`
		TelegramID   string     `json:"telegram_id" example:"1"`
		SubscribedAt *time.Time `json:"subscribed_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		ExpiresAt    *time.Time `json:"expires_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		CreatedAt    time.Time  `json:"created_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		UpdatedAt    time.Time  `json:"updated_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		IsActive     bool       `json:"is_active" example:"true"`
	} `json:"data"`
}

type ExistsByTelegramIDSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    bool   `json:"data" example:"true"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
