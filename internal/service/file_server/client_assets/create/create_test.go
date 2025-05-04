package create

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	clientassets "github.com/go-jedi/lingvogramm_backend/internal/domain/file_server/client_assets"
	clientassetsrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/file_server/client_assets"
	createmocks "github.com/go-jedi/lingvogramm_backend/internal/repository/file_server/client_assets/create/mocks"
	fileserver "github.com/go-jedi/lingvogramm_backend/pkg/file_server"
	clientassetsmocks "github.com/go-jedi/lingvogramm_backend/pkg/file_server/client_assets/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	poolsmocks "github.com/go-jedi/lingvogramm_backend/pkg/postgres/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		ctx  context.Context
		file *multipart.FileHeader
	}

	type want struct {
		clientAsset clientassets.ClientAssets
		err         error
	}

	var (
		ctx            = context.TODO()
		file           = &multipart.FileHeader{}
		nameFile       = gofakeit.Name()
		serverPathFile = gofakeit.URL()
		clientPathFile = gofakeit.URL()
		extension      = gofakeit.FileExtension()
		quality        = gofakeit.IntRange(1, 100)
		oldNameFile    = gofakeit.Name()
		oldExtension   = gofakeit.FileExtension()
		imageData      = clientassets.UploadAndConvertToWebpResponse{
			NameFile:       nameFile,
			ServerPathFile: serverPathFile,
			ClientPathFile: clientPathFile,
			Extension:      extension,
			Quality:        quality,
			OldNameFile:    oldNameFile,
			OldExtension:   oldExtension,
		}
		testClientAsset = clientassets.ClientAssets{
			ID:             gofakeit.Int64(),
			NameFile:       nameFile,
			ServerPathFile: serverPathFile,
			ClientPathFile: clientPathFile,
			Extension:      extension,
			Quality:        quality,
			OldNameFile:    oldNameFile,
			OldExtension:   oldExtension,
			CreatedAt:      gofakeit.Date(),
			UpdatedAt:      gofakeit.Date(),
		}
		txOptions = pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		}
		queryTimeout = int64(2)
	)

	tests := []struct {
		name                     string
		mockPoolBehavior         func(m *poolsmocks.IPool, tx *poolsmocks.ITx)
		mockTxBehavior           func(tx *poolsmocks.ITx)
		mockLoggerBehavior       func(m *loggermocks.ILogger)
		mockClientAssetsBehavior func(m *clientassetsmocks.IClientAssets)
		mockCreateBehavior       func(m *createmocks.ICreate, tx *poolsmocks.ITx)
		in                       in
		want                     want
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
				m.On("Debug", "[create a client assets] execute service")
			},
			mockClientAssetsBehavior: func(m *clientassetsmocks.IClientAssets) {
				m.On("UploadAndConvertToWebP", ctx, file).Return(imageData, nil)
			},
			mockCreateBehavior: func(m *createmocks.ICreate, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					imageData,
				).Return(testClientAsset, nil)
			},
			in: in{
				ctx:  ctx,
				file: file,
			},
			want: want{
				clientAsset: testClientAsset,
				err:         nil,
			},
		},
		{
			name: "begin_transaction_error",
			mockPoolBehavior: func(m *poolsmocks.IPool, _ *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(nil, errors.New("begin transaction error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute service")
			},
			in: in{
				ctx:  ctx,
				file: file,
			},
			want: want{
				clientAsset: clientassets.ClientAssets{},
				err:         errors.New("begin transaction error"),
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
				m.On("Debug", "[create a client assets] execute service")
				m.On("Warn", "failed to remove image after db error", "warn", mock.Anything)
			},
			mockClientAssetsBehavior: func(m *clientassetsmocks.IClientAssets) {
				m.On("UploadAndConvertToWebP", ctx, file).Return(imageData, nil)
			},
			mockCreateBehavior: func(m *createmocks.ICreate, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					imageData,
				).Return(clientassets.ClientAssets{}, errors.New("some error"))
			},
			in: in{
				ctx:  ctx,
				file: file,
			},
			want: want{
				clientAsset: clientassets.ClientAssets{},
				err:         errors.New("some error"),
			},
		},
		{
			name: "rollback_transaction_on_upload_and_convert_to_webp",
			mockPoolBehavior: func(m *poolsmocks.IPool, tx *poolsmocks.ITx) {
				m.On("BeginTx", mock.Anything, txOptions).Return(tx, nil)
			},
			mockTxBehavior: func(tx *poolsmocks.ITx) {
				tx.On("Rollback", mock.Anything).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute service")
			},
			mockClientAssetsBehavior: func(m *clientassetsmocks.IClientAssets) {
				m.On("UploadAndConvertToWebP", ctx, file).Return(clientassets.UploadAndConvertToWebpResponse{}, errors.New("some error"))
			},
			in: in{
				ctx:  ctx,
				file: file,
			},
			want: want{
				clientAsset: clientassets.ClientAssets{},
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
				m.On("Debug", "[create a client assets] execute service")
			},
			mockClientAssetsBehavior: func(m *clientassetsmocks.IClientAssets) {
				m.On("UploadAndConvertToWebP", ctx, file).Return(imageData, nil)
			},
			mockCreateBehavior: func(m *createmocks.ICreate, tx *poolsmocks.ITx) {
				m.On(
					"Execute",
					ctx,
					tx,
					imageData,
				).Return(testClientAsset, nil)
			},
			in: in{
				ctx:  ctx,
				file: file,
			},
			want: want{
				clientAsset: clientassets.ClientAssets{},
				err:         errors.New("commit error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockPool := poolsmocks.NewIPool(t)
			mockTx := poolsmocks.NewITx(t)
			mockLogger := loggermocks.NewILogger(t)
			mockClientAssets := clientassetsmocks.NewIClientAssets(t)
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
			if test.mockClientAssetsBehavior != nil {
				test.mockClientAssetsBehavior(mockClientAssets)
			}
			if test.mockCreateBehavior != nil {
				test.mockCreateBehavior(mockCreate, mockTx)
			}

			car := &clientassetsrepository.Repository{
				Create: mockCreate,
			}

			pg := &postgres.Postgres{
				Pool:         mockPool,
				QueryTimeout: queryTimeout,
			}

			fs := &fileserver.FileServer{
				ClientAssets: mockClientAssets,
			}

			create := New(car, mockLogger, pg, fs)

			result, err := create.Execute(test.in.ctx, test.in.file)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.clientAsset, result)

			mockPool.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockClientAssets.AssertExpectations(t)
			mockCreate.AssertExpectations(t)
		})
	}
}
