package hasstreakdaysincrementtoday

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IHasStreakDaysIncrementToday --output=mocks --case=underscore
type IHasStreakDaysIncrementToday interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error)
}

type HasStreakDaysIncrementToday struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *HasStreakDaysIncrementToday {
	r := &HasStreakDaysIncrementToday{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *HasStreakDaysIncrementToday) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *HasStreakDaysIncrementToday) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
	r.logger.Debug("[check streak days increment today by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM user_stats
			WHERE telegram_id = $1
			AND last_active_at = CURRENT_DATE
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check streak days increment today by telegram id", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check streak days increment today by telegram id", "err", err)
		return false, fmt.Errorf("could not check streak days increment today by telegram id: %w", err)
	}

	return ie, nil
}
