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

// Execute checks if a subscription exists by Telegram ID.
// @Summary Check subscription existence by Telegram ID
// @Description Returns true if the specified Telegram ID has an active subscription, false otherwise.
// @Tags Subscription
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param telegramID path string true "Telegram ID"
// @Success 200 {object} subscription.ExistsByTelegramIDSwaggerResponse "Successful response"
// @Failure 400 {object} subscription.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} subscription.ErrorSwaggerResponse "Internal server error"
// @Router /v1/subscription/exists/telegram/{telegramID} [get]
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
