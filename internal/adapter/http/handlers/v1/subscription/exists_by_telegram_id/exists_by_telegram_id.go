package existsbytelegramid

import (
	"context"
	"time"

	subscriptionservice "github.com/go-jedi/lingramm_backend/internal/service/v1/subscription"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type ExistsByTelegramID struct {
	subscriptionService *subscriptionservice.Service
	logger              logger.ILogger
}

func New(
	subscriptionService *subscriptionservice.Service,
	logger logger.ILogger,
) *ExistsByTelegramID {
	return &ExistsByTelegramID{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

func (h *ExistsByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[check subscription exists by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.subscriptionService.ExistsByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to check subscription exists by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to check subscription exists by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[bool](true, "success", "", result))
}
