package experiencepoint

import (
	createxpevents "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/experience_point/create_xp_events"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	experiencepointservice "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	createXPEvents *createxpevents.CreateXPEvents
}

func New(
	experiencePointService *experiencepointservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		createXPEvents: createxpevents.New(experiencePointService, logger, validator),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/experience_point",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Post("/events", h.createXPEvents.Execute)
	}
}
