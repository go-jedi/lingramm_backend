package check

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingramm_backend/internal/domain/auth"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/auth"
	servicemocks "github.com/go-jedi/lingramm_backend/internal/service/auth/check/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	validatormocks "github.com/go-jedi/lingramm_backend/pkg/validator/mocks"
	"github.com/gofiber/fiber/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
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
		ctx        = context.TODO()
		telegramID = gofakeit.UUID()
		token      = gofakeit.UUID()
		dto        = auth.CheckDTO{
			TelegramID: telegramID,
			Token:      token,
		}
		testResponse = auth.CheckResponse{
			TelegramID: telegramID,
			Token:      token,
			ExpAt:      gofakeit.Date(),
		}
	)

	tests := []struct {
		name                  string
		mockLoggerBehavior    func(m *loggermocks.ILogger)
		mockValidatorBehavior func(m *validatormocks.IValidator)
		mockCheckBehavior     func(m *servicemocks.ICheck)
		in                    in
		want                  want
	}{
		{
			name: "ok",
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[check user token] execute handler")
			},
			mockValidatorBehavior: func(m *validatormocks.IValidator) {
				m.On("Struct", dto).Return(nil)
			},
			mockCheckBehavior: func(m *servicemocks.ICheck) {
				m.On("Execute", ctx, dto).Return(testResponse, nil)
			},
			in: in{
				body: dto,
			},
			want: want{
				statusCode: http.StatusOK,
				response:   testResponse,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockLogger := loggermocks.NewILogger(t)
			mockValidator := validatormocks.NewIValidator(t)
			mockCheck := servicemocks.NewICheck(t)

			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}
			if test.mockValidatorBehavior != nil {
				test.mockValidatorBehavior(mockValidator)
			}
			if test.mockCheckBehavior != nil {
				test.mockCheckBehavior(mockCheck)
			}

			app := fiber.New()

			as := &authservice.Service{
				Check: mockCheck,
			}

			check := New(as, mockLogger, mockValidator)

			app.Post("/v1/auth/check", check.Execute)

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
				"/v1/auth/check",
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
				var result response.Response[auth.CheckResponse]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.NotNil(t, result)
			default:
				var result response.Response[any]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.Nil(t, result)
			}

			mockLogger.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
			mockCheck.AssertExpectations(t)
		})
	}
}
