package event

import (
	eventtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	internalcurrency "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency"
	levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userachievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement"
	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	createevents "github.com/go-jedi/lingramm_backend/internal/service/v1/event/create_events"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/rabbitmq"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
)

type Service struct {
	CreateEvents createevents.ICreateEvents
}

func New(
	experiencePointRepository *experiencepointrepository.Repository,
	userRepository *userrepository.Repository,
	userStatsRepository *userstatsrepository.Repository,
	eventTypeRepository *eventtyperepository.Repository,
	levelRepository *levelrepository.Repository,
	internalCurrencyRepository *internalcurrency.Repository,
	userAchievementRepository *userachievementrepository.Repository,
	userDailyTaskRepository *userdailytaskrepository.Repository,
	notificationRepository *notificationrepository.Repository,
	logger logger.ILogger,
	rabbitMQ *rabbitmq.RabbitMQ,
	postgres *postgres.Postgres,
	redis *redis.Redis,
) *Service {
	return &Service{
		CreateEvents: createevents.New(
			experiencePointRepository,
			userRepository,
			userStatsRepository,
			eventTypeRepository,
			levelRepository,
			internalCurrencyRepository,
			userAchievementRepository,
			userDailyTaskRepository,
			notificationRepository,
			logger,
			rabbitMQ,
			postgres,
			redis,
		),
	}
}
