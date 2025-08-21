package create

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type Create struct {
	notificationService *notificationservice.Service
	logger              logger.ILogger
	validator           validator.IValidator
}

func New(
	notificationService *notificationservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Create {
	return &Create{
		notificationService: notificationService,
		logger:              logger,
		validator:           validator,
	}
}

func (h *Create) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a new notification] execute handler")

	var dto notification.CreateDTO
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

	result, err := h.notificationService.Create.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to create a new notification", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a new notification", err.Error(), nil))
	}

	return c.JSON(response.New[notification.Notification](true, "success", "", result))
}
