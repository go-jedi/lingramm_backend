package signin

import (
	"context"
	"fmt"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/user"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ISignIn --output=mocks --case=underscore
type ISignIn interface {
	Execute(ctx context.Context, dto auth.SignInDTO) (user.User, error)
}

type SignIn struct {
	userRepository *userrepository.Repository
	logger         logger.ILogger
	postgres       *postgres.Postgres
	bigCache       *bigcachepkg.BigCache
	uuid           uuid.IUUID
}

func New(
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	bigCache *bigcachepkg.BigCache,
	uuid uuid.IUUID,
) *SignIn {
	return &SignIn{
		userRepository: userRepository,
		logger:         logger,
		postgres:       postgres,
		bigCache:       bigCache,
		uuid:           uuid,
	}
}

func (s *SignIn) Execute(ctx context.Context, dto auth.SignInDTO) (user.User, error) {
	s.logger.Debug("[sign in user] execute service")

	var (
		err error
		u   user.User
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return user.User{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	ie, err := s.checkExistsUser(ctx, tx, dto.TelegramID, dto.Username)
	if err != nil {
		return user.User{}, err
	}

	if ie {
		u, err = s.findOrReturnExisting(ctx, tx, dto.TelegramID)
	} else {
		u, err = s.createUser(ctx, tx, dto)
	}
	if err != nil {
		return user.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// checkExistsUser checks whether a user exists either in the cache or the database.
// First, it attempts to find the user by Telegram ID in the cache.
// If not found (or if an error occurs other than "entry not found"), it queries the database using Telegram ID and username.
// Returns true if the user exists, otherwise false.
// Any unexpected error (e.g., cache failure or database error) will be returned.
func (s *SignIn) checkExistsUser(ctx context.Context, tx pgx.Tx, telegramID string, username string) (bool, error) {
	// Check if the user exists in the cache by Telegram ID.
	// If found and no error occurred, return true immediately.
	ieFromCache, err := s.bigCache.User.Exists(telegramID, s.bigCache.User.GetPrefixTelegramID())
	if err == nil && ieFromCache {
		return true, nil
	}

	// If the user is not found in the cache (or an error occurred),
	// query the database to check if the user exists.
	ieFromDB, err := s.userRepository.Exists.Execute(ctx, tx, telegramID, username)
	if err != nil {
		return false, err
	}

	// Return the result from the database.
	return ieFromDB, nil
}

// createUser generates a UUID and creates a new user in the database.
// After creation, the user is cached using the Telegram ID as the key.
func (s *SignIn) createUser(ctx context.Context, tx pgx.Tx, dto auth.SignInDTO) (user.User, error) {
	// generate a unique UUID for the new user.
	newUUID, err := s.uuid.Generate()
	if err != nil {
		return user.User{}, err
	}

	createDTO := user.CreateDTO{
		UUID:       newUUID,
		TelegramID: dto.TelegramID,
		Username:   dto.Username,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
	}

	// create the user in the database.
	u, err := s.userRepository.Create.Execute(ctx, tx, createDTO)
	if err != nil {
		return user.User{}, err
	}

	// save the newly created user in the cache (prefix: telegram_id:).
	if err := s.bigCache.User.Set(u.TelegramID, u, s.bigCache.User.GetPrefixTelegramID()); err != nil {
		s.logger.Warn(fmt.Sprintf("failed to cache new user: %v", err))
	}

	// save the newly created user in the cache (prefix: uuid:).
	if err := s.bigCache.User.Set(u.TelegramID, u, s.bigCache.User.GetPrefixUUID()); err != nil {
		s.logger.Warn(fmt.Sprintf("failed to cache new user: %v", err))
	}

	return u, nil
}

// findOrReturnExisting attempts to retrieve a user from the cache by Telegram ID.
// If the user is found in the cache and the data is valid, it returns the cached user.
// Otherwise, it queries the database to retrieve the user by Telegram ID.
func (s *SignIn) findOrReturnExisting(ctx context.Context, tx pgx.Tx, telegramID string) (user.User, error) {
	// Try to get the user from the cache.
	userFromCache, err := s.bigCache.User.Get(telegramID, s.bigCache.User.GetPrefixTelegramID())
	if err == nil && userFromCache.TelegramID == telegramID {
		// Cache hit and valid data — return the cached user.
		return userFromCache, nil
	}

	// Cache miss or invalid data — fallback to the database.
	return s.userRepository.GetByTelegramID.Execute(ctx, tx, telegramID)
}
