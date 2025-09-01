package dependencies

import (
	userachievementhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_achievement"
	userachievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement"
	userachievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_achievement"
)

func (d *Dependencies) UserAchievementRepository() *userachievementrepository.Repository {
	if d.userAchievementRepository == nil {
		d.userAchievementRepository = userachievementrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.userAchievementRepository
}

func (d *Dependencies) UserAchievementService() *userachievementservice.Service {
	if d.userAchievementService == nil {
		d.userAchievementService = userachievementservice.New(
			d.UserAchievementRepository(),
			d.UserRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.userAchievementService
}

func (d *Dependencies) UserAchievementHandler() *userachievementhandler.Handler {
	if d.userAchievementHandler == nil {
		d.userAchievementHandler = userachievementhandler.New(
			d.UserAchievementService(),
			d.app,
			d.logger,
			d.middleware,
		)
	}

	return d.userAchievementHandler
}
