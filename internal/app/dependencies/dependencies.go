package dependencies

import (
	adminhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/admin"
	authhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth"
	bigcachehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/bigcache"
	clientassetshandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/file_server/client_assets"
	internalcurrencyhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/internal_currency"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	adminrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/admin"
	clientassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets"
	internalcurrencyrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
	bigcacheservice "github.com/go-jedi/lingramm_backend/internal/service/v1/bigcache"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	internalcurrencyservice "github.com/go-jedi/lingramm_backend/internal/service/v1/internal_currency"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/go-jedi/lingramm_backend/pkg/uuid"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	app        *fiber.App
	logger     *logger.Logger
	validator  *validator.Validator
	uuid       *uuid.UUID
	jwt        *jwt.JWT
	middleware *middleware.Middleware
	postgres   *postgres.Postgres
	redis      *redis.Redis
	bigCache   *bigcachepkg.BigCache
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

	// bigcache
	bigCacheService *bigcacheservice.Service
	bigCacheHandler *bigcachehandler.Handler

	// internal currency
	internalCurrencyRepository *internalcurrencyrepository.Repository
	internalCurrencyService    *internalcurrencyservice.Service
	internalCurrencyHandler    *internalcurrencyhandler.Handler

	// admin
	adminRepository *adminrepository.Repository
	adminService    *adminservice.Service
	adminHandler    *adminhandler.Handler
}

func New(
	app *fiber.App,
	logger *logger.Logger,
	validator *validator.Validator,
	uuid *uuid.UUID,
	jwt *jwt.JWT,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	bigCache *bigcachepkg.BigCache,
	fileServer *fileserver.FileServer,
) *Dependencies {
	d := &Dependencies{
		app:        app,
		logger:     logger,
		validator:  validator,
		uuid:       uuid,
		jwt:        jwt,
		postgres:   postgres,
		redis:      redis,
		bigCache:   bigCache,
		fileServer: fileServer,
	}

	d.initMiddleware()
	d.initHandler()

	return d
}

// initMiddleware initialize middlewares.
func (d *Dependencies) initMiddleware() {
	d.middleware = middleware.New(
		d.AdminService(),
		d.jwt,
		d.redis,
	)
}

// initHandler initialize handlers.
func (d *Dependencies) initHandler() {
	_ = d.AuthHandler()
	_ = d.ClientAssetsHandler()
	_ = d.BigCacheHandler()
	_ = d.InternalCurrencyHandler()
	_ = d.AdminHandler()
}
