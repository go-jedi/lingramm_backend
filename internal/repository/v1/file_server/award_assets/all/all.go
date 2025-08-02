package all

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

//go:generate mockery --name=IAll --output=mocks --case=underscore
type IAll interface {
	Execute(ctx context.Context, tx pgx.Tx) ([]awardassets.AwardAssets, error)
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

func (r *All) Execute(ctx context.Context, tx pgx.Tx) ([]awardassets.AwardAssets, error) {
	r.logger.Debug("[get all award assets] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM award_assets
		ORDER BY id;
	`

	rows, err := tx.Query(ctxTimeout, q)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get all award assets", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get all award assets", "err", err)
		return nil, fmt.Errorf("could not get all award assets: %w", err)
	}
	defer rows.Close()

	var awardAssets []awardassets.AwardAssets

	for rows.Next() {
		var awardAsset awardassets.AwardAssets

		if err := rows.Scan(
			&awardAsset.ID, &awardAsset.NameFile,
			&awardAsset.ServerPathFile, &awardAsset.ClientPathFile,
			&awardAsset.Extension, &awardAsset.Quality,
			&awardAsset.OldNameFile, &awardAsset.OldExtension,
			&awardAsset.CreatedAt, &awardAsset.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row to get all award assets", "err", err)
			return nil, fmt.Errorf("failed to scan row to get all award assets: %w", err)
		}

		awardAssets = append(awardAssets, awardAsset)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed to get all award assets", "err", rows.Err())
		return nil, fmt.Errorf("failed to get all award assets: %w", err)
	}

	return awardAssets, nil
}
