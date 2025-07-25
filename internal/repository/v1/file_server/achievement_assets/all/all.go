package all

import (
	"context"
	"errors"
	"fmt"
	"time"

	achievementassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/achievement_assets"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAll --output=mocks --case=underscore
type IAll interface {
	Execute(ctx context.Context, tx pgx.Tx) ([]achievementassets.AchievementAssets, error)
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

func (r *All) Execute(ctx context.Context, tx pgx.Tx) ([]achievementassets.AchievementAssets, error) {
	r.logger.Debug("[get all achievement assets] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT *
		FROM achievement_assets
		ORDER BY id;
	`

	rows, err := tx.Query(ctxTimeout, q)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get all achievement assets", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get all achievement assets", "err", err)
		return nil, fmt.Errorf("could not get all achievement assets: %w", err)
	}
	defer rows.Close()

	var achievementAssets []achievementassets.AchievementAssets

	for rows.Next() {
		var achievementAsset achievementassets.AchievementAssets

		if err := rows.Scan(
			&achievementAsset.ID, &achievementAsset.NameFile,
			&achievementAsset.ServerPathFile, &achievementAsset.ClientPathFile,
			&achievementAsset.Extension, &achievementAsset.Quality,
			&achievementAsset.OldNameFile, &achievementAsset.OldExtension,
			&achievementAsset.CreatedAt, &achievementAsset.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row to get all achievement assets", "err", err)
			return nil, fmt.Errorf("failed to scan row to get all achievement assets: %w", err)
		}

		achievementAssets = append(achievementAssets, achievementAsset)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed to get all achievement assets", "err", rows.Err())
		return nil, fmt.Errorf("failed to get all achievement assets: %w", err)
	}

	return achievementAssets, nil
}
