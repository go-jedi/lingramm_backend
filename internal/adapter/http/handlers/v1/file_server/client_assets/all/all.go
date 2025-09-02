package all

import (
	"context"
	"time"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type All struct {
	clientAssetsService *clientassetsservice.Service
	logger              logger.ILogger
	validator           validator.IValidator
}

func New(
	clientAssetsService *clientassetsservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *All {
	return &All{
		clientAssetsService: clientAssetsService,
		logger:              logger,
		validator:           validator,
	}
}

// Execute returns all client assets (admin).
// @Summary Get all client assets (admin)
// @Description Returns a full list of uploaded client asset images and their metadata.
// @Tags Client asset
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Success 200 {object} clientassets.AllSwaggerResponse "Successful response"
// @Failure 500 {object} clientassets.ErrorSwaggerResponse "Internal server error"
// @Router /v1/fs/client_assets/all [get]
func (h *All) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all client assets] execute handler")

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.clientAssetsService.All.Execute(ctxTimeout)
	if err != nil {
		h.logger.Error("failed to get all client assets", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all client assets", err.Error(), nil))
	}

	return c.JSON(response.New[[]clientassets.ClientAssets](true, "success", "", result))
}
