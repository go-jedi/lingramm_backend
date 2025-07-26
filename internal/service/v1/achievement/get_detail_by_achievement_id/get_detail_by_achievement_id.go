package getdetailbyachievementid

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetDetailByAchievementID --output=mocks --case=underscore
type IGetDetailByAchievementID interface {
	Execute(ctx context.Context, achievementID int64) (achievement.Detail, error)
}

type GetDetailByAchievementID struct {
	achievementRepository *achievementrepository.Repository
	logger                logger.ILogger
	postgres              *postgres.Postgres
}

func New(
	achievementRepository *achievementrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *GetDetailByAchievementID {
	return &GetDetailByAchievementID{
		achievementRepository: achievementRepository,
		logger:                logger,
		postgres:              postgres,
	}
}

func (s *GetDetailByAchievementID) Execute(ctx context.Context, achievementID int64) (achievement.Detail, error) {
	s.logger.Debug("[get detail by achievement id] execute service")

	var (
		err    error
		result achievement.Detail
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return achievement.Detail{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	result, err = s.achievementRepository.GetDetailByAchievementID.Execute(ctx, tx, achievementID)
	if err != nil {
		return achievement.Detail{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return achievement.Detail{}, err
	}

	return result, nil
}
