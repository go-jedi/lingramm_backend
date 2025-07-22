package existsbytelegramid

import (
	"context"
	"log"

	adminrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/admin"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IExistsByTelegramID --output=mocks --case=underscore
type IExistsByTelegramID interface {
	Execute(ctx context.Context, telegramID string) (bool, error)
}

type ExistsByTelegramID struct {
	adminRepository *adminrepository.Repository
	logger          logger.ILogger
	postgres        *postgres.Postgres
	bigCache        *bigcachepkg.BigCache
}

func New(
	adminRepository *adminrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
) *ExistsByTelegramID {
	return &ExistsByTelegramID{
		adminRepository: adminRepository,
		logger:          logger,
		postgres:        postgres,
		bigCache:        bigCache,
	}
}

func (s *ExistsByTelegramID) Execute(ctx context.Context, telegramID string) (bool, error) {
	s.logger.Debug("[check admin exists by telegram id] execute service")

	var err error

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

	ie, err := s.checkExistsAdmin(ctx, tx, telegramID)
	if err != nil {
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return ie, nil
}

// checkExistsAdmin checks whether admin exists either in the cache or the database.
// First, it attempts to find the admin by Telegram ID in the cache.
// If not found (or if an error occurs other than "entry not found"), it queries the database using Telegram ID.
// Returns true if the admin exists, otherwise false.
// Any unexpected error (e.g., cache failure or database error) will be returned.
func (s *ExistsByTelegramID) checkExistsAdmin(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
	// Check if the admin exists in the cache by Telegram ID.
	// If found and no error occurred, return true immediately.
	ieFromCache, err := s.bigCache.Admin.Exists(telegramID, s.bigCache.Admin.GetPrefixTelegramID())
	if err == nil && ieFromCache {
		return true, nil
	}

	// If the admin is not found in the cache (or an error occurred),
	// query the database to check if the admin exists.
	ieFromDB, err := s.adminRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return false, err
	}

	// Return the result from the database.
	return ieFromDB, nil
}
