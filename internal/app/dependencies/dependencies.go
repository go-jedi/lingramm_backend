package dependencies

import (
	"context"

	"github.com/go-jedi/lingramm_backend/config"
	undeletefileachievementcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_achievement_cleaner"
	undeletefileawardcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_award_cleaner"
	undeletefileclientcleaner "github.com/go-jedi/lingramm_backend/internal/adapter/cron/jobs/v1/un_delete_file_client_cleaner"
	achievementhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/achievement"
	adminhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/admin"
	authhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/auth"
	bigcachehandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/bigcache"
	experiencepointhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/experience_point"
	clientassetshandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/file_server/client_assets"
	internalcurrencyhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/internal_currency"
	localizedtexthandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/localized_text"
	notificationhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/notification"
	subscriptionhandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/subscription"
	userstatshandler "github.com/go-jedi/lingramm_backend/internal/adapter/http/handlers/v1/user_stats"
	notificationwebsockethandler "github.com/go-jedi/lingramm_backend/internal/adapter/websocket/handlers/v1/notification"
	"github.com/go-jedi/lingramm_backend/internal/middleware"
	achievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/achievement"
	adminrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/admin"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	achievementassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/achievement_assets"
	awardassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/award_assets"
	clientassetsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/file_server/client_assets"
	internalcurrencyrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency"
	levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"
	localizedtextepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/localized_text"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	subscriptionrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	achievementservice "github.com/go-jedi/lingramm_backend/internal/service/v1/achievement"
	adminservice "github.com/go-jedi/lingramm_backend/internal/service/v1/admin"
	authservice "github.com/go-jedi/lingramm_backend/internal/service/v1/auth"
	bigcacheservice "github.com/go-jedi/lingramm_backend/internal/service/v1/bigcache"
	experiencepointservice "github.com/go-jedi/lingramm_backend/internal/service/v1/experience_point"
	clientassetsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/file_server/client_assets"
	internalcurrencyservice "github.com/go-jedi/lingramm_backend/internal/service/v1/internal_currency"
	localizedtextservice "github.com/go-jedi/lingramm_backend/internal/service/v1/localized_text"
	notificationservice "github.com/go-jedi/lingramm_backend/internal/service/v1/notification"
	subscriptionservice "github.com/go-jedi/lingramm_backend/internal/service/v1/subscription"
	userstatsservice "github.com/go-jedi/lingramm_backend/internal/service/v1/user_stats"
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

	// achievement.
	achievementRepository *achievementrepository.Repository
	achievementService    *achievementservice.Service
	achievementHandler    *achievementhandler.Handler

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

	// level.
	levelRepository *levelrepository.Repository

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
	_ = d.AuthHandler()
	_ = d.UserStatsHandler()
	_ = d.ClientAssetsHandler()
	_ = d.BigCacheHandler()
	_ = d.InternalCurrencyHandler()
	_ = d.AchievementHandler()
	_ = d.LocalizedTextHandler()
	_ = d.NotificationHandler()
	_ = d.SubscriptionHandler()
	_ = d.ExperiencePointHandler()
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
}
