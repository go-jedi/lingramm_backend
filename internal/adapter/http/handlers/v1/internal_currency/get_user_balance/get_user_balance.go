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

// Execute returns a user's internal currency balance.
// @Summary Get user balance
// @Description Returns the current internal currency balance for the user identified by Telegram ID.
// @Tags Internal currency
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param telegramID path string true "Telegram ID"
// @Success 200 {object} userbalance.GetUserBalanceSwaggerResponse "Successful response"
// @Failure 400 {object} userbalance.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} userbalance.ErrorSwaggerResponse "Internal server error"
// @Router /v1/internal_currency/user/balance/telegram/{telegramID} [get]
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
