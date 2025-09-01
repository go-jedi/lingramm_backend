package syncuserstatsfromxpeventsbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/event"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ISyncUserStatsFromXPEventsByTelegramID --output=mocks --case=underscore
type ISyncUserStatsFromXPEventsByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string, actions event.Actions) error
}

type SyncUserStatsFromXPEventsByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *SyncUserStatsFromXPEventsByTelegramID {
	r := &SyncUserStatsFromXPEventsByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *SyncUserStatsFromXPEventsByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *SyncUserStatsFromXPEventsByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string, actions event.Actions) error {
	r.logger.Debug("[sync user stats from xp events by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.sync_user_stats_from_xp_events($1, $2);`

	commandTag, err := tx.Exec(
		ctxTimeout, q,
		telegramID, actions,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while sync user stats from xp events by telegram id", "err", err)
			return fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to sync user stats from xp events by telegram id", "err", err)
		return fmt.Errorf("could not sync user stats from xp events by telegram id: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrNoRowsWereAffected
	}

	return nil
}
