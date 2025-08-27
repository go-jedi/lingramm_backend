package achievementtype

import (
	existsbyname "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement_type/exists_by_name"
	getbyname "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement_type/get_by_name"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	ExistsByName existsbyname.IExistsByName
	GetByName    getbyname.IGetByName
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		ExistsByName: existsbyname.New(queryTimeout, logger),
		GetByName:    getbyname.New(queryTimeout, logger),
	}
}
