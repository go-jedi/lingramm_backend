package dependencies

import (
	subscriptionhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/subscription"
	subscriptionrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription"
	subscriptionservice "github.com/go-jedi/lingramm_backend/internal/service/v1/subscription"
)

func (d *Dependencies) SubscriptionRepository() *subscriptionrepository.Repository {
	if d.subscriptionRepository == nil {
		d.subscriptionRepository = subscriptionrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.subscriptionRepository
}

func (d *Dependencies) SubscriptionService() *subscriptionservice.Service {
	if d.subscriptionService == nil {
		d.subscriptionService = subscriptionservice.New(
			d.SubscriptionRepository(),
			d.UserRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.subscriptionService
}

func (d *Dependencies) SubscriptionHandler() *subscriptionhandler.Handler {
	if d.subscriptionHandler == nil {
		d.subscriptionHandler = subscriptionhandler.New(
			d.SubscriptionService(),
			d.app,
			d.logger,
			d.middleware,
		)
	}

	return d.subscriptionHandler
}
