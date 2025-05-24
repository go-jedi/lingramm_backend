package admin

import (
	addadminuser "github.com/go-jedi/lingramm_backend/internal/repository/admin/add_admin_user"
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/admin/exists_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	AddAdminUser       addadminuser.IAddAdminUser
	ExistsByTelegramID existsbytelegramid.IExistsByTelegramID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		AddAdminUser:       addadminuser.New(queryTimeout, logger),
		ExistsByTelegramID: existsbytelegramid.New(queryTimeout, logger),
	}
}
