package createxpevents

import (
	"context"
	"log"

	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICreateXPEvents --output=mocks --case=underscore
type ICreateXPEvents interface {
	Execute(ctx context.Context, dto experiencepoint.CreateXPEventDTO) error
}

type CreateXPEvents struct {
	experiencePointRepository *experiencepointrepository.Repository
	userRepository            *userrepository.Repository
	userStatsRepository       *userstatsrepository.Repository
	levelRepository           *levelrepository.Repository
	logger                    logger.ILogger
	postgres                  *postgres.Postgres
}

func New(
	experiencePointRepository *experiencepointrepository.Repository,
	userRepository *userrepository.Repository,
	userStatsRepository *userstatsrepository.Repository,
	levelRepository *levelrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
) *CreateXPEvents {
	return &CreateXPEvents{
		experiencePointRepository: experiencePointRepository,
		userRepository:            userRepository,
		userStatsRepository:       userStatsRepository,
		levelRepository:           levelRepository,
		logger:                    logger,
		postgres:                  postgres,
	}
}

func (s *CreateXPEvents) Execute(ctx context.Context, dto experiencepoint.CreateXPEventDTO) error {
	s.logger.Debug("[create a new xp events] execute service")

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
	userExists, err = s.userRepository.ExistsByTelegramID.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return err
	}

	// check user stats exists by telegram id.
	userStatsExists, err = s.userStatsRepository.ExistsByTelegramID.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	if !userStatsExists { // if user stats does not exist.
		err = apperrors.ErrUserStatsDoesNotExist
		return err
	}

	// create a new xp events.
	err = s.experiencePointRepository.CreateXPEvents.Execute(ctx, tx, dto)
	if err != nil {
		return err
	}

	// sync user stats from xp events by telegram id.
	err = s.userStatsRepository.SyncUserStatsFromXPEventsByTelegramID.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	// backfill missing level history by telegram id.
	err = s.levelRepository.BackFillMissingLevelHistoryByTelegramID.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	// check has streak days increment today by telegram id.
	isStreakDaysIncrementToday, err = s.userStatsRepository.HasStreakDaysIncrementToday.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	if !isStreakDaysIncrementToday { // has streak days is not increment today.
		// ensure streak days increment today.
		err = s.userStatsRepository.EnsureStreakDaysIncrementToday.Execute(ctx, tx, dto.TelegramID)
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
