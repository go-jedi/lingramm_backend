package createuserlevelhistory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/level"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreateUserLevelHistory --output=mocks --case=underscore
type ICreateUserLevelHistory interface {
	Execute(ctx context.Context, tx pgx.Tx, dto level.CreateUserLevelHistoryDTO) (level.UserLevelHistory, error)
}

type CreateUserLevelHistory struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateUserLevelHistory {
	r := &CreateUserLevelHistory{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateUserLevelHistory) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateUserLevelHistory) Execute(ctx context.Context, tx pgx.Tx, dto level.CreateUserLevelHistoryDTO) (level.UserLevelHistory, error) {
	r.logger.Debug("[create a user level history] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO user_level_history(
			telegram_id,
		    level_number,
		    xp_event_id,
		    xp_at_reach,
		    reached_at
		) VALUES($1, $2, $3, $4, NOW())
		RETURNING *;
	`

	var nulh level.UserLevelHistory

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.TelegramID, dto.LevelNumber,
		dto.XPEventID, dto.XPAtReach,
	).Scan(
		&nulh.ID, &nulh.TelegramID, &nulh.LevelNumber,
		&nulh.XPEventID, &nulh.XPAtReach, &nulh.ReachedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a user level history", "err", err)
			return level.UserLevelHistory{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a user level history", "err", err)
		return level.UserLevelHistory{}, fmt.Errorf("could not create a user level history: %w", err)
	}

	return nulh, nil
}
