package dependencies

import (
	"context"

	leaderboardweeksprocessbatch "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/leaderboard_weeks_process_batch"
)

func (d *Dependencies) LeaderboardWeeksProcessBatchCron(ctx context.Context) *leaderboardweeksprocessbatch.LeaderboardWeeksProcessBatch {
	if d.leaderboardWeeksProcessBatch == nil {
		d.leaderboardWeeksProcessBatch = leaderboardweeksprocessbatch.New(
			ctx,
			d.ExperiencePointService(),
			d.cfg.Cron,
			d.logger,
		)
	}

	return d.leaderboardWeeksProcessBatch
}
