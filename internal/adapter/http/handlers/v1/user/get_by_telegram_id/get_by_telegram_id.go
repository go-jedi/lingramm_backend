package getbytelegramid

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	userservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetByTelegramID struct {
	userService *userservice.Service
	logger      logger.ILogger
}

func New(
	userService *userservice.Service,
	logger logger.ILogger,
) *GetByTelegramID {
	return &GetByTelegramID{
		userService: userService,
		logger:      logger,
	}
}

func (h *GetByTelegramID) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get user by telegram id] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.userService.GetByTelegramID.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get user by telegram id", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get user by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[user.User](true, "success", "", result))
}
