package dependencies

import (
	userstudiedlanguagehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_studied_language"
	userstudiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language"
	userstudiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language"
)

func (d *Dependencies) UserStudiedLanguageRepository() *userstudiedlanguagerepository.Repository {
	if d.userStudiedLanguageRepository == nil {
		d.userStudiedLanguageRepository = userstudiedlanguagerepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.userStudiedLanguageRepository
}

func (d *Dependencies) UserStudiedLanguageService() *userstudiedlanguageservice.Service {
	if d.userStudiedLanguageService == nil {
		d.userStudiedLanguageService = userstudiedlanguageservice.New(
			d.UserStudiedLanguageRepository(),
			d.UserRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.userStudiedLanguageService
}

func (d *Dependencies) UserStudiedLanguageHandler() *userstudiedlanguagehandler.Handler {
	if d.userStudiedLanguageHandler == nil {
		d.userStudiedLanguageHandler = userstudiedlanguagehandler.New(
			d.UserStudiedLanguageService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.userStudiedLanguageHandler
}
