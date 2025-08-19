package dependencies

import (
	experiencepointhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/experience_point"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	experiencepointservice "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point"
)

func (d *Dependencies) ExperiencePointRepository() *experiencepointrepository.Repository {
	if d.experiencePointRepository == nil {
		d.experiencePointRepository = experiencepointrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.experiencePointRepository
}

func (d *Dependencies) ExperiencePointService() *experiencepointservice.Service {
	if d.experiencePointService == nil {
		d.experiencePointService = experiencepointservice.New(
			d.ExperiencePointRepository(),
			d.UserRepository(),
			d.UserStatsRepository(),
			d.LevelRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.experiencePointService
}

func (d *Dependencies) ExperiencePointHandler() *experiencepointhandler.Handler {
	if d.experiencePointHandler == nil {
		d.experiencePointHandler = experiencepointhandler.New(
			d.ExperiencePointService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.experiencePointHandler
}
