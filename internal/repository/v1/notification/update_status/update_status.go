package updatestatus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IUpdateStatus --output=mocks --case=underscore
type IUpdateStatus interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64, status string) error
}

type UpdateStatus struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *UpdateStatus {
	r := &UpdateStatus{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	return r
}

func (r *UpdateStatus) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *UpdateStatus) Execute(ctx context.Context, tx pgx.Tx, id int64, status string) error {
	r.logger.Debug("[update status notification] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		UPDATE notifications SET
			status = $1
		WHERE id = $2;
	`

	commandTag, err := tx.Exec(
		ctxTimeout, q,
		status, id,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while update status notification", "err", err)
			return fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to update status notification", "err", err)
		return fmt.Errorf("could not update status notification: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrNoRowsWereAffected
	}

	return nil
}
