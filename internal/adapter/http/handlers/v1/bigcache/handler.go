package bigcache

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/bigcache/iterator"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
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
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		iterator: iterator.New(bigCacheService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/bigcache",
		middleware.Auth.AuthMiddleware,
		middleware.AdminGuard.AdminGuardMiddleware,
	)
	{
		api.Get("/info", h.iterator.Execute)
	}
}
