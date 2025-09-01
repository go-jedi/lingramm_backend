package dependencies

import (
	eventhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/event"
	eventservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event"
)

func (d *Dependencies) EventService() *eventservice.Service {
	if d.eventService == nil {
		d.eventService = eventservice.New(
			d.ExperiencePointRepository(),
			d.UserRepository(),
			d.UserStatsRepository(),
			d.EventTypeRepository(),
			d.LevelRepository(),
			d.InternalCurrencyRepository(),
			d.UserAchievementRepository(),
			d.NotificationRepository(),
			d.logger,
			d.rabbitMQ,
			d.postgres,
			d.redis,
		)
	}

	return d.eventService
}

func (d *Dependencies) EventHandler() *eventhandler.Handler {
	if d.eventHandler == nil {
		d.eventHandler = eventhandler.New(
			d.EventService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.eventHandler
}
