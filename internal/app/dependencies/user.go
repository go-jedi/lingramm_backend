package dependencies

import (
	userhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user"
)

func (d *Dependencies) UserRepository() *userrepository.Repository {
	if d.userRepository == nil {
		d.userRepository = userrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.userRepository
}

func (d *Dependencies) UserService() *userservice.Service {
	if d.userService == nil {
		d.userService = userservice.New(
			d.UserRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.userService
}

func (d *Dependencies) UserHandler() *userhandler.Handler {
	if d.userHandler == nil {
		d.userHandler = userhandler.New(
			d.UserService(),
			d.app,
			d.logger,
			d.middleware,
		)
	}

	return d.userHandler
}
