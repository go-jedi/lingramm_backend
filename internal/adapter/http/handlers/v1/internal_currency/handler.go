package internalcurrency

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/internal_currency/get_user_balance"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/internal_currency"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	getUserBalance *getuserbalance.GetUserBalance
}

func New(
	internalCurrencyService *internalcurrency.Service,
	app *fiber.App,
	logger logger.ILogger,
	_ validator.IValidator,
) *Handler {
	h := &Handler{
		getUserBalance: getuserbalance.New(internalCurrencyService, logger),
	}

	h.initRoutes(app)

	return h
}

func (h *Handler) initRoutes(app *fiber.App) {
	api := app.Group("/v1/internal_currency")
	{
		api.Get("/user/balance/:telegramID", h.getUserBalance.Execute)
	}
}
