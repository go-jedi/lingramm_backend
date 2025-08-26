package localizedtext

import (
	createtextcontent "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/localized_text/create_text_content"
	createtexttranslation "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/localized_text/create_text_translation"
	gettextsbylanguage "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/localized_text/get_texts_by_language"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	localizedtextservice "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	createTextContent     *createtextcontent.CreateTextContent
	createTextTranslation *createtexttranslation.CreateTextTranslation
	getTextsByLanguage    *gettextsbylanguage.GetTextsByLanguage
}

func New(
	localizedTextService *localizedtextservice.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		createTextContent:     createtextcontent.New(localizedTextService, logger, validator),
		createTextTranslation: createtexttranslation.New(localizedTextService, logger, validator),
		getTextsByLanguage:    gettextsbylanguage.New(localizedTextService, logger, validator),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group("/v1/localized_text")
	{
		api.Post("/content", middleware.Auth.AuthMiddleware, h.createTextContent.Execute)
		api.Post("/translation", middleware.Auth.AuthMiddleware, h.createTextTranslation.Execute)
		api.Get("/texts/language/:language", middleware.Auth.AuthMiddleware, h.getTextsByLanguage.Execute)
	}
}
