package addadminuser

import (
	"github.com/go-jedi/lingramm_backend/internal/domain/admin"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type AddAdminUser struct {
	adminService *adminservice.Service
	logger       logger.ILogger
}

func New(
	adminService *adminservice.Service,
	logger logger.ILogger,
) *AddAdminUser {
	return &AddAdminUser{
		adminService: adminService,
		logger:       logger,
	}
}

func (h *AddAdminUser) Execute(c fiber.Ctx) error {
	h.logger.Debug("[add a new admin user] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	result, err := h.adminService.AddAdminUser.Execute(c.Context(), telegramID)
	if err != nil {
		h.logger.Error("failed to add admin user", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to add admin user", err.Error(), nil))
	}

	return c.JSON(response.New[admin.Admin](true, "success", "", result))
}
