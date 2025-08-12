package allpendingbeforebytelegramid

import (
	"context"
	"log"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAllPendingBeforeByTelegramID --output=mocks --case=underscore
type IAllPendingBeforeByTelegramID interface {
	Execute(ctx context.Context, telegramID string, t0 time.Time) ([]notification.Notification, error)
}

type AllPendingBeforeByTelegramID struct {
	notificationRepository *notificationrepository.Repository
	userRepository         *userrepository.Repository
	logger                 logger.ILogger
	postgres               *postgres.Postgres
}

func New(
	notificationRepository *notificationrepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *AllPendingBeforeByTelegramID {
	return &AllPendingBeforeByTelegramID{
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
		logger:                 logger,
		postgres:               postgres,
	}
}

func (s *AllPendingBeforeByTelegramID) Execute(ctx context.Context, telegramID string, t0 time.Time) ([]notification.Notification, error) {
	s.logger.Debug("[get all pending before notifications by telegram id] execute service")

	var (
		err        error
		result     []notification.Notification
		userExists bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check user exists by telegram id.
	userExists, err = s.userRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return nil, err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return nil, err
	}

	// get all pending before notifications by telegram id.
	result, err = s.notificationRepository.AllPendingBeforeByTelegramID.Execute(ctx, tx, telegramID, t0)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
