package user

import "time"

// User represents a user in the system.
type User struct {
	ID         int64     `json:"id"`
	UUID       string    `json:"uuid"`
	TelegramID string    `json:"telegram_id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

//
// CREATE
//

// CreateDTO represents the data required to create a new user.
type CreateDTO struct {
	UUID       string `json:"uuid"`
	TelegramID string `json:"telegram_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}
