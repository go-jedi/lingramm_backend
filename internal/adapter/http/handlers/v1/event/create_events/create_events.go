package createevents

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/event"
	eventservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type CreateEvents struct {
	eventService *eventservice.Service
	logger       logger.ILogger
	validator    validator.IValidator
}

func New(
	eventService *eventservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *CreateEvents {
	return &CreateEvents{
		eventService: eventService,
		logger:       logger,
		validator:    validator,
	}
}

func (h *CreateEvents) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a new events] execute handler")

	var dto event.CreateEventsDTO
	if err := c.Bind().Body(&dto); err != nil {
		h.logger.Error("failed to bind body", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to bind body", err.Error(), nil))
	}

	if err := h.validator.StructCtx(c.RequestCtx(), dto); err != nil {
		h.logger.Error("failed to validate struct", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate struct", err.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	if err := h.eventService.CreateEvents.Execute(ctxTimeout, dto); err != nil {
		h.logger.Error("failed to create a new events", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a new events", err.Error(), nil))
	}

	return c.JSON(response.New[any](true, "success", "", nil))
}
