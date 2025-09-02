package assigndailytaskbytelegramid

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

//go:generate mockery --name=IAssignDailyTaskByTelegramID --output=mocks --case=underscore
type IAssignDailyTaskByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userdailytask.AssignDailyTaskByTelegramIDResponse, error)
}

type AssignDailyTaskByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *AssignDailyTaskByTelegramID {
	r := &AssignDailyTaskByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *AssignDailyTaskByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *AssignDailyTaskByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userdailytask.AssignDailyTaskByTelegramIDResponse, error) {
	r.logger.Debug("[assign daily task by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.assign_daily_task($1);`

	var result userdailytask.AssignDailyTaskByTelegramIDResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(&result); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while assign daily task by telegram id", "err", err)
			return userdailytask.AssignDailyTaskByTelegramIDResponse{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to assign daily task by telegram id", "err", err)
		return userdailytask.AssignDailyTaskByTelegramIDResponse{}, fmt.Errorf("could not assign daily task by telegram id: %w", err)
	}

	return result, nil
}
