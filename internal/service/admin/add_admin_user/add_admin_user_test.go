package addadminuser

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/admin"
	adminrepository "github.com/go-jedi/lingramm_backend/internal/repository/admin"
	addadminusermocks "github.com/go-jedi/lingramm_backend/internal/repository/admin/add_admin_user/mocks"
	existsbytelegramidmocks "github.com/go-jedi/lingramm_backend/internal/repository/admin/exists_by_telegram_id/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	adminbigcachemocks "github.com/go-jedi/lingramm_backend/pkg/bigcache/admin/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	poolsmocks "github.com/go-jedi/lingramm_backend/pkg/postgres/mocks"
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
		admin admin.Admin
		err   error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
		testAdmin  = admin.Admin{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
		}
		txOptions = pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		}
		queryTimeout = int64(2)
		errCache     = errors.New("cache set failed")
	)

	tests := []struct {
		name                           string
		mockPoolBehavior               func(m *poolsmocks.IPool, tx *poolsmocks.ITx)
		mockTxBehavior                 func(tx *poolsmocks.ITx)
		mockLoggerBehavior             func(m *loggermocks.ILogger)
		mockAdminBigCacheBehavior      func(m *adminbigcachemocks.IAdmin)
		mockExistsByTelegramIDBehavior func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx)
		mockAddAdminUserBehavior       func(m *addadminusermocks.IAddAdminUser, tx *poolsmocks.ITx)
		in                             in
		want                           want
	}{
		{
			name: "ok_admin_does_not_exists_create_admin_cache_miss_db_hit",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute service")
			},
			mockAdminBigCacheBehavior: func(m *adminbigcachemocks.IAdmin) {
				prefixTelegramID := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefixTelegramID)
				m.On("Exists", telegramID, prefixTelegramID).Return(false, bigcache.ErrEntryNotFound)
				m.On("Set", telegramID, testAdmin, prefixTelegramID).Return(nil)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(false, nil)
			},
			mockAddAdminUserBehavior: func(m *addadminusermocks.IAddAdminUser, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(testAdmin, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				admin: testAdmin,
				err:   nil,
			},
		},
		{
			name: "ok_admin_exists_cache_hit",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute service")
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
				admin: admin.Admin{},
				err:   apperrors.ErrAdminAlreadyExists,
			},
		},
		{
			name: "ok_admin_created_cache_set_failed",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute service")
				m.On("Warn", fmt.Sprintf("failed to cache new admin: %v", errCache))
			},
			mockAdminBigCacheBehavior: func(m *adminbigcachemocks.IAdmin) {
				prefixTelegramID := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefixTelegramID)
				m.On("Exists", telegramID, prefixTelegramID).Return(false, bigcache.ErrEntryNotFound)
				m.On("Set", telegramID, testAdmin, prefixTelegramID).Return(errCache)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(false, nil)
			},
			mockAddAdminUserBehavior: func(m *addadminusermocks.IAddAdminUser, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(testAdmin, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				admin: testAdmin,
				err:   nil,
			},
		},
		{
			name: "begin_transaction_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, _ *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(nil, errors.New("begin transaction error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute service")
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				admin: admin.Admin{},
				err:   errors.New("begin transaction error"),
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
				m.On("Debug", "[add a new admin user] execute service")
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
				admin: admin.Admin{},
				err:   errors.New("some error"),
			},
		},
		{
			name: "rollback_transaction_on_create_admin_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute service")
			},
			mockAdminBigCacheBehavior: func(m *adminbigcachemocks.IAdmin) {
				prefixTelegramID := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefixTelegramID)
				m.On("Exists", telegramID, prefixTelegramID).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(false, nil)
			},
			mockAddAdminUserBehavior: func(m *addadminusermocks.IAddAdminUser, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(admin.Admin{}, errors.New("create admin in database error"))
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				admin: admin.Admin{},
				err:   errors.New("create admin in database error"),
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
				m.On("Debug", "[add a new admin user] execute service")
			},
			mockAdminBigCacheBehavior: func(m *adminbigcachemocks.IAdmin) {
				prefixTelegramID := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefixTelegramID)
				m.On("Exists", telegramID, prefixTelegramID).Return(false, bigcache.ErrEntryNotFound)
				m.On("Set", telegramID, testAdmin, prefixTelegramID).Return(nil)
			},
			mockExistsByTelegramIDBehavior: func(m *existsbytelegramidmocks.IExistsByTelegramID, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(false, nil)
			},
			mockAddAdminUserBehavior: func(m *addadminusermocks.IAddAdminUser, tx *poolsmocks.ITx) {
				m.On("Execute", ctx, tx, telegramID).Return(testAdmin, nil)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				admin: admin.Admin{},
				err:   errors.New("commit error"),
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
			mockAddAdminUser := addadminusermocks.NewIAddAdminUser(t)

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
			if test.mockAddAdminUserBehavior != nil {
				test.mockAddAdminUserBehavior(mockAddAdminUser, mockTx)
			}

			ar := &adminrepository.Repository{
				AddAdminUser:       mockAddAdminUser,
				ExistsByTelegramID: mockExistsByTelegramID,
			}

			pg := &postgres.Postgres{
				Pool:         mockPool,
				QueryTimeout: queryTimeout,
			}

			bc := &bigcachepkg.BigCache{
				Admin: mockAdminBigCache,
			}

			addAdminUser := New(ar, mockLogger, pg, bc)

			result, err := addAdminUser.Execute(test.in.ctx, test.in.telegramID)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.admin, result)

			mockPool.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockAdminBigCache.AssertExpectations(t)
			mockExistsByTelegramID.AssertExpectations(t)
			mockAddAdminUser.AssertExpectations(t)
		})
	}
}
