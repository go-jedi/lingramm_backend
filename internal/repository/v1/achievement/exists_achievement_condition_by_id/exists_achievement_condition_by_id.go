package existsachievementconditionbyid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsAchievementConditionByID --output=mocks --case=underscore
type IExistsAchievementConditionByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (bool, error)
}

type ExistsAchievementConditionByID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsAchievementConditionByID {
	r := &ExistsAchievementConditionByID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsAchievementConditionByID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsAchievementConditionByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (bool, error) {
	r.logger.Debug("[check achievement condition exists by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM achievement_conditions
			WHERE id = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check achievement condition exists by id", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check achievement condition exists by id", "err", err)
		return false, fmt.Errorf("could not check achievement condition exists by id: %w", err)
	}

	return ie, nil
}
