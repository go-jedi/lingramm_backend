package all

import (
	"context"
	"time"

	eventtype "github.com/go-jedi/lingramm_backend/internal/domain/event_type"
	eventtypeservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event_type"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type All struct {
	eventTypeService *eventtypeservice.Service
	logger           logger.ILogger
}

func New(
	eventTypeService *eventtypeservice.Service,
	logger logger.ILogger,
) *All {
	return &All{
		eventTypeService: eventTypeService,
		logger:           logger,
	}
}

func (h *All) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all event types] execute handler")

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.eventTypeService.All.Execute(ctxTimeout)
	if err != nil {
		h.logger.Error("failed to get all event types", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all event types", err.Error(), nil))
	}

	return c.JSON(response.New[[]eventtype.EventType](true, "success", "", result))
}
