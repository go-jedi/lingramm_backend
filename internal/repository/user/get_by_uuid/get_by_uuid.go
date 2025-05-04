package getbyuuid

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

//go:generate mockery --name=IGetByUUID --output=mocks --case=underscore
type IGetByUUID interface {
	Execute(ctx context.Context, tx pgx.Tx, uuid string) (user.User, error)
}

type GetByUUID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetByUUID {
	gbui := &GetByUUID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	gbui.init()

	return gbui
}

func (gbui *GetByUUID) init() {
	if gbui.queryTimeout == 0 {
		gbui.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (gbui *GetByUUID) Execute(ctx context.Context, tx pgx.Tx, uuid string) (user.User, error) {
	gbui.logger.Debug("[get user by uuid] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(gbui.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM users
		WHERE uuid = $1;
	`

	var u user.User

	if err := tx.QueryRow(
		ctxTimeout, q, uuid,
	).Scan(
		&u.ID, &u.UUID, &u.TelegramID,
		&u.Username, &u.FirstName, &u.LastName,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			gbui.logger.Error("request timed out while get user by uuid", "err", err)
			return user.User{}, fmt.Errorf("the request timed out: %w", err)
		}
		gbui.logger.Error("failed to get user by uuid", "err", err)
		return user.User{}, fmt.Errorf("could not get user by uuid: %w", err)
	}

	return u, nil
}
