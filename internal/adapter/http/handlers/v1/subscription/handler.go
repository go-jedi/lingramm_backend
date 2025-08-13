package subscription

import (
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/subscription/exists_by_telegram_id"
	getbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/subscription/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	subscriptionservice "github.com/go-jedi/lingramm_backend/internal/service/v1/subscription"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	existsByTelegramID *existsbytelegramid.ExistsByTelegramID
	getByTelegramID    *getbytelegramid.GetByTelegramID
}

func New(
	subscriptionService *subscriptionservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		existsByTelegramID: existsbytelegramid.New(subscriptionService, logger),
		getByTelegramID:    getbytelegramid.New(subscriptionService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/subscription",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Get("/telegram/:telegramID", h.getByTelegramID.Execute)
		api.Get("/exists/telegram/:telegramID", h.existsByTelegramID.Execute)
	}
}
