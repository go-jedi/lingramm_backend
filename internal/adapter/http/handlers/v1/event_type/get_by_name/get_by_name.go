package getbyname

import (
	"context"
	"time"

	eventtype "github.com/go-jedi/lingramm_backend/internal/domain/event_type"
	eventtypeservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event_type"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetByName struct {
	eventTypeService *eventtypeservice.Service
	logger           logger.ILogger
}

func New(
	eventTypeService *eventtypeservice.Service,
	logger logger.ILogger,
) *GetByName {
	return &GetByName{
		eventTypeService: eventTypeService,
		logger:           logger,
	}
}

func (h *GetByName) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get event type by name] execute handler")

	name := c.Query("name")
	if name == "" {
		h.logger.Error("failed to get query name", "error", apperrors.ErrQueryIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get query name", apperrors.ErrQueryIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.eventTypeService.GetByName.Execute(ctxTimeout, name)
	if err != nil {
		h.logger.Error("failed to get event type by name", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get event type by name", err.Error(), nil))
	}

	return c.JSON(response.New[eventtype.EventType](true, "success", "", result))
}
