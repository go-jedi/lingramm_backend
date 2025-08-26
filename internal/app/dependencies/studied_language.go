package dependencies

import (
	studiedlanguagehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/studied_language"
	studiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/studied_language"
	studiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/studied_language"
)

func (d *Dependencies) StudiedLanguageRepository() *studiedlanguagerepository.Repository {
	if d.studiedLanguageRepository == nil {
		d.studiedLanguageRepository = studiedlanguagerepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.studiedLanguageRepository
}

func (d *Dependencies) StudiedLanguageService() *studiedlanguageservice.Service {
	if d.studiedLanguageService == nil {
		d.studiedLanguageService = studiedlanguageservice.New(
			d.StudiedLanguageRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.studiedLanguageService
}

func (d *Dependencies) StudiedLanguageHandler() *studiedlanguagehandler.Handler {
	if d.studiedLanguageHandler == nil {
		d.studiedLanguageHandler = studiedlanguagehandler.New(
			d.StudiedLanguageService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.studiedLanguageHandler
}
