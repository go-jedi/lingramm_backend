package getdailytaskweeksummarybytelegramid

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

//go:generate mockery --name=IGetDailyTaskWeekSummaryByTelegramID --output=mocks --case=underscore
type IGetDailyTaskWeekSummaryByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) ([]userdailytask.GetDailyTaskWeekSummaryByTelegramIDResponse, error)
}

type GetDailyTaskWeekSummaryByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetDailyTaskWeekSummaryByTelegramID {
	r := &GetDailyTaskWeekSummaryByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetDailyTaskWeekSummaryByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetDailyTaskWeekSummaryByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) ([]userdailytask.GetDailyTaskWeekSummaryByTelegramIDResponse, error) {
	r.logger.Debug("[get daily task week summary by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.daily_task_week_summary_get($1);`

	var result []userdailytask.GetDailyTaskWeekSummaryByTelegramIDResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(&result); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get daily task week summary by telegram id", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get daily task week summary by telegram id", "err", err)
		return nil, fmt.Errorf("could not get daily task week summary by telegram id: %w", err)
	}

	return result, nil
}
