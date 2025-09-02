package getbytelegramid

import (
	"context"
	"time"

	userstudiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/user_studied_language"
	userstudiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetByTelegramID struct {
	userStudiedLanguageService *userstudiedlanguageservice.Service
	logger                     logger.ILogger
}

func New(
	userStudiedLanguageService *userstudiedlanguageservice.Service,
	logger logger.ILogger,
) *GetByTelegramID {
	return &GetByTelegramID{
		userStudiedLanguageService: userStudiedLanguageService,
		logger:                     logger,
	}
}

// Execute returns a user's studied language by Telegram ID.
// @Summary Get user studied language by Telegram ID
// @Description Returns the studied language record linked to the specified Telegram ID, including language metadata.
// @Tags User studied language
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param telegramID path string true "Telegram ID"
// @Success 200 {object} userstudiedlanguage.GetByTelegramIDSwaggerResponse "Successful response"
// @Failure 400 {object} userstudiedlanguage.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} userstudiedlanguage.ErrorSwaggerResponse "Internal server error"
// @Router /v1/user_studied_language/telegram/{telegramID} [get]
func (h *GetByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get user studied language by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.userStudiedLanguageService.GetByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get user studied language by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get user studied language by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[userstudiedlanguage.GetByTelegramIDResponse](true, "success", "", result))
}
