package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
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
	r := &Create{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *Create) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *Create) Execute(ctx context.Context, tx pgx.Tx, dto user.CreateDTO) (user.User, error) {
	r.logger.Debug("[create a new user] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	rawData, err := jsoniter.Marshal(dto)
	if err != nil {
		r.logger.Error("failed to marshal user data", "err", err)
		return user.User{}, err
	}

	q := `SELECT * FROM public.user_create($1);`

	var nu user.User

	if err := tx.QueryRow(
		ctxTimeout, q,
		rawData,
	).Scan(
		&nu.ID, &nu.TelegramID,
		&nu.Username, &nu.FirstName, &nu.LastName,
		&nu.CreatedAt, &nu.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a new user", "err", err)
			return user.User{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a new user", "err", err)
		return user.User{}, fmt.Errorf("could not create a new user: %w", err)
	}

	return nu, nil
}
