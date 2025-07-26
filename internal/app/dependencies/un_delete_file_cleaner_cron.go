package dependencies

import (
	"context"

	undelfilecleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_cleaner"
)

func (d *Dependencies) UnDeleteFileCleaner(ctx context.Context) *undelfilecleaner.UnDeleteFileCleaner {
	if d.unDeleteFileCleaner == nil {
		d.unDeleteFileCleaner = undelfilecleaner.New(
			ctx,
			d.cfg.Cron,
			d.logger,
			d.redis,
		)
	}

	return d.unDeleteFileCleaner
}
