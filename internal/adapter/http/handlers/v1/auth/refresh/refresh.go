package refresh

import (
	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Refresh struct {
	authService *authservice.Service
	logger      logger.ILogger
	validator   validator.IValidator
}

func New(
	authService *authservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Refresh {
	return &Refresh{
		authService: authService,
		logger:      logger,
		validator:   validator,
	}
}

func (h *Refresh) Execute(c fiber.Ctx) error {
	h.logger.Debug("[refresh user token] execute handler")

	var dto auth.RefreshDTO
	if err := c.Bind().Body(&dto); err != nil {
		h.logger.Error("failed to bind body", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to bind body", err.Error(), nil))
	}

	if err := h.validator.StructCtx(c, dto); err != nil {
		h.logger.Error("failed to validate struct", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate struct", err.Error(), nil))
	}

	result, err := h.authService.Refresh.Execute(c, dto)
	if err != nil {
		h.logger.Error("failed to refresh tokens", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to refresh tokens", err.Error(), nil))
	}

	return c.JSON(response.New[auth.RefreshResponse](true, "success", "", result))
}
