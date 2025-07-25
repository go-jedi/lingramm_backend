package create

import (
	"context"
	"log"
	"os"

	"github.com/go-jedi/lingramm_backend/internal/domain/achievement"
	achievementassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/achievement_assets"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, dto achievement.CreateDTO) (achievement.CreateResponse, error)
}

type Create struct {
	achievementRepository       *achievementrepository.Repository
	achievementAssetsRepository *achievementassetsrepository.Repository
	logger                      logger.ILogger
	postgres                    *postgres.Postgres
	fileServer                  *fileserver.FileServer
}

func New(
	achievementRepository *achievementrepository.Repository,
	achievementAssetsRepository *achievementassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	fileServer *fileserver.FileServer,
) *Create {
	return &Create{
		achievementRepository:       achievementRepository,
		achievementAssetsRepository: achievementAssetsRepository,
		logger:                      logger,
		postgres:                    postgres,
		fileServer:                  fileServer,
	}
}

func (s *Create) Execute(ctx context.Context, dto achievement.CreateDTO) (achievement.CreateResponse, error) {
	s.logger.Debug("[create a achievement] execute service")

	var (
		err       error
		imageData achievementassets.UploadAndConvertToWebpResponse
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return achievement.CreateResponse{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
			s.cleanupTempFile(imageData.ServerPathFile)
		}
	}()

	// check achievement exists by code
	existsAchievementByCode, err := s.achievementRepository.ExistsAchievementByCode.Execute(ctx, tx, dto.Code)
	if err != nil {
		return achievement.CreateResponse{}, err
	}

	if existsAchievementByCode {
		return achievement.CreateResponse{}, apperrors.ErrAchievementAlreadyExists
	}

	// check achievement condition exists by condition type.
	existsAchievementConditionByConditionType, err := s.achievementRepository.ExistsAchievementConditionByConditionType.Execute(ctx, tx, dto.ConditionType)
	if err != nil {
		return achievement.CreateResponse{}, err
	}

	if existsAchievementConditionByConditionType {
		return achievement.CreateResponse{}, apperrors.ErrAchievementConditionAlreadyExists
	}

	// convert png or jpg image to webp and upload.
	imageData, err = s.fileServer.AchievementAssets.UploadAndConvertToWebP(ctx, dto.FileHeader)
	if err != nil {
		return achievement.CreateResponse{}, err
	}

	// create achievement asset.
	resultAchievementAsset, err := s.createAchievementAsset(ctx, tx, imageData)
	if err != nil {
		return achievement.CreateResponse{}, err
	}

	// create achievement.
	resultAchievement, err := s.createAchievement(ctx, tx, dto, resultAchievementAsset.ID)
	if err != nil {
		return achievement.CreateResponse{}, err
	}

	// create achievement condition.
	resultAchievementCondition, err := s.createAchievementCondition(ctx, tx, dto, resultAchievement.ID)
	if err != nil {
		return achievement.CreateResponse{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return achievement.CreateResponse{}, err
	}

	return achievement.CreateResponse{
		Achievement:       resultAchievement,
		Condition:         resultAchievementCondition,
		AchievementAssets: resultAchievementAsset,
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

// createAchievement create achievement.
func (s *Create) createAchievement(ctx context.Context, tx pgx.Tx, dto achievement.CreateDTO, achievementAssetsID int64) (achievement.Achievement, error) {
	createAchievementDTO := achievement.CreateAchievementDTO{
		AchievementAssetsID: achievementAssetsID,
		Description:         dto.Description,
		Code:                dto.Code,
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
func (s *Create) createAchievementCondition(ctx context.Context, tx pgx.Tx, dto achievement.CreateDTO, achievementID int64) (achievement.Condition, error) {
	createAchievementConditionDTO := achievement.CreateAchievementConditionDTO{
		AchievementID: achievementID,
		Value:         dto.Value,
		ConditionType: dto.ConditionType,
		Operator:      dto.Operator,
	}

	// create achievement condition.
	result, err := s.achievementRepository.CreateAchievementCondition.Execute(ctx, tx, createAchievementConditionDTO)
	if err != nil {
		return achievement.Condition{}, err
	}

	return result, nil
}

// cleanupTempFile cleanup temp file.
func (s *Create) cleanupTempFile(path string) {
	if path == "" {
		return
	}

	if err := os.Remove(path); err != nil {
		s.logger.Warn("failed to remove temporary file", "path", path, "error", err)
	} else {
		s.logger.Debug("successfully removed temporary file", "path", path)
	}
}
