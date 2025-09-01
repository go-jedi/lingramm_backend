package alldetailbytelegramid

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

//go:generate mockery --name=IAllDetailByTelegramID --output=mocks --case=underscore
type IAllDetailByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) ([]userachievement.Detail, error)
}

type AllDetailByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *AllDetailByTelegramID {
	r := &AllDetailByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *AllDetailByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *AllDetailByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) ([]userachievement.Detail, error) {
	r.logger.Debug("[get all user achievements detail] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT
			JSONB_AGG(
				JSONB_BUILD_OBJECT(
					'id', ua.id,
					'telegram_id', ua.telegram_id,
					'name', a.name,
					'description', a.description,
					'achievement_path_file', aa.client_path_file,
					'award_path_file', awa.client_path_file
				)
			)
		FROM user_achievements ua
		INNER JOIN achievements a ON ua.achievement_id = a.id
		INNER JOIN achievement_assets aa ON a.achievement_assets_id = aa.id
		INNER JOIN award_assets awa ON a.award_assets_id = awa.id
		WHERE ua.telegram_id = $1;
	`

	var d []userachievement.Detail

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
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
