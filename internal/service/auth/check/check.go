package check

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ICheck --output=mocks --case=underscore
type ICheck interface {
	Execute(ctx context.Context, dto auth.CheckDTO) (auth.CheckResponse, error)
}

type Check struct {
	userRepository *userrepository.Repository
	logger         logger.ILogger
	postgres       *postgres.Postgres
	bigCache       *bigcachepkg.BigCache
	jwt            *jwt.JWT
}

func New(
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
	jwt *jwt.JWT,
) *Check {
	return &Check{
		userRepository: userRepository,
		logger:         logger,
		postgres:       postgres,
		bigCache:       bigCache,
		jwt:            jwt,
	}
}

func (s *Check) Execute(ctx context.Context, dto auth.CheckDTO) (auth.CheckResponse, error) {
	s.logger.Debug("[check user token] execute service")

	var err error

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return auth.CheckResponse{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check exists user from cache or database.
	ie, err := s.checkExistsUser(ctx, tx, dto.TelegramID)
	if err != nil {
		return auth.CheckResponse{}, err
	}

	if !ie {
		return auth.CheckResponse{}, apperrors.ErrUserDoesNotExist
	}

	// check verify token.
	vr, err := s.jwt.Verify(dto.TelegramID, dto.Token)
	if err != nil {
		return auth.CheckResponse{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return auth.CheckResponse{}, err
	}

	return auth.CheckResponse{
		TelegramID: vr.TelegramID,
		Token:      dto.Token,
		ExpAt:      vr.ExpAt,
	}, nil
}

// checkExistsUser checks whether a user exists either in the cache or the database.
// First, it attempts to find the user by Telegram ID in the cache.
// If not found (or if an error occurs other than "entry not found"), it queries the database using Telegram ID.
// Returns true if the user exists, otherwise false.
// Any unexpected error (e.g., cache failure or database error) will be returned.
func (s *Check) checkExistsUser(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
	// Check if the user exists in the cache by Telegram ID.
	// If found and no error occurred, return true immediately.
	ieFromCache, err := s.bigCache.User.Exists(telegramID, s.bigCache.User.GetPrefixTelegramID())
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
