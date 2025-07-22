package getbytelegramid

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/user"
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
		user user.User
		err  error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
		testUser   = user.User{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			Username:   gofakeit.Username(),
			FirstName:  gofakeit.FirstName(),
			LastName:   gofakeit.LastName(),
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
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testUser.ID

					tgID := args.Get(1).(*string)
					*tgID = testUser.TelegramID

					un := args.Get(2).(*string)
					*un = testUser.Username

					fn := args.Get(3).(*string)
					*fn = testUser.FirstName

					ln := args.Get(4).(*string)
					*ln = testUser.LastName

					ca := args.Get(5).(*time.Time)
					*ca = testUser.CreatedAt

					ua := args.Get(6).(*time.Time)
					*ua = testUser.UpdatedAt
				}).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user by telegram id] execute repository")
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				user: testUser,
				err:  nil,
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
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Return(context.DeadlineExceeded)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user by telegram id] execute repository")
				m.On("Error", "request timed out while get user by telegram id", "err", context.DeadlineExceeded)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				user: user.User{},
				err:  errors.New("the request timed out: context deadline exceeded"),
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
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Return(errors.New("database error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user by telegram id] execute repository")
				m.On("Error", "failed to get user by telegram id", "err", errors.New("database error"))
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
			},
			want: want{
				user: user.User{},
				err:  errors.New("could not get user by telegram id: database error"),
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

			getByTelegramID := New(queryTimeout, mockLogger)

			result, err := getByTelegramID.Execute(test.in.ctx, mockTx, test.in.telegramID)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.user, result)

			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}
