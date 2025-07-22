package clientassets

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/file_server/client_assets/all"
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/file_server/client_assets/create"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	create *create.Create
	all    *all.All
}

func New(
	clientAssetsService *clientassetsservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
) *Handler {
	h := &Handler{
		create: create.New(clientAssetsService, logger, validator),
		all:    all.New(clientAssetsService, logger, validator),
	}

	h.initRoutes(app)

	return h
}

func (h *Handler) initRoutes(app *fiber.App) {
	api := app.Group("/v1/fs/client_assets")
	{
		api.Post("", h.create.Execute)
		api.Get("/all", h.all.Execute)
	}
}
