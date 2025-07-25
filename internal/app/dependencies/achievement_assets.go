package dependencies

import (
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
)

func (d *Dependencies) AchievementAssetsRepository() *achievementassetsrepository.Repository {
	if d.achievementAssetsRepository == nil {
		d.achievementAssetsRepository = achievementassetsrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.achievementAssetsRepository
}
