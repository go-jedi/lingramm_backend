package dependencies

import (
	userstatshandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_stats"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	userstatsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats"
)

func (d *Dependencies) UserStatsRepository() *userstatsrepository.Repository {
	if d.userStatsRepository == nil {
		d.userStatsRepository = userstatsrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.userStatsRepository
}

func (d *Dependencies) UserStatsService() *userstatsservice.Service {
	if d.userStatsService == nil {
		d.userStatsService = userstatsservice.New(
			d.UserStatsRepository(),
			d.UserRepository(),
			d.UserAchievementRepository(),
			d.NotificationRepository(),
			d.logger,
			d.rabbitMQ,
			d.postgres,
			d.redis,
		)
	}

	return d.userStatsService
}

func (d *Dependencies) UserStatsHandler() *userstatshandler.Handler {
	if d.userStatsHandler == nil {
		d.userStatsHandler = userstatshandler.New(
			d.UserStatsService(),
			d.app,
			d.logger,
			d.middleware,
		)
	}

	return d.userStatsHandler
}
