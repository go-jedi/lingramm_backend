package dailytask

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/daily_task/create"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	Create create.ICreate
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		Create: create.New(queryTimeout, logger),
	}
}
