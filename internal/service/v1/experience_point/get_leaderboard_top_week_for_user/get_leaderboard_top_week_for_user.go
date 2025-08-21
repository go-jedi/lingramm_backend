package getleaderboardtopweekforuser

import (
	"context"
	"log"

	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetLeaderboardTopWeekForUser --output=mocks --case=underscore
type IGetLeaderboardTopWeekForUser interface {
	Execute(ctx context.Context, dto experiencepoint.GetLeaderboardTopWeekForUserDTO) ([]experiencepoint.GetLeaderboardTopWeekForUserResponse, error)
}

type GetLeaderboardTopWeekForUser struct {
	experiencePointRepository *experiencepointrepository.Repository
	logger                    logger.ILogger
	postgres                  *postgres.Postgres
}

func New(
	experiencePointRepository *experiencepointrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *GetLeaderboardTopWeekForUser {
	return &GetLeaderboardTopWeekForUser{
		experiencePointRepository: experiencePointRepository,
		logger:                    logger,
		postgres:                  postgres,
	}
}

func (s *GetLeaderboardTopWeekForUser) Execute(ctx context.Context, dto experiencepoint.GetLeaderboardTopWeekForUserDTO) ([]experiencepoint.GetLeaderboardTopWeekForUserResponse, error) {
	s.logger.Debug("[get leaderboard top week for user] execute service")

	var (
		err    error
		result []experiencepoint.GetLeaderboardTopWeekForUserResponse
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

	// get leaderboard top week for user.
	result, err = s.experiencePointRepository.GetLeaderboardTopWeekForUser.Execute(ctx, tx, dto)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
