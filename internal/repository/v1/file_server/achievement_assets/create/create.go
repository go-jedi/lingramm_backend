package create

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

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, tx pgx.Tx, data achievementassets.UploadAndConvertToWebpResponse) (achievementassets.AchievementAssets, error)
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

func (r *Create) Execute(ctx context.Context, tx pgx.Tx, data achievementassets.UploadAndConvertToWebpResponse) (achievementassets.AchievementAssets, error) {
	r.logger.Debug("[create a achievement assets] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO achievement_assets(
		    name_file,
		    server_path_file,
		    client_path_file,
		    extension,
		    quality,
		    old_name_file,
		    old_extension
		) VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING *;
	`

	var aa achievementassets.AchievementAssets

	if err := tx.QueryRow(
		ctxTimeout, q,
		data.NameFile, data.ServerPathFile, data.ClientPathFile,
		data.Extension, data.Quality, data.OldNameFile, data.OldExtension,
	).Scan(
		&aa.ID, &aa.NameFile, &aa.ServerPathFile,
		&aa.ClientPathFile, &aa.Extension, &aa.Quality,
		&aa.OldNameFile, &aa.OldExtension, &aa.CreatedAt, &aa.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create a achievement assets", "err", err)
			return achievementassets.AchievementAssets{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create a achievement assets", "err", err)
		return achievementassets.AchievementAssets{}, fmt.Errorf("could not create a achievement assets: %w", err)
	}

	return aa, nil
}
