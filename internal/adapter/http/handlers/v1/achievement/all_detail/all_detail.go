package alldetail

import (
	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

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

func (h *AllDetail) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all detail] execute handler")

	result, err := h.achievementService.All.Execute(c)
	if err != nil {
		h.logger.Error("failed to get all detail", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all detail", err.Error(), nil))
	}

	return c.JSON(response.New[[]achievement.Detail](true, "success", "", result))
}
