package experiencepoint

import (
	createxpevents "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point/create_xp_events"
	getleaderboardtopweek "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point/get_leaderboard_top_week"
	getleaderboardtopweekforuser "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point/get_leaderboard_top_week_for_user"
	leaderboardweeksprocessbatch "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point/leaderboard_weeks_process_batch"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	CreateXPEvents               createxpevents.ICreateXPEvents
	GetLeaderboardTopWeek        getleaderboardtopweek.IGetLeaderboardTopWeek
	GetLeaderboardTopWeekForUser getleaderboardtopweekforuser.IGetLeaderboardTopWeekForUser
	LeaderboardWeeksProcessBatch leaderboardweeksprocessbatch.ILeaderboardWeeksProcessBatch
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		CreateXPEvents:               createxpevents.New(queryTimeout, logger),
		GetLeaderboardTopWeek:        getleaderboardtopweek.New(queryTimeout, logger),
		GetLeaderboardTopWeekForUser: getleaderboardtopweekforuser.New(queryTimeout, logger),
		LeaderboardWeeksProcessBatch: leaderboardweeksprocessbatch.New(queryTimeout, logger),
	}
}
