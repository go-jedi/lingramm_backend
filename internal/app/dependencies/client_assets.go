package dependencies

import (
	clientassetshandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/file_server/client_assets"
	clientassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
)

func (d *Dependencies) ClientAssetsRepository() *clientassetsrepository.Repository {
	if d.clientAssetsRepository == nil {
		d.clientAssetsRepository = clientassetsrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.clientAssetsRepository
}

func (d *Dependencies) ClientAssetsService() *clientassetsservice.Service {
	if d.clientAssetsService == nil {
		d.clientAssetsService = clientassetsservice.New(
			d.ClientAssetsRepository(),
			d.logger,
			d.postgres,
			d.redis,
			d.fileServer,
		)
	}

	return d.clientAssetsService
}

func (d *Dependencies) ClientAssetsHandler() *clientassetshandler.Handler {
	if d.clientAssetsHandler == nil {
		d.clientAssetsHandler = clientassetshandler.New(
			d.ClientAssetsService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.clientAssetsHandler
}
