package notification

import (
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/notification/create"
	updatestatus "github.com/go-jedi/lingramm_backend/internal/service/v1/notification/update_status"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	Create       create.ICreate
	UpdateStatus updatestatus.IUpdateStatus
}

func New(
	notificationRepository *notificationrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		Create:       create.New(notificationRepository, logger, postgres),
		UpdateStatus: updatestatus.New(notificationRepository, logger, postgres),
	}
}
