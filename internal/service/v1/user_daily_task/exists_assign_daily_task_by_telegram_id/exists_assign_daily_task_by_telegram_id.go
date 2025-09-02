package existsassigndailytaskbytelegramid

import (
	"context"
	"log"

	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsAssignDailyTaskByTelegramID --output=mocks --case=underscore
type IExistsAssignDailyTaskByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (bool, error)
}

type ExistsAssignDailyTaskByTelegramID struct {
	userDailyTaskRepository *userdailytaskrepository.Repository
	logger                  logger.ILogger
	postgres                *postgres.Postgres
}

func New(
	userDailyTaskRepository *userdailytaskrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *ExistsAssignDailyTaskByTelegramID {
	return &ExistsAssignDailyTaskByTelegramID{
		userDailyTaskRepository: userDailyTaskRepository,
		logger:                  logger,
		postgres:                postgres,
	}
}

func (s *ExistsAssignDailyTaskByTelegramID) Execute(ctx context.Context, telegramID string) (bool, error) {
	s.logger.Debug("[check assign daily task exists by telegram id] execute service")

	var (
		err    error
		result bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check assign daily task exists by telegram id.
	result, err = s.userDailyTaskRepository.ExistsAssignDailyTaskByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return false, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return result, nil
}
