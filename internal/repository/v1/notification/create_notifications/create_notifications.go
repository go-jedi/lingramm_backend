package createnotifications

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

//go:generate mockery --name=ICreateNotifications --output=mocks --case=underscore
type ICreateNotifications interface {
	Execute(ctx context.Context, tx pgx.Tx, dto []notification.CreateDTO) ([]notification.Notification, error)
}

type CreateNotifications struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateNotifications {
	r := &CreateNotifications{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateNotifications) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateNotifications) Execute(ctx context.Context, tx pgx.Tx, dto []notification.CreateDTO) ([]notification.Notification, error) {
	r.logger.Debug("[create a new notifications] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	rawData, err := jsoniter.Marshal(dto)
	if err != nil {
		return nil, err
	}

	q := `SELECT * FROM public.notifications_create($1);`

	var nn []notification.Notification

	if err := tx.QueryRow(
		ctxTimeout, q,
		rawData,
	).Scan(&nn); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new notifications", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new notifications", "err", err)
		return nil, fmt.Errorf("could not create a new notifications: %w", err)
	}

	return nn, nil
}
