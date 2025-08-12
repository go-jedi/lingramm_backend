package allbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAllByTelegramID --output=mocks --case=underscore
type IAllByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) ([]notification.Notification, error)
}

type AllByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *AllByTelegramID {
	r := &AllByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *AllByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *AllByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) ([]notification.Notification, error) {
	r.logger.Debug("[get all notifications by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM notifications
		WHERE telegram_id = $1;
	`

	rows, err := tx.Query(
		ctxTimeout, q,
		telegramID,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get all notifications by telegram id", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get all notifications by telegram id", "err", err)
		return nil, fmt.Errorf("could not get all notifications by telegram id: %w", err)
	}
	defer rows.Close()

	var notifications []notification.Notification

	for rows.Next() {
		var n notification.Notification

		if err := rows.Scan(
			&n.ID, &n.Type, &n.TelegramID,
			&n.Status, &n.Message,
			&n.CreatedAt, &n.SentAt,
		); err != nil {
			r.logger.Error("failed to scan row to get all notifications by telegram id", "err", err)
			return nil, fmt.Errorf("failed to scan row to get all notifications by telegram id: %w", err)
		}

		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed to get all notifications by telegram id", "err", rows.Err())
		return nil, fmt.Errorf("failed to get all notifications by telegram id: %w", err)
	}

	return notifications, nil
}
