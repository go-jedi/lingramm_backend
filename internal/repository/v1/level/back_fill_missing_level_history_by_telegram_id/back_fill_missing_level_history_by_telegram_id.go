package backfillmissinglevelhistorybytelegramid

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

//go:generate mockery --name=IBackFillMissingLevelHistoryByTelegramID --output=mocks --case=underscore
type IBackFillMissingLevelHistoryByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) error
}

type BackFillMissingLevelHistoryByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *BackFillMissingLevelHistoryByTelegramID {
	r := &BackFillMissingLevelHistoryByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *BackFillMissingLevelHistoryByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *BackFillMissingLevelHistoryByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) error {
	r.logger.Debug("[back fill missing level history by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.back_fill_missing_level_history($1);`

	commandTag, err := tx.Exec(ctxTimeout, q, telegramID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while back fill missing level history by telegram id", "err", err)
			return fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to back fill missing level history by telegram id", "err", err)
		return fmt.Errorf("could not back fill missing level history by telegram id: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrNoRowsWereAffected
	}

	return nil
}
