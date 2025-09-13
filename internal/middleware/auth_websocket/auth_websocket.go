package authwebsocket

import (
	"errors"
	"strings"

	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const (
	authorizationHeader       = "Authorization"
	authorizationType         = "Bearer"
	telegramIDCtx             = "telegramID"
	connectionHeader          = "Connection"
	upgradeHeader             = "Upgrade"
	secWebSocketVersionHeader = "Sec-WebSocket-Version"
	secWebSocketKeyHeader     = "Sec-WebSocket-Key"
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
	AuthWebSocketMiddleware(c fiber.Ctx) error
	GetTelegramIDFromContext(c fiber.Ctx) (string, error)
}

type Middleware struct {
	jwt   *jwt.JWT
	redis *redis.Redis
}

func New(
	jwt *jwt.JWT,
	redis *redis.Redis,
) *Middleware {
	return &Middleware{
		jwt:   jwt,
		redis: redis,
	}
}

func (m *Middleware) AuthWebSocketMiddleware(c fiber.Ctx) error {
	if !m.isWebSocketHandshake(c) {
		c.Status(fiber.StatusUpgradeRequired)
		return c.JSON(response.New[any](false, "upgrade to websocket required", "not a websocket handshake", nil))
	}

	token, err := m.extractTokenFromHeader(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(response.New[any](false, "unauthorized: invalid or missing token", err.Error(), nil))
	}

	vr, err := m.jwt.ParseToken(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(response.New[any](false, "unauthorized: invalid token signature", err.Error(), nil))
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

// isWebSocketHandshake check is websocket handshake.
func (m *Middleware) isWebSocketHandshake(c fiber.Ctx) bool {
	conn := strings.ToLower(c.Get(connectionHeader))
	upg := strings.ToLower(c.Get(upgradeHeader))
	ver := strings.TrimSpace(c.Get(secWebSocketVersionHeader))
	key := strings.TrimSpace(c.Get(secWebSocketKeyHeader))

	return strings.Contains(conn, "upgrade") &&
		upg == "websocket" &&
		ver != "" &&
		key != ""
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
