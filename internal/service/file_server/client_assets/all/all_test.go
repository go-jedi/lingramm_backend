package all

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/file_server/client_assets"
	allmocks "github.com/go-jedi/lingramm_backend/internal/repository/file_server/client_assets/all/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	poolsmocks "github.com/go-jedi/lingramm_backend/pkg/postgres/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		ctx context.Context
	}

	type want struct {
		clientAssets []clientassets.ClientAssets
		err          error
	}

	var (
		ctx              = context.TODO()
		testClientAssets = []clientassets.ClientAssets{
			{
				ID:             gofakeit.Int64(),
				NameFile:       gofakeit.Name(),
				ServerPathFile: gofakeit.URL(),
				ClientPathFile: gofakeit.URL(),
				Extension:      gofakeit.FileExtension(),
				Quality:        gofakeit.IntRange(1, 100),
				OldNameFile:    gofakeit.Name(),
				OldExtension:   gofakeit.FileExtension(),
				CreatedAt:      gofakeit.Date(),
				UpdatedAt:      gofakeit.Date(),
			},
		}
		txOptions = pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		}
		queryTimeout = int64(2)
	)

	tests := []struct {
		name               string
		mockPoolBehavior   func(m *poolsmocks.IPool, tx *poolsmocks.ITx)
		mockTxBehavior     func(tx *poolsmocks.ITx)
		mockLoggerBehavior func(m *loggermocks.ILogger)
		mockAllBehavior    func(m *allmocks.IAll, tx *poolsmocks.ITx)
		in                 in
		want               want
	}{
		{
			name: "ok",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Commit", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get all client assets] execute service")
			},
			mockAllBehavior: func(m *allmocks.IAll, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
				).Return(testClientAssets, nil)
			},
			in: in{
				ctx: ctx,
			},
			want: want{
				clientAssets: testClientAssets,
				err:          nil,
			},
		},
		{
			name: "begin_transaction_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, _ *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(nil, errors.New("begin transaction error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get all client assets] execute service")
			},
			in: in{
				ctx: ctx,
			},
			want: want{
				clientAssets: nil,
				err:          errors.New("begin transaction error"),
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
				m.On("Debug", "[get all client assets] execute service")
			},
			mockAllBehavior: func(m *allmocks.IAll, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
				).Return(nil, errors.New("some error"))
			},
			in: in{
				ctx: ctx,
			},
			want: want{
				clientAssets: nil,
				err:          errors.New("some error"),
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
				m.On("Debug", "[get all client assets] execute service")
			},
			mockAllBehavior: func(m *allmocks.IAll, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
				).Return(testClientAssets, nil)
			},
			in: in{
				ctx: ctx,
			},
			want: want{
				clientAssets: nil,
				err:          errors.New("commit error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockPool := poolsmocks.NewIPool(t)
			mockTx := poolsmocks.NewITx(t)
			mockLogger := loggermocks.NewILogger(t)
			mockAll := allmocks.NewIAll(t)

			if test.mockPoolBehavior != nil {
				test.mockPoolBehavior(mockPool, mockTx)
			}
			if test.mockTxBehavior != nil {
				test.mockTxBehavior(mockTx)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}
			if test.mockAllBehavior != nil {
				test.mockAllBehavior(mockAll, mockTx)
			}

			car := &clientassetsrepository.Repository{
				All: mockAll,
			}

			pg := &postgres.Postgres{
				Pool:         mockPool,
				QueryTimeout: queryTimeout,
			}

			all := New(car, mockLogger, pg)

			result, err := all.Execute(test.in.ctx)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.clientAssets, result)

			mockPool.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockAll.AssertExpectations(t)
		})
	}
}
