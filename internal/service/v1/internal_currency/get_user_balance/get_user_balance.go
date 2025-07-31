package getuserbalance

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/internal_currency/user_balance"
	"github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IGetUserBalance --output=mocks --case=underscore
type IGetUserBalance interface {
	Execute(ctx context.Context, telegramID string) (userbalance.UserBalance, error)
}

type GetUserBalance struct {
	internalCurrencyRepository *internalcurrency.Repository
	userRepository             *userrepository.Repository
	logger                     logger.ILogger
	postgres                   *postgres.Postgres
	bigCache                   *bigcachepkg.BigCache
}

func New(
	internalCurrencyRepository *internalcurrency.Repository,
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
) *GetUserBalance {
	return &GetUserBalance{
		internalCurrencyRepository: internalCurrencyRepository,
		userRepository:             userRepository,
		logger:                     logger,
		postgres:                   postgres,
		bigCache:                   bigCache,
	}
}

func (s *GetUserBalance) Execute(ctx context.Context, telegramID string) (userbalance.UserBalance, error) {
	s.logger.Debug("[get user balance] execute service")

	var err error

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return userbalance.UserBalance{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	ie, err := s.checkExistsUser(ctx, tx, telegramID)
	if err != nil {
		return userbalance.UserBalance{}, err
	}

	if !ie {
		return userbalance.UserBalance{}, apperrors.ErrUserDoesNotExist
	}

	result, err := s.internalCurrencyRepository.GetUserBalance.Execute(ctx, tx, telegramID)
	if err != nil {
		return userbalance.UserBalance{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return userbalance.UserBalance{}, err
	}

	return result, nil
}

// checkExistsUser checks whether a user exists either in the cache or the database.
// First, it attempts to find the user by Telegram ID in the cache.
// If not found (or if an error occurs other than "entry not found"), it queries the database using Telegram ID.
// Returns true if the user exists, otherwise false.
// Any unexpected error (e.g., cache failure or database error) will be returned.
func (s *GetUserBalance) checkExistsUser(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
	// Check if the user exists in the cache by Telegram ID.
	// If found and no error occurred, return true immediately.
	ieFromCache, err := s.bigCache.User.Exists(telegramID)
	if err == nil && ieFromCache {
		return true, nil
	}

	// If the user is not found in the cache (or an error occurred),
	// query the database to check if the user exists.
	ieFromDB, err := s.userRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return false, err
	}

	// Return the result from the database.
	return ieFromDB, nil
}
