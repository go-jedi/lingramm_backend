package dailytask

import (
	dailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/daily_task"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/daily_task/create"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	Create create.ICreate
}

func New(
	dailyTaskRepository *dailytaskrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		Create: create.New(dailyTaskRepository, logger, postgres),
	}
}
