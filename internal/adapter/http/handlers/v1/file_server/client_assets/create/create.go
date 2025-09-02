package create

import (
	"context"
	"fmt"
	"time"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

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

// Execute uploads a client asset image (admin).
// @Summary Upload client asset (admin)
// @Description Uploads a single image file (multipart/form-data) to create a client asset. Only supported image content types are accepted.
// @Tags Client asset
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param file formData file true "Image file to upload"
// @Success 200 {object} clientassets.CreateSwaggerResponse "Successful response"
// @Failure 400 {object} clientassets.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} clientassets.ErrorSwaggerResponse "Internal server error"
// @Router /v1/fs/client_assets [post]
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

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.clientAssetsService.Create.Execute(ctxTimeout, fileHeader)
	if err != nil {
		h.logger.Error("failed to create a client assets", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a client assets", err.Error(), nil))
	}

	return c.JSON(response.New[clientassets.ClientAssets](true, "success", "", result))
}
