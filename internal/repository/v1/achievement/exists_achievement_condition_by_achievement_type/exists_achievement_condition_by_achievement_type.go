package existsachievementconditionbyachievementtype

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsAchievementConditionByAchievementType --output=mocks --case=underscore
type IExistsAchievementConditionByAchievementType interface {
	Execute(ctx context.Context, tx pgx.Tx, achievementType string) (bool, error)
}

type ExistsAchievementConditionByAchievementType struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsAchievementConditionByAchievementType {
	r := &ExistsAchievementConditionByAchievementType{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsAchievementConditionByAchievementType) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsAchievementConditionByAchievementType) Execute(ctx context.Context, tx pgx.Tx, achievementType string) (bool, error) {
	r.logger.Debug("[check achievement condition exists by achievement type] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM achievement_conditions ac
			INNER JOIN achievement_types at ON ac.achievement_type_id = at.id
			WHERE at.name = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		achievementType,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check achievement condition exists by achievement type", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check achievement condition exists by achievement type", "err", err)
		return false, fmt.Errorf("could not check achievement condition exists by achievement type: %w", err)
	}

	return ie, nil
}
