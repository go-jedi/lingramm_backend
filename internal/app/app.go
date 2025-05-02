package app

import (
	"context"

	"github.com/go-jedi/lingvogramm_backend/config"
	"github.com/go-jedi/lingvogramm_backend/internal/app/dependencies"
	"github.com/go-jedi/lingvogramm_backend/pkg/httpserver"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	"github.com/go-jedi/lingvogramm_backend/pkg/uuid"
	"github.com/go-jedi/lingvogramm_backend/pkg/validator"
)

type App struct {
	cfg          config.Config
	logger       *logger.Logger
	validator    *validator.Validator
	uuid         *uuid.UUID
	postgres     *postgres.Postgres
	hs           *httpserver.HTTPServer
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
		a.initPostgres,
		a.initHTTPServer,
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

// initPostgres initialize postgres.
func (a *App) initPostgres(ctx context.Context) (err error) {
	a.postgres, err = postgres.New(ctx, a.cfg.Postgres, a.logger)
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

// initDependencies initialize dependencies.
func (a *App) initDependencies(_ context.Context) error {
	a.dependencies = dependencies.New(
		a.hs.App,
		a.logger,
		a.validator,
		a.uuid,
		a.postgres,
	)

	return nil
}

// runHTTPServer run http server.
func (a *App) runHTTPServer() error {
	return a.hs.Start()
}
