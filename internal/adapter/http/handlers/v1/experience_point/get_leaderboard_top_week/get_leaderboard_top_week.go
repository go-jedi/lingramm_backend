package getleaderboardtopweek

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

type GetLeaderboardTopWeek struct {
	experiencePointService *experiencepointservice.Service
	logger                 logger.ILogger
	validator              validator.IValidator
}

func New(
	experiencePointService *experiencepointservice.Service,
	logger logger.ILogger,
	validator validator.IValidator,
) *GetLeaderboardTopWeek {
	return &GetLeaderboardTopWeek{
		experiencePointService: experiencePointService,
		logger:                 logger,
		validator:              validator,
	}
}

// Execute returns weekly XP leaderboard.
// @Summary Get weekly leaderboard (XP)
// @Description Returns the top users by XP for the current week in the given timezone.
// @Description Rules:
// @Description • `limit` is required, must be > 0 and ≤ 30
// @Description • `tz` is required and must be `Europe/Moscow`
// @Tags Experience point
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer <token>)
// @Param payload body experiencepoint.GetLeaderboardTopWeekDTO true "Leaderboard request"
// @Success 200 {object} experiencepoint.GetLeaderboardTopWeekSwaggerResponse "Successful response"
// @Failure 400 {object} experiencepoint.ErrorSwaggerResponse "Bad request error"
// @Failure 500 {object} experiencepoint.ErrorSwaggerResponse "Internal server error"
// @Router /v1/experience_point/leaderboard/week_top [post]
func (h *GetLeaderboardTopWeek) Execute(c fiber.Ctx) error {
	h.logger.Debug("[get leaderboard top week] execute handler")

	var dto experiencepoint.GetLeaderboardTopWeekDTO
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

	result, err := h.experiencePointService.GetLeaderboardTopWeek.Execute(ctxTimeout, dto)
	if err != nil {
		h.logger.Error("failed to get leaderboard top week", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.New[any](false, "failed to get leaderboard top week", err.Error(), nil))
	}

	return c.JSON(response.New[[]experiencepoint.GetLeaderboardTopWeekResponse](true, "success", "", result))
}
