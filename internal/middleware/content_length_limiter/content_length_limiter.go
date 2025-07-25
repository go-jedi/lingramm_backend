package contentlengthlimiter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeValue  = "multipart/form-data"
)

var (
	ErrUnsupportedContentType = errors.New("unsupported content type")
	ErrMissingContentLength   = errors.New("missing content length")
	ErrFileTooLarge           = errors.New("file too large")
)

type Middleware struct {
	maxBodySize int
}

func New(maxBodySize int) *Middleware {
	return &Middleware{
		maxBodySize: maxBodySize,
	}
}

func (m *Middleware) ContentLengthLimiterMiddleware(c fiber.Ctx) error {
	contentType := c.Get(contentTypeHeader)
	if !strings.HasPrefix(contentType, contentTypeValue) {
		c.Status(fiber.StatusUnsupportedMediaType)
		return c.JSON(response.New[any](false, "Content-Type must be multipart/form-data", ErrUnsupportedContentType.Error(), nil))
	}

	contentLength := c.Request().Header.ContentLength()
	if contentLength == -1 {
		c.Status(fiber.StatusLengthRequired)
		return c.JSON(response.New[any](false, "Content-Length header is required", ErrMissingContentLength.Error(), nil))
	}

	if contentLength > m.maxBodySize {
		c.Status(fiber.StatusRequestEntityTooLarge)
		return c.JSON(response.New[any](false, fmt.Sprintf("uploaded file exceeds limit (%d bytes)", m.maxBodySize), ErrFileTooLarge.Error(), nil))
	}

	return c.Next()
}
