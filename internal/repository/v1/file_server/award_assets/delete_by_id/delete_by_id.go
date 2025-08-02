package deletebyid

import (
	"context"
	"errors"
	"fmt"
	"time"

	awardassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/award_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IDeleteByID --output=mocks --case=underscore
type IDeleteByID interface {
	Execute(ctx context.Context, tx pgx.Tx, id int64) (awardassets.AwardAssets, error)
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

func (r *DeleteByID) Execute(ctx context.Context, tx pgx.Tx, id int64) (awardassets.AwardAssets, error) {
	r.logger.Debug("[delete award assets by id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		DELETE FROM award_assets
		WHERE id = $1
		RETURNING *;
	`

	var daa awardassets.AwardAssets

	if err := tx.QueryRow(
		ctxTimeout, q,
		id,
	).Scan(
		&daa.ID, &daa.NameFile, &daa.ServerPathFile,
		&daa.ClientPathFile, &daa.Extension, &daa.Quality,
		&daa.OldNameFile, &daa.OldExtension, &daa.CreatedAt, &daa.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while delete award assets by id", "err", err)
			return awardassets.AwardAssets{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to delete award assets by id", "err", err)
		return awardassets.AwardAssets{}, fmt.Errorf("could not delete award assets by id: %w", err)
	}

	return daa, nil
}
