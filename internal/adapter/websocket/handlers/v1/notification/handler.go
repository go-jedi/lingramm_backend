package notification

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/websocket/handlers/v1/notification/stream"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/rabbitmq"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	wsmanager "github.com/go-jedi/lingramm_backend/pkg/ws_manager"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	stream *stream.Stream
}

func New(
	notificationService *notificationservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	rabbitMQ *rabbitmq.RabbitMQ,
	redis *redis.Redis,
	wsManager *wsmanager.WSManager,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		stream: stream.New(notificationService, logger, rabbitMQ, redis, wsManager.NotificationHUB),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group("/v1/ws/notification", middleware.AuthWebSocket.AuthWebSocketMiddleware)
	{
		api.Get("/stream/:telegramID", h.stream.Execute)
	}
}
