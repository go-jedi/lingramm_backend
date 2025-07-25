package existsachievementconditionbyconditiontype

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsAchievementConditionByConditionType --output=mocks --case=underscore
type IExistsAchievementConditionByConditionType interface {
	Execute(ctx context.Context, tx pgx.Tx, conditionType string) (bool, error)
}

type ExistsAchievementConditionByConditionType struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsAchievementConditionByConditionType {
	r := &ExistsAchievementConditionByConditionType{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsAchievementConditionByConditionType) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsAchievementConditionByConditionType) Execute(ctx context.Context, tx pgx.Tx, conditionType string) (bool, error) {
	r.logger.Debug("[check achievement condition exists by condition type] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM achievement_conditions
			WHERE condition_type = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		conditionType,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check achievement condition exists by condition type", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check achievement condition exists by condition type", "err", err)
		return false, fmt.Errorf("could not check achievement condition exists by condition type: %w", err)
	}

	return ie, nil
}
