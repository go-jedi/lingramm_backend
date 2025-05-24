package getuserbalance

import (
	userbalance "github.com/go-jedi/lingramm_backend/internal/domain/user_balance"
	internalcurrency "github.com/go-jedi/lingramm_backend/internal/service/internal_currency"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

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

func (gub *GetUserBalance) Execute(c fiber.Ctx) error {
	gub.logger.Debug("[get user balance] execute handler")

	telegramID := c.Params("telegramID")
	if telegramID == "" {
		gub.logger.Error("failed to get param telegramID", "error", apperrors.ErrParamIsRequired)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to get param telegramID", apperrors.ErrParamIsRequired.Error(), nil))
	}

	result, err := gub.internalCurrencyService.GetUserBalance.Execute(c.Context(), telegramID)
	if err != nil {
		gub.logger.Error("failed to get user balance", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to exists user by telegram id", err.Error(), nil))
	}

	return c.JSON(response.New[userbalance.UserBalance](true, "success", "", result))
}
