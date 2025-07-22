package getbalancetransactioneventbycode

import (
	"context"
	"errors"
	"fmt"
	"time"

	balancetransactionevent "github.com/go-jedi/lingramm_backend/internal/domain/internal_currency/balance_transaction_event"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetBalanceTransactionEventByCode --output=mocks --case=underscore
type IGetBalanceTransactionEventByCode interface {
	Execute(ctx context.Context, tx pgx.Tx, code string) (balancetransactionevent.BalanceTransactionEvent, error)
}

type GetBalanceTransactionEventByCode struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetBalanceTransactionEventByCode {
	r := &GetBalanceTransactionEventByCode{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetBalanceTransactionEventByCode) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetBalanceTransactionEventByCode) Execute(ctx context.Context, tx pgx.Tx, code string) (balancetransactionevent.BalanceTransactionEvent, error) {
	r.logger.Debug("[get user balance] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM balance_transaction_events
		WHERE code = $1;
	`

	var bte balancetransactionevent.BalanceTransactionEvent

	if err := tx.QueryRow(
		ctxTimeout, q, code,
	).Scan(
		&bte.ID, &bte.Code, &bte.Description, &bte.Amount,
		&bte.IsActive, &bte.CreatedAt, &bte.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get balance transaction event", "err", err)
			return balancetransactionevent.BalanceTransactionEvent{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get balance transaction event", "err", err)
		return balancetransactionevent.BalanceTransactionEvent{}, fmt.Errorf("could not get balance transaction event: %w", err)
	}

	return bte, nil
}
