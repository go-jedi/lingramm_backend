package auth

import (
	"errors"
	"strings"

	"github.com/go-jedi/lingvogramm_backend/pkg/jwt"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const (
	authorizationHeader = "Authorization"
	authorizationType   = "Bearer"
	telegramIDCtx       = "telegramID"
)

var (
	ErrEmptyAuthorizationHeader              = errors.New("empty authorization header")
	ErrInvalidAuthorizationHeader            = errors.New("invalid authorization header")
	ErrTokenIsEmpty                          = errors.New("token is empty")
	ErrTelegramIDMakingRequestNotFound       = errors.New("telegram id making request not found")
	ErrTelegramIDMakingRequestHasInvalidType = errors.New("telegram id making request has invalid type")
)

//go:generate mockery --name=IMiddleware --output=mocks --case=underscore
type IMiddleware interface {
	AuthMiddleware(c fiber.Ctx) error
	GetTelegramIDFromContext(c fiber.Ctx) (string, error)
}

type Middleware struct {
	jwt *jwt.JWT
}

func New(jwt *jwt.JWT) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}

func (m *Middleware) AuthMiddleware(c fiber.Ctx) error {
	token, err := m.extractTokenFromHeader(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(response.New[any](false, "failed to extract token from header", err.Error(), nil))
	}

	vr, err := m.jwt.ParseToken(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(response.New[any](false, "failed to parse token", err.Error(), nil))
	}

	c.Locals(telegramIDCtx, vr.TelegramID)

	return c.Next()
}

// GetTelegramIDFromContext get telegram id making request from context.
func (m *Middleware) GetTelegramIDFromContext(c fiber.Ctx) (string, error) {
	val := c.Locals(telegramIDCtx)
	if val == nil {
		return "", ErrTelegramIDMakingRequestNotFound
	}

	telegramID, ok := val.(string)
	if !ok {
		return "", ErrTelegramIDMakingRequestHasInvalidType
	}

	return telegramID, nil
}

// extractTokenFromHeader extract token.
func (m *Middleware) extractTokenFromHeader(c fiber.Ctx) (string, error) {
	header := c.Get(authorizationHeader)
	if header == "" {
		return "", ErrEmptyAuthorizationHeader
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != authorizationType {
		return "", ErrInvalidAuthorizationHeader
	}

	if len(headerParts[1]) == 0 {
		return "", ErrTokenIsEmpty
	}

	return headerParts[1], nil
}
