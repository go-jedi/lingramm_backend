package admin

import (
	addadminuserhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/admin/add_admin_user"
	"github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/admin/exists_by_telegram_id"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	addAdminUser       *addadminuserhandler.AddAdminUser
	existsByTelegramID *existsbytelegramid.ExistsByTelegramID
}

func New(
	adminService *adminservice.Service,
	app *fiber.App,
	logger logger.ILogger,
) *Handler {
	h := &Handler{
		addAdminUser:       addadminuserhandler.New(adminService, logger),
		existsByTelegramID: existsbytelegramid.New(adminService, logger),
	}

	h.initRoutes(app)

	return h
}

func (h *Handler) initRoutes(app *fiber.App) {
	api := app.Group("/v1/admin")
	{
		api.Get("/add/:telegramID", h.addAdminUser.Execute)
		api.Get("/exists/:telegramID", h.existsByTelegramID.Execute)
	}
}
