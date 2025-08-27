package dependencies

import achievementtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement_type"

func (d *Dependencies) AchievementTypeRepository() *achievementtyperepository.Repository {
	if d.achievementTypeRepository == nil {
		d.achievementTypeRepository = achievementtyperepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.achievementTypeRepository
}
