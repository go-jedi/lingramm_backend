package deletebyid

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IDeleteByID --output=mocks --case=underscore
type IDeleteByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (clientassets.ClientAssets, error)
}

type DeleteByID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *DeleteByID {
	r := &DeleteByID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *DeleteByID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *DeleteByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (clientassets.ClientAssets, error) {
	r.logger.Debug("[delete client assets by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		DELETE FROM client_assets
		WHERE id = $1
		RETURNING *;
	`

	var dca clientassets.ClientAssets

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(
		&dca.ID, &dca.NameFile, &dca.NameFileWithoutExtension,
		&dca.ServerPathFile, &dca.ClientPathFile, &dca.Extension,
		&dca.Quality, &dca.OldNameFile, &dca.OldExtension,
		&dca.CreatedAt, &dca.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while delete client assets by id", "err", err)
			return clientassets.ClientAssets{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to delete client assets by id", "err", err)
		return clientassets.ClientAssets{}, fmt.Errorf("could not delete client assets by id: %w", err)
	}

	return dca, nil
}
