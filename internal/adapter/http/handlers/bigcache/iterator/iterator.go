package iterator

import (
	bigcacheservice "github.com/go-jedi/lingvogramm_backend/internal/service/bigcache"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

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

func (i *Iterator) Execute(c fiber.Ctx) error {
	i.logger.Debug("[iterator for show data in bigcache] execute handler")

	result, err := i.bigCacheService.Iterator.Execute(c.Context())
	if err != nil {
		i.logger.Error("failed to show data in bigcache", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to show data in bigcache", err.Error(), nil))
	}

	return c.JSON(response.New[map[string]any](true, "success", "", result))
}
