package dependencies

import (
	userdailytaskhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_daily_task"
	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	userdailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task"
)

func (d *Dependencies) UserDailyTaskRepository() *userdailytaskrepository.Repository {
	if d.userDailyTaskRepository == nil {
		d.userDailyTaskRepository = userdailytaskrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.userDailyTaskRepository
}

func (d *Dependencies) UserDailyTaskService() *userdailytaskservice.Service {
	if d.userDailyTaskService == nil {
		d.userDailyTaskService = userdailytaskservice.New(
			d.UserDailyTaskRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.userDailyTaskService
}

func (d *Dependencies) UserDailyTaskHandler() *userdailytaskhandler.Handler {
	if d.userDailyTaskHandler == nil {
		d.userDailyTaskHandler = userdailytaskhandler.New(
			d.UserDailyTaskService(),
			d.app,
			d.logger,
			d.middleware,
		)
	}

	return d.userDailyTaskHandler
}
