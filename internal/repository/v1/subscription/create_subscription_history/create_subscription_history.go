package createsubscriptionhistory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/subscription"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreateSubscriptionHistory --output=mocks --case=underscore
type ICreateSubscriptionHistory interface {
	Execute(ctx context.Context, tx pgx.Tx, dto subscription.CreateSubscriptionHistoryDTO) error
}

type CreateSubscriptionHistory struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateSubscriptionHistory {
	r := &CreateSubscriptionHistory{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateSubscriptionHistory) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateSubscriptionHistory) Execute(ctx context.Context, tx pgx.Tx, dto subscription.CreateSubscriptionHistoryDTO) error {
	r.logger.Debug("[create a new subscription history history] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO subscription_history(
			telegram_id,
		    action_time,
		    expires_at
		) VALUES($1, $2, $3);
	`

	commandTag, err := tx.Exec(
		ctxTimeout, q,
		dto.TelegramID, dto.ActionTime,
		dto.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new subscription history history", "err", err)
			return fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new subscription history history", "err", err)
		return fmt.Errorf("could not create a new subscription history history: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrNoRowsWereAffected
	}

	return nil
}
