package getbytelegramid

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/subscription"
	subscriptionrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetByTelegramID --output=mocks --case=underscore
type IGetByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (subscription.Subscription, error)
}

type GetByTelegramID struct {
	subscriptionRepository *subscriptionrepository.Repository
	userRepository         *userrepository.Repository
	logger                 logger.ILogger
	postgres               *postgres.Postgres
}

func New(
	subscriptionRepository *subscriptionrepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *GetByTelegramID {
	return &GetByTelegramID{
		subscriptionRepository: subscriptionRepository,
		userRepository:         userRepository,
		logger:                 logger,
		postgres:               postgres,
	}
}

func (s *GetByTelegramID) Execute(ctx context.Context, telegramID string) (subscription.Subscription, error) {
	s.logger.Debug("[get subscription by telegram id] execute service")

	var (
		err                error
		result             subscription.Subscription
		userExists         bool
		subscriptionExists bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return subscription.Subscription{}, err
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
		return subscription.Subscription{}, err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return subscription.Subscription{}, err
	}

	// check subscription exists by telegram id.
	subscriptionExists, err = s.subscriptionRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return subscription.Subscription{}, err
	}

	if !subscriptionExists { // if subscription does not exist.
		err = apperrors.ErrSubscriptionDoesNotExist
		return subscription.Subscription{}, err
	}

	//  get subscription by telegram id.
	result, err = s.subscriptionRepository.GetByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return subscription.Subscription{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return subscription.Subscription{}, err
	}

	return result, nil
}
