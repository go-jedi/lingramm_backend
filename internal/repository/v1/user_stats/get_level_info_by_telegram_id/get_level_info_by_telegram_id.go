package getlevelinfobytelegramid

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

//go:generate mockery --name=IGetLevelInfoByTelegramID --output=mocks --case=underscore
type IGetLevelInfoByTelegramID interface {
	Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userstats.GetLevelInfoByTelegramIDResponse, error)
}

type GetLevelInfoByTelegramID struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *GetLevelInfoByTelegramID {
	r := &GetLevelInfoByTelegramID{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *GetLevelInfoByTelegramID) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *GetLevelInfoByTelegramID) Execute(ctx context.Context, tx pgx.Tx, telegramID string) (userstats.GetLevelInfoByTelegramIDResponse, error) {
	r.logger.Debug("[get level info by telegram id] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.get_level_info($1);`

	var li userstats.GetLevelInfoByTelegramIDResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		telegramID,
	).Scan(
		&li,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while get level info by telegram id", "err", err)
			return userstats.GetLevelInfoByTelegramIDResponse{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to get level info by telegram id", "err", err)
		return userstats.GetLevelInfoByTelegramIDResponse{}, fmt.Errorf("could not get level info by telegram id: %w", err)
	}

	return li, nil
}
