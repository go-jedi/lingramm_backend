package achievementassets

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets/all"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets/create"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	Create create.ICreate
	All    all.IAll
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		Create: create.New(queryTimeout, logger),
		All:    all.New(queryTimeout, logger),
	}
}
