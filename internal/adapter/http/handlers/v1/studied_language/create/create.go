package create

import (
	"context"
	"time"

	studiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/studied_language"
	studiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type Create struct {
	studiedLanguageService *studiedlanguageservice.Service
	logger                 logger.ILogger
	validator              validator.IValidator
}

func New(
	studiedLanguageService *studiedlanguageservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Create {
	return &Create{
		studiedLanguageService: studiedLanguageService,
		logger:                 logger,
		validator:              validator,
	}
}

// Execute creates a new studied language (admin).
// @Summary Create studied language (admin)
// @Description Creates a studied language with required `name`, `description`, and a 2-letter `lang` code.
// @Tags Studied language
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param payload body studiedlanguage.CreateDTO true "Studied language data"
// @Success 200 {object} studiedlanguage.CreateSwaggerResponse "Successful response"
// @Failure 400 {object} studiedlanguage.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} studiedlanguage.ErrorSwaggerResponse "Internal server error"
// @Router /v1/studied_language [post]
func (h *Create) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a new studied language] execute handler")

	var dto studiedlanguage.CreateDTO
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

	result, err := h.studiedLanguageService.Create.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to create a new studied language", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a new studied language", err.Error(), nil))
	}

	return c.JSON(response.New[studiedlanguage.StudiedLanguage](true, "success", "", result))
}
