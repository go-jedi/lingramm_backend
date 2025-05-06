package dependencies

import (
	authhandler "github.com/go-jedi/lingvogramm_backend/internal/adapter/http/handlers/auth"
	authservice "github.com/go-jedi/lingvogramm_backend/internal/service/auth"
)

func (d *Dependencies) AuthService() *authservice.Service {
	if d.authService == nil {
		d.authService = authservice.New(
			d.UserRepository(),
			d.logger,
			d.postgres,
			d.bigCache,
			d.uuid,
		)
	}

	return d.authService
}

func (d *Dependencies) AuthHandler() *authhandler.Handler {
	if d.authHandler == nil {
		d.authHandler = authhandler.New(
			d.AuthService(),
			d.app,
			d.logger,
			d.validator,
		)
	}

	return d.authHandler
}
