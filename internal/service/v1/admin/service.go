package admin

import (
	adminrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/admin"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/admin/add_admin_user"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/admin/exists_by_telegram_id"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	AddAdminUser       addadminuser.IAddAdminUser
	ExistsByTelegramID existsbytelegramid.IExistsByTelegramID
}

func New(
	adminRepository *adminrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
) *Service {
	return &Service{
		AddAdminUser:       addadminuser.New(adminRepository, logger, postgres, bigCache),
		ExistsByTelegramID: existsbytelegramid.New(adminRepository, logger, postgres, bigCache),
	}
}
