package getbyname

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

//go:generate mockery --name=IGetByName --output=mocks --case=underscore
type IGetByName interface {
	Execute(ctx context.Context, tx pgx.Tx, name string) (eventtype.EventType, error)
}

type GetByName struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetByName {
	r := &GetByName{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetByName) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetByName) Execute(ctx context.Context, tx pgx.Tx, name string) (eventtype.EventType, error) {
	r.logger.Debug("[get event type by name] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT * 
		FROM event_types
		WHERE name = $1;
	`

	var et eventtype.EventType

	if err := tx.QueryRow(
		ctxTimeout, q,
		name,
	).Scan(
		&et.ID, &et.Name,
		&et.Description, &et.XP,
		&et.Amount, &et.NotificationMessage,
		&et.IsSendNotification, &et.IsActive,
		&et.CreatedAt, &et.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get event type by name", "err", err)
			return eventtype.EventType{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get event type by name", "err", err)
		return eventtype.EventType{}, fmt.Errorf("could not get event type by name: %w", err)
	}

	return et, nil
}
