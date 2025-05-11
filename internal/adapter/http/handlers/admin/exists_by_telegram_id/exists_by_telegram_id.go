package existsbytelegramid

import (
	adminservice "github.com/go-jedi/lingvogramm_backend/internal/service/admin"
	"github.com/go-jedi/lingvogramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

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

func (ebt *ExistsByTelegramID) Execute(c fiber.Ctx) error {
	ebt.logger.Debug("[check admin exists by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		ebt.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	result, err := ebt.adminService.ExistsByTelegramID.Execute(c.Context(), telegramID)
	if err != nil {
		ebt.logger.Error("failed to exists admin by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to exists admin by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[bool](true, "success", "", result))
}
