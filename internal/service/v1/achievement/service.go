package achievement

import (
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/achievement/create"
	getdetailbyachievementid "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement/get_detail_by_achievement_id"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	Create                   create.ICreate
	GetDetailByAchievementID getdetailbyachievementid.IGetDetailByAchievementID
}

func New(
	achievementRepository *achievementrepository.Repository,
	achievementAssetsRepository *achievementassetsrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	fileServer *fileserver.FileServer,
) *Service {
	return &Service{
		Create:                   create.New(achievementRepository, achievementAssetsRepository, logger, postgres, fileServer),
		GetDetailByAchievementID: getdetailbyachievementid.New(achievementRepository, logger, postgres),
	}
}
