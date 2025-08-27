package existsachievementbyachievementtype

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsAchievementByAchievementType --output=mocks --case=underscore
type IExistsAchievementByAchievementType interface {
	Execute(ctx context.Context, tx pgx.Tx, achievementType string) (bool, error)
}

type ExistsAchievementByAchievementType struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsAchievementByAchievementType {
	r := &ExistsAchievementByAchievementType{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsAchievementByAchievementType) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsAchievementByAchievementType) Execute(ctx context.Context, tx pgx.Tx, achievementType string) (bool, error) {
	r.logger.Debug("[check achievement exists by achievement type] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM achievements a
			INNER JOIN achievement_types at ON a.achievement_type_id = at.id
			WHERE at.name = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		achievementType,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check achievement exists by achievement type", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check achievement exists by achievement type", "err", err)
		return false, fmt.Errorf("could not check achievement exists by achievement type: %w", err)
	}

	return ie, nil
}
