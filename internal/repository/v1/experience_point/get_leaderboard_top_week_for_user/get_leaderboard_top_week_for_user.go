package getleaderboardtopweekforuser

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

//go:generate mockery --name=IGetLeaderboardTopWeekForUser --output=mocks --case=underscore
type IGetLeaderboardTopWeekForUser interface {
	Execute(ctx context.Context, tx pgx.Tx, dto experiencepoint.GetLeaderboardTopWeekForUserDTO) ([]experiencepoint.GetLeaderboardTopWeekForUserResponse, error)
}

type GetLeaderboardTopWeekForUser struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetLeaderboardTopWeekForUser {
	r := &GetLeaderboardTopWeekForUser{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetLeaderboardTopWeekForUser) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetLeaderboardTopWeekForUser) Execute(ctx context.Context, tx pgx.Tx, dto experiencepoint.GetLeaderboardTopWeekForUserDTO) ([]experiencepoint.GetLeaderboardTopWeekForUserResponse, error) {
	r.logger.Debug("[get leaderboard top week for user] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.leaderboard_weeks_top_week_for_user_get($1, $2, $3);`

	var lbtwfu []experiencepoint.GetLeaderboardTopWeekForUserResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.TelegramID, dto.Limit,
		dto.TZ,
	).Scan(&lbtwfu); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get leaderboard top week for user", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get leaderboard top week for user", "err", err)
		return nil, fmt.Errorf("could not get leaderboard top week for user: %w", err)
	}

	return lbtwfu, nil
}
