package create

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
		ctx context.Context
		dto user.CreateDTO
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
		dto        = user.CreateDTO{
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
					dto.UUID, dto.TelegramID, dto.Username,
					dto.FirstName, dto.LastName,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Run(func(args mock.Arguments) {
					id := args.Get(0).(*int64)
					*id = testUser.ID

					uuid := args.Get(1).(*string)
					*uuid = testUser.UUID

					tgID := args.Get(2).(*string)
					*tgID = testUser.TelegramID

					un := args.Get(3).(*string)
					*un = testUser.Username

					fn := args.Get(4).(*string)
					*fn = testUser.FirstName

					ln := args.Get(5).(*string)
					*ln = testUser.LastName

					ca := args.Get(6).(*time.Time)
					*ca = testUser.CreatedAt

					ua := args.Get(7).(*time.Time)
					*ua = testUser.UpdatedAt
				}).Return(nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a new user] execute repository")
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
			name: "timeout error",
			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
				tx.On(
					"QueryRow",
					mock.Anything,
					mock.Anything,
					dto.UUID, dto.TelegramID, dto.Username,
					dto.FirstName, dto.LastName,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Return(context.DeadlineExceeded)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a new user] execute repository")
				m.On("Error", "request timed out while creating the user", "err", context.DeadlineExceeded)
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: user.User{},
				err:  errors.New("the request timed out"),
			},
		},
		{
			name: "database error",
			mockTxBehavior: func(tx *poolsmocks.ITx, row *poolsmocks.RowMock) {
				tx.On(
					"QueryRow",
					mock.Anything,
					mock.Anything,
					dto.UUID, dto.TelegramID, dto.Username,
					dto.FirstName, dto.LastName,
				).Return(row)

				row.On("Scan",
					mock.AnythingOfType("*int64"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*string"),
					mock.AnythingOfType("*time.Time"),
					mock.AnythingOfType("*time.Time"),
				).Return(errors.New("database error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a new user] execute repository")
				m.On("Error", "failed to create user", "err", errors.New("database error"))
			},
			in: in{
				ctx: ctx,
				dto: dto,
			},
			want: want{
				user: user.User{},
				err:  errors.New("could not create user: database error"),
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

			result, err := create.Execute(test.in.ctx, mockTx, test.in.dto)

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
