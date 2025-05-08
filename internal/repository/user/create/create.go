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
	jsoniter "github.com/json-iterator/go"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, dto user.CreateDTO) (user.User, error)
}

type Create struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Create {
	c := &Create{
		queryTimeout: queryTimeout,
		logger:       logger,
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

	rawData, err := jsoniter.Marshal(dto)
	if err != nil {
		c.logger.Error("failed to marshal user data", "err", err)
		return user.User{}, err
	}

	q := `SELECT * FROM public.user_create($1);`

	var nu user.User

	if err := tx.QueryRow(
		ctxTimeout, q,
		rawData,
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
