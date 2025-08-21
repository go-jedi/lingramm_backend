package getleaderboardtopweek

import (
	"context"
	"errors"
	"fmt"
	"time"

	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetLeaderboardTopWeek --output=mocks --case=underscore
type IGetLeaderboardTopWeek interface {
	Execute(ctx context.Context, tx pgx.Tx, dto experiencepoint.GetLeaderboardTopWeekDTO) ([]experiencepoint.GetLeaderboardTopWeekResponse, error)
}

type GetLeaderboardTopWeek struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetLeaderboardTopWeek {
	r := &GetLeaderboardTopWeek{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetLeaderboardTopWeek) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetLeaderboardTopWeek) Execute(ctx context.Context, tx pgx.Tx, dto experiencepoint.GetLeaderboardTopWeekDTO) ([]experiencepoint.GetLeaderboardTopWeekResponse, error) {
	r.logger.Debug("[get leaderboard top week] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.leaderboard_weeks_top_week_get($1, $2);`

	var lbtw []experiencepoint.GetLeaderboardTopWeekResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.Limit, dto.TZ,
	).Scan(&lbtw); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get leaderboard top week", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get leaderboard top week", "err", err)
		return nil, fmt.Errorf("could not get leaderboard top week: %w", err)
	}

	return lbtw, nil
}
