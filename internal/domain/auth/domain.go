package auth

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
