package achievement

import (
	alldetail "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/all_detail"
	createachievement "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/create_achievement"
	deleteachievementbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/delete_achievement_by_id"
	existsachievementbyachievementtype "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/exists_achievement_by_achievement_type"
	existsachievementbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/exists_achievement_by_id"
	getdetailbyachievementid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/get_detail_by_achievement_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	AllDetail                          alldetail.IAllDetail
	CreateAchievement                  createachievement.ICreateAchievement
	DeleteAchievementByID              deleteachievementbyid.IDeleteAchievementByID
	ExistsAchievementByAchievementType existsachievementbyachievementtype.IExistsAchievementByAchievementType
	ExistsAchievementByID              existsachievementbyid.IExistsAchievementByID
	GetDetailByAchievementID           getdetailbyachievementid.IGetDetailByAchievementID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		AllDetail:                          alldetail.New(queryTimeout, logger),
		CreateAchievement:                  createachievement.New(queryTimeout, logger),
		DeleteAchievementByID:              deleteachievementbyid.New(queryTimeout, logger),
		ExistsAchievementByAchievementType: existsachievementbyachievementtype.New(queryTimeout, logger),
		ExistsAchievementByID:              existsachievementbyid.New(queryTimeout, logger),
		GetDetailByAchievementID:           getdetailbyachievementid.New(queryTimeout, logger),
	}
}
