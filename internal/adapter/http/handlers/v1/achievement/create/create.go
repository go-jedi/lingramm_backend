package create

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/achievement_assets"
	awardassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/award_assets"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

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

// Execute creates a new achievement with metadata and two images.
// @Summary Create achievement (admin)
// @Description Creates an achievement with name, type, optional description, and uploads for achievement & award images (multipart/form-data).
// @Tags Achievement
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param name formData string true "Achievement name"
// @Param achievement_type formData string true "Achievement type identifier"
// @Param description formData string false "Optional description"
// @Param file_achievement formData file true "Achievement image file"
// @Param file_award formData file true "Award image file"
// @Success 200 {object} achievement.DetailSwaggerResponse "Successful response"
// @Failure 400 {object} achievement.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} achievement.ErrorSwaggerResponse "Internal server error"
// @Router /v1/achievement [post]
func (h *Create) Execute(c fiber.Ctx) error {
	h.logger.Debug("[create a achievement] execute handler")

	var (
		description     = c.FormValue("description")
		name            = c.FormValue("name")
		achievementType = c.FormValue("achievement_type")
	)

	fileAchievementHeader, err := c.FormFile("file_achievement")
	if err != nil {
		h.logger.Error("failed to get the file achievement for the provided form key", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get the file achievement for the provided form key", err.Error(), nil))
	}

	fileAwardHeader, err := c.FormFile("file_award")
	if err != nil {
		h.logger.Error("failed to get the file award for the provided form key", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get the file award for the provided form key", err.Error(), nil))
	}

	contentTypeAchievement := strings.ToLower(fileAchievementHeader.Header.Get("Content-Type"))
	if _, ok := achievementassets.SupportedImageTypes[contentTypeAchievement]; !ok {
		h.logger.Error(fmt.Sprintf("unsupported file achievement type: %s", contentTypeAchievement), "error")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "unsupported file achievement type", fmt.Errorf("%w: %s", apperrors.ErrUnsupportedFormat, contentTypeAchievement).Error(), nil))
	}

	contentTypeAward := strings.ToLower(fileAwardHeader.Header.Get("Content-Type"))
	if _, ok := awardassets.SupportedImageTypes[contentTypeAward]; !ok {
		h.logger.Error(fmt.Sprintf("unsupported file award type: %s", contentTypeAward), "error")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "unsupported file award type", fmt.Errorf("%w: %s", apperrors.ErrUnsupportedFormat, contentTypeAward).Error(), nil))
	}

	dto := achievement.CreateDTO{
		FileAchievementHeader: fileAchievementHeader,
		FileAwardHeader:       fileAwardHeader,
		Description:           &description,
		Name:                  name,
		AchievementType:       achievementType,
	}

	if err := h.validator.StructCtx(c.RequestCtx(), dto); err != nil {
		h.logger.Error("failed to validate struct", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate struct", err.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.achievementService.Create.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to create a achievement", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to create a achievement", err.Error(), nil))
	}

	return c.JSON(response.New[achievement.Detail](true, "success", "", result))
}
