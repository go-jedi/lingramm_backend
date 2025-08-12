package notification

import (
	allbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/all_by_telegram_id"
	allpendingbeforebytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/all_pending_before_by_telegram_id"
	allpendingbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/all_pending_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/create"
	existsbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/exists_by_id"
	updatestatus "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/update_status"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	AllByTelegramID              allbytelegramid.IAllByTelegramID
	AllPendingBeforeByTelegramID allpendingbeforebytelegramid.IAllPendingBeforeByTelegramID
	AllPendingByTelegramID       allpendingbytelegramid.IAllPendingByTelegramID
	Create                       create.ICreate
	ExistsByID                   existsbyid.IExistsByID
	UpdateStatus                 updatestatus.IUpdateStatus
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		AllByTelegramID:              allbytelegramid.New(queryTimeout, logger),
		AllPendingBeforeByTelegramID: allpendingbeforebytelegramid.New(queryTimeout, logger),
		AllPendingByTelegramID:       allpendingbytelegramid.New(queryTimeout, logger),
		Create:                       create.New(queryTimeout, logger),
		ExistsByID:                   existsbyid.New(queryTimeout, logger),
		UpdateStatus:                 updatestatus.New(queryTimeout, logger),
	}
}
