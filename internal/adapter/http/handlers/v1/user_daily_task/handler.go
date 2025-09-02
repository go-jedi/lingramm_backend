package userdailytask

import (
	getcurrentdailytaskbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_daily_task/get_current_daily_task_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	userdailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	getCurrentDailyTaskByTelegramID *getcurrentdailytaskbytelegramid.GetCurrentDailyTaskByTelegramID
}

func New(
	userDailyTaskService *userdailytaskservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		getCurrentDailyTaskByTelegramID: getcurrentdailytaskbytelegramid.New(userDailyTaskService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/user_daily_task",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Get("/telegram/:telegramID", h.getCurrentDailyTaskByTelegramID.Execute)
	}
}
