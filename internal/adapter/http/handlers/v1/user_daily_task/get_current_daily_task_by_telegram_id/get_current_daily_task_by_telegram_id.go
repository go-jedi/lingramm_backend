package getcurrentdailytaskbytelegramid

import (
	"context"
	"time"

	userdailytask "github.com/go-jedi/lingramm_backend/internal/domain/user_daily_task"
	userdailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetCurrentDailyTaskByTelegramID struct {
	userDailyTaskService *userdailytaskservice.Service
	logger               logger.ILogger
}

func New(
	userDailyTaskService *userdailytaskservice.Service,
	logger logger.ILogger,
) *GetCurrentDailyTaskByTelegramID {
	return &GetCurrentDailyTaskByTelegramID{
		userDailyTaskService: userDailyTaskService,
		logger:               logger,
	}
}

func (h *GetCurrentDailyTaskByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get current daily task by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.userDailyTaskService.GetCurrentDailyTaskByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get current daily task by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get current daily task by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[userdailytask.GetCurrentDailyTaskByTelegramIDResponse](true, "success", "", result))
}
