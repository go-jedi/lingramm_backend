package all

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

//go:generate mockery --name=IAll --output=mocks --case=underscore
type IAll interface {
	Execute(ctx context.Context, tx pgx.Tx) ([]clientassets.ClientAssets, error)
}

type All struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *All {
	r := &All{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *All) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *All) Execute(ctx context.Context, tx pgx.Tx) ([]clientassets.ClientAssets, error) {
	r.logger.Debug("[get all client assets] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM client_assets
		ORDER BY id;
	`

	rows, err := tx.Query(ctxTimeout, q)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get all client assets", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get all client assets", "err", err)
		return nil, fmt.Errorf("could not get all client assets: %w", err)
	}
	defer rows.Close()

	var clientAssets []clientassets.ClientAssets

	for rows.Next() {
		var clientAsset clientassets.ClientAssets

		if err := rows.Scan(
			&clientAsset.ID, &clientAsset.NameFile, &clientAsset.ServerPathFile,
			&clientAsset.ClientPathFile, &clientAsset.Extension, &clientAsset.Quality,
			&clientAsset.OldNameFile, &clientAsset.OldExtension, &clientAsset.CreatedAt, &clientAsset.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row to get all client assets", "err", err)
			return nil, fmt.Errorf("failed to scan row to get all client assets: %w", err)
		}

		clientAssets = append(clientAssets, clientAsset)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed to get all client assets", "err", rows.Err())
		return nil, fmt.Errorf("failed to get all client assets: %w", err)
	}

	return clientAssets, nil
}
