package getdailytaskweeksummarybytelegramid

import (
	"context"
	"log"

	userdailytask "github.com/go-jedi/lingramm_backend/internal/domain/user_daily_task"
	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetDailyTaskWeekSummaryByTelegramID --output=mocks --case=underscore
type IGetDailyTaskWeekSummaryByTelegramID interface {
	Execute(ctx context.Context, telegramID string) ([]userdailytask.GetDailyTaskWeekSummaryByTelegramIDResponse, error)
}

type GetDailyTaskWeekSummaryByTelegramID struct {
	userDailyTaskRepository *userdailytaskrepository.Repository
	logger                  logger.ILogger
	postgres                *postgres.Postgres
}

func New(
	userDailyTaskRepository *userdailytaskrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *GetDailyTaskWeekSummaryByTelegramID {
	return &GetDailyTaskWeekSummaryByTelegramID{
		userDailyTaskRepository: userDailyTaskRepository,
		logger:                  logger,
		postgres:                postgres,
	}
}

func (s *GetDailyTaskWeekSummaryByTelegramID) Execute(ctx context.Context, telegramID string) ([]userdailytask.GetDailyTaskWeekSummaryByTelegramIDResponse, error) {
	s.logger.Debug("[get daily task week summary by telegram id] execute service")

	var (
		err    error
		result []userdailytask.GetDailyTaskWeekSummaryByTelegramIDResponse
		ie     bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if !ie { // if daily task does not assign.
		// assign daily task by telegram id.
		if _, err := s.userDailyTaskRepository.AssignDailyTaskByTelegramID.Execute(ctx, tx, telegramID); err != nil {
			return nil, err
		}
	}

	// get daily task week summary by telegram id.
	result, err = s.userDailyTaskRepository.GetDailyTaskWeekSummaryByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return nil, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
