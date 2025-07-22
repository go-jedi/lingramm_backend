package existsbytelegramid

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	existsbytelegramidservicemocks "github.com/go-jedi/lingramm_backend/internal/service/v1/admin/exists_by_telegram_id/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		telegramID string
	}

	type want struct {
		statusCode int
		response   interface{}
	}

	var (
		telegramID = gofakeit.UUID()
		exists     = gofakeit.Bool()
	)

	tests := []struct {
		name                   string
		mockExistsByTelegramID func(m *existsbytelegramidservicemocks.IExistsByTelegramID)
		mockLoggerBehavior     func(m *loggermocks.ILogger)
		in                     in
		want                   want
	}{
		{
			name: "ok",
			mockExistsByTelegramID: func(m *existsbytelegramidservicemocks.IExistsByTelegramID) {
				m.On("Execute", mock.Anything, telegramID).Return(exists, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check admin exists by telegram id] execute handler")
			},
			in: in{
				telegramID: telegramID,
			},
			want: want{
				statusCode: fiber.StatusOK,
				response:   response.New[bool](true, "success", "", exists),
			},
		},
		{
			name: "service_error",
			mockExistsByTelegramID: func(m *existsbytelegramidservicemocks.IExistsByTelegramID) {
				m.On("Execute", mock.Anything, telegramID).Return(false, errors.New("service error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check admin exists by telegram id] execute handler")
				m.On("Error", "failed to exists admin by telegram id", "error", errors.New("service error"))
			},
			in: in{
				telegramID: telegramID,
			},
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   response.New[any](false, "failed to exists admin by telegram id", "service error", nil),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockExistsByTelegramID := existsbytelegramidservicemocks.NewIExistsByTelegramID(t)
			mockLogger := loggermocks.NewILogger(t)

			if test.mockExistsByTelegramID != nil {
				test.mockExistsByTelegramID(mockExistsByTelegramID)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}

			app := fiber.New()

			as := &adminservice.Service{
				ExistsByTelegramID: mockExistsByTelegramID,
			}

			existsByTelegramID := New(as, mockLogger)

			app.Get("/v1/admin/exists/:telegramID", existsByTelegramID.Execute)
			req := httptest.NewRequest(fiber.MethodGet, "/v1/admin/exists/"+test.in.telegramID, nil)
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
				var result response.Response[bool]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.NotNil(t, result.Data)
			default:
				var result response.Response[any]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.Nil(t, result.Data)
			}

			mockExistsByTelegramID.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
