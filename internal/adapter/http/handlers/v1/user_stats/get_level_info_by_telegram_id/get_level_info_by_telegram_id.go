package getlevelinfobytelegramid

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

type GetLevelInfoByTelegramID struct {
	userStatsService *userstatsservice.Service
	logger           logger.ILogger
}

func New(
	userStatsService *userstatsservice.Service,
	logger logger.ILogger,
) *GetLevelInfoByTelegramID {
	return &GetLevelInfoByTelegramID{
		userStatsService: userStatsService,
		logger:           logger,
	}
}

// Execute returns detailed level/XP info for a user by Telegram ID.
// @Summary Get level info by Telegram ID
// @Description Returns detailed level progress data: total XP, current level, level floor/ceil XP, next level, XP within level, XP to next level, progress ratio, and level name.
// @Tags User stats
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param telegramID path string true "Telegram ID"
// @Success 200 {object} userstats.GetLevelInfoByTelegramIDSwaggerResponse "Successful response"
// @Failure 400 {object} userstats.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} userstats.ErrorSwaggerResponse "Internal server error"
// @Router /v1/user_stats/level_info/telegram/{telegramID} [get]
func (h *GetLevelInfoByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get level info by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.userStatsService.GetLevelInfoByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get level info by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get level info by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[userstats.GetLevelInfoByTelegramIDResponse](true, "success", "", result))
}
