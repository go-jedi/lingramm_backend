package signin

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/auth"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/user"
	authservice "github.com/go-jedi/lingvogramm_backend/internal/service/auth"
	servicemocks "github.com/go-jedi/lingvogramm_backend/internal/service/auth/sign_in/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
	validatormocks "github.com/go-jedi/lingvogramm_backend/pkg/validator/mocks"
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
		uuid       = gofakeit.UUID()
		telegramID = gofakeit.UUID()
		username   = gofakeit.Username()
		firstname  = gofakeit.FirstName()
		lastname   = gofakeit.LastName()
		createdAt  = time.Now()
		updatedAt  = time.Now()
		dto        = auth.SignInDTO{
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
	)

	tests := []struct {
		name                  string
		mockSignInBehavior    func(m *servicemocks.ISignIn)
		mockLoggerBehavior    func(m *loggermocks.ILogger)
		mockValidatorBehavior func(m *validatormocks.IValidator)
		in                    in
		want                  want
	}{
		{
			name: "ok",
			mockSignInBehavior: func(m *servicemocks.ISignIn) {
				m.On("Execute", mock.Anything, dto).Return(testUser, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute handler")
			},
			mockValidatorBehavior: func(m *validatormocks.IValidator) {
				m.On("Struct", dto).Return(nil)
			},
			in: in{
				body: dto,
			},
			want: want{
				statusCode: fiber.StatusOK,
				response:   response.New[user.User](true, "success", "", testUser),
			},
		},
		{
			name: "bind error",
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute handler")
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
				m.On("Debug", "[sign in user] execute handler")
				m.On("Error", "failed to validate struct", "error", mock.Anything)
			},
			mockValidatorBehavior: func(m *validatormocks.IValidator) {
				m.On("Struct", mock.Anything).Return(errors.New("validation error"))
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
			mockSignInBehavior: func(m *servicemocks.ISignIn) {
				m.On("Execute", mock.Anything, dto).Return(user.User{}, errors.New("service error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[sign in user] execute handler")
				m.On("Error", "failed to sign in user", "error", mock.Anything)
			},
			mockValidatorBehavior: func(m *validatormocks.IValidator) {
				m.On("Struct", dto).Return(nil)
			},
			in: in{
				body: dto,
			},
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   response.New[any](false, "failed to sign in user", "service error", nil),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSignIn := servicemocks.NewISignIn(t)
			mockLogger := loggermocks.NewILogger(t)
			mockValidator := validatormocks.NewIValidator(t)

			if test.mockSignInBehavior != nil {
				test.mockSignInBehavior(mockSignIn)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}
			if test.mockValidatorBehavior != nil {
				test.mockValidatorBehavior(mockValidator)
			}

			app := fiber.New()

			as := &authservice.Service{
				SignIn: mockSignIn,
			}

			signIn := New(as, mockLogger, mockValidator)

			app.Post("/v1/auth/signin", signIn.Execute)

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
				"/v1/auth/signin",
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
				var result response.Response[user.User]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.NotNil(t, result.Data)
			default:
				var result response.Response[any]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.Nil(t, result.Data)
			}

			mockSignIn.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}
