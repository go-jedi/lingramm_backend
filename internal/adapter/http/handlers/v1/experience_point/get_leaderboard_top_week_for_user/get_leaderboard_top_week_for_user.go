package getleaderboardtopweekforuser

import (
	"context"
	"time"

	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	experiencepointservice "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/response"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

const timeout = 5 * time.Second

type GetLeaderboardTopWeekForUser struct {
	experiencePointService *experiencepointservice.Service
	logger                 logger.ILogger
	validator              validator.IValidator
}

func New(
	experiencePointService *experiencepointservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *GetLeaderboardTopWeekForUser {
	return &GetLeaderboardTopWeekForUser{
		experiencePointService: experiencePointService,
		logger:                 logger,
		validator:              validator,
	}
}

// Execute returns weekly XP leaderboard scoped around a user.
// @Summary Get weekly leaderboard for user (XP)
// @Description Returns the weekly XP leaderboard centered around the specified user (by Telegram ID).
// @Description Rules:
// @Description • `limit` is required, must be > 0 and ≤ 30
// @Description • `telegram_id` is required
// @Description • `tz` is required and must be `Europe/Moscow`
// @Tags Experience point
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param payload body experiencepoint.GetLeaderboardTopWeekForUserDTO true "Leaderboard request for user"
// @Success 200 {object} experiencepoint.GetLeaderboardTopWeekForUserSwaggerResponse "Successful response"
// @Failure 400 {object} experiencepoint.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} experiencepoint.ErrorSwaggerResponse "Internal server error"
// @Router /v1/experience_point/leaderboard/week_top/user [post]
func (h *GetLeaderboardTopWeekForUser) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get leaderboard top week for user] execute handler")

	var dto experiencepoint.GetLeaderboardTopWeekForUserDTO
	if err := c.Bind().Body(&dto); err != nil {
		h.logger.Error("failed to bind body", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to bind body", err.Error(), nil))
	}

	if err := h.validator.StructCtx(c.RequestCtx(), dto); err != nil {
		h.logger.Error("failed to validate struct", "error", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.New[any](false, "failed to validate struct", err.Error(), nil))
	}

	ctxTimeout, cancel := context.WithTimeout(c.RequestCtx(), timeout)
	defer cancel()

	result, err := h.experiencePointService.GetLeaderboardTopWeekForUser.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to get leaderboard top week for user", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get leaderboard top week for user", err.Error(), nil))
	}

	return c.JSON(response.New[[]experiencepoint.GetLeaderboardTopWeekForUserResponse](true, "success", "", result))
}
