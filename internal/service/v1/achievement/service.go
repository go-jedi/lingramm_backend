package achievement

import (
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
	awardassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/award_assets"
	alldetail "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement/all_detail"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/achievement/create"
	deletedetailbyachievementid "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement/delete_detail_by_achievement_id"
	getdetailbyachievementid "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement/get_detail_by_achievement_id"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Service struct {
	All                         alldetail.IAllDetail
	Create                      create.ICreate
	DeleteDetailByAchievementID deletedetailbyachievementid.IDeleteDetailByAchievementID
	GetDetailByAchievementID    getdetailbyachievementid.IGetDetailByAchievementID
}

func New(
	achievementRepository *achievementrepository.Repository,
	achievementAssetsRepository *achievementassetsrepository.Repository,
	awardAssetsRepository *awardassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	fileServer *fileserver.FileServer,
) *Service {
	return &Service{
		All:                         alldetail.New(achievementRepository, logger, postgres),
		Create:                      create.New(achievementRepository, achievementAssetsRepository, awardAssetsRepository, logger, postgres, redis, fileServer),
		DeleteDetailByAchievementID: deletedetailbyachievementid.New(achievementRepository, achievementAssetsRepository, awardAssetsRepository, logger, postgres, redis),
		GetDetailByAchievementID:    getdetailbyachievementid.New(achievementRepository, logger, postgres),
	}
}
