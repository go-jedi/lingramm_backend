package create

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

//go:generate mockery --name=ICreate --output=mocks --case=underscore
type ICreate interface {
	Execute(ctx context.Context, dto userstudiedlanguage.CreateDTO) (userstudiedlanguage.UserStudiedLanguage, error)
}

type Create struct {
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
) *Create {
	return &Create{
		userStudiedLanguageRepository: userStudiedLanguageRepository,
		userRepository:                userRepository,
		logger:                        logger,
		postgres:                      postgres,
	}
}

func (s *Create) Execute(ctx context.Context, dto userstudiedlanguage.CreateDTO) (userstudiedlanguage.UserStudiedLanguage, error) {
	s.logger.Debug("[create a new user studied language] execute service")

	var (
		err               error
		result            userstudiedlanguage.UserStudiedLanguage
		userExists        bool
		userStudiedExists bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return userstudiedlanguage.UserStudiedLanguage{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check user exists by telegram id.
	userExists, err = s.userRepository.ExistsByTelegramID.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return userstudiedlanguage.UserStudiedLanguage{}, err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return userstudiedlanguage.UserStudiedLanguage{}, err
	}

	// check user studied language exists by telegram id.
	userStudiedExists, err = s.userStudiedLanguageRepository.ExistsByTelegramID.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return userstudiedlanguage.UserStudiedLanguage{}, err
	}

	if userStudiedExists { // if user studied language already exist.
		err = apperrors.ErrUserStudiedLanguageAlreadyExists
		return userstudiedlanguage.UserStudiedLanguage{}, err
	}

	// create new user studied language.
	result, err = s.userStudiedLanguageRepository.Create.Execute(ctx, tx, dto)
	if err != nil {
		return userstudiedlanguage.UserStudiedLanguage{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return userstudiedlanguage.UserStudiedLanguage{}, err
	}

	return result, nil
}
