package update

import (
	"context"
	"time"

	userstudiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/user_studied_language"
	userstudiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type Update struct {
	userStudiedLanguageService *userstudiedlanguageservice.Service
	logger                     logger.ILogger
	validator                  validator.IValidator
}

func New(
	userStudiedLanguageService *userstudiedlanguageservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Update {
	return &Update{
		userStudiedLanguageService: userStudiedLanguageService,
		logger:                     logger,
		validator:                  validator,
	}
}

func (h *Update) Execute(c fiber.Ctx) error {
	h.logger.Debug("[update user studied language] execute handler")

	var dto userstudiedlanguage.UpdateDTO
	if err := c.Bind().Body(&dto); err != nil {
		h.logger.Error("failed to bind body", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to bind body", err.Error(), nil))
	}

	if err := h.validator.StructCtx(c.RequestCtx(), dto); err != nil {
		h.logger.Error("failed to validate struct", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate struct", err.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.userStudiedLanguageService.Update.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to update user studied language", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to update user studied language", err.Error(), nil))
	}

	return c.JSON(response.New[userstudiedlanguage.UserStudiedLanguage](true, "success", "", result))
}
