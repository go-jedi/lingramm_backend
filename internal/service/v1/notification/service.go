package notification

import (
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/notification/create"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/nats"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Service struct {
	Create create.ICreate
}

func New(
	notificationRepository *notificationrepository.Repository,
	natsTimeout int64,
	logger logger.ILogger,
	nats *nats.Nats,
	postgres *postgres.Postgres,
) *Service {
	return &Service{
		Create: create.New(notificationRepository, natsTimeout, logger, nats, postgres),
	}
}
