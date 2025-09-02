package create

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

type Create struct {
	userStudiedLanguageService *userstudiedlanguageservice.Service
	logger                     logger.ILogger
	validator                  validator.IValidator
}

func New(
	userStudiedLanguageService *userstudiedlanguageservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Create {
	return &Create{
		userStudiedLanguageService: userStudiedLanguageService,
		logger:                     logger,
		validator:                  validator,
	}
}

// Execute creates a new user–studied-language relation.
// @Summary Create user studied language
// @Description Links a user (by Telegram ID) to a studied language. Rules:
// @Description • `studied_languages_id` is required and must be > 0
// @Description • `telegram_id` is required
// @Tags User studied language
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param payload body userstudiedlanguage.CreateDTO true "User studied language data"
// @Success 200 {object} userstudiedlanguage.CreateSwaggerResponse "Successful response"
// @Failure 400 {object} userstudiedlanguage.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} userstudiedlanguage.ErrorSwaggerResponse "Internal server error"
// @Router /v1/user_studied_language [post]
func (h *Create) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a new user studied language] execute handler")

	var dto userstudiedlanguage.CreateDTO
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

	result, err := h.userStudiedLanguageService.Create.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to create a new user studied language", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a new user studied language", err.Error(), nil))
	}

	return c.JSON(response.New[userstudiedlanguage.UserStudiedLanguage](true, "success", "", result))
}
