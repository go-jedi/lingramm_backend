package dependencies

import (
	achievementhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/achievement"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
)

func (d *Dependencies) AchievementRepository() *achievementrepository.Repository {
	if d.achievementRepository == nil {
		d.achievementRepository = achievementrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.achievementRepository
}

func (d *Dependencies) AchievementService() *achievementservice.Service {
	if d.achievementService == nil {
		d.achievementService = achievementservice.New(
			d.AchievementRepository(),
			d.AchievementAssetsRepository(),
			d.logger,
			d.postgres,
			d.redis,
			d.fileServer,
		)
	}

	return d.achievementService
}

func (d *Dependencies) AchievementHandler() *achievementhandler.Handler {
	if d.achievementHandler == nil {
		d.achievementHandler = achievementhandler.New(
			d.AchievementService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.achievementHandler
}
