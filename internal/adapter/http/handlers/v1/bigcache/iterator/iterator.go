package iterator

import (
	"context"
	"time"

	bigcacheservice "github.com/go-jedi/lingramm_backend/internal/service/v1/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type Iterator struct {
	bigCacheService *bigcacheservice.Service
	logger          logger.ILogger
}

func New(
	bigCacheService *bigcacheservice.Service,
	logger logger.ILogger,
) *Iterator {
	return &Iterator{
		bigCacheService: bigCacheService,
		logger:          logger,
	}
}

// Execute lists entries from BigCache (admin).
// @Summary Iterate BigCache (admin)
// @Description Iterates over BigCache and returns a map of entries (keys/values) for inspection.
// @Tags BigCache
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Success 200 {object} map[string]interface{} "Successful response with BigCache entries"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/bigcache/info [get]
func (h *Iterator) Execute(c fiber.Ctx) error {
	h.logger.Debug("[iterator for show data in bigcache] execute handler")

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.bigCacheService.Iterator.Execute(ctxTimeout)
	if err != nil {
		h.logger.Error("failed to show data in bigcache", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to show data in bigcache", err.Error(), nil))
	}

	return c.JSON(response.New[map[string]any](true, "success", "", result))
}
