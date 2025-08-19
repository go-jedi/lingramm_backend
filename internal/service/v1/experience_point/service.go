package experiencepoint

import (
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	createxpevents "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point/create_xp_events"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	CreateXPEvents createxpevents.ICreateXPEvents
}

func New(
	experiencePointRepository *experiencepointrepository.Repository,
	userRepository *userrepository.Repository,
	userStatsRepository *userstatsrepository.Repository,
	levelRepository *levelrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		CreateXPEvents: createxpevents.New(experiencePointRepository, userRepository, userStatsRepository, levelRepository, logger, postgres),
	}
}
