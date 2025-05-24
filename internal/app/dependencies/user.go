package dependencies

import userrepository "github.com/go-jedi/lingramm_backend/internal/repository/user"

func (d *Dependencies) UserRepository() *userrepository.Repository {
	if d.userRepository == nil {
		d.userRepository = userrepository.New(d.postgres.QueryTimeout, d.logger)
	}

	return d.userRepository
}
