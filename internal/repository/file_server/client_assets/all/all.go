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
	a := &All{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	a.init()

	return a
}

func (a *All) init() {
	if a.queryTimeout == 0 {
		a.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (a *All) Execute(ctx context.Context, tx pgx.Tx) ([]clientassets.ClientAssets, error) {
	a.logger.Debug("[get all client assets] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(a.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM client_assets
		ORDER BY id;
	`

	rows, err := tx.Query(ctxTimeout, q)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			a.logger.Error("request timed out while get all client assets", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		a.logger.Error("failed to get all client assets", "err", err)
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
			a.logger.Error("failed to scan row to get all client assets", "err", err)
			return nil, fmt.Errorf("failed to scan row to get all client assets: %w", err)
		}

		clientAssets = append(clientAssets, clientAsset)
	}

	if err := rows.Err(); err != nil {
		a.logger.Error("failed to get all client assets", "err", rows.Err())
		return nil, fmt.Errorf("failed to get all client assets: %w", err)
	}

	return clientAssets, nil
}
