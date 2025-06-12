package check

import (
	"context"
	"errors"
	"testing"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/user"
	existsbytelegramidmocks "github.com/go-jedi/lingramm_backend/internal/repository/user/exists_by_telegram_id/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	userbigcachemocks "github.com/go-jedi/lingramm_backend/pkg/bigcache/user/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	jwtmocks "github.com/go-jedi/lingramm_backend/pkg/jwt/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	poolsmocks "github.com/go-jedi/lingramm_backend/pkg/postgres/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheck(t *testing.T) {
	type in struct {
		ctx context.Context
		dto auth.CheckDTO
	}

	type want struct {
		result auth.CheckResponse
		err    error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
		token      = gofakeit.UUID()
		expAt      = gofakeit.Date()
		dto        = auth.CheckDTO{
			TelegramID: telegramID,
			Token:      token,
		}
		jwtVerifyResp = jwt.VerifyResp{
			TelegramID: telegramID,
			ExpAt:      expAt,
		}
		testResult = auth.CheckResponse{
			TelegramID: telegramID,
			Token:      token,
			ExpAt:      expAt,
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
		mockUserBigCache               func(m *userbigcachemocks.IUser)
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
				m.On("Debug", "[check user token] execute service")
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.Token).Return(jwtVerifyResp, nil)
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
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
				m.On("Debug", "[check user token] execute service")
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.Token).Return(jwtVerifyResp, nil)
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(true, nil)
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
			name: "ok_user_does_not_exists",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check user token] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
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
				result: auth.CheckResponse{},
				err:    apperrors.ErrUserDoesNotExist,
			},
		},
		{
			name: "begin_transaction_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, _ *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(nil, errors.New("begin transaction error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check user token] execute service")
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.CheckResponse{},
				err:    errors.New("begin transaction error"),
			},
		},
		{
			name: "err_get_user_from_db",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check user token] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
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
				result: auth.CheckResponse{},
				err:    errors.New("some error"),
			},
		},
		{
			name: "err_check_verify_token",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check user token] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.Token).Return(jwt.VerifyResp{}, errors.New("some error"))
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.CheckResponse{},
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
				m.On("Debug", "[check user token] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockJWTBehavior: func(m *jwtmocks.IJWT) {
				m.On("Verify", dto.TelegramID, dto.Token).Return(jwtVerifyResp, nil)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, dto.TelegramID).Return(true, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				result: auth.CheckResponse{},
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
			if test.mockUserBigCache != nil {
				test.mockUserBigCache(mockUserBigCache)
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

			check := New(ur, mockLogger, pg, bc, mockJWT)

			result, err := check.Execute(test.in.ctx, test.in.dto)
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
			mockJWT.AssertExpectations(t)
			mockExistsByTelegramID.AssertExpectations(t)
		})
	}
}
