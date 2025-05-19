package addadminuser

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/admin"
	adminservice "github.com/go-jedi/lingvogramm_backend/internal/service/admin"
	addadminuserservicemocks "github.com/go-jedi/lingvogramm_backend/internal/service/admin/add_admin_user/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
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
		testAdmin  = admin.Admin{
			ID:         gofakeit.Int64(),
			TelegramID: telegramID,
			CreatedAt:  gofakeit.Date(),
			UpdatedAt:  gofakeit.Date(),
		}
	)

	tests := []struct {
		name                     string
		mockAddAdminUserBehavior func(m *addadminuserservicemocks.IAddAdminUser)
		mockLoggerBehavior       func(m *loggermocks.ILogger)
		in                       in
		want                     want
	}{
		{
			name: "ok",
			mockAddAdminUserBehavior: func(m *addadminuserservicemocks.IAddAdminUser) {
				m.On("Execute", mock.Anything, telegramID).Return(testAdmin, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute handler")
			},
			in: in{
				telegramID: telegramID,
			},
			want: want{
				statusCode: fiber.StatusOK,
				response:   response.New[admin.Admin](true, "success", "", testAdmin),
			},
		},
		{
			name: "service_error",
			mockAddAdminUserBehavior: func(m *addadminuserservicemocks.IAddAdminUser) {
				m.On("Execute", mock.Anything, telegramID).Return(admin.Admin{}, errors.New("service error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[add a new admin user] execute handler")
				m.On("Error", "failed to add admin user", "error", errors.New("service error"))
			},
			in: in{
				telegramID: telegramID,
			},
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   response.New[any](false, "failed to add admin user", "service error", nil),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAddAdminUser := addadminuserservicemocks.NewIAddAdminUser(t)
			mockLogger := loggermocks.NewILogger(t)

			if test.mockAddAdminUserBehavior != nil {
				test.mockAddAdminUserBehavior(mockAddAdminUser)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}

			app := fiber.New()

			as := &adminservice.Service{
				AddAdminUser: mockAddAdminUser,
			}

			addAdminUser := New(as, mockLogger)

			app.Get("/v1/admin/add/:telegramID", addAdminUser.Execute)

			req := httptest.NewRequest(fiber.MethodGet, "/v1/admin/add/"+test.in.telegramID, nil)
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
				var result response.Response[admin.Admin]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.NotNil(t, result.Data)
			default:
				var result response.Response[any]
				err := jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.Nil(t, result.Data)
			}

			mockAddAdminUser.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
