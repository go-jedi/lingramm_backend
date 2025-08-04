package dependencies

import (
	awardassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/award_assets"
)

func (d *Dependencies) AwardAssetsRepository() *awardassetsrepository.Repository {
	if d.awardAssetsRepository == nil {
		d.awardAssetsRepository = awardassetsrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.awardAssetsRepository
}
