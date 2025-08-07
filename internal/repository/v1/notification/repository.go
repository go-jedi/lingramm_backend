package notification

import (
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/create"
	existsbyid "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/exists_by_id"
	updatestatus "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification/update_status"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	Create       create.ICreate
	ExistsByID   existsbyid.IExistsByID
	UpdateStatus updatestatus.IUpdateStatus
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		Create:       create.New(queryTimeout, logger),
		ExistsByID:   existsbyid.New(queryTimeout, logger),
		UpdateStatus: updatestatus.New(queryTimeout, logger),
	}
}
