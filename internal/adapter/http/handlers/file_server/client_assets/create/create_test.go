package create

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	clientassets "github.com/go-jedi/lingvogramm_backend/internal/domain/file_server/client_assets"
	clientassetsservice "github.com/go-jedi/lingvogramm_backend/internal/service/file_server/client_assets"
	servicemocks "github.com/go-jedi/lingvogramm_backend/internal/service/file_server/client_assets/create/mocks"
	loggermocks "github.com/go-jedi/lingvogramm_backend/pkg/logger/mocks"
	"github.com/go-jedi/lingvogramm_backend/pkg/response"
	validatormocks "github.com/go-jedi/lingvogramm_backend/pkg/validator/mocks"
	"github.com/gofiber/fiber/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMultipartFormFile(t *testing.T, fieldName, filename, contentType string, content []byte) (*bytes.Buffer, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	assert.NoError(t, err)

	_, err = part.Write(content)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	return body, writer.FormDataContentType()
}

func TestExecute(t *testing.T) {
	type in struct {
		filename    string
		contentType string
		content     []byte
	}

	type want struct {
		statusCode int
		response   response.Response[any]
	}

	var (
		testClientAsset = clientassets.ClientAssets{
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
		}
	)

	tests := []struct {
		name               string
		mockCreateBehavior func(m *servicemocks.ICreate)
		mockLoggerBehavior func(m *loggermocks.ILogger)
		in                 *in
		want               want
	}{
		{
			name: "ok",
			mockCreateBehavior: func(m *servicemocks.ICreate) {
				m.On("Execute", mock.Anything, mock.Anything).Return(testClientAsset, nil)
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute handler")
			},
			in: &in{
				filename:    "image.png",
				contentType: "image/png",
				content:     []byte("dummy"),
			},
			want: want{
				statusCode: fiber.StatusOK,
				response:   *response.New[any](true, "success", "", testClientAsset),
			},
		},
		{
			name:               "missing file",
			mockCreateBehavior: func(m *servicemocks.ICreate) {},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute handler")
				m.On("Error", "failed to get the first file for the provided form key", "error", mock.Anything)
			},
			in: nil,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   *response.New[any](false, "failed to get the first file for the provided form key", "request Content-Type has bad boundary or is not multipart/form-data", nil),
			},
		},
		{
			name:               "unsupported file type",
			mockCreateBehavior: func(m *servicemocks.ICreate) {},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute handler")
				m.On("Error", "unsupported file type: application/json", "error")
			},
			in: &in{
				filename:    "data.json",
				contentType: "application/json",
				content:     []byte(`{"key": "value"}`),
			},
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   *response.New[any](false, "unsupported file type", "unsupported file format: application/json", nil),
			},
		},
		{
			name: "internal service error",
			mockCreateBehavior: func(m *servicemocks.ICreate) {
				m.On("Execute", mock.Anything, mock.Anything).Return(clientassets.ClientAssets{}, fmt.Errorf("db error"))
			},
			mockLoggerBehavior: func(m *loggermocks.ILogger) {
				m.On("Debug", "[create a client assets] execute handler")
				m.On("Error", "failed to create a client assets", "error", mock.Anything)
			},
			in: &in{
				filename:    "image.png",
				contentType: "image/png",
				content:     []byte("dummy"),
			},
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response: *response.New[any](
					false,
					"failed to create a client assets",
					"db error",
					nil,
				),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCreate := servicemocks.NewICreate(t)
			mockLogger := loggermocks.NewILogger(t)
			mockValidator := validatormocks.NewIValidator(t)

			if test.mockCreateBehavior != nil {
				test.mockCreateBehavior(mockCreate)
			}
			if test.mockLoggerBehavior != nil {
				test.mockLoggerBehavior(mockLogger)
			}

			app := fiber.New()

			cas := &clientassetsservice.Service{
				Create: mockCreate,
			}

			create := New(cas, mockLogger, mockValidator)

			app.Post("/v1/fs/client_assets", create.Execute)

			var req *http.Request
			if test.in != nil {
				body, contentType := createMultipartFormFile(t, "file", test.in.filename, test.in.contentType, test.in.content)
				req = httptest.NewRequest("POST", "/v1/fs/client_assets", body)
				req.Header.Set("Content-Type", contentType)
			} else {
				req = httptest.NewRequest("POST", "/v1/fs/client_assets", nil)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Error(err)
				}
			}(resp.Body)

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			switch test.name {
			case "ok":
				var result response.Response[clientassets.ClientAssets]
				err = jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.Equal(t, testClientAsset, result.Data)
			default:
				var result response.Response[any]
				err = jsoniter.Unmarshal(respBody, &result)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, result)
			}

			mockCreate.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}
