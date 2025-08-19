package dependencies

import levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"

func (d *Dependencies) LevelRepository() *levelrepository.Repository {
	if d.levelRepository == nil {
		d.levelRepository = levelrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.levelRepository
}
