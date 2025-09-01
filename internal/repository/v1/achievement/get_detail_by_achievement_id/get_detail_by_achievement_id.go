package getdetailbyachievementid

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

//go:generate mockery --name=IGetDetailByAchievementID --output=mocks --case=underscore
type IGetDetailByAchievementID interface {
	Execute(ctx context.Context, tx pgx.Tx, achievementID int64) (achievement.Detail, error)
}

type GetDetailByAchievementID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetDetailByAchievementID {
	r := &GetDetailByAchievementID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetDetailByAchievementID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetDetailByAchievementID) Execute(ctx context.Context, tx pgx.Tx, achievementID int64) (achievement.Detail, error) {
	r.logger.Debug("[get detail by achievement id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT
			JSONB_BUILD_OBJECT(
				'achievement', JSONB_BUILD_OBJECT(
					'id', a.id,
					'achievement_assets_id', a.achievement_assets_id,
					'award_assets_id', a.award_assets_id,
					'achievement_type_id', a.achievement_type_id,
					'name', a.name,
					'description', a.description,
					'created_at', a.created_at,
					'updated_at', a.updated_at
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
		FROM achievements a
		INNER JOIN achievement_assets aa ON a.achievement_assets_id = aa.id
		INNER JOIN award_assets awa ON a.award_assets_id = awa.id
		WHERE a.id = $1;
	`

	var d achievement.Detail

	if err := tx.QueryRow(
		ctxTimeout, q, achievementID,
	).Scan(&d); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get detail by achievement id", "err", err)
			return achievement.Detail{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get detail by achievement id", "err", err)
		return achievement.Detail{}, fmt.Errorf("could not get detail by achievement id: %w", err)
	}

	return d, nil
}
