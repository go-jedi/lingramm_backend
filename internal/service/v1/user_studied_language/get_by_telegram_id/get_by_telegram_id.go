package getbytelegramid

import (
	"context"
	"log"

	userstudiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/user_studied_language"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstudiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetByTelegramID --output=mocks --case=underscore
type IGetByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (userstudiedlanguage.GetByTelegramIDResponse, error)
}

type GetByTelegramID struct {
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
) *GetByTelegramID {
	return &GetByTelegramID{
		userStudiedLanguageRepository: userStudiedLanguageRepository,
		userRepository:                userRepository,
		logger:                        logger,
		postgres:                      postgres,
	}
}

func (s *GetByTelegramID) Execute(ctx context.Context, telegramID string) (userstudiedlanguage.GetByTelegramIDResponse, error) {
	s.logger.Debug("[get user studied language by telegram id] execute service")

	var (
		err               error
		result            userstudiedlanguage.GetByTelegramIDResponse
		userExists        bool
		userStudiedExists bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return userstudiedlanguage.GetByTelegramIDResponse{}, err
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
		return userstudiedlanguage.GetByTelegramIDResponse{}, err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return userstudiedlanguage.GetByTelegramIDResponse{}, err
	}

	// check user studied language exists by telegram id.
	userStudiedExists, err = s.userStudiedLanguageRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return userstudiedlanguage.GetByTelegramIDResponse{}, err
	}

	if !userStudiedExists { // user studied language does not exist.
		err = apperrors.ErrUserStudiedLanguageDoesNotExist
		return userstudiedlanguage.GetByTelegramIDResponse{}, err
	}

	// get user studied language by telegram id.
	result, err = s.userStudiedLanguageRepository.GetByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return userstudiedlanguage.GetByTelegramIDResponse{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return userstudiedlanguage.GetByTelegramIDResponse{}, err
	}

	return result, nil
}
