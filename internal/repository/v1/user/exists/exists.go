package exists

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExists --output=mocks --case=underscore
type IExists interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string, username string) (bool, error)
}

type Exists struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Exists {
	r := &Exists{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *Exists) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *Exists) Execute(ctx context.Context, tx pgx.Tx, telegramID string, username string) (bool, error) {
	r.logger.Debug("[check a user exists] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE telegram_id = $1
			OR username = $2
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID, username,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check exists user", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check exists user", "err", err)
		return false, fmt.Errorf("could not check exists user: %w", err)
	}

	return ie, nil
}
