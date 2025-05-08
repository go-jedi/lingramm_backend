package iterator

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/user"
	bigcacheservice "github.com/go-jedi/lingvogramm_backend/internal/service/bigcache"
	servicemocks "github.com/go-jedi/lingvogramm_backend/internal/service/bigcache/iterator/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type want struct {
		statusCode int
		status     bool
		message    string
		error      string
	}

	var (
		telegramID = gofakeit.UUID()
		prefix     = "user:telegram_id:"
		testUser   = user.User{
			ID:         1,
			UUID:       gofakeit.UUID(),
			TelegramID: telegramID,
			Username:   gofakeit.Username(),
			FirstName:  gofakeit.FirstName(),
			LastName:   gofakeit.LastName(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		testIterator = map[string]any{
			prefix + telegramID: testUser,
		}
	)

	tests := []struct {
		name                 string
		mockIteratorBehavior func(m *servicemocks.IIterator)
		mockLoggerBehavior   func(m *loggermocks.ILogger)
		want                 want
	}{
		{
			name: "ok",
			mockIteratorBehavior: func(m *servicemocks.IIterator) {
				m.On("Execute", mock.Anything).Return(testIterator, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[iterator for show data in bigcache] execute handler")
			},
			want: want{
				statusCode: http.StatusOK,
				status:     true,
				message:    "success",
				error:      "",
			},
		},
		{
			name: "service error",
			mockIteratorBehavior: func(m *servicemocks.IIterator) {
				m.On("Execute", mock.Anything).Return(nil, errors.New("service error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[iterator for show data in bigcache] execute handler")
				m.On("Error", "failed to show data in bigcache", "error", mock.Anything)
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				status:     false,
				message:    "failed to show data in bigcache",
				error:      "service error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockIterator := servicemocks.NewIIterator(t)
			mockLogger := loggermocks.NewILogger(t)

			if test.mockIteratorBehavior != nil {
				test.mockIteratorBehavior(mockIterator)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}

			app := fiber.New()

			bcs := &bigcacheservice.Service{
				Iterator: mockIterator,
			}

			iterator := New(bcs, mockLogger)

			app.Get("/v1/bigcache/info", iterator.Execute)

			req := httptest.NewRequest(
				fiber.MethodGet,
				"/v1/bigcache/info",
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			}(resp.Body)

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			switch test.name {
			case "ok":
				var result response.Response[map[string]any]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.True(t, result.Status)
				assert.Equal(t, test.want.message, result.Message)
				assert.Empty(t, result.Error)
				assert.NotNil(t, result.Data)
			default:
				var result response.Response[any]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.False(t, result.Status)
				assert.Equal(t, test.want.message, result.Message)
				assert.NotEmpty(t, result.Error)
				assert.Nil(t, result.Data)
			}

			mockIterator.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
