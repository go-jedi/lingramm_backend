package event

import (
	createevents "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/event/create_events"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	eventservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	createEvents *createevents.CreateEvents
}

func New(
	eventService *eventservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		createEvents: createevents.New(eventService, logger, validator),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/event",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Post("", h.createEvents.Execute)
	}
}
