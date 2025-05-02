package getbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingvogramm_backend/internal/domain/user"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
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

func New(logger logger.ILogger) *GetByTelegramID {
	gbti := &GetByTelegramID{
		logger: logger,
	}

	gbti.init()

	return gbti
}

func (gbti *GetByTelegramID) init() {
	if gbti.queryTimeout == 0 {
		gbti.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (gbti *GetByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (user.User, error) {
	gbti.logger.Debug("[get user by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(gbti.queryTimeout)*time.Second)
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
			gbti.logger.Error("request timed out while get user by telegram id", "err", err)
			return user.User{}, fmt.Errorf("the request timed out: %w", err)
		}
		gbti.logger.Error("failed to get user by telegram id", "err", err)
		return user.User{}, fmt.Errorf("could not get user by telegram id: %w", err)
	}

	return u, nil
}
