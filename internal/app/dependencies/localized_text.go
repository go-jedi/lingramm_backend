package dependencies

import (
	localizedtexthandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/localized_text"
	localizedtextepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text"
	localizedtextservice "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text"
)

func (d *Dependencies) LocalizedTextRepository() *localizedtextepository.Repository {
	if d.localizedTextRepository == nil {
		d.localizedTextRepository = localizedtextepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.localizedTextRepository
}

func (d *Dependencies) LocalizedTextService() *localizedtextservice.Service {
	if d.localizedTextService == nil {
		d.localizedTextService = localizedtextservice.New(
			d.LocalizedTextRepository(),
			d.logger,
			d.postgres,
			d.bigCache,
		)
	}

	return d.localizedTextService
}

func (d *Dependencies) LocalizedTextHandler() *localizedtexthandler.Handler {
	if d.localizedTextHandler == nil {
		d.localizedTextHandler = localizedtexthandler.New(
			d.LocalizedTextService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.localizedTextHandler
}
