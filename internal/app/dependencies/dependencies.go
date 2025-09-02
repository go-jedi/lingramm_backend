package dependencies

import (
	"context"

	"github.com/go-jedi/lingramm_backend/config"
	leaderboardweeksprocessbatch "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/leaderboard_weeks_process_batch"
	undeletefileachievementcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_achievement_cleaner"
	undeletefileawardcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_award_cleaner"
	undeletefileclientcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_client_cleaner"
	achievementhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/achievement"
	adminhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/admin"
	authhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth"
	bigcachehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/bigcache"
	dailytaskhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/daily_task"
	eventhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/event"
	eventtypehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/event_type"
	experiencepointhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/experience_point"
	clientassetshandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/file_server/client_assets"
	internalcurrencyhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/internal_currency"
	localizedtexthandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/localized_text"
	notificationhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/notification"
	studiedlanguagehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/studied_language"
	subscriptionhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/subscription"
	userhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user"
	userachievementhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_achievement"
	userdailytaskhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_daily_task"
	userstatshandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_stats"
	userstudiedlanguagehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_studied_language"
	notificationwebsockethandler "github.com/go-jedi/lingramm_backend/internal/adapter/websocket/handlers/v1/notification"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	achievementtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement_type"
	adminrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/admin"
	dailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/daily_task"
	eventtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
	awardassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/award_assets"
	clientassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets"
	internalcurrencyrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency"
	levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"
	localizedtextepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	studiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/studied_language"
	subscriptionrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userachievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement"
	userdailytaskrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_daily_task"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	userstudiedlanguagerepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_studied_language"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
	bigcacheservice "github.com/go-jedi/lingramm_backend/internal/service/v1/bigcache"
	dailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/daily_task"
	eventservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event"
	eventtypeservice "github.com/go-jedi/lingramm_backend/internal/service/v1/event_type"
	experiencepointservice "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	internalcurrencyservice "github.com/go-jedi/lingramm_backend/internal/service/v1/internal_currency"
	localizedtextservice "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	studiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/studied_language"
	subscriptionservice "github.com/go-jedi/lingramm_backend/internal/service/v1/subscription"
	userservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user"
	userachievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_achievement"
	userdailytaskservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_daily_task"
	userstatsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats"
	userstudiedlanguageservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_studied_language"
	bigcachepkg "github.com/go-jedi/lingramm_backend/pkg/bigcache"
	fileserver "github.com/go-jedi/lingramm_backend/pkg/file_server"
	"github.com/go-jedi/lingramm_backend/pkg/jwt"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/rabbitmq"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/go-jedi/lingramm_backend/pkg/uuid"
	"github.com/go-jedi/lingramm_backend/pkg/validator"
	wsmanager "github.com/go-jedi/lingramm_backend/pkg/ws_manager"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	cfg        config.Config
	app        *fiber.App
	logger     *logger.Logger
	validator  *validator.Validator
	uuid       *uuid.UUID
	jwt        *jwt.JWT
	middleware *middleware.Middleware
	rabbitMQ   *rabbitmq.RabbitMQ
	postgres   *postgres.Postgres
	redis      *redis.Redis
	bigCache   *bigcachepkg.BigCache
	wsManager  *wsmanager.WSManager
	fileServer *fileserver.FileServer

	// auth.
	authService *authservice.Service
	authHandler *authhandler.Handler

	// user.
	userRepository *userrepository.Repository
	userService    *userservice.Service
	userHandler    *userhandler.Handler

	// user stats.
	userStatsRepository *userstatsrepository.Repository
	userStatsService    *userstatsservice.Service
	userStatsHandler    *userstatshandler.Handler

	// client assets.
	clientAssetsRepository *clientassetsrepository.Repository
	clientAssetsService    *clientassetsservice.Service
	clientAssetsHandler    *clientassetshandler.Handler

	// bigcache.
	bigCacheService *bigcacheservice.Service
	bigCacheHandler *bigcachehandler.Handler

	// internal currency.
	internalCurrencyRepository *internalcurrencyrepository.Repository
	internalCurrencyService    *internalcurrencyservice.Service
	internalCurrencyHandler    *internalcurrencyhandler.Handler

	// achievement assets.
	achievementAssetsRepository *achievementassetsrepository.Repository

	// award assets.
	awardAssetsRepository *awardassetsrepository.Repository

	// achievement type.
	achievementTypeRepository *achievementtyperepository.Repository

	// achievement.
	achievementRepository *achievementrepository.Repository
	achievementService    *achievementservice.Service
	achievementHandler    *achievementhandler.Handler

	// user achievement.
	userAchievementRepository *userachievementrepository.Repository
	userAchievementService    *userachievementservice.Service
	userAchievementHandler    *userachievementhandler.Handler

	// studied language.
	studiedLanguageRepository *studiedlanguagerepository.Repository
	studiedLanguageService    *studiedlanguageservice.Service
	studiedLanguageHandler    *studiedlanguagehandler.Handler

	// user studied language.
	userStudiedLanguageRepository *userstudiedlanguagerepository.Repository
	userStudiedLanguageService    *userstudiedlanguageservice.Service
	userStudiedLanguageHandler    *userstudiedlanguagehandler.Handler

	// localized text.
	localizedTextRepository *localizedtextepository.Repository
	localizedTextService    *localizedtextservice.Service
	localizedTextHandler    *localizedtexthandler.Handler

	// notification.
	notificationRepository *notificationrepository.Repository
	notificationService    *notificationservice.Service
	notificationHandler    *notificationhandler.Handler

	// subscription.
	subscriptionRepository *subscriptionrepository.Repository
	subscriptionService    *subscriptionservice.Service
	subscriptionHandler    *subscriptionhandler.Handler

	// experience point.
	experiencePointRepository *experiencepointrepository.Repository
	experiencePointService    *experiencepointservice.Service
	experiencePointHandler    *experiencepointhandler.Handler

	// event.
	eventService *eventservice.Service
	eventHandler *eventhandler.Handler

	// level.
	levelRepository *levelrepository.Repository

	// event type.
	eventTypeRepository *eventtyperepository.Repository
	eventTypeService    *eventtypeservice.Service
	eventTypeHandler    *eventtypehandler.Handler

	// daily task.
	dailyTaskRepository *dailytaskrepository.Repository
	dailyTaskService    *dailytaskservice.Service
	dailyTaskHandler    *dailytaskhandler.Handler

	// user daily task.
	userDailyTaskRepository *userdailytaskrepository.Repository
	userDailyTaskService    *userdailytaskservice.Service
	userDailyTaskHandler    *userdailytaskhandler.Handler

	// admin.
	adminRepository *adminrepository.Repository
	adminService    *adminservice.Service
	adminHandler    *adminhandler.Handler

	// websocket.
	notificationWebSocketHandler *notificationwebsockethandler.Handler

	// cron.
	unDeleteFileAchievementCleaner *undeletefileachievementcleaner.UnDeleteFileAchievementCleaner
	unDeleteFileAwardCleaner       *undeletefileawardcleaner.UnDeleteFileAwardCleaner
	unDeleteFileClientCleaner      *undeletefileclientcleaner.UnDeleteFileClientCleaner
	leaderboardWeeksProcessBatch   *leaderboardweeksprocessbatch.LeaderboardWeeksProcessBatch
}

func New(
	ctx context.Context,
	cfg config.Config,
	app *fiber.App,
	logger *logger.Logger,
	validator *validator.Validator,
	uuid *uuid.UUID,
	jwt *jwt.JWT,
	rabbitMQ *rabbitmq.RabbitMQ,
	postgres *postgres.Postgres,
	redis *redis.Redis,
	bigCache *bigcachepkg.BigCache,
	wsManager *wsmanager.WSManager,
	fileServer *fileserver.FileServer,
) *Dependencies {
	d := &Dependencies{
		cfg:        cfg,
		app:        app,
		logger:     logger,
		validator:  validator,
		uuid:       uuid,
		jwt:        jwt,
		rabbitMQ:   rabbitMQ,
		postgres:   postgres,
		redis:      redis,
		bigCache:   bigCache,
		wsManager:  wsManager,
		fileServer: fileServer,
	}

	d.initMiddleware()
	d.initHandler()
	d.initWebSocket()
	d.initCron(ctx)

	return d
}

// initMiddleware initialize middlewares.
func (d *Dependencies) initMiddleware() {
	d.middleware = middleware.New(
		d.cfg.Middleware,
		d.AdminService(),
		d.jwt,
		d.redis,
	)
}

// initHandler initialize handlers.
func (d *Dependencies) initHandler() {
	_ = d.UserHandler()
	_ = d.AuthHandler()
	_ = d.UserStatsHandler()
	_ = d.ClientAssetsHandler()
	_ = d.BigCacheHandler()
	_ = d.InternalCurrencyHandler()
	_ = d.AchievementHandler()
	_ = d.UserAchievementHandler()
	_ = d.StudiedLanguageHandler()
	_ = d.UserStudiedLanguageHandler()
	_ = d.LocalizedTextHandler()
	_ = d.NotificationHandler()
	_ = d.SubscriptionHandler()
	_ = d.ExperiencePointHandler()
	_ = d.EventHandler()
	_ = d.EventTypeHandler()
	_ = d.DailyTaskHandler()
	_ = d.UserDailyTaskHandler()
	_ = d.AdminHandler()
}

// initWebSocket initialize web sockets.
func (d *Dependencies) initWebSocket() {
	_ = d.NotificationWebSocket()
}

// initCron initialize cron.
func (d *Dependencies) initCron(ctx context.Context) {
	_ = d.UnDeleteFileAchievementCleanerCron(ctx)
	_ = d.UnDeleteFileAwardCleanerCron(ctx)
	_ = d.UnDeleteFileClientCleanerCron(ctx)
	_ = d.LeaderboardWeeksProcessBatchCron(ctx)
}
