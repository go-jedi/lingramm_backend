package ensurestreakdaysincrementtoday

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IEnsureStreakDaysIncrementToday --output=mocks --case=underscore
type IEnsureStreakDaysIncrementToday interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) error
}

type EnsureStreakDaysIncrementToday struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *EnsureStreakDaysIncrementToday {
	r := &EnsureStreakDaysIncrementToday{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *EnsureStreakDaysIncrementToday) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *EnsureStreakDaysIncrementToday) Execute(ctx context.Context, tx pgx.Tx, telegramID string) error {
	r.logger.Debug("[ensure streak days increment today] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.ensure_streak_days_increment_today($1);`

	commandTag, err := tx.Exec(
		ctxTimeout, q,
		telegramID,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while ensure streak days increment today", "err", err)
			return fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to ensure streak days increment today", "err", err)
		return fmt.Errorf("could not ensure streak days increment today: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrNoRowsWereAffected
	}

	return nil
}
