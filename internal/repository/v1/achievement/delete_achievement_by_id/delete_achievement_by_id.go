package deleteachievementbyid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IDeleteAchievementByID --output=mocks --case=underscore
type IDeleteAchievementByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (achievement.Achievement, error)
}

type DeleteAchievementByID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *DeleteAchievementByID {
	r := &DeleteAchievementByID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *DeleteAchievementByID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *DeleteAchievementByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (achievement.Achievement, error) {
	r.logger.Debug("[delete achievement by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		DELETE FROM achievements
		WHERE id = $1
		RETURNING *;
	`

	var da achievement.Achievement

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(
		&da.ID, &da.AchievementAssetsID,
		&da.Code, &da.Name, &da.Description,
		&da.CreatedAt, &da.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while delete achievement by id", "err", err)
			return achievement.Achievement{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to delete achievement by id", "err", err)
		return achievement.Achievement{}, fmt.Errorf("could not delete achievement by id: %w", err)
	}

	return da, nil
}
