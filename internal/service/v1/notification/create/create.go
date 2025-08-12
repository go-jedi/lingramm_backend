package create

import (
	"context"
	"fmt"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/rabbitmq"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, dto notification.CreateDTO) (notification.Notification, error)
}

type Create struct {
	notificationRepository *notificationrepository.Repository
	logger                 logger.ILogger
	rabbitMQ               *rabbitmq.RabbitMQ
	postgres               *postgres.Postgres
	redis                  *redis.Redis
}

func New(
	notificationRepository *notificationrepository.Repository,
	logger logger.ILogger,
	rabbitMQ *rabbitmq.RabbitMQ,
	postgres *postgres.Postgres,
	redis *redis.Redis,
) *Create {
	return &Create{
		notificationRepository: notificationRepository,
		logger:                 logger,
		rabbitMQ:               rabbitMQ,
		postgres:               postgres,
		redis:                  redis,
	}
}

func (s *Create) Execute(ctx context.Context, dto notification.CreateDTO) (notification.Notification, error) {
	s.logger.Debug("[create a new notification] execute service")

	var (
		err            error
		result         notification.Notification
		isUserPresence bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return notification.Notification{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// create notification.
	result, err = s.notificationRepository.Create.Execute(ctx, tx, dto)
	if err != nil {
		return notification.Notification{}, err
	}

	// check exists user is online for send notification with message broker.
	isUserPresence, err = s.redis.UserPresence.Exists(ctx, dto.TelegramID)
	if err != nil {
		return notification.Notification{}, err
	}

	if isUserPresence { // if the user is online.
		// send notification user by rabbitmq.
		s.sendNotification(ctx, result)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return notification.Notification{}, err
	}

	return result, nil
}

// sendNotification send notification.
func (s *Create) sendNotification(ctx context.Context, n notification.Notification) {
	data := notification.SendNotificationDTO{
		ID:         n.ID,
		Message:    n.Message,
		Type:       n.Type,
		TelegramID: n.TelegramID,
		CreatedAt:  n.CreatedAt,
	}

	// send notification in rabbitmq.
	if err := s.rabbitMQ.Notification.Publisher.Execute(ctx, data.TelegramID, data); err != nil {
		s.logger.Warn(fmt.Sprintf("failed to publish notification by rabbitmq: %v", err))
	}
}
