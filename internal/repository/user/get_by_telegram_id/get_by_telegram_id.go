package getbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetByTelegramID --output=mocks --case=underscore
type IGetByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (user.User, error)
}

type GetByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetByTelegramID {
	r := &GetByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (user.User, error) {
	r.logger.Debug("[get user by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM users
		WHERE telegram_id = $1;
	`

	var u user.User

	if err := tx.QueryRow(
		ctxTimeout, q, telegramID,
	).Scan(
		&u.ID, &u.UUID, &u.TelegramID,
		&u.Username, &u.FirstName, &u.LastName,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get user by telegram id", "err", err)
			return user.User{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get user by telegram id", "err", err)
		return user.User{}, fmt.Errorf("could not get user by telegram id: %w", err)
	}

	return u, nil
}
