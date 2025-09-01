package createevents

import (
	"context"
	"fmt"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/domain/event"
	eventtype "github.com/go-jedi/lingramm_backend/internal/domain/event_type"
	experiencepoint "github.com/go-jedi/lingramm_backend/internal/domain/experience_point"
	userbalance "github.com/go-jedi/lingramm_backend/internal/domain/internal_currency/user_balance"
	"github.com/go-jedi/lingramm_backend/internal/domain/level"
	"github.com/go-jedi/lingramm_backend/internal/domain/notification"
	userachievement "github.com/go-jedi/lingramm_backend/internal/domain/user_achievement"
	eventtyperepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/event_type"
	experiencepointrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/experience_point"
	internalcurrency "github.com/go-jedi/lingramm_backend/internal/repository/v1/internal_currency"
	levelrepository "github.com/go-jedi/lingramm_backend/internal/repository/v1/level"
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
	"github.com/shopspring/decimal"
)

//go:generate mockery --name=ICreateEvents --output=mocks --case=underscore
type ICreateEvents interface {
	Execute(ctx context.Context, dto event.CreateEventsDTO) error
}

type CreateEvents struct {
	experiencePointRepository  *experiencepointrepository.Repository
	userRepository             *userrepository.Repository
	userStatsRepository        *userstatsrepository.Repository
	eventTypeRepository        *eventtyperepository.Repository
	levelRepository            *levelrepository.Repository
	internalCurrencyRepository *internalcurrency.Repository
	userAchievementRepository  *userachievementrepository.Repository
	notificationRepository     *notificationrepository.Repository
	logger                     logger.ILogger
	rabbitMQ                   *rabbitmq.RabbitMQ
	postgres                   *postgres.Postgres
	redis                      *redis.Redis
}

func New(
	experiencePointRepository *experiencepointrepository.Repository,
	userRepository *userrepository.Repository,
	userStatsRepository *userstatsrepository.Repository,
	eventTypeRepository *eventtyperepository.Repository,
	levelRepository *levelrepository.Repository,
	internalCurrencyRepository *internalcurrency.Repository,
	userAchievementRepository *userachievementrepository.Repository,
	notificationRepository *notificationrepository.Repository,
	logger logger.ILogger,
	rabbitMQ *rabbitmq.RabbitMQ,
	postgres *postgres.Postgres,
	redis *redis.Redis,
) *CreateEvents {
	return &CreateEvents{
		experiencePointRepository:  experiencePointRepository,
		userRepository:             userRepository,
		userStatsRepository:        userStatsRepository,
		eventTypeRepository:        eventTypeRepository,
		levelRepository:            levelRepository,
		internalCurrencyRepository: internalCurrencyRepository,
		userAchievementRepository:  userAchievementRepository,
		notificationRepository:     notificationRepository,
		logger:                     logger,
		rabbitMQ:                   rabbitMQ,
		postgres:                   postgres,
		redis:                      redis,
	}
}

func (s *CreateEvents) Execute(ctx context.Context, dto event.CreateEventsDTO) error {
	s.logger.Debug("[create a new events] execute service")

	var (
		err                         error
		eventTypeData               eventtype.EventType
		backFillMissingLevelHistory level.BackFillMissingLevelHistoryByTelegramIDResponse
		unlockAvailableAchievements []userachievement.UnlockAvailableAchievementsResponse
		notifications               []notification.Notification
		isStreakDaysIncrementToday  bool
		isAccrualInternalCurrency   bool
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

	// check user exist by telegram id.
	err = s.checkUserExistByTelegramID(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	// check user stats exist by telegram id.
	err = s.checkUserStatsExistByTelegramID(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	// get event type data.
	eventTypeData, err = s.getEventTypeData(ctx, tx, dto.EventType)
	if err != nil {
		return err
	}

	// create a new xp events.
	err = s.createXPEvents(ctx, tx, dto.TelegramID, eventTypeData.Name, eventTypeData.XP)
	if err != nil {
		return err
	}

	// sync user stats from xp events by telegram id.
	err = s.userStatsRepository.SyncUserStatsFromXPEventsByTelegramID.Execute(ctx, tx, dto.TelegramID, dto.Actions)
	if err != nil {
		return err
	}

	// backfill missing level history by telegram id.
	backFillMissingLevelHistory, err = s.levelRepository.BackFillMissingLevelHistoryByTelegramID.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	// if amount is not nil and amount is positive number.
	if eventTypeData.Amount != nil {
		// check and accrual internal currency.
		isAccrualInternalCurrency, err = s.checkAndAccrualInternalCurrency(ctx, tx, dto.TelegramID, eventTypeData.ID, *eventTypeData.Amount, eventTypeData.Description)
		if err != nil {
			return err
		}
	}

	// check has streak days increment today by telegram id.
	isStreakDaysIncrementToday, err = s.userStatsRepository.HasStreakDaysIncrementToday.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	if !isStreakDaysIncrementToday { // has streak days is not increment today.
		// ensure streak days increment today.
		err = s.userStatsRepository.EnsureStreakDaysIncrementToday.Execute(ctx, tx, dto.TelegramID)
		if err != nil {
			return err
		}
	}

	// unlock available achievements.
	unlockAvailableAchievements, err = s.userAchievementRepository.UnlockAvailableAchievements.Execute(ctx, tx, dto.TelegramID)
	if err != nil {
		return err
	}

	// здесь будем проверять выполнил ли пользователь ежедневное задание.

	// create notifications in database.
	notifications, err = s.createNotifications(ctx, tx, dto.TelegramID, backFillMissingLevelHistory, unlockAvailableAchievements, isAccrualInternalCurrency)
	if err != nil {
		return err
	}

	// check exists user is online for send notification with message broker.
	isUserPresence, err = s.redis.UserPresence.Exists(ctx, dto.TelegramID)
	if err != nil {
		return err
	}

	if isUserPresence {
		s.sendNotifications(ctx, notifications)
	}

	// commit transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// checkUserExistByTelegramID check user exist by telegram id.
func (s *CreateEvents) checkUserExistByTelegramID(ctx context.Context, tx pgx.Tx, telegramID string) error {
	// check user exists by telegram id.
	ie, err := s.userRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return err
	}

	if !ie { // if user does not exist.
		return apperrors.ErrUserDoesNotExist
	}

	return nil
}

// checkUserStatsExistByTelegramID check user stats exist by telegram id.
func (s *CreateEvents) checkUserStatsExistByTelegramID(ctx context.Context, tx pgx.Tx, telegramID string) error {
	// check user stats exists by telegram id.
	ie, err := s.userStatsRepository.ExistsByTelegramID.Execute(ctx, tx, telegramID)
	if err != nil {
		return err
	}

	if !ie { // if user stats does not exist.
		return apperrors.ErrUserStatsDoesNotExist
	}

	return nil
}

// getEventTypeData get event type data.
func (s *CreateEvents) getEventTypeData(ctx context.Context, tx pgx.Tx, eventType string) (eventtype.EventType, error) {
	// check event type exist by name.
	ie, err := s.eventTypeRepository.ExistsByName.Execute(ctx, tx, eventType)
	if err != nil {
		return eventtype.EventType{}, err
	}

	if !ie { // if event type does not exist.
		return eventtype.EventType{}, apperrors.ErrEventTypeDoesNotExist
	}

	// get event type by name.
	eventTypeData, err := s.eventTypeRepository.GetByName.Execute(ctx, tx, eventType)
	if err != nil {
		return eventtype.EventType{}, err
	}

	return eventTypeData, nil
}

// createXPEvents create xp events.
func (s *CreateEvents) createXPEvents(
	ctx context.Context,
	tx pgx.Tx,
	telegramID string,
	eventType string,
	deltaXP int64,
) error {
	dto := experiencepoint.CreateXPEventDTO{
		TelegramID: telegramID,
		EventType:  eventType,
		DeltaXP:    deltaXP,
	}

	// create a new xp events.
	if err := s.experiencePointRepository.CreateXPEvents.Execute(ctx, tx, dto); err != nil {
		return err
	}

	return nil
}

// checkAndAccrualInternalCurrency check and accrual internal currency.
func (s *CreateEvents) checkAndAccrualInternalCurrency(
	ctx context.Context,
	tx pgx.Tx,
	telegramID string,
	eventTypeID int64,
	amount decimal.Decimal,
	description *string,
) (bool, error) {
	if !amount.IsPositive() {
		s.logger.Debug(fmt.Sprintf("skip add user balance because amount is not positive: %s", amount.String()))
		return false, nil
	}

	dto := userbalance.AddUserBalanceDTO{
		EventTypeID: eventTypeID,
		Amount:      amount,
		TelegramID:  telegramID,
		Description: description,
	}

	// add user balance.
	if _, err := s.internalCurrencyRepository.AddUserBalance.Execute(ctx, tx, dto); err != nil {
		return false, err
	}

	return true, nil
}

// createNotifications create notifications.
func (s *CreateEvents) createNotifications(
	ctx context.Context,
	tx pgx.Tx,
	telegramID string,
	backFillMissingLevelHistory level.BackFillMissingLevelHistoryByTelegramIDResponse,
	unlockAvailableAchievements []userachievement.UnlockAvailableAchievementsResponse,
	isAccrualInternalCurrency bool,
) ([]notification.Notification, error) {
	dto := []notification.CreateDTO{
		{
			Message: notification.Message{
				Title: "Уведомление",
				Text:  "Победа! Вы успешно выполнили мини-игру!",
			},
			Type:       notification.MiniGameType,
			TelegramID: telegramID,
		},
	}

	if backFillMissingLevelHistory.IsLevelUp {
		dto = append(dto, notification.CreateDTO{
			Message: notification.Message{
				Title: "Уведомление",
				Text:  fmt.Sprintf("Поздравляем! Вы перешли на %d уровень!", backFillMissingLevelHistory.NewLevel),
			},
			Type:       notification.LevelType,
			TelegramID: telegramID,
		})
	}

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

	if isAccrualInternalCurrency {
		dto = append(dto, notification.CreateDTO{
			Message: notification.Message{
				Title: "Уведомление",
				Text:  "Поздравляем! Баланс пополнен!",
			},
			Type:       notification.InternalCurrencyType,
			TelegramID: telegramID,
		})
	}

	// create notifications.
	notifications, err := s.notificationRepository.CreateNotifications.Execute(ctx, tx, dto)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// sendNotifications send notifications.
func (s *CreateEvents) sendNotifications(ctx context.Context, notifications []notification.Notification) {
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
