package notification

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/websocket/handlers/v1/notification/stream"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	userdailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task"
	userstatsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats"
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
	userStatsService *userstatsservice.Service,
	userDailyTaskService *userdailytaskservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	rabbitMQ *rabbitmq.RabbitMQ,
	redis *redis.Redis,
	wsManager *wsmanager.WSManager,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		stream: stream.New(
			notificationService,
			userStatsService,
			userDailyTaskService,
			logger,
			rabbitMQ,
			redis,
			middleware,
			wsManager.NotificationHUB,
		),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group("/v1/ws/notification", middleware.AuthWebSocket.AuthWebSocketMiddleware)
	{
		api.Get("/stream", h.stream.Execute)
	}
}
