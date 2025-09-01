package deletedetailbyachievementid

import (
	"context"
	"log"
	"os"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/achievement_assets"
	awardassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/award_assets"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
	awardassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/award_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IDeleteDetailByAchievementID --output=mocks --case=underscore
type IDeleteDetailByAchievementID interface {
	Execute(ctx context.Context, achievementID int64) (achievement.Detail, error)
}

type DeleteDetailByAchievementID struct {
	achievementRepository       *achievementrepository.Repository
	achievementAssetsRepository *achievementassetsrepository.Repository
	awardAssetsRepository       *awardassetsrepository.Repository
	logger                      logger.ILogger
	postgres                    *postgres.Postgres
	redis                       *redis.Redis
}

func New(
	achievementRepository *achievementrepository.Repository,
	achievementAssetsRepository *achievementassetsrepository.Repository,
	awardAssetsRepository *awardassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
) *DeleteDetailByAchievementID {
	return &DeleteDetailByAchievementID{
		achievementRepository:       achievementRepository,
		achievementAssetsRepository: achievementAssetsRepository,
		awardAssetsRepository:       awardAssetsRepository,
		logger:                      logger,
		postgres:                    postgres,
		redis:                       redis,
	}
}

func (s *DeleteDetailByAchievementID) Execute(ctx context.Context, achievementID int64) (achievement.Detail, error) {
	s.logger.Debug("[delete detail by achievement id] execute service")

	var (
		err                         error
		resultAchievement           achievement.Achievement
		resultAchievementAsset      achievementassets.AchievementAssets
		resultAwardAssets           awardassets.AwardAssets
		existsAchievementByID       bool
		existsAchievementAssetsByID bool
		existsAwardAssetsByID       bool
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

	// check achievement exists by id.
	existsAchievementByID, err = s.achievementRepository.ExistsAchievementByID.Execute(ctx, tx, achievementID)
	if err != nil {
		return achievement.Detail{}, err
	}

	if !existsAchievementByID { // if achievement does not exist.
		err = apperrors.ErrAchievementDoesNotExist
		return achievement.Detail{}, err
	}

	// delete achievement by id.
	resultAchievement, err = s.achievementRepository.DeleteAchievementByID.Execute(ctx, tx, achievementID)
	if err != nil {
		return achievement.Detail{}, err
	}

	// check achievement assets exists by id.
	existsAchievementAssetsByID, err = s.achievementAssetsRepository.ExistsByID.Execute(ctx, tx, resultAchievement.AchievementAssetsID)
	if err != nil {
		return achievement.Detail{}, err
	}

	if !existsAchievementAssetsByID { // if achievement assets does not exist.
		err = apperrors.ErrAchievementAssetsDoesNotExist
		return achievement.Detail{}, err
	}

	// delete achievement assets by id.
	resultAchievementAsset, err = s.achievementAssetsRepository.DeleteByID.Execute(ctx, tx, resultAchievement.AchievementAssetsID)
	if err != nil {
		return achievement.Detail{}, err
	}

	// check award assets exists by id.
	existsAwardAssetsByID, err = s.awardAssetsRepository.ExistsByID.Execute(ctx, tx, resultAchievement.AwardAssetsID)
	if err != nil {
		return achievement.Detail{}, err
	}

	if !existsAwardAssetsByID { // if award assets does not exist.
		err = apperrors.ErrAwardAssetsDoesNotExist
		return achievement.Detail{}, err
	}

	// delete award assets by id.
	resultAwardAssets, err = s.awardAssetsRepository.DeleteByID.Execute(ctx, tx, resultAchievement.AwardAssetsID)
	if err != nil {
		return achievement.Detail{}, err
	}

	// remove file achievement.
	s.deleteAchievementFile(ctx, resultAchievementAsset.NameFileWithoutExtension, resultAchievementAsset.ServerPathFile)
	// remove file award.
	s.deleteAwardFile(ctx, resultAwardAssets.NameFileWithoutExtension, resultAwardAssets.ServerPathFile)

	err = tx.Commit(ctx)
	if err != nil {
		return achievement.Detail{}, err
	}

	return achievement.Detail{
		Achievement:       resultAchievement,
		AchievementAssets: resultAchievementAsset,
		AwardAssets:       resultAwardAssets,
	}, nil
}

// deleteAchievementFile delete achievement file.
func (s *DeleteDetailByAchievementID) deleteAchievementFile(ctx context.Context, nameFileWithoutExtension string, path string) {
	if err := os.Remove(path); err != nil {
		s.logger.Warn("failed to remove achievement file", "path", path, "error", err)

		if err := s.redis.UnDeleteFileAchievement.Set(ctx, nameFileWithoutExtension, path); err != nil {
			s.logger.Warn("failed to set un delete achievement file", "path", path, "error", err)
		}
	} else {
		s.logger.Debug("achievement file removed", "path", path)
	}
}

// deleteAwardFile delete file.
func (s *DeleteDetailByAchievementID) deleteAwardFile(ctx context.Context, nameFileWithoutExtension string, path string) {
	if err := os.Remove(path); err != nil {
		s.logger.Warn("failed to remove award file", "path", path, "error", err)

		if err := s.redis.UnDeleteFileAward.Set(ctx, nameFileWithoutExtension, path); err != nil {
			s.logger.Warn("failed to set un delete award file", "path", path, "error", err)
		}
	} else {
		s.logger.Debug("award file removed", "path", path)
	}
}
