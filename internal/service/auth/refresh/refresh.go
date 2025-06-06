package refresh

import (
	"context"
	"log"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/user"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IRefresh --output=mocks --case=underscore
type IRefresh interface {
	Execute(ctx context.Context, dto auth.RefreshDTO) (auth.RefreshResponse, error)
}

type Refresh struct {
	userRepository *userrepository.Repository
	logger         logger.ILogger
	postgres       *postgres.Postgres
	redis          *redis.Redis
	bigCache       *bigcachepkg.BigCache
	jwt            jwt.IJWT
}

func New(
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	bigCache *bigcachepkg.BigCache,
	jwt jwt.IJWT,
) *Refresh {
	return &Refresh{
		userRepository: userRepository,
		logger:         logger,
		postgres:       postgres,
		redis:          redis,
		bigCache:       bigCache,
		jwt:            jwt,
	}
}

func (s *Refresh) Execute(ctx context.Context, dto auth.RefreshDTO) (auth.RefreshResponse, error) {
	s.logger.Debug("[refresh user token] execute service")

	var err error

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return auth.RefreshResponse{}, err
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
		return auth.RefreshResponse{}, err
	}

	if !ie {
		return auth.RefreshResponse{}, apperrors.ErrUserDoesNotExist
	}

	// check verify token.
	vr, err := s.jwt.Verify(dto.TelegramID, dto.RefreshToken)
	if err != nil {
		return auth.RefreshResponse{}, err
	}

	// get refresh token from cache.
	err = s.validateRefreshTokenFromCache(ctx, dto.TelegramID, dto.RefreshToken)
	if err != nil {
		return auth.RefreshResponse{}, err
	}

	// generate access, refresh tokens.
	tokens, err := s.jwt.Generate(vr.TelegramID)
	if err != nil {
		return auth.RefreshResponse{}, err
	}

	// set new refresh token in cache.
	err = s.redis.RefreshToken.SetWithExpiration(ctx, dto.TelegramID, tokens.RefreshToken, time.Until(tokens.RefreshExpAt))
	if err != nil {
		return auth.RefreshResponse{}, err
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return auth.RefreshResponse{}, err
	}

	return auth.RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		AccessExpAt:  tokens.AccessExpAt,
		RefreshExpAt: tokens.RefreshExpAt,
	}, nil
}

// checkExistsUser checks whether a user exists either in the cache or the database.
// First, it attempts to find the user by Telegram ID in the cache.
// If not found (or if an error occurs other than "entry not found"), it queries the database using Telegram ID.
// Returns true if the user exists, otherwise false.
// Any unexpected error (e.g., cache failure or database error) will be returned.
func (s *Refresh) checkExistsUser(ctx context.Context, tx pgx.Tx, telegramID string) (bool, error) {
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

// validateRefreshTokenFromCache checks whether the provided refresh token matches.
// the one stored in Redis for the given Telegram ID.
func (s *Refresh) validateRefreshTokenFromCache(ctx context.Context, telegramID string, token string) error {
	refreshTokenFromCache, err := s.redis.RefreshToken.Get(ctx, telegramID)
	if err != nil {
		return err // redis access error.
	}
	if refreshTokenFromCache == "" {
		return apperrors.ErrNoActiveSessionFound // no token stored — session is invalid.
	}
	if refreshTokenFromCache != token {
		return apperrors.ErrTokenMismatchOrExpired // tokens do not match.
	}

	return nil // token is valid.
}
