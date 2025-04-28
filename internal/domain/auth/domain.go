package auth

import "time"

//
// SIGN IN
//

// SignInDTO represents the request body for signing in a user.
// @param telegram_id string true "Telegram ID of the user"
// @param username string true "Username of the user"
// @param first_name string true "First name of the user"
// @param last_name string true "Last name of the user"
type SignInDTO struct {
	TelegramID string `json:"telegram_id" validate:"required,min=1"`
	Username   string `json:"username" validate:"omitempty,min=1"`
	FirstName  string `json:"first_name" validate:"omitempty,min=1"`
	LastName   string `json:"last_name" validate:"omitempty,min=1"`
}

// SignInResp represents the response body for a successful sign-in.
// @param access_token string true "Access token"
// @param refresh_token string true "Refresh token"
// @param access_exp_at string true "Expiry time of the access token"
// @param refresh_exp_at string true "Expiry time of the refresh token"
type SignInResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AccessExpAt  time.Time `json:"access_exp_at"`
	RefreshExpAt time.Time `json:"refresh_exp_at"`
}

//
// CHECK
//

// CheckDTO represents the request body for checking a user token.
// @param telegram_id string true "Telegram ID of the user"
// @param token string true "Token to validate"
type CheckDTO struct {
	TelegramID string `json:"telegram_id" validate:"required,min=1"`
	Token      string `validate:"required,min=1"`
}

// CheckResponse represents the response body for a successful token check.
// @param telegram_id string true "Telegram ID of the user"
// @param token string true "Validated token"
// @param exp_at string true "Expiry time of the token"
type CheckResponse struct {
	TelegramID string    `json:"telegram_id"`
	Token      string    `json:"token"`
	ExpAt      time.Time `json:"exp_at"`
}

//
// REFRESH
//

// RefreshDTO represents the request body for refreshing a user token.
// @param telegram_id string true "Telegram ID of the user"
// @param refresh_token string true "Refresh token to validate"
type RefreshDTO struct {
	TelegramID   string `json:"telegram_id" validate:"required,min=1"`
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
}

// RefreshResponse represents the response body for a successful token refresh.
// @param access_token string true "New access token"
// @param refresh_token string true "New refresh token"
// @param access_exp_at string true "Expiry time of the access token"
// @param refresh_exp_at string true "Expiry time of the refresh token"
type RefreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AccessExpAt  time.Time `json:"access_exp_at"`
	RefreshExpAt time.Time `json:"refresh_exp_at"`
}
