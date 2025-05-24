package auth

import "time"

//
// SIGN IN
//

// SignInDTO represents the request body for signing in a user.
type SignInDTO struct {
	TelegramID string `json:"telegram_id" validate:"required,min=1"`
	Username   string `json:"username" validate:"omitempty,min=1"`
	FirstName  string `json:"first_name" validate:"omitempty,min=1"`
	LastName   string `json:"last_name" validate:"omitempty,min=1"`
}

//
// CHECK
//

// CheckDTO represents the request body for checking a user token.
// @param telegram_id string true "Telegram ID of the user".
// @param token string true "Token to validate".
type CheckDTO struct {
	TelegramID string `json:"telegram_id" validate:"required,min=1"`
	Token      string `validate:"required,min=1"`
}

// CheckResponse represents the response body for a successful token check.
// @param telegram_id string true "Telegram ID of the user".
// @param token string true "Validated token".
// @param exp_at string true "Expiry time of the token".
type CheckResponse struct {
	TelegramID string    `json:"telegram_id"`
	Token      string    `json:"token"`
	ExpAt      time.Time `json:"exp_at"`
}
