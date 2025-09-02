package getbytelegramid

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/subscription"
	subscriptionservice "github.com/go-jedi/lingramm_backend/internal/service/v1/subscription"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetByTelegramID struct {
	subscriptionService *subscriptionservice.Service
	logger              logger.ILogger
}

func New(
	subscriptionService *subscriptionservice.Service,
	logger logger.ILogger,
) *GetByTelegramID {
	return &GetByTelegramID{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

// Execute returns a user's subscription by Telegram ID.
// @Summary Get subscription by Telegram ID
// @Description Returns the subscription record for the specified Telegram ID.
// @Tags Subscription
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param telegramID path string true "Telegram ID"
// @Success 200 {object} subscription.GetByTelegramIDSwaggerResponse "Successful response"
// @Failure 400 {object} subscription.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} subscription.ErrorSwaggerResponse "Internal server error"
// @Router /v1/subscription/telegram/{telegramID} [get]
func (h *GetByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get subscription by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.subscriptionService.GetByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get subscription by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get subscription by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[subscription.Subscription](true, "success", "", result))
}
