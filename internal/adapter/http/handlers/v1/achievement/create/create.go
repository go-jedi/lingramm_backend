package create

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/achievement_assets"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Create struct {
	achievementService *achievementservice.Service
	logger             logger.ILogger
	validator          validator.IValidator
}

func New(
	achievementService *achievementservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *Create {
	return &Create{
		achievementService: achievementService,
		logger:             logger,
		validator:          validator,
	}
}

func (h *Create) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a achievement] execute handler")

	var (
		valueStr      = c.FormValue("value")
		description   = c.FormValue("description")
		code          = c.FormValue("code")
		name          = c.FormValue("name")
		conditionType = c.FormValue("condition_type")
		operator      = c.FormValue("operator")
	)

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		h.logger.Error("failed parse string to int64", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed parse string to int64", err.Error(), nil))
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get the first file for the provided form key", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get the first file for the provided form key", err.Error(), nil))
	}

	contentType := strings.ToLower(fileHeader.Header.Get("Content-Type"))
	if _, ok := achievementassets.SupportedImageTypes[contentType]; !ok {
		h.logger.Error(fmt.Sprintf("unsupported file type: %s", contentType), "error")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "unsupported file type", fmt.Errorf("%w: %s", apperrors.ErrUnsupportedFormat, contentType).Error(), nil))
	}

	dto := achievement.CreateDTO{
		Value:         value,
		FileHeader:    fileHeader,
		Description:   &description,
		Code:          code,
		Name:          name,
		ConditionType: conditionType,
		Operator:      operator,
	}

	if err := h.validator.StructCtx(c, dto); err != nil {
		h.logger.Error("failed to validate struct", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate struct", err.Error(), nil))
	}

	result, err := h.achievementService.Create.Execute(c, dto)
	if err != nil {
		h.logger.Error("failed to create a achievement", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a achievement", err.Error(), nil))
	}

	return c.JSON(response.New[achievement.Detail](true, "success", "", result))
}
