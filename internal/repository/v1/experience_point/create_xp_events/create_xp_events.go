package createxpevents

import (
	"context"
	"errors"
	"fmt"
	"time"

	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreateXPEvents --output=mocks --case=underscore
type ICreateXPEvents interface {
	Execute(ctx context.Context, tx pgx.Tx, dto experiencepoint.CreateXPEventDTO) error
}

type CreateXPEvents struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateXPEvents {
	r := &CreateXPEvents{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateXPEvents) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateXPEvents) Execute(ctx context.Context, tx pgx.Tx, dto experiencepoint.CreateXPEventDTO) error {
	r.logger.Debug("[create a new xp events] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.xp_event_create($1, $2, $3);`

	commandTag, err := tx.Exec(
		ctxTimeout, q,
		dto.TelegramID, dto.EventType, dto.DeltaXP,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new xp events", "err", err)
			return fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new xp events", "err", err)
		return fmt.Errorf("could not create a new xp events: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrNoRowsWereAffected
	}

	return nil
}
