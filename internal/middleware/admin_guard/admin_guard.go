package adminguard

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-jedi/lingvogramm_backend/pkg/jwt"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
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
	jwt *jwt.JWT
}

func New(jwt *jwt.JWT) *Middleware {
	return &Middleware{
		jwt: jwt,
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

	fmt.Println(vr.TelegramID)
	// вызываем сервис проверяющий является ли пользователь администратором.
	// если является администратором, то пропускаем дальше, а иначе ошибку.
	/*

		ie, err := m.adminService.ExistsByTelegramID(c.Request.Context(), vr.TelegramID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "internal server error",
			})
			return
		}

		if !ie {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": ErrAccessDenied.Error(),
			})
			return
		}
	*/

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
