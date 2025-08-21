package leaderboardweeksprocessbatch

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

//go:generate mockery --name=ILeaderboardWeeksProcessBatch --output=mocks --case=underscore
type ILeaderboardWeeksProcessBatch interface {
	Execute(ctx context.Context, tx pgx.Tx, dto experiencepoint.LeaderboardWeeksProcessBatchDTO) (experiencepoint.LeaderboardWeeksProcessBatchResponse, error)
}

type LeaderboardWeeksProcessBatch struct {
	queryTimeout int64
	logger       logger.ILogger
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *LeaderboardWeeksProcessBatch {
	r := &LeaderboardWeeksProcessBatch{
		queryTimeout: queryTimeout,
		logger:       logger,
	}

	r.init()

	return r
}

func (r *LeaderboardWeeksProcessBatch) init() {
	if r.queryTimeout == 0 {
		r.queryTimeout = postgres.DefaultQueryTimeout
	}
}

func (r *LeaderboardWeeksProcessBatch) Execute(ctx context.Context, tx pgx.Tx, dto experiencepoint.LeaderboardWeeksProcessBatchDTO) (experiencepoint.LeaderboardWeeksProcessBatchResponse, error) {
	r.logger.Debug("[execute a leaderboard weeks process batch] execute repository")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.queryTimeout)*time.Second)
	defer cancel()

	q := `SELECT * FROM public.leaderboard_weeks_process_batch($1, $2, $3, $4);`

	var lbwpb experiencepoint.LeaderboardWeeksProcessBatchResponse

	if err := tx.QueryRow(
		ctxTimeout, q,
		dto.WorkerName, dto.BatchSize,
		dto.StatementTimeoutMS, dto.LockTimeoutMS,
	).Scan(
		&lbwpb,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.logger.Error("request timed out while execute a leaderboard weeks process batch", "err", err)
			return experiencepoint.LeaderboardWeeksProcessBatchResponse{}, fmt.Errorf("the request timed out: %w", err)
		}
		r.logger.Error("failed to execute a leaderboard weeks process batch", "err", err)
		return experiencepoint.LeaderboardWeeksProcessBatchResponse{}, fmt.Errorf("could not execute a leaderboard weeks process batch: %w", err)
	}

	return lbwpb, nil
}
