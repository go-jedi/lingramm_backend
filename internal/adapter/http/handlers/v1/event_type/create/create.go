package create

import (
	"context"
	"time"

	eventtype "github.com/go-jedi/lingramm_backend/internal/domain/event_type"
	eventtypeservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event_type"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type Create struct {
	eventTypeService *eventtypeservice.Service
	logger           logger.ILogger
	validator        validator.IValidator
}

func New(
	eventTypeService *eventtypeservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Create {
	return &Create{
		eventTypeService: eventTypeService,
		logger:           logger,
		validator:        validator,
	}
}

func (h *Create) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a new event type] execute handler")

	var dto eventtype.CreateDTO
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

	if dto.Amount != nil && !dto.Amount.IsPositive() {
		h.logger.Error("failed to validate amount", "error", "amount must be greater than zero")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate amount", "amount must be greater than zero", nil))
	}

	if dto.IsSendNotification && dto.NotificationMessage == nil {
		h.logger.Error("failed to validate notification_message", "error", "notification_message must be provided when is_send_notification is true")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate notification_message", "notification_message must be provided when is_send_notification is true", nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.eventTypeService.Create.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to create a new event type", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a new event type", err.Error(), nil))
	}

	return c.JSON(response.New[eventtype.EventType](true, "success", "", result))
}
