package refresh

import (
	"context"
	"errors"
	"testing"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	existsbytelegramidmocks "github.com/go-jedi/lingramm_backend/internal/repository/v1/user/exists_by_telegram_id/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	userbigcachemocks "github.com/go-jedi/lingramm_backend/pkg/bigcache/user/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	jwtmocks "github.com/go-jedi/lingramm_backend/pkg/jwt/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	poolsmocks "github.com/go-jedi/lingramm_backend/pkg/postgres/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	refreshtokenredismocks "github.com/go-jedi/lingramm_backend/pkg/redis/refresh_token/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRefresh(t *testing.T) {
	type in struct {
		ctx context.Context
		dto auth.RefreshDTO
	}

	type want struct {
		result auth.RefreshResponse
		err    error
	}

	var (
		ctx          = context.TODO()
		telegramID   = gofakeit.UUID()
		refreshToken = gofakeit.UUID()
		dto          = auth.RefreshDTO{
			TelegramID:   telegramID,
			RefreshToken: refreshToken,
		}
		jwtVerifyResp = jwt.VerifyResp{
			TelegramID: telegramID,
			ExpAt:      gofakeit.Date(),
		}
		jwtGenerateResp = jwt.GenerateResp{
			AccessToken:  gofakeit.UUID(),
			RefreshToken: gofakeit.UUID(),
			AccessExpAt:  gofakeit.Date(),
			RefreshExpAt: gofakeit.Date(),
		}
		testResult = auth.RefreshResponse{
			AccessToken:  jwtGenerateResp.AccessToken,
			RefreshToken: jwtGenerateResp.RefreshToken,
			AccessExpAt:  jwtGenerateResp.AccessExpAt,
			RefreshExpAt: jwtGenerateResp.RefreshExpAt,
		}
		txOptions = pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		}
		queryTimeout = int64(2)
	)

	tests := []struct {
		name                           string
		mockPoolBehavior               func(m *poolsmocks.IPool, tx *poolsmocks.ITx)
		mockTxBehavior                 func(tx *poolsmocks.ITx)
		mockLoggerBehavior             func(m *loggermocks.ILogger)
		mockUserBigCacheBehavior       func(m *userbigcachemocks.IUser)
		mockRefreshTokenRedisBehavior  func(m *refreshtokenredismocks.IRefreshToken)
		mockJWTBehavior                func(m *jwtmocks.IJWT)
		mockExistsByTelegramIDBehavior func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx)
		in                             in
		want                           want
	}{
		{
			name: "ok_user_exists_cache_miss_db_hit",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockRefreshTokenRedisBehavior: func(m *refreshtokenredismocks.IRefreshToken) {
				m.On("Get", ctx, dto.TelegramID).Return(refreshToken, nil)
				m.On("SetWithExpiration", ctx, dto.TelegramID, jwtGenerateResp.RefreshToken, mock.Anything).Return(nil)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwtVerifyResp, nil)
				m.On("Generate", jwtVerifyResp.TelegramID).Return(jwtGenerateResp, nil)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: testResult,
				err:    nil,
			},
		},
		{
			name: "ok_user_exists_cache_hit",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(true, nil)
			},
			mockRefreshTokenRedisBehavior: func(m *refreshtokenredismocks.IRefreshToken) {
				m.On("Get", ctx, dto.TelegramID).Return(refreshToken, nil)
				m.On("SetWithExpiration", ctx, dto.TelegramID, jwtGenerateResp.RefreshToken, mock.Anything).Return(nil)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwtVerifyResp, nil)
				m.On("Generate", jwtVerifyResp.TelegramID).Return(jwtGenerateResp, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: testResult,
				err:    nil,
			},
		},
		{
			name: "err_exists_by_telegram_id_from_db",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(false, errors.New("some error"))
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    errors.New("some error"),
			},
		},
		{
			name: "err_user_already_exists",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(false, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    apperrors.ErrUserDoesNotExist,
			},
		},
		{
			name: "err_verify_token",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwt.VerifyResp{}, errors.New("some error"))
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    errors.New("some error"),
			},
		},
		{
			name: "err_access_get_refresh_token_from_redis",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			mockRefreshTokenRedisBehavior: func(m *refreshtokenredismocks.IRefreshToken) {
				m.On("Get", ctx, dto.TelegramID).Return("", errors.New("some error"))
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwtVerifyResp, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    errors.New("some error"),
			},
		},
		{
			name: "err_session_invalid_get_refresh_token_from_redis",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			mockRefreshTokenRedisBehavior: func(m *refreshtokenredismocks.IRefreshToken) {
				m.On("Get", ctx, dto.TelegramID).Return("", nil)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwtVerifyResp, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    apperrors.ErrNoActiveSessionFound,
			},
		},
		{
			name: "err_tokens_do_not_match_get_refresh_token_from_redis",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			mockRefreshTokenRedisBehavior: func(m *refreshtokenredismocks.IRefreshToken) {
				m.On("Get", ctx, dto.TelegramID).Return(gofakeit.UUID(), nil)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwtVerifyResp, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    apperrors.ErrTokenMismatchOrExpired,
			},
		},
		{
			name: "err_generate_tokens",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			mockRefreshTokenRedisBehavior: func(m *refreshtokenredismocks.IRefreshToken) {
				m.On("Get", ctx, dto.TelegramID).Return(refreshToken, nil)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwtVerifyResp, nil)
				m.On("Generate", jwtVerifyResp.TelegramID).Return(jwt.GenerateResp{}, errors.New("some error"))
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    errors.New("some error"),
			},
		},
		{
			name: "err_save_new_refresh_token_in_cache",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			mockRefreshTokenRedisBehavior: func(m *refreshtokenredismocks.IRefreshToken) {
				m.On("Get", ctx, dto.TelegramID).Return(refreshToken, nil)
				m.On("SetWithExpiration", ctx, dto.TelegramID, jwtGenerateResp.RefreshToken, mock.Anything).Return(errors.New("some error"))
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwtVerifyResp, nil)
				m.On("Generate", jwtVerifyResp.TelegramID).Return(jwtGenerateResp, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    errors.New("some error"),
			},
		},
		{
			name: "commit_transaction_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(errors.New("commit error"))
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute service")
			},
			mockUserBigCacheBehavior: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			mockRefreshTokenRedisBehavior: func(m *refreshtokenredismocks.IRefreshToken) {
				m.On("Get", ctx, dto.TelegramID).Return(refreshToken, nil)
				m.On("SetWithExpiration", ctx, dto.TelegramID, jwtGenerateResp.RefreshToken, mock.Anything).Return(nil)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.RefreshToken).Return(jwtVerifyResp, nil)
				m.On("Generate", jwtVerifyResp.TelegramID).Return(jwtGenerateResp, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.RefreshResponse{},
				err:    errors.New("commit error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockPool := poolsmocks.NewIPool(t)
			mockTx := poolsmocks.NewITx(t)
			mockLogger := loggermocks.NewILogger(t)
			mockUserBigCache := userbigcachemocks.NewIUser(t)
			mockRefreshTokenRedis := refreshtokenredismocks.NewIRefreshToken(t)
			mockJWT := jwtmocks.NewIJWT(t)
			mockExistsByTelegramID := existsbytelegramidmocks.NewIExistsByTelegramID(t)

			if test.mockPoolBehavior != nil {
				test.mockPoolBehavior(mockPool, mockTx)
			}
			if test.mockTxBehavior != nil {
				test.mockTxBehavior(mockTx)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}
			if test.mockUserBigCacheBehavior != nil {
				test.mockUserBigCacheBehavior(mockUserBigCache)
			}
			if test.mockRefreshTokenRedisBehavior != nil {
				test.mockRefreshTokenRedisBehavior(mockRefreshTokenRedis)
			}
			if test.mockJWTBehavior != nil {
				test.mockJWTBehavior(mockJWT)
			}
			if test.mockExistsByTelegramIDBehavior != nil {
				test.mockExistsByTelegramIDBehavior(mockExistsByTelegramID, mockTx)
			}

			ur := &userrepository.Repository{
				ExistsByTelegramID: mockExistsByTelegramID,
			}
			pg := &postgres.Postgres{
				Pool:         mockPool,
				QueryTimeout: queryTimeout,
			}
			bc := &bigcachepkg.BigCache{
				User: mockUserBigCache,
			}
			r := &redis.Redis{
				RefreshToken: mockRefreshTokenRedis,
			}

			refresh := New(ur, mockLogger, pg, r, bc, mockJWT)

			result, err := refresh.Execute(test.in.ctx, test.in.dto)
			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.result, result)

			mockPool.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockUserBigCache.AssertExpectations(t)
			mockRefreshTokenRedis.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
			mockExistsByTelegramID.AssertExpectations(t)
		})
	}
}
