package notification

import (
	allbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/notification/all_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/notification/create"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	create          *create.Create
	allByTelegramID *allbytelegramid.AllByTelegramID
}

func New(
	notificationService *notificationservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		create:          create.New(notificationService, logger, validator),
		allByTelegramID: allbytelegramid.New(notificationService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/notification",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Post("", middleware.AdminGuard.AdminGuardMiddleware, h.create.Execute)
		api.Get("/all/telegram/:telegramID", h.allByTelegramID.Execute)
	}
}
