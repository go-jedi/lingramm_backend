package createxpevents

import (
	"context"
	"time"

	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	experiencepointservice "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const (
	tagTelegramID = "required,min=1"
	tagEvents     = "required,min=1,max=250"
	timeout       = 5 * time.Second
)

type CreateXPEvents struct {
	experiencePointService *experiencepointservice.Service
	logger                 logger.ILogger
	validator              validator.IValidator
}

func New(
	experiencePointService *experiencepointservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *CreateXPEvents {
	return &CreateXPEvents{
		experiencePointService: experiencePointService,
		logger:                 logger,
		validator:              validator,
	}
}

func (h *CreateXPEvents) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a new xp events] execute handler")

	var dto experiencepoint.CreateXPEventDTO
	if err := c.Bind().Body(&dto); err != nil {
		h.logger.Error("failed to bind body", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to bind body", err.Error(), nil))
	}

	if err := h.validator.VarCtx(c.RequestCtx(), dto.TelegramID, tagTelegramID); err != nil {
		h.logger.Error("failed to validate telegram_id", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate telegram_id", err.Error(), nil))
	}

	if err := h.validator.VarCtx(c.RequestCtx(), dto.Events, tagEvents); err != nil {
		h.logger.Error("failed to validate slice", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate slice", err.Error(), nil))
	}

	for i := range dto.Events {
		if err := h.validator.StructCtx(c.RequestCtx(), dto.Events[i]); err != nil {
			h.logger.Error("failed to validate struct", "error", err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.New[any](false, "failed to validate struct", err.Error(), nil))
		}
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	if err := h.experiencePointService.CreateXPEvents.Execute(ctxTimeout, dto); err != nil {
		h.logger.Error("failed to create a new xp events", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a new xp events", err.Error(), nil))
	}

	return c.JSON(response.New[any](true, "success", "", nil))
}
