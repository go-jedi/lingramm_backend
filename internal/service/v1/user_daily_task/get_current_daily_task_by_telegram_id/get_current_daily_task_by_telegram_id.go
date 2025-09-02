package getcurrentdailytaskbytelegramid

import (
	"context"
	"log"

	userdailytask "github.com/go-jedi/lingramm_backend/internal/domain/user_daily_task"
	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetCurrentDailyTaskByTelegramID --output=mocks --case=underscore
type IGetCurrentDailyTaskByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (userdailytask.GetCurrentDailyTaskByTelegramIDResponse, error)
}

type GetCurrentDailyTaskByTelegramID struct {
	userDailyTaskRepository *userdailytaskrepository.Repository
	logger                  logger.ILogger
	postgres                *postgres.Postgres
}

func New(
	userDailyTaskRepository *userdailytaskrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *GetCurrentDailyTaskByTelegramID {
	return &GetCurrentDailyTaskByTelegramID{
		userDailyTaskRepository: userDailyTaskRepository,
		logger:                  logger,
		postgres:                postgres,
	}
}

func (s *GetCurrentDailyTaskByTelegramID) Execute(ctx context.Context, telegramID string) (userdailytask.GetCurrentDailyTaskByTelegramIDResponse, error) {
	s.logger.Debug("[get current daily task by telegram id] execute service")

	var (
		err             error
		assignDailyTask userdailytask.AssignDailyTaskByTelegramIDResponse
		result          userdailytask.GetCurrentDailyTaskByTelegramIDResponse
		ie              bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return userdailytask.GetCurrentDailyTaskByTelegramIDResponse{}, err
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
		return userdailytask.GetCurrentDailyTaskByTelegramIDResponse{}, err
	}

	if !ie { // if daily task does not assign.
		// assign daily task by telegram id.
		assignDailyTask, err = s.userDailyTaskRepository.AssignDailyTaskByTelegramID.Execute(ctx, tx, telegramID)
		if err != nil {
			return userdailytask.GetCurrentDailyTaskByTelegramIDResponse{}, err
		}

		return assignDailyTask.ConvertToGetCurrentDailyTask(), nil
	}

	// get current daily task by telegram id.
	result, err = s.userDailyTaskRepository.GetCurrentDailyTaskByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return userdailytask.GetCurrentDailyTaskByTelegramIDResponse{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return userdailytask.GetCurrentDailyTaskByTelegramIDResponse{}, err
	}

	return result, nil
}
