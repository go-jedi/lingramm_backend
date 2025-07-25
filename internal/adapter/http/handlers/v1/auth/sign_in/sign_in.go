package signin

import (
	"github.com/go-jedi/lingramm_backend/config"
	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type SignIn struct {
	authService *authservice.Service
	cookie      config.CookieConfig
	logger      logger.ILogger
	validator   validator.IValidator
}

func New(
	authService *authservice.Service,
	cookie config.CookieConfig,
	logger logger.ILogger,
	validator validator.IValidator,
) *SignIn {
	return &SignIn{
		authService: authService,
		cookie:      cookie,
		logger:      logger,
		validator:   validator,
	}
}

func (h *SignIn) Execute(c fiber.Ctx) error {
	h.logger.Debug("[sign in user] execute handler")

	var dto auth.SignInDTO
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

	result, err := h.authService.SignIn.Execute(c, dto)
	if err != nil {
		h.logger.Error("failed to sign in user", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to sign in user", err.Error(), nil))
	}

	c.Cookie(&fiber.Cookie{
		Expires:     result.RefreshExpAt,
		Name:        h.cookie.Refresh.Name,
		Value:       result.RefreshToken,
		Path:        h.cookie.Refresh.Path,
		Domain:      h.cookie.Refresh.Domain,
		SameSite:    h.cookie.Refresh.SameSite,
		MaxAge:      h.cookie.Refresh.MaxAge,
		Secure:      h.cookie.Refresh.Secure,
		HTTPOnly:    h.cookie.Refresh.HTTPOnly,
		Partitioned: h.cookie.Refresh.Partitioned,
		SessionOnly: h.cookie.Refresh.SessionOnly,
	})

	return c.JSON(response.New[auth.SignInResp](true, "success", "", result))
}
