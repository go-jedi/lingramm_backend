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
	ebt := &ExistsByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	ebt.init()

	return ebt
}

func (ebt *ExistsByTelegramID) init() {
	if ebt.queryTimeout == 0 {
		ebt.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (ebt *ExistsByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
	ebt.logger.Debug("[check admin exists by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(ebt.queryTimeout)*time.Second)
	defer cancel()

	ie := false

	q := `
		SELECT EXISTS(
			SELECT 1
			FROM admins
			WHERE telegram_id = $1
		);
	`

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(&ie); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			ebt.logger.Error("request timed out while check exists admin by telegram id", "err", err)
			return false, fmt.Errorf("the request timed out: %w", err)
		}
		ebt.logger.Error("failed to check exists admin by telegram id", "err", err)
		return false, fmt.Errorf("could not check exists admin by telegram id: %w", err)
	}

	return ie, nil
}
