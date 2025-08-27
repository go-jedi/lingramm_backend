package deleteachievementconditionbyid

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

//go:generate mockery --name=IDeleteAchievementConditionByID --output=mocks --case=underscore
type IDeleteAchievementConditionByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (achievement.Condition, error)
}

type DeleteAchievementConditionByID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *DeleteAchievementConditionByID {
	r := &DeleteAchievementConditionByID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *DeleteAchievementConditionByID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *DeleteAchievementConditionByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (achievement.Condition, error) {
	r.logger.Debug("[delete achievement condition by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		DELETE FROM achievement_conditions
		WHERE id = $1
		RETURNING *;
	`

	var dac achievement.Condition

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(
		&dac.ID, &dac.AchievementID,
		&dac.AchievementTypeID, &dac.Operator,
		&dac.Value, &dac.CreatedAt, &dac.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while delete achievement condition by id", "err", err)
			return achievement.Condition{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to delete achievement condition by id", "err", err)
		return achievement.Condition{}, fmt.Errorf("could not delete achievement condition by id: %w", err)
	}

	return dac, nil
}
