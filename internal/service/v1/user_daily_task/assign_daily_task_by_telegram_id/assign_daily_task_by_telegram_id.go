package assigndailytaskbytelegramid

import (
	"context"
	"log"

	userdailytask "github.com/go-jedi/lingramm_backend/internal/domain/user_daily_task"
	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAssignDailyTaskByTelegramID --output=mocks --case=underscore
type IAssignDailyTaskByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (userdailytask.AssignDailyTaskByTelegramIDResponse, error)
}

type AssignDailyTaskByTelegramID struct {
	userDailyTaskRepository *userdailytaskrepository.Repository
	logger                  logger.ILogger
	postgres                *postgres.Postgres
}

func New(
	userDailyTaskRepository *userdailytaskrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *AssignDailyTaskByTelegramID {
	return &AssignDailyTaskByTelegramID{
		userDailyTaskRepository: userDailyTaskRepository,
		logger:                  logger,
		postgres:                postgres,
	}
}

func (s *AssignDailyTaskByTelegramID) Execute(ctx context.Context, telegramID string) (userdailytask.AssignDailyTaskByTelegramIDResponse, error) {
	s.logger.Debug("[assign daily task by telegram id] execute service")

	var (
		err    error
		result userdailytask.AssignDailyTaskByTelegramIDResponse
		ie     bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return userdailytask.AssignDailyTaskByTelegramIDResponse{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check assign daily task exists by telegram id.
	ie, err = s.userDailyTaskRepository.ExistsAssignDailyTaskByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return userdailytask.AssignDailyTaskByTelegramIDResponse{}, err
	}

	if ie { // if assign daily task already exist.
		err = apperrors.ErrDailyTaskAlreadyAssign
		return userdailytask.AssignDailyTaskByTelegramIDResponse{}, err
	}

	// assign daily task by telegram id.
	result, err = s.userDailyTaskRepository.AssignDailyTaskByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return userdailytask.AssignDailyTaskByTelegramIDResponse{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return userdailytask.AssignDailyTaskByTelegramIDResponse{}, err
	}

	return result, nil
}
