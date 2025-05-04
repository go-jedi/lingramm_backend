package dependencies

import (
	authhandler "github.com/go-jedi/lingvogramm_backend/internal/adapter/http/handlers/auth"
	clientassetshandler "github.com/go-jedi/lingvogramm_backend/internal/adapter/http/handlers/file_server/client_assets"
	"github.com/go-jedi/lingvogramm_backend/internal/middleware"
	clientassetsrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/file_server/client_assets"
	userrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/user"
	authservice "github.com/go-jedi/lingvogramm_backend/internal/service/auth"
	clientassetsservice "github.com/go-jedi/lingvogramm_backend/internal/service/file_server/client_assets"
	fileserver "github.com/go-jedi/lingvogramm_backend/pkg/file_server"
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
	fileServer *fileserver.FileServer

	// auth
	authService *authservice.Service
	authHandler *authhandler.Handler

	// user
	userRepository *userrepository.Repository

	// client_assets
	clientAssetsRepository *clientassetsrepository.Repository
	clientAssetsService    *clientassetsservice.Service
	clientAssetsHandler    *clientassetshandler.Handler
}

func New(
	app *fiber.App,
	logger *logger.Logger,
	validator *validator.Validator,
	uuid *uuid.UUID,
	postgres *postgres.Postgres,
	fileServer *fileserver.FileServer,
) *Dependencies {
	d := &Dependencies{
		app:        app,
		logger:     logger,
		validator:  validator,
		uuid:       uuid,
		postgres:   postgres,
		fileServer: fileServer,
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
	_ = d.ClientAssetsHandler()
}
