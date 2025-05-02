package dependencies

import (
	authhandler "github.com/go-jedi/lingvogramm_backend/internal/adapter/http/handlers/auth"
	"github.com/go-jedi/lingvogramm_backend/internal/middleware"
	userrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/user"
	authservice "github.com/go-jedi/lingvogramm_backend/internal/service/auth"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	"github.com/go-jedi/lingvogramm_backend/pkg/uuid"
	"github.com/go-jedi/lingvogramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	app        *fiber.App
	logger     *logger.Logger
	validator  *validator.Validator
	uuid       *uuid.UUID
	middleware *middleware.Middleware
	postgres   *postgres.Postgres

	//	auth
	authService *authservice.Service
	authHandler *authhandler.Handler

	// user
	userRepository *userrepository.Repository
}

func New(
	app *fiber.App,
	logger *logger.Logger,
	validator *validator.Validator,
	uuid *uuid.UUID,
	postgres *postgres.Postgres,
) *Dependencies {
	d := &Dependencies{
		app:       app,
		logger:    logger,
		validator: validator,
		uuid:      uuid,
		postgres:  postgres,
	}

	d.initMiddleware()
	d.initHandler()

	return d
}

// initMiddleware initialize middlewares.
func (d *Dependencies) initMiddleware() {
	d.middleware = middleware.New()
}

// initHandler initialize handlers.
func (d *Dependencies) initHandler() {
	_ = d.AuthHandler()
}
