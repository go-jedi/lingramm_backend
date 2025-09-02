package auth

import "time"

//
// SIGN IN
//

// SignInDTO represents the request body for signing in a user.
// @param telegram_id string true "Telegram ID of the user".
// @param username string true "Username of the user".
// @param first_name string true "First name of the user".
// @param last_name string true "Last name of the user".
type SignInDTO struct {
	TelegramID string `json:"telegram_id" validate:"required,min=1"`
	Username   string `json:"username" validate:"omitempty,min=1"`
	FirstName  string `json:"first_name" validate:"omitempty,min=1"`
	LastName   string `json:"last_name" validate:"omitempty,min=1"`
}

// SignInResp represents the response body for a successful sign-in.
// @param access_token string true "Access token".
// @param refresh_token string true "Refresh token".
// @param access_exp_at string true "Expiry time of the access token".
// @param refresh_exp_at string true "Expiry time of the refresh token".
type SignInResp struct {
	AccessExpAt  time.Time `json:"access_exp_at"`
	RefreshExpAt time.Time `json:"refresh_exp_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
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

//
// REFRESH
//

// RefreshDTO represents the request body for refreshing a user token.
// @param telegram_id string true "Telegram ID of the user".
// @param refresh_token string true "Refresh token to validate".
type RefreshDTO struct {
	TelegramID   string `json:"telegram_id" validate:"required,min=1"`
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
}

// RefreshResponse represents the response body for a successful token refresh.
// @param access_token string true "New access token".
// @param refresh_token string true "New refresh token".
// @param access_exp_at string true "Expiry time of the access token".
// @param refresh_exp_at string true "Expiry time of the refresh token".
type RefreshResponse struct {
	AccessExpAt  time.Time `json:"access_exp_at"`
	RefreshExpAt time.Time `json:"refresh_exp_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

//
// SWAGGER
//

type SignInSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		AccessExpAt  time.Time `json:"access_exp_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		RefreshExpAt time.Time `json:"refresh_exp_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWxlZ3JhbV9pZCI6IjEyMzQ1Njc4OTAiLCJleHAiOjE3NTY4MTYyMjAsImlhdCI6MTc1NjgxMjYyMH0.F39iYi6QaEEeQ9pTBO4HL_sOyfaAjRJs7IQiCSQihGE"`
		RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWxlZ3JhbV9pZCI6IjEyMzQ1Njc4OTAiLCJleHAiOjE3NTc0MTc0MjAsImlhdCI6MTc1NjgxMjYyMH0.BTdkTO_6gMG_zmGi1kse2bXoTJ9iU5SZUSl6Pp4upYk"`
	} `json:"data"`
}

type CheckSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		TelegramID string    `json:"telegram_id" example:"1"`
		Token      string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWxlZ3JhbV9pZCI6IjEyMzQ1Njc4OTAiLCJleHAiOjE3NTY4MTYyMjAsImlhdCI6MTc1NjgxMjYyMH0.F39iYi6QaEEeQ9pTBO4HL_sOyfaAjRJs7IQiCSQihGE"`
		ExpAt      time.Time `json:"exp_at" example:"2025-09-02T15:30:20.095307198+03:00"`
	} `json:"data"`
}

type RefreshSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		AccessExpAt  time.Time `json:"access_exp_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		RefreshExpAt time.Time `json:"refresh_exp_at" example:"2025-09-02T15:30:20.095307198+03:00"`
		AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWxlZ3JhbV9pZCI6IjEyMzQ1Njc4OTAiLCJleHAiOjE3NTY4MTYyMjAsImlhdCI6MTc1NjgxMjYyMH0.F39iYi6QaEEeQ9pTBO4HL_sOyfaAjRJs7IQiCSQihGE"`
		RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWxlZ3JhbV9pZCI6IjEyMzQ1Njc4OTAiLCJleHAiOjE3NTc0MTc0MjAsImlhdCI6MTc1NjgxMjYyMH0.BTdkTO_6gMG_zmGi1kse2bXoTJ9iU5SZUSl6Pp4upYk"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
