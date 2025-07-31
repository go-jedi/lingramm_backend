package gettextsbylanguage

import (
	localizedtext "github.com/go-jedi/lingramm_backend/internal/domain/localized_text"
	localizedtextservice "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
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

func (h *GetTextsByLanguage) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get texts by language] execute handler")

	const languageLength = 2

	language := c.Params("language")
	if len(language) != languageLength {
		h.logger.Error("failed to get param language", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param language", apperrors.ErrParamIsRequired.Error(), nil))
	}

	result, err := h.localizedTextService.GetTextsByLanguage.Execute(c, language)
	if err != nil {
		h.logger.Error("failed to get texts by language", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get texts by language", err.Error(), nil))
	}

	return c.JSON(response.New[map[string][]localizedtext.LocalizedTexts](true, "success", "", result))
}
