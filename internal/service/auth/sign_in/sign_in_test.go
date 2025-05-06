package signin

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/auth"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/user"
	userrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/user"
	createmocks "github.com/go-jedi/lingvogramm_backend/internal/repository/user/create/mocks"
	existsmocks "github.com/go-jedi/lingvogramm_backend/internal/repository/user/exists/mocks"
	getbytelegramidmocks "github.com/go-jedi/lingvogramm_backend/internal/repository/user/get_by_telegram_id/mocks"
	bigcachepkg "github.com/go-jedi/lingvogramm_backend/pkg/bigcache"
	userbigcachemocks "github.com/go-jedi/lingvogramm_backend/pkg/bigcache/user/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	poolsmocks "github.com/go-jedi/lingvogramm_backend/pkg/postgres/mocks"
	uuidmocks "github.com/go-jedi/lingvogramm_backend/pkg/uuid/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		ctx context.Context
		dto auth.SignInDTO
	}

	type want struct {
		user user.User
		err  error
	}

	var (
		ctx        = context.TODO()
		uuid       = gofakeit.UUID()
		telegramID = gofakeit.UUID()
		username   = gofakeit.Username()
		firstname  = gofakeit.FirstName()
		lastname   = gofakeit.LastName()
		createdAt  = time.Now()
		updatedAt  = time.Now()
		dto        = auth.SignInDTO{
			TelegramID: telegramID,
			Username:   username,
			FirstName:  firstname,
			LastName:   lastname,
		}
		createDTO = user.CreateDTO{
			UUID:       uuid,
			TelegramID: telegramID,
			Username:   username,
			FirstName:  firstname,
			LastName:   lastname,
		}
		testUser = user.User{
			ID:         gofakeit.Int64(),
			UUID:       uuid,
			TelegramID: telegramID,
			Username:   username,
			FirstName:  firstname,
			LastName:   lastname,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		}
		txOptions = pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		}
		queryTimeout = int64(2)
		errCache     = errors.New("cache set failed")
	)

	tests := []struct {
		name                        string
		mockPoolBehavior            func(m *poolsmocks.IPool, tx *poolsmocks.ITx)
		mockTxBehavior              func(tx *poolsmocks.ITx)
		mockLoggerBehavior          func(m *loggermocks.ILogger)
		mockUUIDBehavior            func(m *uuidmocks.IUUID)
		mockUserBigCache            func(m *userbigcachemocks.IUser)
		mockExistsBehavior          func(mock *existsmocks.IExists, tx *poolsmocks.ITx)
		mockGetByTelegramIDBehavior func(mock *getbytelegramidmocks.IGetByTelegramID, tx *poolsmocks.ITx)
		mockCreateBehavior          func(mock *createmocks.ICreate, tx *poolsmocks.ITx)
		in                          in
		want                        want
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
				m.On("Debug", "[sign in user] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
				m.On("Get", dto.TelegramID, prefix).Return(user.User{}, bigcache.ErrEntryNotFound)
			},
			mockExistsBehavior: func(m *existsmocks.IExists, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					dto.TelegramID,
					dto.Username,
				).Return(true, nil)
			},
			mockGetByTelegramIDBehavior: func(m *getbytelegramidmocks.IGetByTelegramID, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					mock.Anything,
					tx,
					dto.TelegramID,
				).Return(testUser, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: testUser,
				err:  nil,
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
				m.On("Debug", "[sign in user] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(true, nil)
				m.On("Get", dto.TelegramID, prefix).Return(testUser, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: testUser,
				err:  nil,
			},
		},
		{
			name: "ok_user_does_not_exists_create_user",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute service")
			},
			mockUUIDBehavior: func(m *uuidmocks.IUUID) {
				m.On("Generate").Return(uuid, nil)
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
				m.On("Set", testUser.TelegramID, testUser, prefix).Return(nil)
			},
			mockExistsBehavior: func(m *existsmocks.IExists, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					dto.TelegramID,
					dto.Username,
				).Return(false, nil)
			},
			mockCreateBehavior: func(m *createmocks.ICreate, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					createDTO,
				).Return(testUser, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: testUser,
				err:  nil,
			},
		},
		{
			name: "ok_user_created_cache_set_failed",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute service")
				m.On("Warn", fmt.Sprintf("failed to cache new user: %v", errCache))
			},
			mockUUIDBehavior: func(m *uuidmocks.IUUID) {
				m.On("Generate").Return(uuid, nil)
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
				m.On("Set", testUser.TelegramID, testUser, prefix).Return(errCache)
			},
			mockExistsBehavior: func(m *existsmocks.IExists, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					dto.TelegramID,
					dto.Username,
				).Return(false, nil)
			},
			mockCreateBehavior: func(m *createmocks.ICreate, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					createDTO,
				).Return(testUser, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: testUser,
				err:  nil,
			},
		},
		{
			name: "begin_transaction_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, _ *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(nil, errors.New("begin transaction error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute service")
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: user.User{},
				err:  errors.New("begin transaction error"),
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
				m.On("Debug", "[sign in user] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsBehavior: func(m *existsmocks.IExists, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					dto.TelegramID,
					dto.Username,
				).Return(true, errors.New("some error"))
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: user.User{},
				err:  errors.New("some error"),
			},
		},
		{
			name: "rollback_transaction_on_create_user_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute service")
			},
			mockUUIDBehavior: func(m *uuidmocks.IUUID) {
				m.On("Generate").Return(uuid, nil)
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsBehavior: func(m *existsmocks.IExists, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					dto.TelegramID,
					dto.Username,
				).Return(false, nil)
			},
			mockCreateBehavior: func(m *createmocks.ICreate, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					createDTO,
				).Return(user.User{}, errors.New("create user in database error"))
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: user.User{},
				err:  errors.New("create user in database error"),
			},
		},
		{
			name: "rollback_transaction_on_get_by_telegram_id_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
				m.On("Get", dto.TelegramID, prefix).Return(user.User{}, bigcache.ErrEntryNotFound)
			},
			mockExistsBehavior: func(m *existsmocks.IExists, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					dto.TelegramID,
					dto.Username,
				).Return(true, nil)
			},
			mockGetByTelegramIDBehavior: func(m *getbytelegramidmocks.IGetByTelegramID, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					mock.Anything,
					tx,
					dto.TelegramID,
				).Return(user.User{}, errors.New("get by telegram id error"))
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: user.User{},
				err:  errors.New("get by telegram id error"),
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
				m.On("Debug", "[sign in user] execute service")
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
				m.On("Get", dto.TelegramID, prefix).Return(user.User{}, bigcache.ErrEntryNotFound)
			},
			mockExistsBehavior: func(m *existsmocks.IExists, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					dto.TelegramID,
					dto.Username,
				).Return(true, nil)
			},
			mockGetByTelegramIDBehavior: func(m *getbytelegramidmocks.IGetByTelegramID, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					mock.Anything,
					tx,
					dto.TelegramID,
				).Return(testUser, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: user.User{},
				err:  errors.New("commit error"),
			},
		},
		{
			name: "generate_uuid_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute service")
			},
			mockUUIDBehavior: func(m *uuidmocks.IUUID) {
				m.On("Generate").Return("", errors.New("generate error"))
			},
			mockUserBigCache: func(m *userbigcachemocks.IUser) {
				prefix := "telegram_id:"
				m.On("GetPrefixTelegramID").Return(prefix)
				m.On("Exists", dto.TelegramID, prefix).Return(false, bigcache.ErrEntryNotFound)
			},
			mockExistsBehavior: func(m *existsmocks.IExists, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					dto.TelegramID,
					dto.Username,
				).Return(false, nil)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: user.User{},
				err:  errors.New("generate error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockPool := poolsmocks.NewIPool(t)
			mockTx := poolsmocks.NewITx(t)
			mockLogger := loggermocks.NewILogger(t)
			mockUUID := uuidmocks.NewIUUID(t)
			mockUserBigCache := userbigcachemocks.NewIUser(t)
			mockExists := existsmocks.NewIExists(t)
			mockGetByTelegramID := getbytelegramidmocks.NewIGetByTelegramID(t)
			mockCreate := createmocks.NewICreate(t)

			if test.mockPoolBehavior != nil {
				test.mockPoolBehavior(mockPool, mockTx)
			}
			if test.mockTxBehavior != nil {
				test.mockTxBehavior(mockTx)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}
			if test.mockUUIDBehavior != nil {
				test.mockUUIDBehavior(mockUUID)
			}
			if test.mockUserBigCache != nil {
				test.mockUserBigCache(mockUserBigCache)
			}
			if test.mockExistsBehavior != nil {
				test.mockExistsBehavior(mockExists, mockTx)
			}
			if test.mockGetByTelegramIDBehavior != nil {
				test.mockGetByTelegramIDBehavior(mockGetByTelegramID, mockTx)
			}
			if test.mockCreateBehavior != nil {
				test.mockCreateBehavior(mockCreate, mockTx)
			}

			ur := &userrepository.Repository{
				Create:          mockCreate,
				Exists:          mockExists,
				GetByTelegramID: mockGetByTelegramID,
			}

			pg := &postgres.Postgres{
				Pool:         mockPool,
				QueryTimeout: queryTimeout,
			}

			bc := &bigcachepkg.BigCache{
				User: mockUserBigCache,
			}

			signIn := New(ur, mockLogger, pg, bc, mockUUID)

			result, err := signIn.Execute(test.in.ctx, test.in.dto)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.user, result)

			mockPool.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockUUID.AssertExpectations(t)
			mockExists.AssertExpectations(t)
			mockGetByTelegramID.AssertExpectations(t)
			mockCreate.AssertExpectations(t)
		})
	}
}
