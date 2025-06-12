package refresh

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/auth"
	servicemocks "github.com/go-jedi/lingramm_backend/internal/service/auth/refresh/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	validatormocks "github.com/go-jedi/lingramm_backend/pkg/validator/mocks"
	"github.com/gofiber/fiber/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type in struct {
		body interface{}
	}

	type want struct {
		statusCode int
		response   interface{}
	}

	var (
		dto = auth.RefreshDTO{
			TelegramID:   gofakeit.UUID(),
			RefreshToken: gofakeit.UUID(),
		}
		testResponse = auth.RefreshResponse{
			AccessToken:  gofakeit.UUID(),
			RefreshToken: gofakeit.UUID(),
			AccessExpAt:  gofakeit.Date(),
			RefreshExpAt: gofakeit.Date(),
		}
	)

	tests := []struct {
		name                  string
		mockRefreshBehavior   func(m *servicemocks.IRefresh)
		mockLoggerBehavior    func(m *loggermocks.ILogger)
		mockValidatorBehavior func(m *validatormocks.IValidator)
		in                    in
		want                  want
	}{
		{
			name: "ok",
			mockRefreshBehavior: func(m *servicemocks.IRefresh) {
				m.On("Execute", mock.Anything, dto).Return(testResponse, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute handler")
			},
			mockValidatorBehavior: func(m *validatormocks.IValidator) {
				m.On("Struct", dto).Return(nil)
			},
			in: in{
				body: dto,
			},
			want: want{
				statusCode: http.StatusOK,
				response:   testResponse,
			},
		},
		{
			name: "bind error",
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute handler")
				m.On("Error", "failed to bind body", "error", mock.Anything)
			},
			in: in{
				body: `"invalid"`,
			},
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   response.New[any](false, "failed to bind body", mock.Anything, nil),
			},
		},
		{
			name: "validation error",
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute handler")
				m.On("Error", "failed to validate struct", "error", mock.Anything)
			},
			mockValidatorBehavior: func(m *validatormocks.IValidator) {
				m.On("Struct", dto).Return(errors.New("some error"))
			},
			in: in{
				body: dto,
			},
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   response.New[any](false, "failed to validate struct", "validation error", nil),
			},
		},
		{
			name: "service error",
			mockRefreshBehavior: func(m *servicemocks.IRefresh) {
				m.On("Execute", mock.Anything, dto).Return(auth.RefreshResponse{}, errors.New("service error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[refresh user token] execute handler")
				m.On("Error", "failed to refresh tokens", "error", mock.Anything)
			},
			mockValidatorBehavior: func(m *validatormocks.IValidator) {
				m.On("Struct", dto).Return(nil)
			},
			in: in{
				body: dto,
			},
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   response.New[any](false, "failed to refresh tokens", "service error", nil),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRefresh := servicemocks.NewIRefresh(t)
			mockLogger := loggermocks.NewILogger(t)
			mockValidator := validatormocks.NewIValidator(t)

			if test.mockRefreshBehavior != nil {
				test.mockRefreshBehavior(mockRefresh)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}
			if test.mockValidatorBehavior != nil {
				test.mockValidatorBehavior(mockValidator)
			}

			app := fiber.New()

			as := &authservice.Service{
				Refresh: mockRefresh,
			}

			refresh := New(as, mockLogger, mockValidator)

			app.Post("/v1/auth/refresh", refresh.Execute)

			var rawData []byte
			var err error

			switch body := test.in.body.(type) {
			case string:
				rawData = []byte(body)
			default:
				rawData, err = jsoniter.Marshal(test.in.body)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(
				fiber.MethodPost,
				"/v1/auth/refresh",
				bytes.NewBuffer(rawData),
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
				var result response.Response[auth.RefreshResponse]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.NotNil(t, result.Data)
			default:
				var result response.Response[any]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.Nil(t, result.Data)
			}

			mockRefresh.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}
