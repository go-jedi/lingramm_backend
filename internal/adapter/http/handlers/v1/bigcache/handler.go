package bigcache

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/bigcache/iterator"
	bigcacheservice "github.com/go-jedi/lingramm_backend/internal/service/v1/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	iterator *iterator.Iterator
}

func New(
	bigCacheService *bigcacheservice.Service,
	app *fiber.App,
	logger logger.ILogger,
) *Handler {
	h := &Handler{
		iterator: iterator.New(bigCacheService, logger),
	}

	h.initRoutes(app)

	return h
}

func (h *Handler) initRoutes(app *fiber.App) {
	api := app.Group("/v1/bigcache")
	{
		api.Get("/info", h.iterator.Execute)
	}
}
