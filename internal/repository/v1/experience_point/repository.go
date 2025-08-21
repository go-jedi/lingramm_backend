package experiencepoint

import (
	createxpevents "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point/create_xp_events"
	leaderboardweeksprocessbatch "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point/leaderboard_weeks_process_batch"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	CreateXPEvents               createxpevents.ICreateXPEvents
	LeaderboardWeeksProcessBatch leaderboardweeksprocessbatch.ILeaderboardWeeksProcessBatch
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		CreateXPEvents:               createxpevents.New(queryTimeout, logger),
		LeaderboardWeeksProcessBatch: leaderboardweeksprocessbatch.New(queryTimeout, logger),
	}
}
