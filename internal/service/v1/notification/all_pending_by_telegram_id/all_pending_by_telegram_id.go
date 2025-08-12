package allpendingbytelegramid

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAllPendingByTelegramID --output=mocks --case=underscore
type IAllPendingByTelegramID interface {
	Execute(ctx context.Context, telegramID string) ([]notification.Notification, error)
}

type AllPendingByTelegramID struct {
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
) *AllPendingByTelegramID {
	return &AllPendingByTelegramID{
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
		logger:                 logger,
		postgres:               postgres,
	}
}

func (s *AllPendingByTelegramID) Execute(ctx context.Context, telegramID string) ([]notification.Notification, error) {
	s.logger.Debug("[get all pending notifications by telegram id] execute service")

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

	// get all pending notifications by telegram id.
	result, err = s.notificationRepository.AllPendingByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
