package signin

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	"github.com/go-jedi/lingramm_backend/internal/domain/level"
	"github.com/go-jedi/lingramm_backend/internal/domain/user"
	levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ISignIn --output=mocks --case=underscore
type ISignIn interface {
	Execute(ctx context.Context, dto auth.SignInDTO) (auth.SignInResp, error)
}

type SignIn struct {
	userRepository          *userrepository.Repository
	levelRepository         *levelrepository.Repository
	userDailyTaskRepository *userdailytaskrepository.Repository
	logger                  logger.ILogger
	postgres                *postgres.Postgres
	redis                   *redis.Redis
	bigCache                *bigcachepkg.BigCache
	jwt                     jwt.IJWT
}

func New(
	userRepository *userrepository.Repository,
	levelRepository *levelrepository.Repository,
	userDailyTaskRepository *userdailytaskrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	bigCache *bigcachepkg.BigCache,
	jwt jwt.IJWT,
) *SignIn {
	return &SignIn{
		userRepository:          userRepository,
		levelRepository:         levelRepository,
		userDailyTaskRepository: userDailyTaskRepository,
		logger:                  logger,
		postgres:                postgres,
		redis:                   redis,
		bigCache:                bigCache,
		jwt:                     jwt,
	}
}

func (s *SignIn) Execute(ctx context.Context, dto auth.SignInDTO) (auth.SignInResp, error) {
	s.logger.Debug("[sign in user] execute service")

	var (
		err error
		u   auth.SignInResp
		ie  bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return auth.SignInResp{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	ie, err = s.checkExistsUser(ctx, tx, dto.TelegramID, dto.Username)
	if err != nil {
		return auth.SignInResp{}, err
	}

	if ie {
		u, err = s.getUserAndGenerateTokens(ctx, tx, dto.TelegramID)
	} else {
		u, err = s.createUser(ctx, tx, dto)
	}
	if err != nil {
		return auth.SignInResp{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return auth.SignInResp{}, err
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
	ieFromCache, err := s.bigCache.User.Exists(telegramID)
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

// createUser creates a new user in the database and generates JWT access, refresh tokens.
// After creation, the user is cached using the Telegram ID as the key.
func (s *SignIn) createUser(ctx context.Context, tx pgx.Tx, dto auth.SignInDTO) (auth.SignInResp, error) {
	createDTO := user.CreateDTO{
		TelegramID: dto.TelegramID,
		Username:   dto.Username,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
	}

	// create new user in the database.
	nu, err := s.userRepository.Create.Execute(ctx, tx, createDTO)
	if err != nil {
		return auth.SignInResp{}, err
	}

	// create user level history (set 1 level for new user).
	if err := s.createUserLevelHistory(ctx, tx, nu.TelegramID); err != nil {
		return auth.SignInResp{}, err
	}

	// assign daily task by telegram id.
	if _, err := s.userDailyTaskRepository.AssignDailyTaskByTelegramID.Execute(ctx, tx, nu.TelegramID); err != nil {
		return auth.SignInResp{}, err
	}

	// generate access, refresh tokens.
	tokens, err := s.jwt.Generate(nu.TelegramID)
	if err != nil {
		return auth.SignInResp{}, err
	}

	// set new refresh token in cache.
	if err := s.redis.RefreshToken.SetWithExpiration(ctx, nu.TelegramID, tokens.RefreshToken, time.Until(tokens.RefreshExpAt)); err != nil {
		return auth.SignInResp{}, err
	}

	// save the newly created user in the cache.
	if err := s.bigCache.User.Set(nu.TelegramID, nu); err != nil {
		s.logger.Warn(fmt.Sprintf("failed to cache new user: %v", err))
	}

	return auth.SignInResp{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		AccessExpAt:  tokens.AccessExpAt,
		RefreshExpAt: tokens.RefreshExpAt,
	}, nil
}

// createUserLevelHistory create user level history (level 1).
func (s *SignIn) createUserLevelHistory(ctx context.Context, tx pgx.Tx, telegramID string) error {
	const (
		levelNumber = 1
		xpAtReach   = 0
	)

	createUserLevelHistoryDTO := level.CreateUserLevelHistoryDTO{
		TelegramID:  telegramID,
		LevelNumber: levelNumber,
		XPAtReach:   xpAtReach,
	}

	if _, err := s.levelRepository.CreateUserLevelHistory.Execute(ctx, tx, createUserLevelHistoryDTO); err != nil {
		return err
	}

	return nil
}

// getUserAndGenerateTokens get user from cache or database.
// If the user is found to generate tokens.
func (s *SignIn) getUserAndGenerateTokens(ctx context.Context, tx pgx.Tx, telegramID string) (auth.SignInResp, error) {
	// get user from cache or database.
	u, err := s.findOrReturnExisting(ctx, tx, telegramID)
	if err != nil {
		return auth.SignInResp{}, err
	}

	// generate access, refresh tokens.
	tokens, err := s.jwt.Generate(u.TelegramID)
	if err != nil {
		return auth.SignInResp{}, err
	}

	// set new refresh token in cache.
	if err := s.redis.RefreshToken.SetWithExpiration(ctx, u.TelegramID, tokens.RefreshToken, time.Until(tokens.RefreshExpAt)); err != nil {
		return auth.SignInResp{}, err
	}

	// save the user in the cache.
	if err := s.bigCache.User.Set(u.TelegramID, u); err != nil {
		s.logger.Warn(fmt.Sprintf("failed to cache new user: %v", err))
	}

	return auth.SignInResp{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		AccessExpAt:  tokens.AccessExpAt,
		RefreshExpAt: tokens.RefreshExpAt,
	}, nil
}

// findOrReturnExisting attempts to retrieve a user from the cache by Telegram ID.
// If the user is found in the cache and the data is valid, it returns the cached user.
// Otherwise, it queries the database to retrieve the user by Telegram ID.
func (s *SignIn) findOrReturnExisting(ctx context.Context, tx pgx.Tx, telegramID string) (user.User, error) {
	// Try to get the user from the cache.
	userFromCache, err := s.bigCache.User.Get(telegramID)
	if err == nil && userFromCache.TelegramID == telegramID {
		// Cache hit and valid data — return the cached user.
		return userFromCache, nil
	}

	// Cache miss or invalid data — fallback to the database.
	return s.userRepository.GetByTelegramID.Execute(ctx, tx, telegramID)
}
