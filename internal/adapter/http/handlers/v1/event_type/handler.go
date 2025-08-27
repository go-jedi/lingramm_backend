package eventtype

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/event_type/all"
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/event_type/create"
	getbyname "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/event_type/get_by_name"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	eventtypeservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event_type"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	all       *all.All
	create    *create.Create
	getByName *getbyname.GetByName
}

func New(
	eventTypeService *eventtypeservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		all:       all.New(eventTypeService, logger),
		create:    create.New(eventTypeService, logger, validator),
		getByName: getbyname.New(eventTypeService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/event_type",
		middleware.Auth.AuthMiddleware,
		middleware.AdminGuard.AdminGuardMiddleware,
	)
	{
		api.Post("", h.create.Execute)
		api.Get("/all", h.all.Execute)
		api.Get("/name", h.getByName.Execute)
	}
}
