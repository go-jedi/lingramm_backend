package existsbyname

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsByName --output=mocks --case=underscore
type IExistsByName interface {
	Execute(ctx context.Context, tx pgx.Tx, name string) (bool, error)
}

type ExistsByName struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsByName {
	r := &ExistsByName{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsByName) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsByName) Execute(ctx context.Context, tx pgx.Tx, name string) (bool, error) {
	r.logger.Debug("[check achievement type exists by name] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM achievement_types
			WHERE name = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		name,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check achievement type exists by name", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check achievement type exists by name", "err", err)
		return false, fmt.Errorf("could not check achievement type exists by name: %w", err)
	}

	return ie, nil
}
