package app

import (
	"context"
	"os"
	"strings"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/go-jedi/lingramm_backend/internal/app/dependencies"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/httpserver"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/go-jedi/lingramm_backend/pkg/uuid"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type App struct {
	cfg          config.Config
	logger       *logger.Logger
	validator    *validator.Validator
	uuid         *uuid.UUID
	jwt          *jwt.JWT
	postgres     *postgres.Postgres
	redis        *redis.Redis
	bigCache     *bigcachepkg.BigCache
	hs           *httpserver.HTTPServer
	fileServer   *fileserver.FileServer
	dependencies *dependencies.Dependencies
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run application.
func (a *App) Run() error {
	return a.runHTTPServer()
}

// initDeps initialize deps.
func (a *App) initDeps(ctx context.Context) error {
	fn := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initValidator,
		a.initUUID,
		a.initJWT,
		a.initPostgres,
		a.initRedis,
		a.initBigCache,
		a.initHTTPServer,
		a.initFileServer,
		a.initDependencies,
	}

	for i := range fn {
		if err := fn[i](ctx); err != nil {
			return err
		}
	}

	return nil
}

// initConfig initialize config.
func (a *App) initConfig(_ context.Context) (err error) {
	a.cfg, err = config.GetConfig()
	if err != nil {
		return err
	}

	return
}

// initLogger initialize logger.
func (a *App) initLogger(_ context.Context) error {
	a.logger = logger.New(a.cfg.Logger)
	return nil
}

// initValidator initialize validator.
func (a *App) initValidator(_ context.Context) error {
	a.validator = validator.New()
	return nil
}

// initUUID initialize uuid.
func (a *App) initUUID(_ context.Context) error {
	a.uuid = uuid.New()
	return nil
}

// initJWT initialize jwt.
func (a *App) initJWT(_ context.Context) (err error) {
	a.jwt, err = jwt.New(a.cfg.JWT, a.uuid)
	if err != nil {
		return err
	}

	return
}

// initPostgres initialize postgres.
func (a *App) initPostgres(ctx context.Context) (err error) {
	a.postgres, err = postgres.New(ctx, a.cfg.Postgres, a.logger)
	if err != nil {
		return err
	}

	return
}

// initRedis initialize redis.
func (a *App) initRedis(ctx context.Context) (err error) {
	a.redis, err = redis.New(ctx, a.cfg.Redis)
	if err != nil {
		return err
	}

	return
}

// initBigCache initialize big cache.
func (a *App) initBigCache(_ context.Context) (err error) {
	a.bigCache, err = bigcachepkg.New(a.cfg.BigCache)
	if err != nil {
		return err
	}

	return
}

// initHTTPServer initialize http server.
func (a *App) initHTTPServer(_ context.Context) (err error) {
	a.hs, err = httpserver.New(a.cfg.HTTPServer)
	if err != nil {
		return err
	}

	return
}

// initFileServer initialize file server.
func (a *App) initFileServer(_ context.Context) error {
	a.fileServer = fileserver.New(a.cfg.FileServer, a.uuid)

	clientAssetsStaticCfg := static.Config{
		FS:       os.DirFS(a.cfg.FileServer.ClientAssets.Dir),
		Browse:   a.cfg.FileServer.ClientAssets.Browse,
		Compress: a.cfg.FileServer.ClientAssets.Compress,
	}
	achievementAssetsStaticCfg := static.Config{
		FS:       os.DirFS(a.cfg.FileServer.AchievementAssets.Dir),
		Browse:   a.cfg.FileServer.AchievementAssets.Browse,
		Compress: a.cfg.FileServer.AchievementAssets.Compress,
	}

	if a.cfg.FileServer.ClientAssets.IsNext {
		clientAssetsStaticCfg.Next = func(c fiber.Ctx) bool { // need don't show any format.
			return strings.HasSuffix(c.Path(), a.cfg.FileServer.ClientAssets.IsNextIgnoreFormat)
		}
	}
	if a.cfg.FileServer.AchievementAssets.IsNext {
		achievementAssetsStaticCfg.Next = func(c fiber.Ctx) bool { // need don't show any format.
			return strings.HasSuffix(c.Path(), a.cfg.FileServer.AchievementAssets.IsNextIgnoreFormat)
		}
	}

	a.hs.App.Get(a.cfg.FileServer.ClientAssets.Path, static.New("", clientAssetsStaticCfg))           // initialize static for client assets.
	a.hs.App.Get(a.cfg.FileServer.AchievementAssets.Path, static.New("", achievementAssetsStaticCfg)) // initialize static for achievement assets.

	return nil
}

// initDependencies initialize dependencies.
func (a *App) initDependencies(ctx context.Context) error {
	a.dependencies = dependencies.New(
		ctx,
		a.cfg,
		a.hs.App,
		a.logger,
		a.validator,
		a.uuid,
		a.jwt,
		a.postgres,
		a.redis,
		a.bigCache,
		a.fileServer,
	)

	return nil
}

// runHTTPServer run http server.
func (a *App) runHTTPServer() error {
	return a.hs.Start()
}
