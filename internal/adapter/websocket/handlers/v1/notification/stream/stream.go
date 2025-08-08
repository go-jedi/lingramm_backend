package stream

import (
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type Stream struct {
	notificationService *notificationservice.Service
	logger              logger.ILogger
}

func New(
	notificationService *notificationservice.Service,
	logger logger.ILogger,
) *Stream {
	return &Stream{
		notificationService: notificationService,
		logger:              logger,
	}
}

func (h *Stream) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get notifications stream] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	return nil
}
