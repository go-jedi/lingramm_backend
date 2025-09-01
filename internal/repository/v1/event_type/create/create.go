package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	eventtype "github.com/go-jedi/lingramm_backend/internal/domain/event_type"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/utils/nullify"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, dto eventtype.CreateDTO) (eventtype.EventType, error)
}

type Create struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Create {
	r := &Create{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *Create) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *Create) Execute(ctx context.Context, tx pgx.Tx, dto eventtype.CreateDTO) (eventtype.EventType, error) {
	r.logger.Debug("[create a new event type] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO event_types(
		    name,
		    description,
		    xp,
		    amount,
		    notification_message,
		    is_send_notification,
		    is_active
		) VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING *;
	`

	var newEventType eventtype.EventType

	if err := tx.QueryRow(
		ctxTimeout, q,
		r.getArgs(dto)...,
	).Scan(
		&newEventType.ID, &newEventType.Name,
		&newEventType.Description, &newEventType.XP,
		&newEventType.Amount, &newEventType.NotificationMessage,
		&newEventType.IsSendNotification, &newEventType.IsActive,
		&newEventType.CreatedAt, &newEventType.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new event type", "err", err)
			return eventtype.EventType{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new event type", "err", err)
		return eventtype.EventType{}, fmt.Errorf("could not create a new event type: %w", err)
	}

	return newEventType, nil
}

// getArgs get args.
func (r *Create) getArgs(dto eventtype.CreateDTO) []interface{} {
	return []interface{}{
		dto.Name,
		nullify.EmptyString(dto.Description),
		dto.XP,
		nullify.EmptyDecimalWithDefault(dto.Amount),
		nullify.EmptyString(dto.NotificationMessage),
		dto.IsSendNotification,
		dto.IsActive,
	}
}
