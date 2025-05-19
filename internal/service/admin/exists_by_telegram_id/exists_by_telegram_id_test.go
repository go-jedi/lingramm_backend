package existsbytelegramid

import (
	"context"
	"errors"
	"testing"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	adminrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/admin"
	existsbytelegramidmocks "github.com/go-jedi/lingvogramm_backend/internal/repository/admin/exists_by_telegram_id/mocks"
	bigcachepkg "github.com/go-jedi/lingvogramm_backend/pkg/bigcache"
	adminbigcachemocks "github.com/go-jedi/lingvogramm_backend/pkg/bigcache/admin/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	poolsmocks "github.com/go-jedi/lingvogramm_backend/pkg/postgres/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		ctx        context.Context
		telegramID string
	}

	type want struct {
		exists bool
		err    error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
		txOptions  = pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		}
		queryTimeout = int64(2)
		// errCache     = errors.New("cache set failed")
	)

	tests := []struct {
		name                           string
		mockPoolBehavior               func(m *poolsmocks.IPool, tx *poolsmocks.ITx)
		mockTxBehavior                 func(tx *poolsmocks.ITx)
		mockLoggerBehavior             func(m *loggermocks.ILogger)
		mockAdminBigCacheBehavior      func(m *adminbigcachemocks.IAdmin)
		mockExistsByTelegramIDBehavior func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx)
		in                             in
		want                           want
	}{
		{
			name: "ok_admin_exists_cache_miss_db_hit",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check admin exists by telegram id] execute service")
			},
			mockAdminBigCacheBehavior: func(m *adminbigcachemocks.IAdmin) {
				prefixTelegramID := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefixTelegramID)
				m.On("Exists", telegramID, prefixTelegramID).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(true, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				exists: true,
				err:    nil,
			},
		},
		{
			name: "ok_admin_exists_cache_hit",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check admin exists by telegram id] execute service")
			},
			mockAdminBigCacheBehavior: func(m *adminbigcachemocks.IAdmin) {
				prefixTelegramID := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefixTelegramID)
				m.On("Exists", telegramID, prefixTelegramID).Return(true, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				exists: true,
				err:    nil,
			},
		},
		{
			name: "begin_transaction_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, _ *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(nil, errors.New("begin transaction error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check admin exists by telegram id] execute service")
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				exists: false,
				err:    errors.New("begin transaction error"),
			},
		},
		{
			name: "rollback_transaction_on_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check admin exists by telegram id] execute service")
			},
			mockAdminBigCacheBehavior: func(m *adminbigcachemocks.IAdmin) {
				prefixTelegramID := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefixTelegramID)
				m.On("Exists", telegramID, prefixTelegramID).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(false, errors.New("some error"))
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				exists: false,
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
				m.On("Debug", "[check admin exists by telegram id] execute service")
			},
			mockAdminBigCacheBehavior: func(m *adminbigcachemocks.IAdmin) {
				prefixTelegramID := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefixTelegramID)
				m.On("Exists", telegramID, prefixTelegramID).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(true, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				exists: false,
				err:    errors.New("commit error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockPool := poolsmocks.NewIPool(t)
			mockTx := poolsmocks.NewITx(t)
			mockLogger := loggermocks.NewILogger(t)
			mockAdminBigCache := adminbigcachemocks.NewIAdmin(t)
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
			if test.mockAdminBigCacheBehavior != nil {
				test.mockAdminBigCacheBehavior(mockAdminBigCache)
			}
			if test.mockExistsByTelegramIDBehavior != nil {
				test.mockExistsByTelegramIDBehavior(mockExistsByTelegramID, mockTx)
			}

			ar := &adminrepository.Repository{
				ExistsByTelegramID: mockExistsByTelegramID,
			}

			pg := &postgres.Postgres{
				Pool:         mockPool,
				QueryTimeout: queryTimeout,
			}

			bc := &bigcachepkg.BigCache{
				Admin: mockAdminBigCache,
			}

			existsByTelegramID := New(ar, mockLogger, pg, bc)

			result, err := existsByTelegramID.Execute(test.in.ctx, test.in.telegramID)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.exists, result)

			mockPool.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockAdminBigCache.AssertExpectations(t)
			mockExistsByTelegramID.AssertExpectations(t)
		})
	}
}
