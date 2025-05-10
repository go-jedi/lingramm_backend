package getuserbalance

import (
	"context"
	"errors"
	"fmt"
	"time"

	userbalance "github.com/go-jedi/lingvogramm_backend/internal/domain/user_balance"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetUserBalance --output=mocks --case=underscore
type IGetUserBalance interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userbalance.UserBalance, error)
}

type GetUserBalance struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetUserBalance {
	gub := &GetUserBalance{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	gub.init()

	return gub
}

func (gub *GetUserBalance) init() {
	if gub.queryTimeout == 0 {
		gub.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (gub *GetUserBalance) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userbalance.UserBalance, error) {
	gub.logger.Debug("[get user balance] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(gub.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM user_balances
		WHERE telegram_id = $1;
	`

	var ub userbalance.UserBalance

	if err := tx.QueryRow(
		ctxTimeout, q, telegramID,
	).Scan(
		&ub.ID, &ub.TelegramID, &ub.Balance,
		&ub.CreatedAt, &ub.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			gub.logger.Error("request timed out while get user balance", "err", err)
			return userbalance.UserBalance{}, fmt.Errorf("the request timed out: %w", err)
		}
		gub.logger.Error("failed to get user balance", "err", err)
		return userbalance.UserBalance{}, fmt.Errorf("could not get user balance: %w", err)
	}

	return ub, nil
}
