package create

import (
	"fmt"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Create struct {
	clientAssetsService *clientassetsservice.Service
	logger              logger.ILogger
	validator           validator.IValidator
}

func New(
	clientAssetsService *clientassetsservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Create {
	return &Create{
		clientAssetsService: clientAssetsService,
		logger:              logger,
		validator:           validator,
	}
}

func (h *Create) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a client assets] execute handler")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get the first file for the provided form key", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get the first file for the provided form key", err.Error(), nil))
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if _, ok := clientassets.SupportedImageTypes[contentType]; !ok {
		h.logger.Error(fmt.Sprintf("unsupported file type: %s", contentType), "error")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "unsupported file type", fmt.Errorf("%w: %s", apperrors.ErrUnsupportedFormat, contentType).Error(), nil))
	}

	result, err := h.clientAssetsService.Create.Execute(c.Context(), fileHeader)
	if err != nil {
		h.logger.Error("failed to create a client assets", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a client assets", err.Error(), nil))
	}

	return c.JSON(response.New[clientassets.ClientAssets](true, "success", "", result))
}
