package dependencies

import (
	authhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
)

func (d *Dependencies) AuthService() *authservice.Service {
	if d.authService == nil {
		d.authService = authservice.New(
			d.UserRepository(),
			d.LevelRepository(),
			d.logger,
			d.postgres,
			d.redis,
			d.bigCache,
			d.jwt,
		)
	}

	return d.authService
}

func (d *Dependencies) AuthHandler() *authhandler.Handler {
	if d.authHandler == nil {
		d.authHandler = authhandler.New(
			d.AuthService(),
			d.app,
			d.cfg.Cookie,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.authHandler
}
