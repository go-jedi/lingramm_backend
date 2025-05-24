package addadminuser

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/admin"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAddAdminUser --output=mocks --case=underscore
type IAddAdminUser interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (admin.Admin, error)
}

type AddAdminUser struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *AddAdminUser {
	r := &AddAdminUser{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *AddAdminUser) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *AddAdminUser) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (admin.Admin, error) {
	r.logger.Debug("[add a new admin user] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO admins(
			telegram_id
		) VALUES ($1)
		RETURNING *;
	`

	var na admin.Admin

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(
		&na.ID, &na.TelegramID,
		&na.CreatedAt, &na.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while add admin user", "err", err)
			return admin.Admin{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to add admin user", "err", err)
		return admin.Admin{}, fmt.Errorf("could not add admin user: %w", err)
	}

	return na, nil
}
