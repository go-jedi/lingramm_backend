package userstats

import (
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userachievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	ensurestreakdaysincrementtoday "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats/ensure_streak_days_increment_today"
	getlevelbytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats/get_level_by_telegram_id"
	getlevelinfobytelegramid "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats/get_level_info_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/rabbitmq"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Service struct {
	EnsureStreakDaysIncrementToday ensurestreakdaysincrementtoday.IEnsureStreakDaysIncrementToday
	GetLevelByTelegramID           getlevelbytelegramid.IGetLevelByTelegramID
	GetLevelInfoByTelegramID       getlevelinfobytelegramid.IGetLevelInfoByTelegramID
}

func New(
	userStatsRepository *userstatsrepository.Repository,
	userRepository *userrepository.Repository,
	userAchievementRepository *userachievementrepository.Repository,
	notificationRepository *notificationrepository.Repository,
	logger logger.ILogger,
	rabbitMQ *rabbitmq.RabbitMQ,
	postgres *postgres.Postgres,
	redis *redis.Redis,
) *Service {
	return &Service{
		EnsureStreakDaysIncrementToday: ensurestreakdaysincrementtoday.New(
			userStatsRepository,
			userRepository,
			userAchievementRepository,
			notificationRepository,
			logger,
			rabbitMQ,
			postgres,
			redis,
		),
		GetLevelByTelegramID:     getlevelbytelegramid.New(userStatsRepository, userRepository, logger, postgres),
		GetLevelInfoByTelegramID: getlevelinfobytelegramid.New(userStatsRepository, userRepository, logger, postgres),
	}
}
