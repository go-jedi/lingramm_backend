package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, dto notification.CreateDTO) (notification.Notification, error)
}

type Create struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Create {
	c := &Create{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	c.init()

	return c
}

func (r *Create) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *Create) Execute(ctx context.Context, tx pgx.Tx, dto notification.CreateDTO) (notification.Notification, error) {
	r.logger.Debug("[create a new notification] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	rawMessage, err := jsoniter.Marshal(dto.Message)
	if err != nil {
		return notification.Notification{}, err
	}

	q := `
		INSERT INTO notifications(
		    type,
		    telegram_id,
		    message
		) VALUES($1, $2, $3)
		RETURNING *;
	`

	var nn notification.Notification

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.Type, dto.TelegramID,
		rawMessage,
	).Scan(
		&nn.ID, &nn.Type, &nn.TelegramID,
		&nn.Status, &nn.Message,
		&nn.CreatedAt, &nn.SentAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new notification", "err", err)
			return notification.Notification{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new notification", "err", err)
		return notification.Notification{}, fmt.Errorf("could not create a new notification: %w", err)
	}

	return nn, nil
}
