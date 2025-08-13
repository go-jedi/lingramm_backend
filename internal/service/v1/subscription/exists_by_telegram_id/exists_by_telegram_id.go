package existsbytelegramid

import (
	"context"
	"log"

	subscriptionrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsByTelegramID --output=mocks --case=underscore
type IExistsByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (bool, error)
}

type ExistsByTelegramID struct {
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
) *ExistsByTelegramID {
	return &ExistsByTelegramID{
		subscriptionRepository: subscriptionRepository,
		userRepository:         userRepository,
		logger:                 logger,
		postgres:               postgres,
	}
}

func (s *ExistsByTelegramID) Execute(ctx context.Context, telegramID string) (bool, error) {
	s.logger.Debug("[check subscription exists by telegram id] execute service")

	var (
		err                error
		userExists         bool
		subscriptionExists bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return false, err
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
		return false, err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return false, err
	}

	// check subscription exists by telegram id.
	subscriptionExists, err = s.subscriptionRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return subscriptionExists, nil
}
