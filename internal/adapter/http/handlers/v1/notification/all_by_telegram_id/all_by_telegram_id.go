package allbytelegramid

import (
	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

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

func (h *AllByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all notifications by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	result, err := h.notificationService.AllByTelegramID.Execute(c.RequestCtx(), telegramID)
	if err != nil {
		h.logger.Error("failed to get all notifications by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all notifications by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[[]notification.Notification](true, "success", "", result))
}
