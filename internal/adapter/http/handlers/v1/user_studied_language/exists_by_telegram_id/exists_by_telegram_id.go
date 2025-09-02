package existsbytelegramid

import (
	"context"
	"time"

	userstudiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type ExistsByTelegramID struct {
	userStudiedLanguageService *userstudiedlanguageservice.Service
	logger                     logger.ILogger
}

func New(
	userStudiedLanguageService *userstudiedlanguageservice.Service,
	logger logger.ILogger,
) *ExistsByTelegramID {
	return &ExistsByTelegramID{
		userStudiedLanguageService: userStudiedLanguageService,
		logger:                     logger,
	}
}

// Execute checks if a user has any studied language by Telegram ID.
// @Summary Check user studied language by Telegram ID
// @Description Returns true if the specified Telegram ID has at least one studied language, false otherwise.
// @Tags User studied language
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param telegramID path string true "Telegram ID"
// @Success 200 {object} userstudiedlanguage.ExistsByTelegramIDSwaggerResponse "Successful response"
// @Failure 400 {object} userstudiedlanguage.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} userstudiedlanguage.ErrorSwaggerResponse "Internal server error"
// @Router /v1/user_studied_language/exists/{telegramID} [get]
func (h *ExistsByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[check user studied language exists by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.userStudiedLanguageService.ExistsByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to check user studied language exists by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to check user studied language exists by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[bool](true, "success", "", result))
}
