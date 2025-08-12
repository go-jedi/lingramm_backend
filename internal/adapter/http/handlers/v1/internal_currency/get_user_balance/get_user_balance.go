package getuserbalance

import (
	"context"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/internal_currency/user_balance"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/internal_currency"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetUserBalance struct {
	internalCurrencyService *internalcurrency.Service
	logger                  logger.ILogger
}

func New(
	internalCurrencyService *internalcurrency.Service,
	logger logger.ILogger,
) *GetUserBalance {
	return &GetUserBalance{
		internalCurrencyService: internalCurrencyService,
		logger:                  logger,
	}
}

func (h *GetUserBalance) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get user balance] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		h.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.internalCurrencyService.GetUserBalance.Execute(ctxTimeout, telegramID)
	if err != nil {
		h.logger.Error("failed to get user balance", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get user balance", err.Error(), nil))
	}

	return c.JSON(response.New[userbalance.UserBalance](true, "success", "", result))
}
