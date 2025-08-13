package createsubscription

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

//go:generate mockery --name=ICreateSubscription --output=mocks --case=underscore
type ICreateSubscription interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (subscription.Subscription, error)
}

type CreateSubscription struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateSubscription {
	r := &CreateSubscription{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateSubscription) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateSubscription) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (subscription.Subscription, error) {
	r.logger.Debug("[create a new subscription] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.subscription_create($1);`

	var ns subscription.Subscription

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(
		&ns.ID, &ns.TelegramID, &ns.SubscribedAt,
		&ns.ExpiresAt, &ns.IsActive, &ns.CreatedAt, &ns.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new subscription", "err", err)
			return subscription.Subscription{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new subscription", "err", err)
		return subscription.Subscription{}, fmt.Errorf("could not create a new subscription: %w", err)
	}

	return ns, nil
}
