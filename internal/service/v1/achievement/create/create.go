package create

import (
	"context"
	"log"
	"os"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementtype "github.com/go-jedi/lingramm_backend/internal/domain/achievement_type"
	achievementassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/achievement_assets"
	awardassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/award_assets"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	achievementtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement_type"
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
	awardassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/award_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, dto achievement.CreateDTO) (achievement.Detail, error)
}

type Create struct {
	achievementRepository       *achievementrepository.Repository
	achievementAssetsRepository *achievementassetsrepository.Repository
	awardAssetsRepository       *awardassetsrepository.Repository
	achievementTypeRepository   *achievementtyperepository.Repository
	logger                      logger.ILogger
	postgres                    *postgres.Postgres
	redis                       *redis.Redis
	fileServer                  *fileserver.FileServer
}

func New(
	achievementRepository *achievementrepository.Repository,
	achievementAssetsRepository *achievementassetsrepository.Repository,
	awardAssetsRepository *awardassetsrepository.Repository,
	achievementTypeRepository *achievementtyperepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	fileServer *fileserver.FileServer,
) *Create {
	return &Create{
		achievementRepository:       achievementRepository,
		achievementAssetsRepository: achievementAssetsRepository,
		awardAssetsRepository:       awardAssetsRepository,
		achievementTypeRepository:   achievementTypeRepository,
		logger:                      logger,
		postgres:                    postgres,
		redis:                       redis,
		fileServer:                  fileServer,
	}
}

func (s *Create) Execute(ctx context.Context, dto achievement.CreateDTO) (achievement.Detail, error) {
	s.logger.Debug("[create a achievement] execute service")

	var (
		err                        error
		imageAchievementData       achievementassets.UploadAndConvertToWebpResponse
		achievementAsset           achievementassets.AchievementAssets
		imageAwardData             awardassets.UploadAndConvertToWebpResponse
		awardAsset                 awardassets.AwardAssets
		achievementTypeData        achievementtype.AchievementType
		resultAchievement          achievement.Achievement
		resultAchievementCondition achievement.Condition
		existsAchievement          bool
		existsAchievementCondition bool
		existsAchievementType      bool
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
			s.deleteAchievementTempFile(ctx, imageAchievementData.NameFileWithoutExtension, imageAchievementData.ServerPathFile)
			s.deleteAwardTempFile(ctx, imageAwardData.NameFileWithoutExtension, imageAwardData.ServerPathFile)
		}
	}()

	// check achievement exists by event type.
	existsAchievement, err = s.achievementRepository.ExistsAchievementByAchievementType.Execute(ctx, tx, dto.AchievementType)
	if err != nil {
		return achievement.Detail{}, err
	}

	if existsAchievement { // if achievement already exist.
		err = apperrors.ErrAchievementAlreadyExists
		return achievement.Detail{}, err
	}

	// check achievement condition exists by event type.
	existsAchievementCondition, err = s.achievementRepository.ExistsAchievementConditionByAchievementType.Execute(ctx, tx, dto.AchievementType)
	if err != nil {
		return achievement.Detail{}, err
	}

	if existsAchievementCondition { // if achievement condition already exist.
		err = apperrors.ErrAchievementConditionAlreadyExists
		return achievement.Detail{}, err
	}

	// convert png or jpg image achievement to webp and upload.
	imageAchievementData, err = s.fileServer.AchievementAssets.UploadAndConvertToWebP(ctx, dto.FileAchievementHeader)
	if err != nil {
		return achievement.Detail{}, err
	}

	// create achievement asset.
	achievementAsset, err = s.createAchievementAsset(ctx, tx, imageAchievementData)
	if err != nil {
		return achievement.Detail{}, err
	}

	// convert png or jpg image award to webp and upload.
	imageAwardData, err = s.fileServer.AwardAssets.UploadAndConvertToWebP(ctx, dto.FileAwardHeader)
	if err != nil {
		return achievement.Detail{}, err
	}

	// create award asset.
	awardAsset, err = s.createAwardAsset(ctx, tx, imageAwardData)
	if err != nil {
		return achievement.Detail{}, err
	}

	// check achievement type exists by name.
	existsAchievementType, err = s.achievementTypeRepository.ExistsByName.Execute(ctx, tx, dto.AchievementType)
	if err != nil {
		return achievement.Detail{}, err
	}

	if !existsAchievementType { // if achievement type does not exist.
		err = apperrors.ErrAchievementTypeDoesNotExist
		return achievement.Detail{}, err
	}

	// get achievement type by name.
	achievementTypeData, err = s.achievementTypeRepository.GetByName.Execute(ctx, tx, dto.AchievementType)
	if err != nil {
		return achievement.Detail{}, err
	}

	// create achievement.
	resultAchievement, err = s.createAchievement(ctx, tx, dto, achievementAsset.ID, awardAsset.ID, achievementTypeData.ID)
	if err != nil {
		return achievement.Detail{}, err
	}

	// create achievement condition.
	resultAchievementCondition, err = s.createAchievementCondition(ctx, tx, dto, resultAchievement.ID, achievementTypeData.ID)
	if err != nil {
		return achievement.Detail{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return achievement.Detail{}, err
	}

	return achievement.Detail{
		Achievement:       resultAchievement,
		Condition:         resultAchievementCondition,
		AchievementAssets: achievementAsset,
		AwardAssets:       awardAsset,
	}, nil
}

// createAchievementAsset create achievement assets.
func (s *Create) createAchievementAsset(ctx context.Context, tx pgx.Tx, imageData achievementassets.UploadAndConvertToWebpResponse) (achievementassets.AchievementAssets, error) {
	// create achievement assets.
	result, err := s.achievementAssetsRepository.Create.Execute(ctx, tx, imageData)
	if err != nil {
		return achievementassets.AchievementAssets{}, err
	}

	return result, nil
}

// createAwardAsset create award assets.
func (s *Create) createAwardAsset(ctx context.Context, tx pgx.Tx, imageData awardassets.UploadAndConvertToWebpResponse) (awardassets.AwardAssets, error) {
	// create award assets.
	result, err := s.awardAssetsRepository.Create.Execute(ctx, tx, imageData)
	if err != nil {
		return awardassets.AwardAssets{}, err
	}

	return result, nil
}

// createAchievement create achievement.
func (s *Create) createAchievement(ctx context.Context, tx pgx.Tx, dto achievement.CreateDTO, achievementAssetsID int64, awardAssetsID int64, achievementTypeID int64) (achievement.Achievement, error) {
	createAchievementDTO := achievement.CreateAchievementDTO{
		AchievementAssetsID: achievementAssetsID,
		AwardAssetsID:       awardAssetsID,
		AchievementTypeID:   achievementTypeID,
		Description:         dto.Description,
		Name:                dto.Name,
	}

	// create achievement.
	result, err := s.achievementRepository.CreateAchievement.Execute(ctx, tx, createAchievementDTO)
	if err != nil {
		return achievement.Achievement{}, err
	}

	return result, nil
}

// createAchievementCondition create achievement condition.
func (s *Create) createAchievementCondition(ctx context.Context, tx pgx.Tx, dto achievement.CreateDTO, achievementID int64, achievementTypeID int64) (achievement.Condition, error) {
	createAchievementConditionDTO := achievement.CreateAchievementConditionDTO{
		AchievementID:     achievementID,
		AchievementTypeID: achievementTypeID,
		Value:             dto.Value,
		Operator:          dto.Operator,
	}

	// create achievement condition.
	result, err := s.achievementRepository.CreateAchievementCondition.Execute(ctx, tx, createAchievementConditionDTO)
	if err != nil {
		return achievement.Condition{}, err
	}

	return result, nil
}

// deleteAchievementTempFile delete achievement temp file.
func (s *Create) deleteAchievementTempFile(ctx context.Context, nameFileWithoutExtension string, path string) {
	if path == "" {
		return
	}

	if err := os.Remove(path); err != nil {
		s.logger.Warn("failed to remove temporary achievement file", "path", path, "error", err)

		if err := s.redis.UnDeleteFileAchievement.Set(ctx, nameFileWithoutExtension, path); err != nil {
			s.logger.Warn("failed to set un delete achievement file", "path", path, "error", err)
		}
	} else {
		s.logger.Debug("successfully removed temporary achievement file", "path", path)
	}
}

// deleteAwardTempFile delete award temp file.
func (s *Create) deleteAwardTempFile(ctx context.Context, nameFileWithoutExtension string, path string) {
	if path == "" {
		return
	}

	if err := os.Remove(path); err != nil {
		s.logger.Warn("failed to remove temporary award file", "path", path, "error", err)

		if err := s.redis.UnDeleteFileAward.Set(ctx, nameFileWithoutExtension, path); err != nil {
			s.logger.Warn("failed to set un delete award file", "path", path, "error", err)
		}
	} else {
		s.logger.Debug("successfully removed temporary award file", "path", path)
	}
}
