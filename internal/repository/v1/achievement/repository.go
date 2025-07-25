package achievement

import (
	createachievement "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/create_achievement"
	createachievementcondition "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/create_achievement_condition"
	existsachievementbycode "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/exists_achievement_by_code"
	existsachievementconditionbyconditiontype "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement/exists_achievement_condition_by_condition_type"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	CreateAchievement                         createachievement.ICreateAchievement
	CreateAchievementCondition                createachievementcondition.ICreateAchievementCondition
	ExistsAchievementByCode                   existsachievementbycode.IExistsAchievementByCode
	ExistsAchievementConditionByConditionType existsachievementconditionbyconditiontype.IExistsAchievementConditionByConditionType
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		CreateAchievement:                         createachievement.New(queryTimeout, logger),
		CreateAchievementCondition:                createachievementcondition.New(queryTimeout, logger),
		ExistsAchievementByCode:                   existsachievementbycode.New(queryTimeout, logger),
		ExistsAchievementConditionByConditionType: existsachievementconditionbyconditiontype.New(queryTimeout, logger),
	}
}
