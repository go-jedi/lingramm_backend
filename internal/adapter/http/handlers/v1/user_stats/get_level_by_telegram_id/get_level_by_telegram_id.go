package getlevelbytelegramid

import (
	"context"
	"time"

	userstats "github.com/go-jedi/lingramm_backend/internal/domain/user_stats"
	userstatsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetLevelByTelegramID struct {
	userStatsService *userstatsservice.Service
	logger           logger.ILogger
}

func New(
	userStatsService *userstatsservice.Service,
	logger logger.ILogger,
) *GetLevelByTelegramID {
	return &GetLevelByTelegramID{
		userStatsService: userStatsService,
		logger:           logger,
	}
}

func (h *GetLevelByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get level user by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.userStatsService.GetLevelByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get level user by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get level user by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[userstats.GetLevelByTelegramIDResponse](true, "success", "", result))
}
