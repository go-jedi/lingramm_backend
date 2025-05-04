package exists

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	poolsmocks "github.com/go-jedi/lingvogramm_backend/pkg/postgres/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		ctx        context.Context
		telegramID string
		username   string
	}

	type want struct {
		exists bool
		err    error
	}

	var (
		ctx          = context.TODO()
		telegramID   = gofakeit.UUID()
		username     = gofakeit.Username()
		exists       = gofakeit.Bool()
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
					username,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*bool"),
				).Run(func(args mock.Arguments) {
					ie := args.Get(0).(*bool)
					*ie = exists
				}).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check a user exists] execute repository")
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
				username:   username,
			},
			want: want{
				exists: exists,
				err:    nil,
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
					username,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*bool"),
				).Run(func(args mock.Arguments) {
					ie := args.Get(0).(*bool)
					*ie = exists
				}).Return(context.DeadlineExceeded)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check a user exists] execute repository")
				m.On("Error", "request timed out while check exists user", "err", context.DeadlineExceeded)
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
				username:   username,
			},
			want: want{
				exists: false,
				err:    errors.New("the request timed out"),
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
					username,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*bool"),
				).Run(func(args mock.Arguments) {
					ie := args.Get(0).(*bool)
					*ie = exists
				}).Return(errors.New("database error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check a user exists] execute repository")
				m.On("Error", "failed to check exists user", "err", errors.New("database error"))
			},
			in: in{
				ctx:        ctx,
				telegramID: telegramID,
				username:   username,
			},
			want: want{
				exists: false,
				err:    errors.New("could not check exists user: database error"),
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

			exists := New(queryTimeout, mockLogger)

			result, err := exists.Execute(test.in.ctx, mockTx, test.in.telegramID, test.in.username)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.exists, result)

			mockTx.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}
