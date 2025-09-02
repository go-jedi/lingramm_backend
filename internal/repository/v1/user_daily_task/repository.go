package userdailytask

import (
	assigndailytaskbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task/assign_daily_task_by_telegram_id"
	existsassigndailytaskbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task/exists_assign_daily_task_by_telegram_id"
	getcurrentdailytaskbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task/get_current_daily_task_by_telegram_id"
	syncuserdailytaskprogress "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task/sync_user_daily_task_progress"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	AssignDailyTaskByTelegramID       assigndailytaskbytelegramid.IAssignDailyTaskByTelegramID
	ExistsAssignDailyTaskByTelegramID existsassigndailytaskbytelegramid.IExistsAssignDailyTaskByTelegramID
	GetCurrentDailyTaskByTelegramID   getcurrentdailytaskbytelegramid.IGetCurrentDailyTaskByTelegramID
	SyncUserDailyTaskProgress         syncuserdailytaskprogress.ISyncUserDailyTaskProgress
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		AssignDailyTaskByTelegramID:       assigndailytaskbytelegramid.New(queryTimeout, logger),
		ExistsAssignDailyTaskByTelegramID: existsassigndailytaskbytelegramid.New(queryTimeout, logger),
		GetCurrentDailyTaskByTelegramID:   getcurrentdailytaskbytelegramid.New(queryTimeout, logger),
		SyncUserDailyTaskProgress:         syncuserdailytaskprogress.New(queryTimeout, logger),
	}
}
