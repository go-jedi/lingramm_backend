package experiencepoint

import (
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	createxpevents "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point/create_xp_events"
	getleaderboardtopweek "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point/get_leaderboard_top_week"
	getleaderboardtopweekforuser "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point/get_leaderboard_top_week_for_user"
	leaderboardweeksprocessbatch "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point/leaderboard_weeks_process_batch"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	CreateXPEvents               createxpevents.ICreateXPEvents
	GetLeaderboardTopWeek        getleaderboardtopweek.IGetLeaderboardTopWeek
	GetLeaderboardTopWeekForUser getleaderboardtopweekforuser.IGetLeaderboardTopWeekForUser
	LeaderboardWeeksProcessBatch leaderboardweeksprocessbatch.ILeaderboardWeeksProcessBatch
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
		CreateXPEvents:               createxpevents.New(experiencePointRepository, userRepository, userStatsRepository, levelRepository, logger, postgres),
		GetLeaderboardTopWeek:        getleaderboardtopweek.New(experiencePointRepository, logger, postgres),
		GetLeaderboardTopWeekForUser: getleaderboardtopweekforuser.New(experiencePointRepository, logger, postgres),
		LeaderboardWeeksProcessBatch: leaderboardweeksprocessbatch.New(experiencePointRepository, logger, postgres),
	}
}
