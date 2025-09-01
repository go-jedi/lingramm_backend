package alldetailbytelegramid

import (
	"context"
	"time"

	userachievement "github.com/go-jedi/lingramm_backend/internal/domain/user_achievement"
	userachievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_achievement"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type AllDetailByTelegramID struct {
	userAchievementService *userachievementservice.Service
	logger                 logger.ILogger
}

func New(
	userAchievementService *userachievementservice.Service,
	logger logger.ILogger,
) *AllDetailByTelegramID {
	return &AllDetailByTelegramID{
		userAchievementService: userAchievementService,
		logger:                 logger,
	}
}

func (h *AllDetailByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all user achievements detail] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.userAchievementService.AllDetailByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get all user achievements detail", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all user achievements detail", err.Error(), nil))
	}

	return c.JSON(response.New[[]userachievement.Detail](true, "success", "", result))
}
