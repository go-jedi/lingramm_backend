package userachievement

import (
	alldetailbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_achievement/all_detail_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	userachievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	allDetailByTelegramID *alldetailbytelegramid.AllDetailByTelegramID
}

func New(
	userAchievementService *userachievementservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		allDetailByTelegramID: alldetailbytelegramid.New(userAchievementService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/user_achievement",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Get("/all/telegram/:telegramID", h.allDetailByTelegramID.Execute)
	}
}
