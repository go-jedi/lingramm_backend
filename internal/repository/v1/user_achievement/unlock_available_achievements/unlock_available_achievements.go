package unlockavailableachievements

import (
	"context"
	"errors"
	"fmt"
	"time"

	userachievement "github.com/go-jedi/lingramm_backend/internal/domain/user_achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IUnlockAvailableAchievements --output=mocks --case=underscore
type IUnlockAvailableAchievements interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) ([]userachievement.UnlockAvailableAchievementsResponse, error)
}

type UnlockAvailableAchievements struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *UnlockAvailableAchievements {
	r := &UnlockAvailableAchievements{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *UnlockAvailableAchievements) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *UnlockAvailableAchievements) Execute(ctx context.Context, tx pgx.Tx, telegramID string) ([]userachievement.UnlockAvailableAchievementsResponse, error) {
	r.logger.Debug("[unlock available achievements] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.unlock_available_achievements($1);`

	var uaa []userachievement.UnlockAvailableAchievementsResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(&uaa); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while unlock available achievements", "err", err)
			return nil, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to unlock available achievements", "err", err)
		return nil, fmt.Errorf("could not unlock available achievements: %w", err)
	}

	return uaa, nil
}
