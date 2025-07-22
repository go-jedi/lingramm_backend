package auth

import (
	checkhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth/check"
	refreshhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth/refresh"
	signinhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth/sign_in"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	signIn  *signinhandler.SignIn
	check   *checkhandler.Check
	refresh *refreshhandler.Refresh
}

func New(
	authService *auth.Service,
	app *fiber.App,
	logger logger.ILogger,
	validator validator.IValidator,
) *Handler {
	h := &Handler{
		signIn:  signinhandler.New(authService, logger, validator),
		check:   checkhandler.New(authService, logger, validator),
		refresh: refreshhandler.New(authService, logger, validator),
	}

	h.initRoutes(app)

	return h
}

func (h *Handler) initRoutes(app *fiber.App) {
	api := app.Group("/v1/auth")
	{
		api.Post("/signin", h.signIn.Execute)
		api.Post("/check", h.check.Execute)
		api.Post("/refresh", h.refresh.Execute)
	}
}
