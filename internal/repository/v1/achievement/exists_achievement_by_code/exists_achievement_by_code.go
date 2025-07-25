package existsachievementbycode

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsAchievementByCode --output=mocks --case=underscore
type IExistsAchievementByCode interface {
	Execute(ctx context.Context, tx pgx.Tx, code string) (bool, error)
}

type ExistsAchievementByCode struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsAchievementByCode {
	r := &ExistsAchievementByCode{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsAchievementByCode) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsAchievementByCode) Execute(ctx context.Context, tx pgx.Tx, code string) (bool, error) {
	r.logger.Debug("[check achievement exists by code] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	ie := false

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM achievements
			WHERE code = $1
		);
	`

	if err := tx.QueryRow(
		ctxTimeout, q,
		code,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check achievement exists by code", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check achievement exists by code", "err", err)
		return false, fmt.Errorf("could not check achievement exists by code: %w", err)
	}

	return ie, nil
}
