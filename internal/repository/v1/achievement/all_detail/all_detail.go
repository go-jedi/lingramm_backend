package alldetail

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAllDetail --output=mocks --case=underscore
type IAllDetail interface {
	Execute(ctx context.Context, tx pgx.Tx) ([]achievement.Detail, error)
}

type AllDetail struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *AllDetail {
	r := &AllDetail{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *AllDetail) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *AllDetail) Execute(ctx context.Context, tx pgx.Tx) ([]achievement.Detail, error) {
	r.logger.Debug("[get all detail] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'achievement', JSONB_BUILD_OBJECT(
						'id', a.id,
						'achievement_assets_id', a.achievement_assets_id,
						'award_assets_id', a.award_assets_id,
						'code', a.code,
						'name', a.name,
						'description', a.description,
						'created_at', a.created_at,
						'updated_at', a.updated_at
					),
					'condition', JSONB_BUILD_OBJECT(
						'id', ac.id,
						'achievement_id', ac.achievement_id,
						'condition_type', ac.condition_type,
						'operator', ac.operator,
						'value', ac.value,
						'created_at', ac.created_at,
						'updated_at', ac.updated_at
					),
					'achievement_assets', JSONB_BUILD_OBJECT(
						'id', aa.id,
						'name_file', aa.name_file,
						'name_file_without_extension', aa.name_file_without_extension,
						'server_path_file', aa.server_path_file,
						'client_path_file', aa.client_path_file,
						'extension', aa.extension,
						'quality', aa.quality,
						'old_name_file', aa.old_name_file,
						'old_extension', aa.old_extension,
						'created_at', aa.created_at,
						'updated_at', aa.updated_at
					),
						'award_assets', JSONB_BUILD_OBJECT(
						'id', awa.id,
						'name_file', awa.name_file,
						'name_file_without_extension', awa.name_file_without_extension,
						'server_path_file', awa.server_path_file,
						'client_path_file', awa.client_path_file,
						'extension', awa.extension,
						'quality', awa.quality,
						'old_name_file', awa.old_name_file,
						'old_extension', awa.old_extension,
						'created_at', awa.created_at,
						'updated_at', awa.updated_at
					)
				)
			)
		FROM achievements a
		INNER JOIN achievement_conditions ac ON a.id = ac.achievement_id
		INNER JOIN achievement_assets aa ON a.achievement_assets_id = aa.id
		INNER JOIN award_assets awa ON a.award_assets_id = awa.id;
	`

	var d []achievement.Detail

	if err := tx.QueryRow(
		ctxTimeout, q,
	).Scan(&d); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get all detail", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get all detail", "err", err)
		return nil, fmt.Errorf("could not get all detail: %w", err)
	}

	return d, nil
}
