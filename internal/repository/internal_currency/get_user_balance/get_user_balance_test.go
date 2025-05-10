package getuserbalance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	userbalance "github.com/go-jedi/lingvogramm_backend/internal/domain/user_balance"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	poolsmocks "github.com/go-jedi/lingvogramm_backend/pkg/postgres/mocks"
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
					mock.AnythingOfType("*decimal.Decimal"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testUserBalance.ID

					telegramID := args.Get(1).(*string)
					*telegramID = testUserBalance.TelegramID

					balance := args.Get(2).(*decimal.Decimal)
					*balance = testUserBalance.Balance

					createdAt := args.Get(3).(*time.Time)
					*createdAt = testUserBalance.CreatedAt

					updatedAt := args.Get(4).(*time.Time)
					*updatedAt = testUserBalance.UpdatedAt
				}).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute repository")
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
					mock.AnythingOfType("*decimal.Decimal"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testUserBalance.ID

					telegramID := args.Get(1).(*string)
					*telegramID = testUserBalance.TelegramID

					balance := args.Get(2).(*decimal.Decimal)
					*balance = testUserBalance.Balance

					createdAt := args.Get(3).(*time.Time)
					*createdAt = testUserBalance.CreatedAt

					updatedAt := args.Get(4).(*time.Time)
					*updatedAt = testUserBalance.UpdatedAt
				}).Return(context.DeadlineExceeded)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute repository")
				m.On("Error", "request timed out while get user balance", "err", context.DeadlineExceeded)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: userbalance.UserBalance{},
				err:         errors.New("the request timed out"),
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
					mock.AnythingOfType("*decimal.Decimal"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testUserBalance.ID

					telegramID := args.Get(1).(*string)
					*telegramID = testUserBalance.TelegramID

					balance := args.Get(2).(*decimal.Decimal)
					*balance = testUserBalance.Balance

					createdAt := args.Get(3).(*time.Time)
					*createdAt = testUserBalance.CreatedAt

					updatedAt := args.Get(4).(*time.Time)
					*updatedAt = testUserBalance.UpdatedAt
				}).Return(errors.New("database error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute repository")
				m.On("Error", "failed to get user balance", "err", errors.New("database error"))
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				userBalance: userbalance.UserBalance{},
				err:         errors.New("could not get user balance: database error"),
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

			getUserBalance := New(queryTimeout, mockLogger)

			result, err := getUserBalance.Execute(test.in.ctx, mockTx, test.in.telegramID)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.userBalance, result)

			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}
