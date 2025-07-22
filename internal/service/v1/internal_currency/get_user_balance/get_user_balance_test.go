package getuserbalance

import (
	"context"
	"errors"
	"testing"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/internal_currency/user_balance"
	internalcurrencyrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency"
	getuserbalancemocks "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency/get_user_balance/mocks"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	existsbytelegramidmocks "github.com/go-jedi/lingramm_backend/internal/repository/v1/user/exists_by_telegram_id/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	userbigcachemocks "github.com/go-jedi/lingramm_backend/pkg/bigcache/user/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	poolsmocks "github.com/go-jedi/lingramm_backend/pkg/postgres/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		ctx        context.Context
		telegramID string
	}

	type want struct {
		userBalance userbalance.UserBalance
		err         error
	}

	var (
		ctx             = context.TODO()
		telegramID      = gofakeit.UUID()
		testUserBalance = userbalance.UserBalance{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			Balance:    decimal.NewFromFloat(0.00),
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
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
		mockExistsByTelegramIDBehavior func(mock *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx)
		mockGetUserBalanceBehavior     func(m *getuserbalancemocks.IGetUserBalance, tx *poolsmocks.ITx)
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
				m.On("Debug", "[get user balance] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", telegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(true, nil)
			},
			mockGetUserBalanceBehavior: func(m *getuserbalancemocks.IGetUserBalance, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(testUserBalance, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: testUserBalance,
				err:         nil,
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
				m.On("Debug", "[get user balance] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", telegramID, prefix).Return(true, nil)
			},
			mockGetUserBalanceBehavior: func(m *getuserbalancemocks.IGetUserBalance, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(testUserBalance, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: testUserBalance,
				err:         nil,
			},
		},
		{
			name: "ok_user_does_not_exists",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", telegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(false, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: userbalance.UserBalance{},
				err:         apperrors.ErrUserDoesNotExist,
			},
		},
		{
			name: "begin_transaction_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, _ *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(nil, errors.New("begin transaction error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute service")
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: userbalance.UserBalance{},
				err:         errors.New("begin transaction error"),
			},
		},
		{
			name: "exists_by_telegram_id_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", telegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(false, errors.New("some error"))
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: userbalance.UserBalance{},
				err:         errors.New("some error"),
			},
		},
		{
			name: "get_user_balance_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", telegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(true, nil)
			},
			mockGetUserBalanceBehavior: func(m *getuserbalancemocks.IGetUserBalance, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(userbalance.UserBalance{}, errors.New("some error"))
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: userbalance.UserBalance{},
				err:         errors.New("some error"),
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
				m.On("Debug", "[get user balance] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", telegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(true, nil)
			},
			mockGetUserBalanceBehavior: func(m *getuserbalancemocks.IGetUserBalance, tx *poolsmocks.ITx) {
				m.On("Execute", mock.Anything, tx, telegramID).Return(testUserBalance, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: userbalance.UserBalance{},
				err:         errors.New("commit error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockPool := poolsmocks.NewIPool(t)
			mockTx := poolsmocks.NewITx(t)
			mockLogger := loggermocks.NewILogger(t)
			mockUserBigCache := userbigcachemocks.NewIUser(t)
			mockExistsByTelegramID := existsbytelegramidmocks.NewIExistsByTelegramID(t)
			mockGetUserBalance := getuserbalancemocks.NewIGetUserBalance(t)

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
			if test.mockExistsByTelegramIDBehavior != nil {
				test.mockExistsByTelegramIDBehavior(mockExistsByTelegramID, mockTx)
			}
			if test.mockGetUserBalanceBehavior != nil {
				test.mockGetUserBalanceBehavior(mockGetUserBalance, mockTx)
			}

			ur := &userrepository.Repository{
				ExistsByTelegramID: mockExistsByTelegramID,
			}

			ic := &internalcurrencyrepository.Repository{
				GetUserBalance: mockGetUserBalance,
			}

			pg := &postgres.Postgres{
				Pool:         mockPool,
				QueryTimeout: queryTimeout,
			}

			bc := &bigcachepkg.BigCache{
				User: mockUserBigCache,
			}

			getUserBalance := New(ic, ur, mockLogger, pg, bc)

			result, err := getUserBalance.Execute(test.in.ctx, test.in.telegramID)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.userBalance, result)

			mockPool.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockUserBigCache.AssertExpectations(t)
			mockExistsByTelegramID.AssertExpectations(t)
			mockGetUserBalance.AssertExpectations(t)
		})
	}
}
