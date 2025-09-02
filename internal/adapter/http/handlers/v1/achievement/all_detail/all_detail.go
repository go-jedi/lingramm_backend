package alldetail

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type AllDetail struct {
	achievementService *achievementservice.Service
	logger             logger.ILogger
}

func New(
	achievementService *achievementservice.Service,
	logger logger.ILogger,
) *AllDetail {
	return &AllDetail{
		achievementService: achievementService,
		logger:             logger,
	}
}

// Execute returns all achievements with their conditions and assets.
// @Summary Get all achievement details (admin)
// @Description Returns a full list of achievements with their condition, achievement assets and award assets
// @Tags Achievement
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Success 200 {object} achievement.AllDetailSwaggerResponse "Successful response"
// @Failure 400 {object} achievement.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} achievement.ErrorSwaggerResponse "Internal server error"
// @Router /v1/achievement/all [get]
func (h *AllDetail) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all detail] execute handler")

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.achievementService.All.Execute(ctxTimeout)
	if err != nil {
		h.logger.Error("failed to get all detail", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all detail", err.Error(), nil))
	}

	return c.JSON(response.New[[]achievement.Detail](true, "success", "", result))
}
