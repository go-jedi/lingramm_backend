package createtexttranslation

import (
	"context"
	"time"

	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	localizedtextservice "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type CreateTextTranslation struct {
	localizedTextService *localizedtextservice.Service
	logger               logger.ILogger
	validator            validator.IValidator
}

func New(
	localizedTextService *localizedtextservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *CreateTextTranslation {
	return &CreateTextTranslation{
		localizedTextService: localizedTextService,
		logger:               logger,
		validator:            validator,
	}
}

// Execute creates a localized text translation (admin).
// @Summary Create text translation (admin)
// @Description Creates a translation for a text content entry.
// @Description Rules:
// @Description • `content_id` is required and must be > 0
// @Description • `lang` must be a 2-letter code (e.g., "en", "ru")
// @Description • `value` is required
// @Tags Localized text
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param payload body localizedtext.CreateTextTranslationDTO true "Text translation data"
// @Success 200 {object} localizedtext.CreateTextTranslationSwaggerResponse "Successful response"
// @Failure 400 {object} localizedtext.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} localizedtext.ErrorSwaggerResponse "Internal server error"
// @Router /v1/localized_text/translation [post]
func (h *CreateTextTranslation) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create text translation] execute handler")

	var dto localizedtext.CreateTextTranslationDTO
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

	result, err := h.localizedTextService.CreateTextTranslation.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to create text translation", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create text translation", err.Error(), nil))
	}

	return c.JSON(response.New[localizedtext.TextTranslations](true, "success", "", result))
}
