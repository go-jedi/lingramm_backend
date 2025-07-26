package deletedetailbyachievementid

import (
	"context"
	"log"
	"os"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IDeleteDetailByAchievementID --output=mocks --case=underscore
type IDeleteDetailByAchievementID interface {
	Execute(ctx context.Context, achievementID int64) (achievement.Detail, error)
}

type DeleteDetailByAchievementID struct {
	achievementRepository       *achievementrepository.Repository
	achievementAssetsRepository *achievementassetsrepository.Repository
	logger                      logger.ILogger
	postgres                    *postgres.Postgres
}

func New(
	achievementRepository *achievementrepository.Repository,
	achievementAssetsRepository *achievementassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *DeleteDetailByAchievementID {
	return &DeleteDetailByAchievementID{
		achievementRepository:       achievementRepository,
		achievementAssetsRepository: achievementAssetsRepository,
		logger:                      logger,
		postgres:                    postgres,
	}
}

func (s *DeleteDetailByAchievementID) Execute(ctx context.Context, achievementID int64) (achievement.Detail, error) {
	s.logger.Debug("[delete detail by achievement id] execute service")

	var err error

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

	// check achievement exists by id.
	existsAchievementByID, err := s.achievementRepository.ExistsAchievementByID.Execute(ctx, tx, achievementID)
	if err != nil {
		return achievement.Detail{}, err
	}

	if !existsAchievementByID {
		return achievement.Detail{}, apperrors.ErrAchievementDoesNotExist
	}

	// delete achievement condition by achievement id.
	resultAchievementCondition, err := s.achievementRepository.DeleteAchievementConditionByAchievementID.Execute(ctx, tx, achievementID)
	if err != nil {
		return achievement.Detail{}, err
	}

	// delete achievement by id.
	resultAchievement, err := s.achievementRepository.DeleteAchievementByID.Execute(ctx, tx, achievementID)
	if err != nil {
		return achievement.Detail{}, err
	}

	// delete achievement assets by id.
	resultAchievementAsset, err := s.achievementAssetsRepository.DeleteByID.Execute(ctx, tx, resultAchievement.AchievementAssetsID)
	if err != nil {
		return achievement.Detail{}, err
	}

	// remove file.
	s.deleteAchievementFile(resultAchievementAsset.ServerPathFile)

	err = tx.Commit(ctx)
	if err != nil {
		return achievement.Detail{}, err
	}

	return achievement.Detail{
		Achievement:       resultAchievement,
		Condition:         resultAchievementCondition,
		AchievementAssets: resultAchievementAsset,
	}, nil
}

// deleteAchievementFile delete file.
func (s *DeleteDetailByAchievementID) deleteAchievementFile(path string) {
	if err := os.Remove(path); err != nil {
		s.logger.Warn("failed to remove asset file", "path", path, "error", err)
	} else {
		s.logger.Debug("asset file removed", "path", path)
	}
}
