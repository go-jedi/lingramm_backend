package dependencies

import (
	"github.com/go-jedi/lingvogramm_backend/internal/middleware"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	"github.com/go-jedi/lingvogramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	app        *fiber.App
	logger     *logger.Logger
	validator  *validator.Validator
	middleware *middleware.Middleware
	postgres   *postgres.Postgres
}

func New(
	app *fiber.App,
	logger *logger.Logger,
	validator *validator.Validator,
	postgres *postgres.Postgres,
) *Dependencies {
	d := &Dependencies{
		app:       app,
		logger:    logger,
		validator: validator,
		postgres:  postgres,
	}

	d.initMiddleware()
	d.initHandler()

	return d
}

// initMiddleware initialize middlewares.
func (d *Dependencies) initMiddleware() {}

// initHandler initialize handlers.
func (d *Dependencies) initHandler() {}
