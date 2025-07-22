package addadminuser

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/admin"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	poolsmocks "github.com/go-jedi/lingramm_backend/pkg/postgres/mocks"
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
					mock.Anything,
					mock.Anything,
					telegramID,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testAdmin.ID

					telegramID := args.Get(1).(*string)
					*telegramID = testAdmin.TelegramID

					createdAt := args.Get(2).(*time.Time)
					*createdAt = testAdmin.CreatedAt

					updatedAt := args.Get(3).(*time.Time)
					*updatedAt = testAdmin.UpdatedAt
				}).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute repository")
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
			name: "timeout error",
			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
				tx.On(
					"QueryRow",
					mock.Anything,
					mock.Anything,
					telegramID,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testAdmin.ID

					telegramID := args.Get(1).(*string)
					*telegramID = testAdmin.TelegramID

					createdAt := args.Get(2).(*time.Time)
					*createdAt = testAdmin.CreatedAt

					updatedAt := args.Get(3).(*time.Time)
					*updatedAt = testAdmin.UpdatedAt
				}).Return(context.DeadlineExceeded)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute repository")
				m.On("Error", "request timed out while add admin user", "err", context.DeadlineExceeded)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				admin: admin.Admin{},
				err:   errors.New("the request timed out"),
			},
		},
		{
			name: "database error",
			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
				tx.On(
					"QueryRow",
					mock.Anything,
					mock.Anything,
					telegramID,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testAdmin.ID

					telegramID := args.Get(1).(*string)
					*telegramID = testAdmin.TelegramID

					createdAt := args.Get(2).(*time.Time)
					*createdAt = testAdmin.CreatedAt

					updatedAt := args.Get(3).(*time.Time)
					*updatedAt = testAdmin.UpdatedAt
				}).Return(errors.New("database error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute repository")
				m.On("Error", "failed to add admin user", "err", errors.New("database error"))
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				admin: admin.Admin{},
				err:   errors.New("could not add admin user: database error"),
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

			addAdminUser := New(queryTimeout, mockLogger)

			result, err := addAdminUser.Execute(test.in.ctx, mockTx, test.in.telegramID)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.admin, result)

			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}
