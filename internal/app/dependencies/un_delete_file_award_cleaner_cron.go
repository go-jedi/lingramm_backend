package dependencies

import (
	"context"

	undeletefileawardcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_award_cleaner"
)

func (d *Dependencies) UnDeleteFileAwardCleanerCron(ctx context.Context) *undeletefileawardcleaner.UnDeleteFileAwardCleaner {
	if d.unDeleteFileAwardCleaner == nil {
		d.unDeleteFileAwardCleaner = undeletefileawardcleaner.New(
			ctx,
			d.cfg.Cron,
			d.logger,
			d.redis,
		)
	}

	return d.unDeleteFileAwardCleaner
}
