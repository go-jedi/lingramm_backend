package notification

import (
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	allbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/notification/all_by_telegram_id"
	allpendingbeforebytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/notification/all_pending_before_by_telegram_id"
	allpendingbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/notification/all_pending_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/internal/service/v1/notification/create"
	updatestatus "github.com/go-jedi/lingramm_backend/internal/service/v1/notification/update_status"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/rabbitmq"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Service struct {
	AllByTelegramID              allbytelegramid.IAllByTelegramID
	AllPendingBeforeByTelegramID allpendingbeforebytelegramid.IAllPendingBeforeByTelegramID
	AllPendingByTelegramID       allpendingbytelegramid.IAllPendingByTelegramID
	Create                       create.ICreate
	UpdateStatus                 updatestatus.IUpdateStatus
}

func New(
	notificationRepository *notificationrepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	rabbitMQ *rabbitmq.RabbitMQ,
	postgres *postgres.Postgres,
	redis *redis.Redis,
) *Service {
	return &Service{
		AllByTelegramID:              allbytelegramid.New(notificationRepository, userRepository, logger, postgres),
		AllPendingBeforeByTelegramID: allpendingbeforebytelegramid.New(notificationRepository, userRepository, logger, postgres),
		AllPendingByTelegramID:       allpendingbytelegramid.New(notificationRepository, userRepository, logger, postgres),
		Create:                       create.New(notificationRepository, logger, rabbitMQ, postgres, redis),
		UpdateStatus:                 updatestatus.New(notificationRepository, logger, postgres),
	}
}
