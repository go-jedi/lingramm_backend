package createachievementcondition

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

//go:generate mockery --name=ICreateAchievementCondition --output=mocks --case=underscore
type ICreateAchievementCondition interface {
	Execute(ctx context.Context, tx pgx.Tx, dto achievement.CreateAchievementConditionDTO) (achievement.Condition, error)
}

type CreateAchievementCondition struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateAchievementCondition {
	r := &CreateAchievementCondition{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateAchievementCondition) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateAchievementCondition) Execute(ctx context.Context, tx pgx.Tx, dto achievement.CreateAchievementConditionDTO) (achievement.Condition, error) {
	r.logger.Debug("[create achievement condition] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO achievement_conditions(
		    achievement_id,
		    condition_type,
		    operator,
		    value
		) VALUES($1, $2, $3, $4)
		RETURNING *;
	`

	var nac achievement.Condition

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.AchievementID, dto.ConditionType,
		dto.Operator, dto.Value,
	).Scan(
		&nac.ID, &nac.AchievementID,
		&nac.ConditionType, &nac.Operator,
		&nac.Value, &nac.CreatedAt, &nac.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create achievement condition", "err", err)
			return achievement.Condition{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create achievement condition", "err", err)
		return achievement.Condition{}, fmt.Errorf("could not create achievement condition: %w", err)
	}

	return nac, nil
}
