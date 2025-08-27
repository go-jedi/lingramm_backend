package dependencies

import (
	eventtypehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/event_type"
	eventtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type"
	eventtypeservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event_type"
)

func (d *Dependencies) EventTypeRepository() *eventtyperepository.Repository {
	if d.eventTypeRepository == nil {
		d.eventTypeRepository = eventtyperepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.eventTypeRepository
}

func (d *Dependencies) EventTypeService() *eventtypeservice.Service {
	if d.eventTypeService == nil {
		d.eventTypeService = eventtypeservice.New(
			d.EventTypeRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.eventTypeService
}

func (d *Dependencies) EventTypeHandler() *eventtypehandler.Handler {
	if d.eventTypeHandler == nil {
		d.eventTypeHandler = eventtypehandler.New(
			d.EventTypeService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.eventTypeHandler
}
