package gettextsbylanguage

import (
	"context"
	"time"

	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	localizedtextservice "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const (
	languageLength = 2
	timeout        = 5 * time.Second
)

type GetTextsByLanguage struct {
	localizedTextService *localizedtextservice.Service
	logger               logger.ILogger
	validator            validator.IValidator
}

func New(
	localizedTextService *localizedtextservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *GetTextsByLanguage {
	return &GetTextsByLanguage{
		localizedTextService: localizedTextService,
		logger:               logger,
		validator:            validator,
	}
}

// Execute returns localized texts grouped by code for a given language.
// @Summary Get texts by language
// @Description Returns a map where keys are page codes and values are arrays of localized texts for the specified 2-letter language code.
// @Tags Localized text
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param language path string true "2-letter language code (e.g., en, ru)"
// @Success 200 {object} map[string][]localizedtext.GetTextsByLanguageSwaggerResponse "Successful response"
// @Failure 400 {object} localizedtext.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} localizedtext.ErrorSwaggerResponse "Internal server error"
// @Router /v1/localized_text/texts/language/{language} [get]
func (h *GetTextsByLanguage) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get texts by language] execute handler")

	language := c.Params("language")
	if len(language) != languageLength {
		h.logger.Error("failed to get param language", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param language", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.localizedTextService.GetTextsByLanguage.Execute(ctxTimeout, language)
	if err != nil {
		h.logger.Error("failed to get texts by language", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get texts by language", err.Error(), nil))
	}

	return c.JSON(response.New[map[string][]localizedtext.LocalizedTexts](true, "success", "", result))
}
