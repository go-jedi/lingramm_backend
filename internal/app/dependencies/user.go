package dependencies

import userrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/user"

func (d *Dependencies) UserRepository() *userrepository.Repository {
	if d.userRepository == nil {
		d.userRepository = userrepository.New(d.logger)
	}

	return d.userRepository
}
