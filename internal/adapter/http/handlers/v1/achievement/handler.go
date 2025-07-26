package achievement

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/achievement/create"
	getdetailbyachievementid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/achievement/get_detail_by_achievement_id"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	create                   *create.Create
	getDetailByAchievementID *getdetailbyachievementid.GetDetailByAchievementID
}

func New(
	achievementService *achievementservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		create:                   create.New(achievementService, logger, validator),
		getDetailByAchievementID: getdetailbyachievementid.New(achievementService, logger),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/achievement",
		middleware.Auth.AuthMiddleware,
		middleware.AdminGuard.AdminGuardMiddleware,
	)
	{
		api.Post(
			"",
			middleware.ContentLengthLimiter.ContentLengthLimiterMiddleware,
			h.create.Execute,
		)
		api.Get("/id/:achievementID", h.getDetailByAchievementID.Execute)
	}
}
