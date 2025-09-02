package existsbytelegramid

import (
	"context"
	"time"

	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type ExistsByTelegramID struct {
	adminService *adminservice.Service
	logger       logger.ILogger
}

func New(
	adminService *adminservice.Service,
	logger logger.ILogger,
) *ExistsByTelegramID {
	return &ExistsByTelegramID{
		adminService: adminService,
		logger:       logger,
	}
}

// Execute checks if an admin exists by Telegram ID (admin).
// @Summary Check admin existence by Telegram ID (admin)
// @Description Returns true if a user with the given Telegram ID is an admin, false otherwise.
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param telegramID path string true "Telegram ID"
// @Success 200 {object} admin.ExistsSwaggerResponse "Successful response"
// @Failure 400 {object} admin.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} admin.ErrorSwaggerResponse "Internal server error"
// @Router /v1/admin/exists/{telegramID}/exists [get]
func (h *ExistsByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[check admin exists by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.adminService.ExistsByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to exists admin by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to exists admin by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[bool](true, "success", "", result))
}
