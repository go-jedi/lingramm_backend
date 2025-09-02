package syncuserdailytaskprogress

import (
	"context"
	"errors"
	"fmt"
	"time"

	userdailytask "github.com/go-jedi/lingramm_backend/internal/domain/user_daily_task"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
)

//go:generate mockery --name=ISyncUserDailyTaskProgress --output=mocks --case=underscore
type ISyncUserDailyTaskProgress interface {
	Execute(ctx context.Context, tx pgx.Tx, dto userdailytask.SyncUserDailyTaskProgressDTO) error
}

type SyncUserDailyTaskProgress struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *SyncUserDailyTaskProgress {
	r := &SyncUserDailyTaskProgress{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *SyncUserDailyTaskProgress) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *SyncUserDailyTaskProgress) Execute(ctx context.Context, tx pgx.Tx, dto userdailytask.SyncUserDailyTaskProgressDTO) error {
	r.logger.Debug("[sync user daily task progress] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	rawActions, err := jsoniter.Marshal(dto.Actions)
	if err != nil {
		return err
	}

	q := `SELECT * FROM public.sync_user_daily_task_progress($1, $2);`

	commandTag, err := tx.Exec(
		ctxTimeout, q,
		dto.TelegramID, rawActions,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while sync user daily task progress", "err", err)
			return fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to sync user daily task progress", "err", err)
		return fmt.Errorf("could not sync user daily task progress: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrNoRowsWereAffected
	}

	return nil
}
