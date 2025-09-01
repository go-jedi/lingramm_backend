package adduserbalance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/internal_currency/user_balance"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/utils/nullify"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

var (
	ErrAmountMustNotBeZero     = errors.New("amount must not be zero")
	ErrAmountMustNotBeNegative = errors.New("amount must not be negative")
	ErrNegativeBalance         = errors.New("resulting user balance would be negative")
)

//go:generate mockery --name=IAddUserBalance --output=mocks --case=underscore
type IAddUserBalance interface {
	Execute(ctx context.Context, tx pgx.Tx, dto userbalance.AddUserBalanceDTO) (userbalance.UserBalance, error)
}

type AddUserBalance struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *AddUserBalance {
	r := &AddUserBalance{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *AddUserBalance) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

// Execute performs the entire process of adding amount to a user's balance:
// 1. Validates the input,
// 2. Retrieves current balance with row lock,
// 3. Ensures resulting balance is not negative,
// 4. Creates a transaction record,
// 5. Updates the user balance,
// 6. Returns the updated user balance.
// All actions are performed inside the provided database transaction.
func (r *AddUserBalance) Execute(ctx context.Context, tx pgx.Tx, dto userbalance.AddUserBalanceDTO) (userbalance.UserBalance, error) {
	r.logger.Debug("[add user balance] execute repository")

	// check amount is zero.
	if dto.Amount.IsZero() {
		return userbalance.UserBalance{}, ErrAmountMustNotBeZero
	}
	// check amount is negative.
	if dto.Amount.IsNegative() {
		return userbalance.UserBalance{}, ErrAmountMustNotBeNegative
	}

	// get current balance.
	cb, err := r.getCurrentBalance(ctx, tx, dto.TelegramID)
	if err != nil {
		return userbalance.UserBalance{}, err
	}

	// current balance plus amount.
	nb := cb.Add(dto.Amount)
	if nb.IsNegative() { // check new balance is negative.
		return userbalance.UserBalance{}, fmt.Errorf("%w: %s + %s = %s", ErrNegativeBalance, cb, dto.Amount, nb)
	}

	// create balance transaction.
	if err := r.createBalanceTransaction(ctx, tx, dto, nb); err != nil {
		return userbalance.UserBalance{}, err
	}

	// update balance.
	ub, err := r.updateBalance(ctx, tx, dto.TelegramID, nb)
	if err != nil {
		return userbalance.UserBalance{}, err
	}

	return ub, nil
}

// getCurrentBalance retrieves the user's current balance from the database using SELECT ... FOR UPDATE
// to ensure safe concurrent access.
func (r *AddUserBalance) getCurrentBalance(ctx context.Context, tx pgx.Tx, telegramID string) (decimal.Decimal, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT balance
		FROM user_balances
		WHERE telegram_id = $1
		FOR UPDATE;
	`

	var currentBalance decimal.Decimal

	if err := tx.QueryRow(
		ctxTimeout, q, telegramID,
	).Scan(&currentBalance); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get current user balance", "err", err)
			return decimal.Decimal{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get current user balance", "err", err)
		return decimal.Decimal{}, fmt.Errorf("could not get current user balance: %w", err)
	}

	return currentBalance, nil
}

// createBalanceTransaction inserts a record into the balance_transactions table
// to log the operation for auditing and tracking purposes.
func (r *AddUserBalance) createBalanceTransaction(ctx context.Context, tx pgx.Tx, dto userbalance.AddUserBalanceDTO, newBalance decimal.Decimal) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO balance_transactions(
			event_type_id,
		    telegram_id,
		    amount,
		    description,
		    balance_after
		) VALUES ($1, $2, $3, $4, $5);
	`

	commandTag, err := tx.Exec(
		ctxTimeout, q,
		dto.EventTypeID, dto.TelegramID,
		dto.Amount, nullify.Empty(dto.Description), newBalance,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create balance transaction", "err", err)
			return fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create balance transaction", "err", err)
		return fmt.Errorf("could not create balance transaction: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrNoRowsWereAffected
	}

	return nil
}

// updateBalance updates the user's balance in the user_balances table
// and returns the full updated user balance record.
func (r *AddUserBalance) updateBalance(ctx context.Context, tx pgx.Tx, telegramID string, newBalance decimal.Decimal) (userbalance.UserBalance, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		UPDATE user_balances SET
			balance = $1
		WHERE telegram_id = $2
		RETURNING *;
	`

	var ub userbalance.UserBalance

	if err := tx.QueryRow(
		ctxTimeout, q,
		newBalance, telegramID,
	).Scan(
		&ub.ID, &ub.TelegramID, &ub.Balance,
		&ub.CreatedAt, &ub.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while update user balance", "err", err)
			return userbalance.UserBalance{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to update user balance", "err", err)
		return userbalance.UserBalance{}, fmt.Errorf("could not update user balance: %w", err)
	}

	return ub, nil
}
