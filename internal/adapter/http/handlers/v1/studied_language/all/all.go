package all

import (
	"context"
	"time"

	studiedlanguage "github.com/go-jedi/lingramm_backend/internal/domain/studied_language"
	studiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/studied_language"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type All struct {
	studiedLanguageService *studiedlanguageservice.Service
	logger                 logger.ILogger
}

func New(
	studiedLanguageService *studiedlanguageservice.Service,
	logger logger.ILogger,
) *All {
	return &All{
		studiedLanguageService: studiedLanguageService,
		logger:                 logger,
	}
}

func (h *All) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get all studied languages] execute handler")

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.studiedLanguageService.All.Execute(ctxTimeout)
	if err != nil {
		h.logger.Error("failed to get all studied languages", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get all studied languages", err.Error(), nil))
	}

	return c.JSON(response.New[[]studiedlanguage.StudiedLanguage](true, "success", "", result))
}
