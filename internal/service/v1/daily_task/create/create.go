package create

import (
	"context"
	"log"

	dailytask "github.com/go-jedi/lingramm_backend/internal/domain/daily_task"
	dailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, dto dailytask.CreateDTO) (dailytask.DailyTask, error)
}

type Create struct {
	dailyTaskRepository *dailytaskrepository.Repository
	logger              logger.ILogger
	postgres            *postgres.Postgres
}

func New(
	dailyTaskRepository *dailytaskrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Create {
	return &Create{
		dailyTaskRepository: dailyTaskRepository,
		logger:              logger,
		postgres:            postgres,
	}
}

func (s *Create) Execute(ctx context.Context, dto dailytask.CreateDTO) (dailytask.DailyTask, error) {
	s.logger.Debug("[create a new daily task] execute service")

	var (
		err    error
		result dailytask.DailyTask
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return dailytask.DailyTask{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// create new daily task.
	result, err = s.dailyTaskRepository.Create.Execute(ctx, tx, dto)
	if err != nil {
		return dailytask.DailyTask{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return dailytask.DailyTask{}, err
	}

	return result, nil
}
