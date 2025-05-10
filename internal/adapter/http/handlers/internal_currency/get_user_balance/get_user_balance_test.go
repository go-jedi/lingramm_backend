package getuserbalance

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	userbalance "github.com/go-jedi/lingvogramm_backend/internal/domain/user_balance"
	internalcurrencyservice "github.com/go-jedi/lingvogramm_backend/internal/service/internal_currency"
	servicemocks "github.com/go-jedi/lingvogramm_backend/internal/service/internal_currency/get_user_balance/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
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
		telegramID      = gofakeit.UUID()
		testUserBalance = userbalance.UserBalance{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			Balance:    decimal.NewFromFloat(0.00),
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
		}
	)

	tests := []struct {
		name                       string
		mockGetUserBalanceBehavior func(m *servicemocks.IGetUserBalance)
		mockLoggerBehavior         func(m *loggermocks.ILogger)
		in                         in
		want                       want
	}{
		{
			name: "ok",
			mockGetUserBalanceBehavior: func(m *servicemocks.IGetUserBalance) {
				m.On("Execute", mock.Anything, telegramID).Return(testUserBalance, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute handler")
			},
			in: in{
				telegramID: telegramID,
			},
			want: want{
				statusCode: fiber.StatusOK,
				response:   response.New[userbalance.UserBalance](true, "success", "", testUserBalance),
			},
		},
		{
			name: "service_error",
			mockGetUserBalanceBehavior: func(m *servicemocks.IGetUserBalance) {
				m.On("Execute", mock.Anything, telegramID).Return(userbalance.UserBalance{}, errors.New("service error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get user balance] execute handler")
				m.On("Error", "failed to get user balance", "error", errors.New("service error"))
			},
			in: in{
				telegramID: telegramID,
			},
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   response.New[any](false, "failed to exists user by telegram id", "service error", nil),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockGetUserBalance := servicemocks.NewIGetUserBalance(t)
			mockLogger := loggermocks.NewILogger(t)

			if test.mockGetUserBalanceBehavior != nil {
				test.mockGetUserBalanceBehavior(mockGetUserBalance)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}

			app := fiber.New()

			ics := &internalcurrencyservice.Service{
				GetUserBalance: mockGetUserBalance,
			}

			getUserBalance := New(ics, mockLogger)

			app.Get("/v1/internal_currency/user/balance/:telegramID", getUserBalance.Execute)

			req := httptest.NewRequest(fiber.MethodGet, "/v1/internal_currency/user/balance/"+test.in.telegramID, nil)
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
				var result response.Response[userbalance.UserBalance]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.NotNil(t, result.Data)
			default:
				var result response.Response[any]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.Nil(t, result.Data)
			}

			mockGetUserBalance.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
