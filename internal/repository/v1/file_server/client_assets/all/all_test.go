package all

//import (
//	"context"
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/brianvoe/gofakeit/v7"
//	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
//	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
//	poolsmocks "github.com/go-jedi/lingramm_backend/pkg/postgres/mocks"
//	"github.com/jackc/pgx/v5"
//	"github.com/jackc/pgx/v5/pgconn"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//func TestExecute(t *testing.T) {
//	type in struct {
//		ctx context.Context
//	}
//
//	type want struct {
//		clientAssets []clientassets.ClientAssets
//		err          error
//	}
//
//	var (
//		ctx              = context.TODO()
//		testClientAssets = []clientassets.ClientAssets{
//			{
//				ID:             gofakeit.Int64(),
//				NameFile:       gofakeit.Name(),
//				ServerPathFile: gofakeit.URL(),
//				ClientPathFile: gofakeit.URL(),
//				Extension:      gofakeit.FileExtension(),
//				Quality:        gofakeit.IntRange(1, 100),
//				OldNameFile:    gofakeit.Name(),
//				OldExtension:   gofakeit.FileExtension(),
//				CreatedAt:      gofakeit.Date(),
//				UpdatedAt:      gofakeit.Date(),
//			},
//		}
//		queryTimeout = int64(2)
//	)
//
//	tests := []struct {
//		name               string
//		mockTxBehavior     func(tx *poolsmocks.ITx, row *poolsmocks.RowMock)
//		mockLoggerBehavior func(m *loggermocks.ILogger)
//		in                 in
//		want               want
//	}{
//		{
//			name: "ok",
//			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
//				rows := &poolsmocks.RowsMock{
//					NextFunc: func() bool {
//						return len(testClientAssets) > 0
//					},
//					ScanFunc: func(dest ...any) error {
//						if len(testClientAssets) == 0 {
//							return pgx.ErrNoRows
//						}
//
//						tca := testClientAssets[0]
//						testClientAssets = testClientAssets[1:]
//
//						*dest[0].(*int64) = tca.ID
//						*dest[1].(*string) = tca.NameFile
//						*dest[2].(*string) = tca.ServerPathFile
//						*dest[3].(*string) = tca.ClientPathFile
//						*dest[4].(*string) = tca.Extension
//						*dest[5].(*int) = tca.Quality
//						*dest[6].(*string) = tca.OldNameFile
//						*dest[7].(*string) = tca.OldExtension
//						*dest[8].(*time.Time) = tca.CreatedAt
//						*dest[9].(*time.Time) = tca.UpdatedAt
//
//						return nil
//					},
//					ErrFunc: func() error {
//						return nil
//					},
//					CloseFunc: func() {},
//					CommandTagFunc: func() pgconn.CommandTag {
//						return pgconn.NewCommandTag("SELECT")
//					},
//					FieldDescriptionsFunc: func() []pgconn.FieldDescription {
//						return []pgconn.FieldDescription{
//							{Name: "id", DataTypeOID: 20},
//							{Name: "name_file", DataTypeOID: 20},
//							{Name: "server_path_file", DataTypeOID: 25},
//							{Name: "client_path_file", DataTypeOID: 25},
//							{Name: "extension", DataTypeOID: 25},
//							{Name: "quality", DataTypeOID: 25},
//							{Name: "old_name_file", DataTypeOID: 25},
//							{Name: "old_extension", DataTypeOID: 25},
//							{Name: "created_at", DataTypeOID: 1114},
//							{Name: "updated_at", DataTypeOID: 1114},
//						}
//					},
//					ValuesFunc: func() ([]any, error) {
//						if len(testClientAssets) == 0 {
//							return nil, pgx.ErrNoRows
//						}
//						tca := testClientAssets[0]
//						return []any{
//							tca.ID,
//							tca.NameFile,
//							tca.ServerPathFile,
//							tca.ClientPathFile,
//							tca.Extension,
//							tca.Quality,
//							tca.OldNameFile,
//							tca.OldExtension,
//							tca.CreatedAt,
//							tca.UpdatedAt,
//						}, nil
//					},
//					RawValuesFunc: func() [][]byte {
//						return nil
//					},
//				}
//				tx.On("Query", mock.Anything, mock.Anything).Return(rows, nil)
//			},
//			mockLoggerBehavior: func(m *loggermocks.ILogger) {
//				m.On("Debug", "[get all client assets] execute repository")
//			},
//			in: in{
//				ctx: ctx,
//			},
//			want: want{
//				clientAssets: testClientAssets,
//				err:          nil,
//			},
//		},
//		{
//			name: "query error",
//			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
//				tx.On("Query", mock.Anything, mock.Anything).Return(nil, pgx.ErrNoRows)
//			},
//			mockLoggerBehavior: func(m *loggermocks.ILogger) {
//				m.On("Debug", "[get all client assets] execute repository")
//				m.On("Error", "failed to get all client assets", "err", pgx.ErrNoRows)
//			},
//			in: in{
//				ctx: ctx,
//			},
//			want: want{
//				clientAssets: nil,
//				err:          fmt.Errorf("could not get all client assets: %w", pgx.ErrNoRows),
//			},
//		},
//		{
//			name: "context timeout",
//			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
//				tx.On("Query", mock.Anything, mock.Anything).Return(nil, context.DeadlineExceeded)
//			},
//			mockLoggerBehavior: func(m *loggermocks.ILogger) {
//				m.On("Debug", "[get all client assets] execute repository")
//				m.On("Error", "request timed out while get all client assets", "err", context.DeadlineExceeded)
//			},
//			in: in{
//				ctx: ctx,
//			},
//			want: want{
//				clientAssets: nil,
//				err:          fmt.Errorf("the request timed out: %w", context.DeadlineExceeded),
//			},
//		},
//		{
//			name: "scan error",
//			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
//				rows := &poolsmocks.RowsMock{
//					NextFunc: func() bool {
//						return true
//					},
//					ScanFunc: func(dest ...any) error {
//						return pgx.ErrNoRows
//					},
//					ErrFunc: func() error {
//						return nil
//					},
//					CloseFunc: func() {},
//				}
//				tx.On("Query", mock.Anything, mock.Anything).Return(rows, nil)
//			},
//			mockLoggerBehavior: func(m *loggermocks.ILogger) {
//				m.On("Debug", "[get all client assets] execute repository")
//				m.On("Error", "failed to scan row to get all client assets", "err", pgx.ErrNoRows)
//			},
//			in: in{
//				ctx: ctx,
//			},
//			want: want{
//				clientAssets: nil,
//				err:          fmt.Errorf("failed to scan row to get all client assets: %w", pgx.ErrNoRows),
//			},
//		},
//		{
//			name: "rows error",
//			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
//				rows := &poolsmocks.RowsMock{
//					NextFunc: func() bool {
//						return false
//					},
//					ScanFunc: func(dest ...any) error {
//						return nil
//					},
//					ErrFunc: func() error {
//						return pgx.ErrNoRows
//					},
//					CloseFunc: func() {},
//				}
//				tx.On("Query", mock.Anything, mock.Anything).Return(rows, nil)
//			},
//			mockLoggerBehavior: func(m *loggermocks.ILogger) {
//				m.On("Debug", "[get all client assets] execute repository")
//				m.On("Error", "failed to get all client assets", "err", pgx.ErrNoRows)
//			},
//			in: in{
//				ctx: ctx,
//			},
//			want: want{
//				clientAssets: nil,
//				err:          fmt.Errorf("failed to get all client assets: %w", pgx.ErrNoRows),
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			mockTx := poolsmocks.NewITx(t)
//			mockRow := poolsmocks.NewMockRow(t)
//			mockLogger := loggermocks.NewILogger(t)
//
//			if test.mockTxBehavior != nil {
//				test.mockTxBehavior(mockTx, mockRow)
//			}
//			if test.mockLoggerBehavior != nil {
//				test.mockLoggerBehavior(mockLogger)
//			}
//
//			all := New(queryTimeout, mockLogger)
//
//			result, err := all.Execute(test.in.ctx, mockTx)
//
//			if test.want.err != nil {
//				assert.Error(t, err)
//				assert.Contains(t, err.Error(), test.want.err.Error())
//			} else {
//				assert.NoError(t, err)
//			}
//
//			assert.Equal(t, test.want.clientAssets, result)
//
//			mockTx.AssertExpectations(t)
//			mockLogger.AssertExpectations(t)
//			mockRow.AssertExpectations(t)
//		})
//	}
//}
