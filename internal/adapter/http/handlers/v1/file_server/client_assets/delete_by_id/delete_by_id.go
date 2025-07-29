package deletebyid

import (
	"strconv"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type DeleteByID struct {
	clientAssetsService *clientassetsservice.Service
	logger              logger.ILogger
}

func New(
	clientAssetsService *clientassetsservice.Service,
	logger logger.ILogger,
) *DeleteByID {
	return &DeleteByID{
		clientAssetsService: clientAssetsService,
		logger:              logger,
	}
}

func (h *DeleteByID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[delete client assets by id] execute handler")

	idStr := c.Params("id")
	if idStr == "" {
		h.logger.Error("failed to get param id", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param id", apperrors.ErrParamIsRequired.Error(), nil))
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("failed parse string to int64", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed parse string to int64", err.Error(), nil))
	}

	if id <= 0 {
		h.logger.Error("invalid id", "error", "client assets id must be a positive integer")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "invalid client assets id", "client assets id must be a positive integer", nil))
	}

	result, err := h.clientAssetsService.DeleteByID.Execute(c, id)
	if err != nil {
		h.logger.Error("failed to delete client assets by id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to delete client assets by id", err.Error(), nil))
	}

	return c.JSON(response.New[clientassets.ClientAssets](true, "success", "", result))
}
