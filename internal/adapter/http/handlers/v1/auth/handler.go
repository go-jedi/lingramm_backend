package auth

import (
	"github.com/go-jedi/lingramm_backend/config"
	checkhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth/check"
	refreshhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth/refresh"
	signinhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth/sign_in"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	check   *checkhandler.Check
	refresh *refreshhandler.Refresh
	signIn  *signinhandler.SignIn
}

func New(
	authService *auth.Service,
	app *fiber.App,
	cookie config.CookieConfig,
	logger logger.ILogger,
	validator validator.IValidator,
	middleware *middleware.Middleware,
) *Handler {
	h := &Handler{
		check:   checkhandler.New(authService, logger, validator),
		refresh: refreshhandler.New(authService, logger, validator),
		signIn:  signinhandler.New(authService, cookie, logger, validator),
	}

	h.initRoutes(app, middleware)

	return h
}

func (h *Handler) initRoutes(app *fiber.App, middleware *middleware.Middleware) {
	api := app.Group("/v1/auth")
	{
		api.Post("/signin", h.signIn.Execute)
		api.Post("/check", middleware.Auth.AuthMiddleware, h.check.Execute)
		api.Post("/refresh", middleware.Auth.AuthMiddleware, h.refresh.Execute)
	}
}
