package user

import (
	getbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	userservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	getByTelegramID *getbytelegramid.GetByTelegramID
}

func New(
	userService *userservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		getByTelegramID: getbytelegramid.New(userService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/user",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Get("/telegram/:telegramID", h.getByTelegramID.Execute)
	}
}
