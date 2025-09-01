package dependencies

import (
	dailytaskhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/daily_task"
	dailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/daily_task"
	dailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/daily_task"
)

func (d *Dependencies) DailyTaskRepository() *dailytaskrepository.Repository {
	if d.dailyTaskRepository == nil {
		d.dailyTaskRepository = dailytaskrepository.New(
			d.postgres.QueryTimeout,
			d.logger,
		)
	}

	return d.dailyTaskRepository
}

func (d *Dependencies) DailyTaskService() *dailytaskservice.Service {
	if d.dailyTaskService == nil {
		d.dailyTaskService = dailytaskservice.New(
			d.DailyTaskRepository(),
			d.logger,
			d.postgres,
		)
	}

	return d.dailyTaskService
}

func (d *Dependencies) DailyTaskHandler() *dailytaskhandler.Handler {
	if d.dailyTaskHandler == nil {
		d.dailyTaskHandler = dailytaskhandler.New(
			d.DailyTaskService(),
			d.app,
			d.logger,
			d.validator,
			d.middleware,
		)
	}

	return d.dailyTaskHandler
}
