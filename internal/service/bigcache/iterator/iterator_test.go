package iterator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/user"
	bigcachepkg "github.com/go-jedi/lingvogramm_backend/pkg/bigcache"
	iteratorbigcachemocks "github.com/go-jedi/lingvogramm_backend/pkg/bigcache/iterator/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	type in struct {
		ctx context.Context
	}

	type want struct {
		result map[string]any
		err    error
	}

	var (
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
		key        = "user:telegram_id:" + telegramID
		testUser   = user.User{
			ID:         gofakeit.Int64(),
			UUID:       gofakeit.UUID(),
			TelegramID: telegramID,
			Username:   gofakeit.Username(),
			FirstName:  gofakeit.FirstName(),
			LastName:   gofakeit.LastName(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		testResult = map[string]any{
			key: testUser,
		}
	)

	tests := []struct {
		name                 string
		mockLoggerBehavior   func(m *loggermocks.ILogger)
		mockIteratorBigCache func(m *iteratorbigcachemocks.IIterator)
		in                   in
		want                 want
	}{
		{
			name: "ok",
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[iterator for show data in bigcache] execute service")
			},
			mockIteratorBigCache: func(m *iteratorbigcachemocks.IIterator) {
				m.On("Iterator", ctx).Return(testResult, nil)
			},
			in: in{
				ctx: ctx,
			},
			want: want{
				result: testResult,
				err:    nil,
			},
		},
		{
			name: "no_data_in_cache",
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[iterator for show data in bigcache] execute service")
			},
			mockIteratorBigCache: func(m *iteratorbigcachemocks.IIterator) {
				m.On("Iterator", ctx).Return(map[string]any{}, nil)
			},
			in: in{
				ctx: ctx,
			},
			want: want{
				result: nil,
				err:    nil,
			},
		},
		{
			name: "get_data_from_cache_error",
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[iterator for show data in bigcache] execute service")
			},
			mockIteratorBigCache: func(m *iteratorbigcachemocks.IIterator) {
				m.On("Iterator", ctx).Return(nil, errors.New("some error"))
			},
			in: in{
				ctx: ctx,
			},
			want: want{
				result: nil,
				err:    errors.New("some error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockLogger := loggermocks.NewILogger(t)
			mockIteratorBigCache := iteratorbigcachemocks.NewIIterator(t)

			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}
			if test.mockIteratorBigCache != nil {
				test.mockIteratorBigCache(mockIteratorBigCache)
			}

			bc := &bigcachepkg.BigCache{
				Iterator: mockIteratorBigCache,
			}

			iterator := New(mockLogger, bc)

			result, err := iterator.Execute(test.in.ctx)

			if test.want.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.result, result)

			mockLogger.AssertExpectations(t)
			mockIteratorBigCache.AssertExpectations(t)
		})
	}
}
