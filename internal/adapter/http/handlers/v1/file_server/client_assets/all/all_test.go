package all

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	servicemocks "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets/all/mocks"
	loggermocks "github.com/go-jedi/lingramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	validatormocks "github.com/go-jedi/lingramm_backend/pkg/validator/mocks"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	type want struct {
		statusCode int
		response   interface{}
	}

	var (
		testClientAssets = []clientassets.ClientAssets{
			{
				ID:             gofakeit.Int64(),
				NameFile:       gofakeit.Name(),
				ServerPathFile: gofakeit.URL(),
				ClientPathFile: gofakeit.URL(),
				Extension:      gofakeit.FileExtension(),
				Quality:        gofakeit.IntRange(1, 100),
				OldNameFile:    gofakeit.Name(),
				OldExtension:   gofakeit.FileExtension(),
				CreatedAt:      gofakeit.Date(),
				UpdatedAt:      gofakeit.Date(),
			},
		}
	)

	tests := []struct {
		name               string
		mockAllBehavior    func(m *servicemocks.IAll)
		mockLoggerBehavior func(m *loggermocks.ILogger)
		want               want
	}{
		{
			name: "ok",
			mockAllBehavior: func(m *servicemocks.IAll) {
				m.On("Execute", mock.Anything).Return(testClientAssets, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get all client assets] execute handler")
			},
			want: want{
				statusCode: fiber.StatusOK,
				response:   response.New[[]clientassets.ClientAssets](true, "success", "", testClientAssets),
			},
		},
		{
			name: "service error",
			mockAllBehavior: func(m *servicemocks.IAll) {
				m.On("Execute", mock.Anything).Return(nil, errors.New("service error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[get all client assets] execute handler")
				m.On("Error", "failed to get all client assets", "error", mock.Anything)
			},
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   response.New[any](false, "failed to get all client assets", mock.Anything, nil),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAll := servicemocks.NewIAll(t)
			mockLogger := loggermocks.NewILogger(t)
			mockValidator := validatormocks.NewIValidator(t)

			if test.mockAllBehavior != nil {
				test.mockAllBehavior(mockAll)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}

			app := fiber.New()

			cas := &clientassetsservice.Service{
				All: mockAll,
			}

			all := New(cas, mockLogger, mockValidator)

			app.Get("/v1/auth/signin/all", all.Execute)

			req := httptest.NewRequest(
				fiber.MethodGet,
				"/v1/auth/signin/all",
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

			mockAll.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}
