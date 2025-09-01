package dailytask

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/daily_task/create"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	dailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	create *create.Create
}

func New(
	dailyTaskService *dailytaskservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		create: create.New(dailyTaskService, logger, validator),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/daily_task",
		middleware.Auth.AuthMiddleware,
		middleware.AdminGuard.AdminGuardMiddleware,
	)
	{
		api.Post("", h.create.Execute)
	}
}
