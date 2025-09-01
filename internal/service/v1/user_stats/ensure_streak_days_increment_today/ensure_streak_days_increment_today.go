package ensurestreakdaysincrementtoday

import (
	"context"
	"fmt"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	userachievement "github.com/go-jedi/lingramm_backend/internal/domain/user_achievement"
	notificationrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/notification"
	userrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user"
	userachievementrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_achievement"
	userstatsrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/user_stats"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
	"github.com/go-jedi/lingramm_backend/pkg/postgres"
	"github.com/go-jedi/lingramm_backend/pkg/rabbitmq"
	"github.com/go-jedi/lingramm_backend/pkg/redis"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=IEnsureStreakDaysIncrementToday --output=mocks --case=underscore
type IEnsureStreakDaysIncrementToday interface {
	Execute(ctx context.Context, telegramID string) error
}

type EnsureStreakDaysIncrementToday struct {
	userStatsRepository       *userstatsrepository.Repository
	userRepository            *userrepository.Repository
	userAchievementRepository *userachievementrepository.Repository
	notificationRepository    *notificationrepository.Repository
	logger                    logger.ILogger
	rabbitMQ                  *rabbitmq.RabbitMQ
	postgres                  *postgres.Postgres
	redis                     *redis.Redis
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
) *EnsureStreakDaysIncrementToday {
	return &EnsureStreakDaysIncrementToday{
		userStatsRepository:       userStatsRepository,
		userRepository:            userRepository,
		userAchievementRepository: userAchievementRepository,
		notificationRepository:    notificationRepository,
		logger:                    logger,
		rabbitMQ:                  rabbitMQ,
		postgres:                  postgres,
		redis:                     redis,
	}
}

func (s *EnsureStreakDaysIncrementToday) Execute(ctx context.Context, telegramID string) error {
	s.logger.Debug("[ensure streak days increment today] execute service")

	var (
		err                         error
		unlockAvailableAchievements []userachievement.UnlockAvailableAchievementsResponse
		notifications               []notification.Notification
		userExists                  bool
		userStatsExists             bool
		isStreakDaysIncrementToday  bool
		isUserPresence              bool
	)

	tx, err := s.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	// check user exists by telegram id.
	userExists, err = s.userRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return err
	}

	if !userExists { // if user does not exist.
		err = apperrors.ErrUserDoesNotExist
		return err
	}

	// check user stats exists by telegram id.
	userStatsExists, err = s.userStatsRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return err
	}

	if !userStatsExists { // if user stats does not exist.
		err = apperrors.ErrUserStatsDoesNotExist
		return err
	}

	// check has streak days increment today by telegram id.
	isStreakDaysIncrementToday, err = s.userStatsRepository.HasStreakDaysIncrementToday.Execute(ctx, tx, telegramID)
	if err != nil {
		return err
	}

	if !isStreakDaysIncrementToday { // has streak days is not increment today.
		// ensure streak days increment today.
		err = s.userStatsRepository.EnsureStreakDaysIncrementToday.Execute(ctx, tx, telegramID)
		if err != nil {
			return err
		}
	}

	// unlock available achievements.
	unlockAvailableAchievements, err = s.userAchievementRepository.UnlockAvailableAchievements.Execute(ctx, tx, telegramID)
	if err != nil {
		return err
	}

	if len(unlockAvailableAchievements) > 0 {
		// create notifications in database.
		notifications, err = s.createNotifications(ctx, tx, telegramID, unlockAvailableAchievements)
		if err != nil {
			return err
		}

		// check exists user is online for send notification with message broker.
		isUserPresence, err = s.redis.UserPresence.Exists(ctx, telegramID)
		if err != nil {
			return err
		}

		if isUserPresence {
			s.sendNotifications(ctx, notifications)
		}
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// createNotifications create notifications.
func (s *EnsureStreakDaysIncrementToday) createNotifications(
	ctx context.Context,
	tx pgx.Tx,
	telegramID string,
	unlockAvailableAchievements []userachievement.UnlockAvailableAchievementsResponse,
) ([]notification.Notification, error) {
	var dto []notification.CreateDTO

	if len(unlockAvailableAchievements) > 0 {
		for i := range unlockAvailableAchievements {
			dto = append(dto, notification.CreateDTO{
				Message: notification.Message{
					Title: "Уведомление",
					Text:  fmt.Sprintf("Поздравляем! Вы получили достижение «%s»! Так держать!", unlockAvailableAchievements[i].AchievementName),
				},
				Type:       notification.AchievementType,
				TelegramID: telegramID,
			})
		}
	}

	// create notifications.
	notifications, err := s.notificationRepository.CreateNotifications.Execute(ctx, tx, dto)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// sendNotifications send notifications.
func (s *EnsureStreakDaysIncrementToday) sendNotifications(ctx context.Context, notifications []notification.Notification) {
	for i := range notifications {
		data := notification.SendNotificationDTO{
			ID:         notifications[i].ID,
			Message:    notifications[i].Message,
			Type:       notifications[i].Type,
			TelegramID: notifications[i].TelegramID,
			CreatedAt:  notifications[i].CreatedAt,
		}

		// send notification in rabbitmq.
		if err := s.rabbitMQ.Notification.Publisher.Execute(ctx, data.TelegramID, data); err != nil {
			s.logger.Warn(fmt.Sprintf("failed to publish notification by rabbitmq: %v", err))
		}
	}
}
