package userstudiedlanguage

import (
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_studied_language/create"
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_studied_language/exists_by_telegram_id"
	getbytelegramid "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_studied_language/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_studied_language/update"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	userstudiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	create             *create.Create
	existsByTelegramID *existsbytelegramid.ExistsByTelegramID
	getByTelegramID    *getbytelegramid.GetByTelegramID
	update             *update.Update
}

func New(
	userStudiedLanguageService *userstudiedlanguageservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		create:             create.New(userStudiedLanguageService, logger, validator),
		existsByTelegramID: existsbytelegramid.New(userStudiedLanguageService, logger),
		getByTelegramID:    getbytelegramid.New(userStudiedLanguageService, logger),
		update:             update.New(userStudiedLanguageService, logger, validator),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group(
		"/v1/user_studied_language",
		middleware.Auth.AuthMiddleware,
	)
	{
		api.Post("", h.create.Execute)
		api.Get("/exists/:telegramID", h.existsByTelegramID.Execute)
		api.Get("/telegram/:telegramID", h.getByTelegramID.Execute)
	}
}
