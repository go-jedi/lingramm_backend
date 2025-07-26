package clientassets

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/file_server/client_assets/all"
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/file_server/client_assets/create"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	all    *all.All
	create *create.Create
}

func New(
	clientAssetsService *clientassetsservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		all:    all.New(clientAssetsService, logger, validator),
		create: create.New(clientAssetsService, logger, validator),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/fs/client_assets",
		middleware.Auth.AuthMiddleware,
		middleware.AdminGuard.AdminGuardMiddleware,
	)
	{
		api.Post("", h.create.Execute)
		api.Get("/all", h.all.Execute)
	}
}
