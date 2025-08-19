package getlevelinfobytelegramid

import (
	"context"
	"log"

	userstats "github.com/go-jedi/lingramm_backend/internal/domain/user_stats"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetLevelInfoByTelegramID --output=mocks --case=underscore
type IGetLevelInfoByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (userstats.GetLevelInfoByTelegramIDResponse, error)
}

type GetLevelInfoByTelegramID struct {
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
) *GetLevelInfoByTelegramID {
	return &GetLevelInfoByTelegramID{
		userStatsRepository: userStatsRepository,
		userRepository:      userRepository,
		logger:              logger,
		postgres:            postgres,
	}
}

func (s *GetLevelInfoByTelegramID) Execute(ctx context.Context, telegramID string) (userstats.GetLevelInfoByTelegramIDResponse, error) {
	s.logger.Debug("[get level info by telegram id] execute service")

	var (
		err             error
		result          userstats.GetLevelInfoByTelegramIDResponse
		userExists      bool
		userStatsExists bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return userstats.GetLevelInfoByTelegramIDResponse{}, err
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
		return userstats.GetLevelInfoByTelegramIDResponse{}, err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return userstats.GetLevelInfoByTelegramIDResponse{}, err
	}

	// check user stats exists by telegram id.
	userStatsExists, err = s.userStatsRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return userstats.GetLevelInfoByTelegramIDResponse{}, err
	}

	if !userStatsExists { // if user stats does not exist.
		err = apperrors.ErrUserStatsDoesNotExist
		return userstats.GetLevelInfoByTelegramIDResponse{}, err
	}

	// get user level info by telegram id.
	result, err = s.userStatsRepository.GetLevelInfoByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return userstats.GetLevelInfoByTelegramIDResponse{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return userstats.GetLevelInfoByTelegramIDResponse{}, err
	}

	return result, nil
}
