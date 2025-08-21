package leaderboardweeksprocessbatch

import (
	"context"
	"log"

	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ILeaderboardWeeksProcessBatch --output=mocks --case=underscore
type ILeaderboardWeeksProcessBatch interface {
	Execute(ctx context.Context, dto experiencepoint.LeaderboardWeeksProcessBatchDTO) (experiencepoint.LeaderboardWeeksProcessBatchResponse, error)
}

type LeaderboardWeeksProcessBatch struct {
	experiencePointRepository *experiencepointrepository.Repository
	logger                    logger.ILogger
	postgres                  *postgres.Postgres
}

func New(
	experiencePointRepository *experiencepointrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *LeaderboardWeeksProcessBatch {
	return &LeaderboardWeeksProcessBatch{
		experiencePointRepository: experiencePointRepository,
		logger:                    logger,
		postgres:                  postgres,
	}
}

func (s *LeaderboardWeeksProcessBatch) Execute(ctx context.Context, dto experiencepoint.LeaderboardWeeksProcessBatchDTO) (experiencepoint.LeaderboardWeeksProcessBatchResponse, error) {
	s.logger.Debug("[execute a leaderboard weeks process batch] execute service")

	var (
		err    error
		result experiencepoint.LeaderboardWeeksProcessBatchResponse
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return experiencepoint.LeaderboardWeeksProcessBatchResponse{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// execute leaderboard weeks process batch.
	result, err = s.experiencePointRepository.LeaderboardWeeksProcessBatch.Execute(ctx, tx, dto)
	if err != nil {
		return experiencepoint.LeaderboardWeeksProcessBatchResponse{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return experiencepoint.LeaderboardWeeksProcessBatchResponse{}, err
	}

	return result, nil
}
