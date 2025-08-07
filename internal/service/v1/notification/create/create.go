package create

import (
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/nats"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
)

type Create struct {
	notificationRepository *notificationrepository.Repository
	logger                 logger.ILogger
	nats                   *nats.Nats
	postgres               *postgres.Postgres
}

func New(
	notificationRepository *notificationrepository.Repository,
	logger logger.ILogger,
	nats *nats.Nats,
	postgres *postgres.Postgres,
) *Create {
	return &Create{
		notificationRepository: notificationRepository,
		logger:                 logger,
		nats:                   nats,
		postgres:               postgres,
	}
}
