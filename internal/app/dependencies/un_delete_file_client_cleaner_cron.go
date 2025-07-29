package dependencies

import (
	"context"

	undeletefileclientcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_client_cleaner"
)

func (d *Dependencies) UnDeleteFileClientCleanerCron(ctx context.Context) *undeletefileclientcleaner.UnDeleteFileClientCleaner {
	if d.unDeleteFileClientCleaner == nil {
		d.unDeleteFileClientCleaner = undeletefileclientcleaner.New(
			ctx,
			d.cfg.Cron,
			d.logger,
			d.redis,
		)
	}

	return d.unDeleteFileClientCleaner
}
