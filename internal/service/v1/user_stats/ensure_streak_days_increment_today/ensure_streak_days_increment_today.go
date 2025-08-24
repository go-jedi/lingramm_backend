package ensurestreakdaysincrementtoday

import (
	"context"
	"log"

	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IEnsureStreakDaysIncrementToday --output=mocks --case=underscore
type IEnsureStreakDaysIncrementToday interface {
	Execute(ctx context.Context, telegramID string) error
}

type EnsureStreakDaysIncrementToday struct {
	userStatsRepository *userstatsrepository.Repository
	userRepository      *userrepository.Repository
	logger              logger.ILogger
	postgres            *postgres.Postgres
}

func New(
	userStatsRepository *userstatsrepository.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *EnsureStreakDaysIncrementToday {
	return &EnsureStreakDaysIncrementToday{
		userStatsRepository: userStatsRepository,
		userRepository:      userRepository,
		logger:              logger,
		postgres:            postgres,
	}
}

func (s *EnsureStreakDaysIncrementToday) Execute(ctx context.Context, telegramID string) error {
	s.logger.Debug("[ensure streak days increment today] execute service")

	var (
		err                        error
		userExists                 bool
		userStatsExists            bool
		isStreakDaysIncrementToday bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
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
		return err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return err
	}

	// check user stats exists by telegram id.
	userStatsExists, err = s.userStatsRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return err
	}

	if !userStatsExists { // if user stats does not exist.
		err = apperrors.ErrUserStatsDoesNotExist
		return err
	}

	// check has streak days increment today by telegram id.
	isStreakDaysIncrementToday, err = s.userStatsRepository.HasStreakDaysIncrementToday.Execute(ctx, tx, telegramID)
	if err != nil {
		return err
	}

	if !isStreakDaysIncrementToday { // has streak days is not increment today.
		// ensure streak days increment today.
		err = s.userStatsRepository.EnsureStreakDaysIncrementToday.Execute(ctx, tx, telegramID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
