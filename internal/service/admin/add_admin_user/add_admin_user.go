package addadminuser

import (
	"context"
	"fmt"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/admin"
	adminrepository "github.com/go-jedi/lingramm_backend/internal/repository/admin"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IAddAdminUser --output=mocks --case=underscore
type IAddAdminUser interface {
	Execute(ctx context.Context, telegramID string) (admin.Admin, error)
}

type AddAdminUser struct {
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
) *AddAdminUser {
	return &AddAdminUser{
		adminRepository: adminRepository,
		logger:          logger,
		postgres:        postgres,
		bigCache:        bigCache,
	}
}

func (s *AddAdminUser) Execute(ctx context.Context, telegramID string) (admin.Admin, error) {
	s.logger.Debug("[add a new admin user] execute service")

	var err error

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return admin.Admin{}, err
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
		return admin.Admin{}, err
	}

	if ie {
		return admin.Admin{}, apperrors.ErrAdminAlreadyExists
	}

	na, err := s.createAdmin(ctx, tx, telegramID)
	if err != nil {
		return admin.Admin{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return admin.Admin{}, err
	}

	return na, nil
}

// checkExistsAdmin checks whether admin exists either in the cache or the database.
// First, it attempts to find the admin by Telegram ID in the cache.
// If not found (or if an error occurs other than "entry not found"), it queries the database using Telegram ID.
// Returns true if the admin exists, otherwise false.
// Any unexpected error (e.g., cache failure or database error) will be returned.
func (s *AddAdminUser) checkExistsAdmin(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
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

// createAdmin creates a new admin in the database.
// After creation, the admin is cached using the Telegram ID as the key.
func (s *AddAdminUser) createAdmin(ctx context.Context, tx pgx.Tx, telegramID string) (admin.Admin, error) {
	// add the admin in the database.
	na, err := s.adminRepository.AddAdminUser.Execute(ctx, tx, telegramID)
	if err != nil {
		return admin.Admin{}, err
	}

	// save the newly created admin in the cache (prefix: telegram_id:).
	if err := s.bigCache.Admin.Set(na.TelegramID, na, s.bigCache.Admin.GetPrefixTelegramID()); err != nil {
		s.logger.Warn(fmt.Sprintf("failed to cache new admin: %v", err))
	}

	// Return the result from the database.
	return na, nil
}
