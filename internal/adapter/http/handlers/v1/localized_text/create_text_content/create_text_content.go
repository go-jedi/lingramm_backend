package createtextcontent

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

type CreateTextContent struct {
	localizedTextService *localizedtextservice.Service
	logger               logger.ILogger
	validator            validator.IValidator
}

func New(
	localizedTextService *localizedtextservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *CreateTextContent {
	return &CreateTextContent{
		localizedTextService: localizedTextService,
		logger:               logger,
		validator:            validator,
	}
}

// Execute creates a localized text content record (admin).
// @Summary Create text content (admin)
// @Description Creates a localized text content entry with required `code` and `page`, and optional `description`.
// @Tags Localized text
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param payload body localizedtext.CreateTextContentDTO true "Text content data"
// @Success 200 {object} localizedtext.CreateTextContentSwaggerResponse "Successful response"
// @Failure 400 {object} localizedtext.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} localizedtext.ErrorSwaggerResponse "Internal server error"
// @Router /v1/localized_text/content [post]
func (h *CreateTextContent) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create text content] execute handler")

	var dto localizedtext.CreateTextContentDTO
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

	result, err := h.localizedTextService.CreateTextContent.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to create text content", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create text content", err.Error(), nil))
	}

	return c.JSON(response.New[localizedtext.TextContents](true, "success", "", result))
}
