package createachievement

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

//go:generate mockery --name=ICreateAchievement --output=mocks --case=underscore
type ICreateAchievement interface {
	Execute(ctx context.Context, tx pgx.Tx, dto achievement.CreateAchievementDTO) (achievement.Achievement, error)
}

type CreateAchievement struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *CreateAchievement {
	r := &CreateAchievement{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *CreateAchievement) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *CreateAchievement) Execute(ctx context.Context, tx pgx.Tx, dto achievement.CreateAchievementDTO) (achievement.Achievement, error) {
	r.logger.Debug("[create achievement] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		INSERT INTO achievements(
		    achievement_assets_id,
		    award_assets_id,
		    achievement_type_id,
		    name,
		    description
		) VALUES($1, $2, $3, $4, $5)
		RETURNING *;
	`

	var na achievement.Achievement

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.AchievementAssetsID, dto.AwardAssetsID,
		dto.AchievementTypeID, dto.Name, dto.Description,
	).Scan(
		&na.ID, &na.AchievementAssetsID,
		&na.AwardAssetsID, &na.AchievementTypeID,
		&na.Name, &na.Description,
		&na.CreatedAt, &na.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while create achievement", "err", err)
			return achievement.Achievement{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to create achievement", "err", err)
		return achievement.Achievement{}, fmt.Errorf("could not create achievement: %w", err)
	}

	return na, nil
}
