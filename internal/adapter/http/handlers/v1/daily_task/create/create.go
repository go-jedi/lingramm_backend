package create

import (
	"context"
	"time"

	dailytask "github.com/go-jedi/lingramm_backend/internal/domain/daily_task"
	dailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type Create struct {
	dailyTaskService *dailytaskservice.Service
	logger           logger.ILogger
	validator        validator.IValidator
}

func New(
	dailyTaskService *dailytaskservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Create {
	return &Create{
		dailyTaskService: dailyTaskService,
		logger:           logger,
		validator:        validator,
	}
}

func (h *Create) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a new daily task] execute handler")

	var dto dailytask.CreateDTO
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

	if (dto.WordsLearnedNeed == nil || *dto.WordsLearnedNeed <= 0) &&
		(dto.TasksCompletedNeed == nil || *dto.TasksCompletedNeed <= 0) &&
		(dto.LessonsFinishedNeed == nil || *dto.LessonsFinishedNeed <= 0) &&
		(dto.WordsTranslateNeed == nil || *dto.WordsTranslateNeed <= 0) &&
		(dto.DialogCompletedNeed == nil || *dto.DialogCompletedNeed <= 0) &&
		(dto.ExperiencePointsNeed == nil || *dto.ExperiencePointsNeed <= 0) {
		h.logger.Error("failed to validate need fields", "error", "at least one of the *_need fields must be provided and greater than 0")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate need fields", "at least one of the *_need fields must be provided and greater than 0", nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.dailyTaskService.Create.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to create a new daily task", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a new daily task", err.Error(), nil))
	}

	return c.JSON(response.New[dailytask.DailyTask](true, "success", "", result))
}
