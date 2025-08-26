package existsbytelegramid

import (
	"context"
	"log"

	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstudiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language"
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
	userStudiedLanguageRepository *userstudiedlanguagerepository.Repository
	userRepository                *userrepository.Repository
	logger                        logger.ILogger
	postgres                      *postgres.Postgres
}

func New(
	userStudiedLanguageRepository *userstudiedlanguagerepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *ExistsByTelegramID {
	return &ExistsByTelegramID{
		userStudiedLanguageRepository: userStudiedLanguageRepository,
		userRepository:                userRepository,
		logger:                        logger,
		postgres:                      postgres,
	}
}

func (s *ExistsByTelegramID) Execute(ctx context.Context, telegramID string) (bool, error) {
	s.logger.Debug("[check user studied language exists by telegram id] execute service")

	var (
		err        error
		userExists bool
		result     bool
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

	// check user studied language exists by telegram id.
	result, err = s.userStudiedLanguageRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return false, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return result, nil
}
