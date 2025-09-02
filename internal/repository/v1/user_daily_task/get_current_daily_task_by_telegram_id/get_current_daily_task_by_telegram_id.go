package getcurrentdailytaskbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	userdailytask "github.com/go-jedi/lingramm_backend/internal/domain/user_daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetCurrentDailyTaskByTelegramID --output=mocks --case=underscore
type IGetCurrentDailyTaskByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userdailytask.GetCurrentDailyTaskByTelegramIDResponse, error)
}

type GetCurrentDailyTaskByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetCurrentDailyTaskByTelegramID {
	r := &GetCurrentDailyTaskByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetCurrentDailyTaskByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetCurrentDailyTaskByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userdailytask.GetCurrentDailyTaskByTelegramIDResponse, error) {
	r.logger.Debug("[get current daily task by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.daily_task_current_get($1);`

	var result userdailytask.GetCurrentDailyTaskByTelegramIDResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(&result); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get current daily task by telegram id", "err", err)
			return userdailytask.GetCurrentDailyTaskByTelegramIDResponse{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get current daily task by telegram id", "err", err)
		return userdailytask.GetCurrentDailyTaskByTelegramIDResponse{}, fmt.Errorf("could not get current daily task by telegram id: %w", err)
	}

	return result, nil
}
