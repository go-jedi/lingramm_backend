package create

import (
	"context"
	"log"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/nats"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, dto notification.CreateDTO) (notification.Notification, error)
}

type Create struct {
	notificationRepository *notificationrepository.Repository
	natsTimeout            int64
	logger                 logger.ILogger
	nats                   *nats.Nats
	postgres               *postgres.Postgres
}

func New(
	notificationRepository *notificationrepository.Repository,
	natsTimeout int64,
	logger logger.ILogger,
	nats *nats.Nats,
	postgres *postgres.Postgres,
) *Create {
	return &Create{
		notificationRepository: notificationRepository,
		natsTimeout:            natsTimeout,
		logger:                 logger,
		nats:                   nats,
		postgres:               postgres,
	}
}

func (s *Create) Execute(ctx context.Context, dto notification.CreateDTO) (notification.Notification, error) {
	s.logger.Debug("[create a new notification] execute service")

	var (
		err    error
		result notification.Notification
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

	// send notification user by nats.
	err = s.sendNotification(ctx, result)
	if err != nil {
		return notification.Notification{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return notification.Notification{}, err
	}

	return result, nil
}

// sendNotification send notification.
func (s *Create) sendNotification(ctx context.Context, n notification.Notification) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(s.natsTimeout)*time.Second)
	defer cancel()

	data := notification.SendNotificationDTO{
		ID:         n.ID,
		Message:    n.Message,
		Type:       n.Type,
		TelegramID: n.TelegramID,
		CreatedAt:  n.CreatedAt,
	}

	if err := s.nats.Notification.Publisher.Start(ctxTimeout, data.TelegramID, data); err != nil {
		log.Printf("failed to publish notification: %v", err)
		return err
	}

	return nil
}
