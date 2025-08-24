package getbytelegramid

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetByTelegramID --output=mocks --case=underscore
type IGetByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (user.User, error)
}

type GetByTelegramID struct {
	userRepository *userrepository.Repository
	logger         logger.ILogger
	postgres       *postgres.Postgres
}

func New(
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *GetByTelegramID {
	return &GetByTelegramID{
		userRepository: userRepository,
		logger:         logger,
		postgres:       postgres,
	}
}

func (s *GetByTelegramID) Execute(ctx context.Context, telegramID string) (user.User, error) {
	s.logger.Debug("[get user by telegram id] execute service")

	var (
		err        error
		result     user.User
		userExists bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return user.User{}, err
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
		return user.User{}, err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return user.User{}, err
	}

	// get user by telegram id.
	result, err = s.userRepository.GetByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return user.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return user.User{}, err
	}

	return result, nil
}
