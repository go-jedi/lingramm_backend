package getleaderboardtopweek

import (
	"context"
	"log"

	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetLeaderboardTopWeek --output=mocks --case=underscore
type IGetLeaderboardTopWeek interface {
	Execute(ctx context.Context, dto experiencepoint.GetLeaderboardTopWeekDTO) ([]experiencepoint.GetLeaderboardTopWeekResponse, error)
}

type GetLeaderboardTopWeek struct {
	experiencePointRepository *experiencepointrepository.Repository
	logger                    logger.ILogger
	postgres                  *postgres.Postgres
}

func New(
	experiencePointRepository *experiencepointrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *GetLeaderboardTopWeek {
	return &GetLeaderboardTopWeek{
		experiencePointRepository: experiencePointRepository,
		logger:                    logger,
		postgres:                  postgres,
	}
}

func (s *GetLeaderboardTopWeek) Execute(ctx context.Context, dto experiencepoint.GetLeaderboardTopWeekDTO) ([]experiencepoint.GetLeaderboardTopWeekResponse, error) {
	s.logger.Debug("[get leaderboard top week] execute service")

	var (
		err    error
		result []experiencepoint.GetLeaderboardTopWeekResponse
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// get leaderboard top week.
	result, err = s.experiencePointRepository.GetLeaderboardTopWeek.Execute(ctx, tx, dto)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
