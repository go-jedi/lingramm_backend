package allbytelegramid

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type AllByTelegramID struct {
	notificationService *notificationservice.Service
	logger              logger.ILogger
}

func New(
	notificationService *notificationservice.Service,
	logger logger.ILogger,
) *AllByTelegramID {
	return &AllByTelegramID{
		notificationService: notificationService,
		logger:              logger,
	}
}

// Execute returns all notifications for a user by Telegram ID.
// @Summary Get all notifications by Telegram ID
// @Description Returns a list of notifications for the specified Telegram ID.
// @Tags Notification
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param telegramID path string true "Telegram ID"
// @Success 200 {object} notification.AllSwaggerResponse "Successful response"
// @Failure 400 {object} notification.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} notification.ErrorSwaggerResponse "Internal server error"
// @Router /v1/notification/all/telegram/{telegramID} [get]
func (h *AllByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all notifications by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.notificationService.AllByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get all notifications by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all notifications by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[[]notification.Notification](true, "success", "", result))
}
