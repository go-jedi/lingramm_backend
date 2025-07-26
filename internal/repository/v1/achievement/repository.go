package achievement

import (
	alldetail "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/all_detail"
	createachievement "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/create_achievement"
	createachievementcondition "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/create_achievement_condition"
	deleteachievementbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/delete_achievement_by_id"
	deleteachievementconditionbyachievementid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/delete_achievement_condition_by_achievement_id"
	deleteachievementconditionbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/delete_achievement_condition_by_id"
	existsachievementbycode "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/exists_achievement_by_code"
	existsachievementbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/exists_achievement_by_id"
	existsachievementconditionbyconditiontype "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/exists_achievement_condition_by_condition_type"
	existsachievementconditionbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/exists_achievement_condition_by_id"
	getdetailbyachievementid "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/get_detail_by_achievement_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	AllDetail                                 alldetail.IAllDetail
	CreateAchievement                         createachievement.ICreateAchievement
	CreateAchievementCondition                createachievementcondition.ICreateAchievementCondition
	DeleteAchievementByID                     deleteachievementbyid.IDeleteAchievementByID
	DeleteAchievementConditionByAchievementID deleteachievementconditionbyachievementid.IDeleteAchievementConditionByAchievementID
	DeleteAchievementConditionByID            deleteachievementconditionbyid.IDeleteAchievementConditionByID
	ExistsAchievementByCode                   existsachievementbycode.IExistsAchievementByCode
	ExistsAchievementByID                     existsachievementbyid.IExistsAchievementByID
	ExistsAchievementConditionByConditionType existsachievementconditionbyconditiontype.IExistsAchievementConditionByConditionType
	ExistsAchievementConditionByID            existsachievementconditionbyid.IExistsAchievementConditionByID
	GetDetailByAchievementID                  getdetailbyachievementid.IGetDetailByAchievementID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		AllDetail:                                 alldetail.New(queryTimeout, logger),
		CreateAchievement:                         createachievement.New(queryTimeout, logger),
		CreateAchievementCondition:                createachievementcondition.New(queryTimeout, logger),
		DeleteAchievementByID:                     deleteachievementbyid.New(queryTimeout, logger),
		DeleteAchievementConditionByAchievementID: deleteachievementconditionbyachievementid.New(queryTimeout, logger),
		DeleteAchievementConditionByID:            deleteachievementconditionbyid.New(queryTimeout, logger),
		ExistsAchievementByCode:                   existsachievementbycode.New(queryTimeout, logger),
		ExistsAchievementByID:                     existsachievementbyid.New(queryTimeout, logger),
		ExistsAchievementConditionByConditionType: existsachievementconditionbyconditiontype.New(queryTimeout, logger),
		ExistsAchievementConditionByID:            existsachievementconditionbyid.New(queryTimeout, logger),
		GetDetailByAchievementID:                  getdetailbyachievementid.New(queryTimeout, logger),
	}
}
