package create

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

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, dto user.CreateDTO) (user.User, error)
}

type Create struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(logger logger.ILogger) *Create {
	c := &Create{
		logger: logger,
	}

	c.init()

	return c
}

func (c *Create) init() {
	if c.queryTimeout == 0 {
		c.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (c *Create) Execute(ctx context.Context, tx pgx.Tx, dto user.CreateDTO) (user.User, error) {
	c.logger.Debug("[create a new user] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(c.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO users(
			uuid,
		    telegram_id,
		    username,
		    first_name,
		    last_name
		) VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (telegram_id) DO NOTHING
		RETURNING *;
	`

	var nu user.User

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.UUID, dto.TelegramID, dto.Username,
		dto.FirstName, dto.LastName,
	).Scan(
		&nu.ID, &nu.UUID, &nu.TelegramID,
		&nu.Username, &nu.FirstName, &nu.LastName,
		&nu.CreatedAt, &nu.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.logger.Error("request timed out while creating the user", "err", err)
			return user.User{}, fmt.Errorf("the request timed out: %w", err)
		}
		c.logger.Error("failed to create user", "err", err)
		return user.User{}, fmt.Errorf("could not create user: %w", err)
	}

	return nu, nil
}
