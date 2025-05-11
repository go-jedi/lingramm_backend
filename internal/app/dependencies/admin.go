package dependencies

import (
	adminhandler "github.com/go-jedi/lingvogramm_backend/internal/adapter/http/handlers/admin"
	adminrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/admin"
	adminservice "github.com/go-jedi/lingvogramm_backend/internal/service/admin"
)

func (d *Dependencies) AdminRepository() *adminrepository.Repository {
	if d.adminRepository == nil {
		d.adminRepository = adminrepository.New(d.postgres.QueryTimeout, d.logger)
	}

	return d.adminRepository
}

func (d *Dependencies) AdminService() *adminservice.Service {
	if d.adminService == nil {
		d.adminService = adminservice.New(
			d.AdminRepository(),
			d.logger,
			d.postgres,
			d.bigCache,
		)
	}

	return d.adminService
}

func (d *Dependencies) AdminHandler() *adminhandler.Handler {
	if d.adminHandler == nil {
		d.adminHandler = adminhandler.New(
			d.AdminService(),
			d.app,
			d.logger,
		)
	}

	return d.adminHandler
}
