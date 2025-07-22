package adminguard

import (
	"errors"
	"strings"

	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const (
	authorizationHeader = "Authorization"
	authorizationType   = "Bearer"
)

var (
	ErrEmptyAuthorizationHeader   = errors.New("empty authorization header")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrTokenIsEmpty               = errors.New("token is empty")
	ErrAccessDenied               = errors.New("access denied: you do not have permission to perform this action")
)

type Middleware struct {
	adminService *adminservice.Service
	jwt          *jwt.JWT
}

func New(
	adminService *adminservice.Service,
	jwt *jwt.JWT,
) *Middleware {
	return &Middleware{
		adminService: adminService,
		jwt:          jwt,
	}
}

func (m *Middleware) AdminGuardMiddleware(c fiber.Ctx) error {
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

	ie, err := m.adminService.ExistsByTelegramID.Execute(c.Context(), vr.TelegramID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "internal server error", err.Error(), nil))
	}

	if !ie {
		c.Status(fiber.StatusForbidden)
		return c.JSON(response.New[any](false, "access denied", ErrAccessDenied.Error(), nil))
	}

	return c.Next()
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
