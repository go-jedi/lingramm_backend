package all

import (
	"context"
	"errors"
	"fmt"
	"time"

	eventtype "github.com/go-jedi/lingramm_backend/internal/domain/event_type"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAll --output=mocks --case=underscore
type IAll interface {
	Execute(ctx context.Context, tx pgx.Tx) ([]eventtype.EventType, error)
}

type All struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *All {
	r := &All{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *All) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *All) Execute(ctx context.Context, tx pgx.Tx) ([]eventtype.EventType, error) {
	r.logger.Debug("[get all event types] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT * 
		FROM event_types
		ORDER BY id;
	`

	rows, err := tx.Query(ctxTimeout, q)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get all event types", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get all event types", "err", err)
		return nil, fmt.Errorf("could not get all event types: %w", err)
	}
	defer rows.Close()

	var eventTypes []eventtype.EventType

	for rows.Next() {
		var et eventtype.EventType

		if err := rows.Scan(
			&et.ID, &et.Name,
			&et.Description, &et.XP,
			&et.Amount, &et.NotificationMessage,
			&et.IsSendNotification, &et.IsActive,
			&et.CreatedAt, &et.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row to get all event types", "err", err)
			return nil, fmt.Errorf("failed to scan row to get all event types: %w", err)
		}

		eventTypes = append(eventTypes, et)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed to get all event types", "err", rows.Err())
		return nil, fmt.Errorf("failed to get all event types: %w", err)
	}

	return eventTypes, nil
}
