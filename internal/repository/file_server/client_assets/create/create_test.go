package create

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	poolsmocks "github.com/go-jedi/lingramm_backend/pkg/postgres/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		ctx  context.Context
		data clientassets.UploadAndConvertToWebpResponse
	}

	type want struct {
		clientAsset clientassets.ClientAssets
		err         error
	}

	var (
		ctx            = context.TODO()
		nameFile       = gofakeit.Name()
		serverPathFile = gofakeit.URL()
		clientPathFile = gofakeit.URL()
		extension      = gofakeit.FileExtension()
		quality        = gofakeit.IntRange(1, 100)
		oldNameFile    = gofakeit.Name()
		oldExtension   = gofakeit.FileExtension()
		data           = clientassets.UploadAndConvertToWebpResponse{
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
		queryTimeout = int64(2)
	)

	tests := []struct {
		name               string
		mockTxBehavior     func(tx *poolsmocks.ITx, row *poolsmocks.RowMock)
		mockLoggerBehavior func(m *loggermocks.ILogger)
		in                 in
		want               want
	}{
		{
			name: "ok",
			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
				tx.On(
					"QueryRow",
					mock.Anything, mock.Anything,
					data.NameFile, data.ServerPathFile, data.ClientPathFile,
					data.Extension, data.Quality, data.OldNameFile, data.OldExtension,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*int"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testClientAsset.ID

					nameFile := args.Get(1).(*string)
					*nameFile = testClientAsset.NameFile

					serverPathFile := args.Get(2).(*string)
					*serverPathFile = testClientAsset.ServerPathFile

					clientPathFile := args.Get(3).(*string)
					*clientPathFile = testClientAsset.ClientPathFile

					extension := args.Get(4).(*string)
					*extension = testClientAsset.Extension

					quality := args.Get(5).(*int)
					*quality = testClientAsset.Quality

					oldNameFile := args.Get(6).(*string)
					*oldNameFile = testClientAsset.OldNameFile

					oldExtension := args.Get(7).(*string)
					*oldExtension = testClientAsset.OldExtension

					createdAt := args.Get(8).(*time.Time)
					*createdAt = testClientAsset.CreatedAt

					updatedAt := args.Get(9).(*time.Time)
					*updatedAt = testClientAsset.UpdatedAt
				}).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute repository")
			},
			in: in{
				ctx:  ctx,
				data: data,
			},
			want: want{
				clientAsset: testClientAsset,
				err:         nil,
			},
		},
		{
			name: "timeout error",
			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
				tx.On(
					"QueryRow",
					mock.Anything, mock.Anything,
					data.NameFile, data.ServerPathFile, data.ClientPathFile,
					data.Extension, data.Quality, data.OldNameFile, data.OldExtension,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*int"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Return(context.DeadlineExceeded)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute repository")
				m.On("Error", "request timed out while create a client assets", "err", context.DeadlineExceeded)
			},
			in: in{
				ctx:  ctx,
				data: data,
			},
			want: want{
				clientAsset: clientassets.ClientAssets{},
				err:         errors.New("the request timed out"),
			},
		},
		{
			name: "database error",
			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
				tx.On(
					"QueryRow",
					mock.Anything, mock.Anything,
					data.NameFile, data.ServerPathFile, data.ClientPathFile,
					data.Extension, data.Quality, data.OldNameFile, data.OldExtension,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*int"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Return(errors.New("database error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute repository")
				m.On("Error", "failed to create a client assets", "err", errors.New("database error"))
			},
			in: in{
				ctx:  ctx,
				data: data,
			},
			want: want{
				clientAsset: clientassets.ClientAssets{},
				err:         errors.New("could not create a client assets: database error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockTx := poolsmocks.NewITx(t)
			mockRow := poolsmocks.NewMockRow(t)
			mockLogger := loggermocks.NewILogger(t)

			if test.mockTxBehavior != nil {
				test.mockTxBehavior(mockTx, mockRow)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}

			create := New(queryTimeout, mockLogger)

			result, err := create.Execute(test.in.ctx, mockTx, test.in.data)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.clientAsset, result)

			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}
