package getbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/subscription"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetByTelegramID --output=mocks --case=underscore
type IGetByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (subscription.Subscription, error)
}

type GetByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetByTelegramID {
	r := &GetByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (subscription.Subscription, error) {
	r.logger.Debug("[get subscription by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM subscriptions
		WHERE telegram_id = $1;
	`

	var s subscription.Subscription

	if err := tx.QueryRow(
		ctxTimeout, q, telegramID,
	).Scan(
		&s.ID, &s.TelegramID, &s.SubscribedAt,
		&s.ExpiresAt, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get subscription by telegram id", "err", err)
			return subscription.Subscription{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get subscription by telegram id", "err", err)
		return subscription.Subscription{}, fmt.Errorf("could not get subscription by telegram id: %w", err)
	}

	return s, nil
}
