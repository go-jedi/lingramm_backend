package getlevelbytelegramid

import (
	"context"
	"errors"
	"fmt"
	"time"

	userstats "github.com/go-jedi/lingramm_backend/internal/domain/user_stats"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetLevelByTelegramID --output=mocks --case=underscore
type IGetLevelByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userstats.GetLevelByTelegramIDResponse, error)
}

type GetLevelByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetLevelByTelegramID {
	r := &GetLevelByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetLevelByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetLevelByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userstats.GetLevelByTelegramIDResponse, error) {
	r.logger.Debug("[get level user by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `
		SELECT telegram_id, level
		FROM user_stats
		WHERE telegram_id = $1;
	`

	var l userstats.GetLevelByTelegramIDResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(
		&l.TelegramID, &l.Level,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get level user by telegram id", "err", err)
			return userstats.GetLevelByTelegramIDResponse{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get level user by telegram id", "err", err)
		return userstats.GetLevelByTelegramIDResponse{}, fmt.Errorf("could not get level user by telegram id: %w", err)
	}

	return l, nil
}
