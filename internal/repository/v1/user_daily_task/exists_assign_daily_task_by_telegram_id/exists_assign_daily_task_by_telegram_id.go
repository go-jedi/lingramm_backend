package existsassigndailytaskbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsAssignDailyTaskByTelegramID --output=mocks --case=underscore
type IExistsAssignDailyTaskByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error)
}

type ExistsAssignDailyTaskByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsAssignDailyTaskByTelegramID {
	r := &ExistsAssignDailyTaskByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsAssignDailyTaskByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsAssignDailyTaskByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
	r.logger.Debug("[check assign daily task exists by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM user_daily_tasks
			WHERE telegram_id = $1
			AND (
			    (occurred_at AT TIME ZONE 'Europe/Moscow')::DATE
			) = (
			    (NOW() AT TIME ZONE 'Europe/Moscow')::DATE
			)
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check assign daily task exists by telegram id", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check assign daily task exists by telegram id", "err", err)
		return false, fmt.Errorf("could not check assign daily task exists by telegram id: %w", err)
	}

	return ie, nil
}
