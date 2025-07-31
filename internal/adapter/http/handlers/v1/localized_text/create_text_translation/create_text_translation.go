package createtexttranslation

import (
	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	localizedtextservice "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

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

func (h *CreateTextTranslation) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create text translation] execute handler")

	var dto localizedtext.CreateTextTranslationDTO
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

	result, err := h.localizedTextService.CreateTextTranslation.Execute(c, dto)
	if err != nil {
		h.logger.Error("failed to create text translation", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create text translation", err.Error(), nil))
	}

	return c.JSON(response.New[localizedtext.TextTranslations](true, "success", "", result))
}
