package userdailytask

import (
	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	assigndailytaskbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task/assign_daily_task_by_telegram_id"
	existsassigndailytaskbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task/exists_assign_daily_task_by_telegram_id"
	getcurrentdailytaskbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task/get_current_daily_task_by_telegram_id"
	getdailytaskweeksummarybytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task/get_daily_task_week_summary_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	AssignDailyTaskByTelegramID         assigndailytaskbytelegramid.IAssignDailyTaskByTelegramID
	ExistsAssignDailyTaskByTelegramID   existsassigndailytaskbytelegramid.IExistsAssignDailyTaskByTelegramID
	GetCurrentDailyTaskByTelegramID     getcurrentdailytaskbytelegramid.IGetCurrentDailyTaskByTelegramID
	GetDailyTaskWeekSummaryByTelegramID getdailytaskweeksummarybytelegramid.IGetDailyTaskWeekSummaryByTelegramID
}

func New(
	userDailyTaskRepository *userdailytaskrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		AssignDailyTaskByTelegramID:         assigndailytaskbytelegramid.New(userDailyTaskRepository, logger, postgres),
		ExistsAssignDailyTaskByTelegramID:   existsassigndailytaskbytelegramid.New(userDailyTaskRepository, logger, postgres),
		GetCurrentDailyTaskByTelegramID:     getcurrentdailytaskbytelegramid.New(userDailyTaskRepository, logger, postgres),
		GetDailyTaskWeekSummaryByTelegramID: getdailytaskweeksummarybytelegramid.New(userDailyTaskRepository, logger, postgres),
	}
}
