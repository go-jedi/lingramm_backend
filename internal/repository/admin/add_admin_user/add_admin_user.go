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
	aad := &AddAdminUser{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	aad.init()

	return aad
}

func (aad *AddAdminUser) init() {
	if aad.queryTimeout == 0 {
		aad.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (aad *AddAdminUser) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (admin.Admin, error) {
	aad.logger.Debug("[add a new admin user] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(aad.queryTimeout)*time.Second)
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
			aad.logger.Error("request timed out while add admin user", "err", err)
			return admin.Admin{}, fmt.Errorf("the request timed out: %w", err)
		}
		aad.logger.Error("failed to add admin user", "err", err)
		return admin.Admin{}, fmt.Errorf("could not add admin user: %w", err)
	}

	return na, nil
}
