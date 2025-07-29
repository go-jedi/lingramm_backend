package dependencies

import (
	"context"

	undeletefileachievementcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_achievement_cleaner"
)

func (d *Dependencies) UnDeleteFileAchievementCleanerCron(ctx context.Context) *undeletefileachievementcleaner.UnDeleteFileAchievementCleaner {
	if d.unDeleteFileAchievementCleaner == nil {
		d.unDeleteFileAchievementCleaner = undeletefileachievementcleaner.New(
			ctx,
			d.cfg.Cron,
			d.logger,
			d.redis,
		)
	}

	return d.unDeleteFileAchievementCleaner
}
