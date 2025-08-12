package dependencies

import (
	notificationhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/notification"
	notificationwebsockethandler "github.com/go-jedi/lingramm_backend/internal/adapter/websocket/handlers/v1/notification"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
)

func (d *Dependencies) NotificationRepository() *notificationrepository.Repository {
	if d.notificationRepository == nil {
		d.notificationRepository = notificationrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.notificationRepository
}

func (d *Dependencies) NotificationService() *notificationservice.Service {
	if d.notificationService == nil {
		d.notificationService = notificationservice.New(
			d.NotificationRepository(),
			d.UserRepository(),
			d.logger,
			d.rabbitMQ,
			d.postgres,
			d.redis,
		)
	}

	return d.notificationService
}

func (d *Dependencies) NotificationHandler() *notificationhandler.Handler {
	if d.notificationHandler == nil {
		d.notificationHandler = notificationhandler.New(
			d.NotificationService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.notificationHandler
}

func (d *Dependencies) NotificationWebSocket() *notificationwebsockethandler.Handler {
	if d.notificationWebSocketHandler == nil {
		d.notificationWebSocketHandler = notificationwebsockethandler.New(
			d.NotificationService(),
			d.app,
			d.logger,
			d.rabbitMQ,
			d.redis,
			d.wsManager,
			d.middleware,
		)
	}

	return d.notificationWebSocketHandler
}
