package all

import (
	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

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

func (h *All) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all client assets] execute handler")

	result, err := h.clientAssetsService.All.Execute(c.Context())
	if err != nil {
		h.logger.Error("failed to get all client assets", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all client assets", err.Error(), nil))
	}

	return c.JSON(response.New[[]clientassets.ClientAssets](true, "success", "", result))
}
