package existsbyid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsByID --output=mocks --case=underscore
type IExistsByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (bool, error)
}

type ExistsByID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsByID {
	r := &ExistsByID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsByID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (bool, error) {
	r.logger.Debug("[check notification exists by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM notifications
			WHERE id = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check notification exists by id", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check notification exists by id", "err", err)
		return false, fmt.Errorf("could not check notification exists by id: %w", err)
	}

	return ie, nil
}
