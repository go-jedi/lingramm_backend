package dependencies

import (
	internalcurrencyhandler "github.com/go-jedi/lingvogramm_backend/internal/adapter/http/handlers/internal_currency"
	internalcurrencyrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/internal_currency"
	internalcurrencyservice "github.com/go-jedi/lingvogramm_backend/internal/service/internal_currency"
)

func (d *Dependencies) InternalCurrencyRepository() *internalcurrencyrepository.Repository {
	if d.internalCurrencyRepository == nil {
		d.internalCurrencyRepository = internalcurrencyrepository.New(d.postgres.QueryTimeout, d.logger)
	}

	return d.internalCurrencyRepository
}

func (d *Dependencies) InternalCurrencyService() *internalcurrencyservice.Service {
	if d.internalCurrencyService == nil {
		d.internalCurrencyService = internalcurrencyservice.New(
			d.InternalCurrencyRepository(),
			d.UserRepository(),
			d.logger,
			d.postgres,
			d.bigCache,
		)
	}

	return d.internalCurrencyService
}

func (d *Dependencies) InternalCurrencyHandler() *internalcurrencyhandler.Handler {
	if d.internalCurrencyHandler == nil {
		d.internalCurrencyHandler = internalcurrencyhandler.New(
			d.InternalCurrencyService(),
			d.app,
			d.logger,
			d.validator,
		)
	}

	return d.internalCurrencyHandler
}
