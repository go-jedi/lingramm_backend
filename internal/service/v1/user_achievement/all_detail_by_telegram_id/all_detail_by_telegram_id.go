package alldetailbytelegramid

import (
	"context"
	"log"

	userachievement "github.com/go-jedi/lingramm_backend/internal/domain/user_achievement"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userachievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAllDetailByTelegramID --output=mocks --case=underscore
type IAllDetailByTelegramID interface {
	Execute(ctx context.Context, telegramID string) ([]userachievement.Detail, error)
}

type AllDetailByTelegramID struct {
	userAchievementRepository *userachievementrepository.Repository
	userRepository            *userrepository.Repository
	logger                    logger.ILogger
	postgres                  *postgres.Postgres
}

func New(
	userAchievementRepository *userachievementrepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *AllDetailByTelegramID {
	return &AllDetailByTelegramID{
		userAchievementRepository: userAchievementRepository,
		userRepository:            userRepository,
		logger:                    logger,
		postgres:                  postgres,
	}
}

func (s *AllDetailByTelegramID) Execute(ctx context.Context, telegramID string) ([]userachievement.Detail, error) {
	s.logger.Debug("[get all user achievements detail] execute service")

	var (
		err        error
		result     []userachievement.Detail
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

	// get detail user achievements.
	result, err = s.userAchievementRepository.AllDetailByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return nil, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
