package studiedlanguage

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/studied_language/all"
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/studied_language/create"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	studiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	create *create.Create
	all    *all.All
}

func New(
	studiedLanguageService *studiedlanguageservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		create: create.New(studiedLanguageService, logger, validator),
		all:    all.New(studiedLanguageService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/studied_language",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Post("", middleware.AdminGuard.AdminGuardMiddleware, h.create.Execute)
		api.Get("/all", h.all.Execute)
	}
}
