package existsbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsByTelegramID --output=mocks --case=underscore
type IExistsByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error)
}

type ExistsByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *ExistsByTelegramID {
	r := &ExistsByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *ExistsByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *ExistsByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
	r.logger.Debug("[check user stats exists by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM user_stats
			WHERE telegram_id = $1
		);
	`

	ie := false

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while check user stats exists by telegram id", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to check user stats exists by telegram id", "err", err)
		return false, fmt.Errorf("could not check user stats exists by telegram id: %w", err)
	}

	return ie, nil
}
