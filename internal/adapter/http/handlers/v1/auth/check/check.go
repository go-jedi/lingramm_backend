package check

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type Check struct {
	authService *authservice.Service
	logger      logger.ILogger
	validator   validator.IValidator
}

func New(
	authService *authservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Check {
	return &Check{
		authService: authService,
		logger:      logger,
		validator:   validator,
	}
}

// Execute check user token.
//
// @Summary Check user token
// @Description Check if the provided Telegram ID and token are valid
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param request body auth.CheckDTO true "Check request body"
// @Success 200 {object} auth.CheckSwaggerResponse "Successful check user token"
// @Failure 400 {object} auth.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} auth.ErrorSwaggerResponse "Internal server error"
// @Router /v1/auth/check [post]
func (h *Check) Execute(c fiber.Ctx) error {
	h.logger.Debug("[check user token] execute handler")

	var dto auth.CheckDTO
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

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.authService.Check.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to check user token", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to check user token", err.Error(), nil))
	}

	return c.JSON(response.New[auth.CheckResponse](true, "success", "", result))
}
