package getbyname

import (
	"context"
	"errors"
	"fmt"
	"time"

	achievementtype "github.com/go-jedi/lingramm_backend/internal/domain/achievement_type"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetByName --output=mocks --case=underscore
type IGetByName interface {
	Execute(ctx context.Context, tx pgx.Tx, name string) (achievementtype.AchievementType, error)
}

type GetByName struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetByName {
	r := &GetByName{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetByName) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetByName) Execute(ctx context.Context, tx pgx.Tx, name string) (achievementtype.AchievementType, error) {
	r.logger.Debug("[get achievement type by name] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT * 
		FROM achievement_types
		WHERE name = $1;
	`

	var at achievementtype.AchievementType

	if err := tx.QueryRow(
		ctxTimeout, q,
		name,
	).Scan(
		&at.ID, &at.Name, &at.Description,
		&at.IsActive, &at.CreatedAt, &at.UpdatedAt,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get achievement type by name", "err", err)
			return achievementtype.AchievementType{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get achievement type by name", "err", err)
		return achievementtype.AchievementType{}, fmt.Errorf("could not get achievement type by name: %w", err)
	}

	return at, nil
}
