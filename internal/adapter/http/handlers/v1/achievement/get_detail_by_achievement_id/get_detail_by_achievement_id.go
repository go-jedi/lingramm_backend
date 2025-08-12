package getdetailbyachievementid

import (
	"context"
	"strconv"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetDetailByAchievementID struct {
	achievementService *achievementservice.Service
	logger             logger.ILogger
}

func New(
	achievementService *achievementservice.Service,
	logger logger.ILogger,
) *GetDetailByAchievementID {
	return &GetDetailByAchievementID{
		achievementService: achievementService,
		logger:             logger,
	}
}

func (h *GetDetailByAchievementID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get detail by achievement id] execute handler")

	achievementIDStr := c.Params("achievementID")
	if achievementIDStr == "" {
		h.logger.Error("failed to get param achievementID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param achievementID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	achievementID, err := strconv.ParseInt(achievementIDStr, 10, 64)
	if err != nil {
		h.logger.Error("failed parse string to int64", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed parse string to int64", err.Error(), nil))
	}

	if achievementID <= 0 {
		h.logger.Error("invalid achievementID", "error", "achievement id must be a positive integer")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "invalid achievement id", "achievement id must be a positive integer", nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.achievementService.GetDetailByAchievementID.Execute(ctxTimeout, achievementID)
	if err != nil {
		h.logger.Error("failed to get detail by achievement id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get detail by achievement id", err.Error(), nil))
	}

	return c.JSON(response.New[achievement.Detail](true, "success", "", result))
}
