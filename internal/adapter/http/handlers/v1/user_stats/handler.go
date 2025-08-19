package userstats

import (
	getlevelbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_stats/get_level_by_telegram_id"
	getlevelinfobytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_stats/get_level_info_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	userstatsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	getLevelByTelegramID     *getlevelbytelegramid.GetLevelByTelegramID
	getLevelInfoByTelegramID *getlevelinfobytelegramid.GetLevelInfoByTelegramID
}

func New(
	userStatsService *userstatsservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		getLevelByTelegramID:     getlevelbytelegramid.New(userStatsService, logger),
		getLevelInfoByTelegramID: getlevelinfobytelegramid.New(userStatsService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/user_stats",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Get("/level/:telegramID", h.getLevelByTelegramID.Execute)
		api.Get("/level_info/:telegramID", h.getLevelInfoByTelegramID.Execute)
	}
}
